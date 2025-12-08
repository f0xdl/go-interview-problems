package main

import (
	"fmt"
	"sync"
)

type Client interface {
	Get(address string) (string, error)
}

type Cache struct {
	client Client
	// cache  sync.Map - alternative
	cache      map[string]*Response
	processing map[string]*sync.WaitGroup
	mu         sync.Mutex
}

type Response struct {
	Value string
	Err   error
}

// Don't update signature of NewCache
func NewCache(client Client) *Cache {
	return &Cache{
		client:     client,
		mu:         sync.Mutex{},
		cache:      map[string]*Response{},
		processing: map[string]*sync.WaitGroup{},
	}
}

// Cache Client.Get result
func (c *Cache) Get(address string) (string, error) {
	c.mu.Lock()

	if cache, ok := c.cache[address]; ok {
		c.mu.Unlock()
		return cache.Value, cache.Err
	}

	if wg, ok := c.processing[address]; ok {
		c.mu.Unlock()
		wg.Wait()

		c.mu.Lock()
		defer c.mu.Unlock()
		cache := c.cache[address]
		return cache.Value, cache.Err
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	c.processing[address] = wg
	c.mu.Unlock()

	resp := c.processAddress(address)

	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[address] = resp
	delete(c.processing, address)
	wg.Done()
	return resp.Value, resp.Err
}

func (c *Cache) processAddress(address string) *Response {
	resp := &Response{}
	defer func() {
		if r := recover(); r != nil {
			resp.Err = fmt.Errorf("panic: %v", r)
		}
	}()
	resp.Value, resp.Err = c.client.Get(address)
	return resp
}
