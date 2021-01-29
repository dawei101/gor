package rkafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"sync"
)

type Manager struct {
	Consumers Consumers
	ctx       context.Context
	ctxCancel context.CancelFunc
	runOnce   sync.Once
	stopOnce  sync.Once
	state     sync.Mutex
}

func New() *Manager {
	m := &Manager{
		Consumers: Consumers{
			regs:     map[uuid.UUID]*Consumer{},
			running:  map[uuid.UUID]*Consumer{},
			stopping: map[uuid.UUID]*Consumer{},
		},
	}
	m.ctx, m.ctxCancel = context.WithCancel(context.Background())
	return m
}

func (m *Manager) NewConsumer(conf *ConsumerConfig) (err error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return
	}
	conf.ConfigMap.SetKey("client.id", id.String())
	rawC, err := kafka.NewConsumer(conf.ConfigMap)
	if err != nil {
		return
	}
	var c = &Consumer{
		id:     id,
		topics: conf.Topics,

		Consumer: rawC,

		OnProcess: conf.OnProcess,
		OnError:   conf.OnError,
	}
	c.ctx, c.ctxCancel = context.WithCancel(m.ctx)
	m.Consumers.Set(c)
	return
}

func (m *Manager) Destroy() {
	m.state.Lock()
	defer m.state.Unlock()
	m.stopOnce.Do(func() {
		m.ctxCancel()
		m.runOnce = sync.Once{}
	})
}

func (m *Manager) Run() {
	m.state.Lock()
	defer m.state.Unlock()
	m.runOnce.Do(func() {
		for _, v := range m.Consumers.All() {
			go func() {
				err := v.Run()
				if err != nil {

				}
			}()
		}
		m.stopOnce = sync.Once{}
	})
}
