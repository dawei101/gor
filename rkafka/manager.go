package rkafka

import (
	"context"
	"github.com/dawei101/gor/rconfig"
	"sync"
)

type Manager struct {
	consumers map[string]*Consumer
	ctx       context.Context
	ctxCancel context.CancelFunc
	runOnce   sync.Once
}

type ManagerConf struct {
	Consumers map[string]ConsumerConf
}

var manager = &Manager{}
var managerConf = &ManagerConf{}

func init() {
	rconfig.DefConf().ValTo("kafka", managerConf)
	manager = New(managerConf)
}

func New(managerConf *ManagerConf) *Manager {
	var m = &Manager{
		consumers: map[string]*Consumer{},
	}
	m.ctx, m.ctxCancel = context.WithCancel(context.Background())

	for name, consumerConf := range managerConf.Consumers {
		m.consumers[name] = m.newConsumer(consumerConf)
	}

	return m
}

func Consume(name string, fn OnProcess) {
	manager.Consume(name, fn)
}

func (m *Manager) Consume(name string, fn OnProcess) {
	consumer, ok := manager.consumers[name]
	if ok {
		consumer.onProcess = fn
		consumer.runOnce.Do(func() {
			go consumer.Run()
		})
	}
}

func (m *Manager) newConsumer(conf ConsumerConf) *Consumer {
	var c = &Consumer{
		clientID: conf.ClientID,
		groupID:  conf.GroupID,
		topic:    conf.Topic,
		server:   conf.Server,
		limiter: Limiter{
			n: conf.Limit,
			c: make(chan struct{}, conf.Limit),
		},
	}
	c.wrap()
	c.ctx, c.ctxCancel = context.WithCancel(m.ctx)
	return c
}

func (m *Manager) Stop() {
	m.ctxCancel()
}
