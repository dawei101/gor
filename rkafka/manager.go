package rkafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"roo.bo/rlib"
	"sync"
)

type Manager struct {
	Consumers Consumers
	ctx       context.Context
	ctxCancel context.CancelFunc
	runOnce   sync.Once
	configs   map[string]*ConsumerConfig
}

func New() *Manager {
	m := &Manager{
		Consumers: Consumers{
			regs: map[string]*Consumer{},
		},
	}

	var configs interface{}
	rlib.LoadDefaultConfig("./config.yml")
	rlib.DefaultConfig().Value("kafka", &configs, nil)

	topics, ok := configs.(map[string]map[string]string)

	if ok {
		for topic, config := range topics {
			conf := m.NewConfig(config["groupId"], topic)
			m.configs[topic] = conf
		}
	}

	m.ctx, m.ctxCancel = context.WithCancel(context.Background())
	return m
}

func (m *Manager) OnProcess(topic string, f OnProcess) {
	if m, ok := m.configs[topic]; ok {
		m.OnProcess = f
	}
}

func (m *Manager) NewConfig(groupId kafka.ConfigValue, topic string) *ConsumerConfig {
	return &ConsumerConfig{
		ConfigMap: kafka.ConfigMap{
			"group.id": groupId,
		},

		Topic: topic,
		Limit: 10,
	}
}

func (m *Manager) NewConsumer(conf *ConsumerConfig) (err error) {
	conf.Defaults()
	conf.Wrap()
	rawC, err := kafka.NewConsumer(&conf.ConfigMap)
	if err != nil {
		return
	}
	var c = &Consumer{
		clientId: conf.ClientId,
		topic:    conf.Topic,
		limiter: Limiter{
			n: conf.Limit,
			c: make(chan struct{}, conf.Limit),
		},

		Consumer: rawC,

		OnProcess: conf.OnProcess,
		OnError:   conf.OnError,
		OnLimiter: conf.OnLimiter,
	}
	c.ctx, c.ctxCancel = context.WithCancel(m.ctx)
	m.Consumers.Set(c)
	return
}

func (m *Manager) Run() {
	m.runOnce.Do(func() {
		for _, conf := range m.configs {
			if conf.OnProcess != nil {
				m.NewConsumer(conf)
			}
		}
		for _, v := range m.Consumers.All() {
			go v.Run()
		}
	})
}
