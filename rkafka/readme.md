# rkafka

## Quick Start!

```go
package main
import (
	"bytes"
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/dawei101/gor/rkafka"
	"strconv"
	"strings"
	"time"
)

func main() {
	manager := rkafka.New()
	err := manager.NewConsumer(&rkafka.ConsumerConfig{
		ConfigMap: &kafka.ConfigMap{
			"bootstrap.servers":  "127.0.0.1",
			"group.id":           "test",
			"session.timeout.ms": 6000,
			//"enable.auto.commit": false,
		},
		OnProcess: func(msg *kafka.Message, ctx context.Context) error {
			fmt.Println(string(msg.Value), msg.TopicPartition)
			return nil
		},

		OnError: func(err error, ctx context.Context) {
			fmt.Printf("%+v", err)
		},
		Topics: []string{"testTopic"},
	})

	if err != nil {
	    println(err)
	}
	manager.Run()
}

```


