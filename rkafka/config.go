package rkafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

//var consumerDefault = map[string]interface{}{
//	"bootstrap.servers":  "127.0.0.1",
//	"group.id":           "0",
//	"session.timeout.ms": 6000,
//	//"enable.auto.commit": false,
//}

type OnProcess func(msg *kafka.Message, ctx context.Context) error
type OnError func(err error, ctx context.Context)

type ConsumerConfig struct {
	ConfigMap *kafka.ConfigMap
	OnProcess OnProcess
	OnError   OnError
	Topics    []string
}
