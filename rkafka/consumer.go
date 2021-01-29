package rkafka

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"os"
)

type Consumer struct {
	id        uuid.UUID
	ctx       context.Context
	ctxCancel context.CancelFunc
	topics    []string

	*kafka.Consumer

	OnProcess OnProcess
	OnError   OnError
}

func (c *Consumer) Run() error {
	err := c.SubscribeTopics(c.topics, nil)
	fmt.Printf("Created Consumer %v\n", c)
	c.Listen()
	return err
}

func (c *Consumer) Stop() {
	c.ctxCancel()
}

func (c *Consumer) Listen() {
	for {
		select {
		case <-c.ctx.Done():
			c.ctxCancel()
			fmt.Printf("consumer %s already stop", c.String())
			return
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				go c.OnProcess(e, c.ctx)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}
