package badempty

import "fmt"

type cacheImpl struct {
	data map[string]interface{}
}

func (c *cacheImpl) Get(key string) interface{} {
	return c.data[key]
}

func (c *cacheImpl) Set(key string, value interface{}) {
	c.data[key] = value
}

// newCache retourne l'implémentation au lieu de l'interface
func newCache() *cacheImpl {
	return &cacheImpl{
		data: make(map[string]interface{}),
	}
}

type processorImpl struct{}

func (p *processorImpl) Process(data interface{}) interface{} {
	fmt.Println(data)
	return data
}

func newProcessor() *processorImpl {
	return &processorImpl{}
}
