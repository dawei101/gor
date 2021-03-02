package rkafka

import (
	"context"
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

func (m *Manager) OnProcesses(name string, fn OnProcess) {
	consumer, ok := m.consumers[name]
	if ok {
		consumer.onProcess = fn
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

func (m *Manager) Run() {
	m.runOnce.Do(func() {
		for _, v := range m.consumers {
			go v.Run()
		}
	})
}

func (m *Manager) Stop() {
	m.ctxCancel()
}
