package rkafka

import (
	"github.com/google/uuid"
	"sync"
)

type Consumers struct {
	regs     map[uuid.UUID]*Consumer
	running  map[uuid.UUID]*Consumer
	stopping map[uuid.UUID]*Consumer
	mu       sync.Mutex
}

func (c *Consumers) Set(consumer *Consumer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.regs[consumer.id] = consumer
}

func (c *Consumers) Get(id uuid.UUID) *Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, _ := c.regs[id]
	return v
}

func (c *Consumers) All() map[uuid.UUID]*Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.regs
}
