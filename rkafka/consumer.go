package rkafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/xid"
	"os"
	"sync"
)

type OnProcess func(msg *kafka.Message, ctx context.Context)

type OnError func(err error, ctx context.Context)

type Consumer struct {
	server   string
	clientID string
	groupID  string
	topic    string

	ctx       context.Context
	ctxCancel context.CancelFunc

	limiter Limiter

	consumer  *kafka.Consumer
	configMap kafka.ConfigMap

	onProcess OnProcess
	onError   OnError
	runOnce   sync.Once
}

type ConsumerConf struct {
	ClientID string `yaml:"clientID"`
	GroupID  string `yaml:"groupID"`
	Server   string `yaml:"server"`
	Topic    string `yaml:"topic"`
	Limit    int    `yaml:"limit"`
}

func (c *Consumer) Run() {
	if c.onProcess == nil {
		return
	}
	err := c.consumer.Subscribe(c.topic, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%% Error: %v \n", err)
		return
	}

	for {
		select {
		case <-c.ctx.Done():
			fmt.Printf("consumer `%s` already stop \n", c.consumer.String())
			return
		default:
			ev := c.consumer.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				c.limiter.Run(func() {
					c.onProcess(e, c.ctx)
				})
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}

func (c *Consumer) wrap() {
	if c.limiter.n <= 0 {
		c.limiter.n = 4
	}

	c.configMap = kafka.ConfigMap{
		"client.id": func() kafka.ConfigValue {
			if c.clientID == "" {
				c.clientID = xid.New().String()
			}
			return c.clientID
		}(),
		"bootstrap.servers": func() kafka.ConfigValue {
			if c.server == "" {
				c.server = "127.0.0.1"
			}
			return c.server
		}(),
		"session.timeout.ms": 6000,
		"group.id":           c.groupID,
	}

	consumerRaw, err := kafka.NewConsumer(&c.configMap)
	if err != nil {
		panic(err)
	}
	c.consumer = consumerRaw
}
