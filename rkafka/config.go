package rkafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

type OnProcess func(msg *kafka.Message, ctx context.Context) error

type OnError func(err error, ctx context.Context)

type OnLimiter func(c chan struct{})

type ConsumerConfig struct {
	ClientId  string
	ConfigMap kafka.ConfigMap

	OnProcess OnProcess
	OnError   OnError
	OnLimiter OnLimiter

	Topic string

	Limit int
}

func (c *ConsumerConfig) Wrap() {
	c.wrapValue("client.id", func(v kafka.ConfigValue) {
		clientIdRaw, ok := c.ConfigMap["client.id"]
		if !ok {
			clientId, ok := clientIdRaw.(string)
			if !ok {
				err := "The client.id must be a string"
				panic(err)
			}
			c.ClientId = clientId
		}
	})
}

func (c *ConsumerConfig) Defaults() {
	if c.Limit <= 0 {
		c.Limit = 1
	}

	c.defaultValue("client.id", func() kafka.ConfigValue {
		id, err := uuid.NewUUID()
		if err != nil {
			panic(err)
		}
		return id
	})

	c.defaultValue("bootstrap.servers", func() kafka.ConfigValue {
		return "127.0.0.1"
	})

	//default
	//ip
	//group.id
	//topic  user login

	//user_login
	//onProcess  user.log  map[string]func()

	//other
	//ip
	//client.id
	//group.id
	//onProcess2

	c.defaultValue("session.timeout.ms", func() kafka.ConfigValue {
		return 6000
	})
}

func (c *ConsumerConfig) defaultValue(key string, valueFunc func() kafka.ConfigValue) {
	_, ok := c.ConfigMap[key]
	if !ok {
		c.ConfigMap.SetKey(key, valueFunc())
	}
}

func (c *ConsumerConfig) wrapValue(key string, f func(v kafka.ConfigValue)) {
	v, ok := c.ConfigMap[key]
	if ok {
		f(v)
	}
}
