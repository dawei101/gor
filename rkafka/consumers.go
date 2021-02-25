package rkafka

import (
	"sync"
)

type Consumers struct {
	regs map[string]*Consumer
	mu   sync.Mutex
}

func (c *Consumers) Set(consumer *Consumer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.regs[consumer.clientId] = consumer
}

func (c *Consumers) Get(clientId string) *Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, _ := c.regs[clientId]
	return v
}

func (c *Consumers) All() map[string]*Consumer {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.regs
}
