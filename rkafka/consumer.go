package rkafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

type Consumer struct {
	clientId  string
	ctx       context.Context
	ctxCancel context.CancelFunc
	topic     string
	limiter   Limiter
	*kafka.Consumer

	OnProcess OnProcess
	OnError   OnError
	OnLimiter OnLimiter
}

func (c *Consumer) Run() {
	err := c.Subscribe(c.topic, nil)
	if err != nil {
		fmt.Printf("Consumer: %v create fail. \n", c)
		return
	}
	c.Listen()

	fmt.Printf("Consumer: %v create success. \n", c)
}

func (c *Consumer) Stop() {
	c.ctxCancel()
}

func (c *Consumer) Listen() {
	c.limiter.Watch(c.OnLimiter)

	for {
		select {
		case <-c.ctx.Done():
			c.ctxCancel()
			fmt.Printf("consumer %s already stop. \n", c.String())
			return
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				c.limiter.Run(func() {
					c.OnProcess(e, c.ctx)

				})
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}
