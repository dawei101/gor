package rkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type KafkaMessage struct {
	Action   string                 `json:"action"`
	ClientId string                 `json:"clientId"`
	Data     map[string]interface{} `json:"data"`
}

func (m *KafkaMessage) DataValue(key string) interface{} {
	return m.Data[key]
}

func (m *KafkaMessage) DataAssignTo(data interface{}) {
}

func (m *KafkaMessage) SetData(data map[string]interface{}) {
	m.Data = data
}

func KafkaConsume(name string, consume func(context.Context, *kafka.Message)) {
	if kconfig, ok := DefaultRooboConfig().KafkaConsumer[name]; ok {
		c, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": kconfig.Servers,
			"group.id":          kconfig.Group,
			"auto.offset.reset": "earliest",
		})
		if err != nil {
			msg := fmt.Sprintf("kafka init (name=%s) failed:%v", name, err)
			Error(context.Background(), msg)
			panic(msg)
		}
		c.SubscribeTopics(kconfig.Topics, nil)
		for {
			msg, err := c.ReadMessage(-1)
			ctx := contextWithRequestId(context.Background())
			if err == nil {
				topic := ""
				if msg.TopicPartition.Topic != nil {
					topic = *msg.TopicPartition.Topic
				}
				partition := msg.TopicPartition.Partition
				offset := msg.TopicPartition.Offset
				delayms := MicroTimestamp(time.Now()) - MicroTimestamp(msg.Timestamp)

				m := &KafkaMessage{}
				err = json.Unmarshal(msg.Value, m)
				logstr := fmt.Sprintf("[KAFKA-CONSUMER][%s] topic(%v) partition(%d) offset(%d) delayms(%d) message(%s), error(%v)", name, topic, partition, offset, delayms, string(msg.Value), err)
				reqLog.Info(ctx, logstr)
				if err == nil {
					go consume(ctx, msg)
				}
			} else {
				logstr := fmt.Sprintf("[KAFKA-CONSUMER] no message get, error(%v)", err)
				reqLog.Info(ctx, logstr)
			}
		}
		c.Close()
	} else {
		msg := "no kafka config named:" + name
		Error(context.Background(), msg)
		panic("no kafka config named:" + name)
	}
}

var k_producers map[string]*kafka.Producer
var k_producers_lock sync.RWMutex

func KafkaProduce(name string, msg *KafkaMessage) {
	var p *kafka.Producer
	p = KafkaProducerGet(name)
	msgd, _ := json.Marshal(msg)
	config, _ := DefaultRooboConfig().KafkaProducer[name]
	var done chan kafka.Event
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &config.Topic, Partition: kafka.PartitionAny},
		Value:          msgd,
	}, done)
	<-done
}

func KafkaProducerStop(name string) {
	p := KafkaProducerGet(name)
	p.Flush(5 * 1000)
	k_producers_lock.Lock()
	defer k_producers_lock.Unlock()
	delete(k_producers, name)
}

func KafkaProducerGet(name string) *kafka.Producer {
	k_producers_lock.RLock()
	p, ok := k_producers[name]
	k_producers_lock.RUnlock()
	if ok {
		return p
	}

	if kconfig, ok := DefaultRooboConfig().KafkaProducer[name]; ok {
		p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kconfig.Servers})
		if err != nil {
			panic(fmt.Sprintf("kafka init failed:%v", err))
		}
		k_producers_lock.Lock()
		k_producers[name] = p
		k_producers_lock.Unlock()
	} else {
		panic("no kafka config named:" + name)
	}
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return p
}
