package lock

import "sync"

// https://stackoverflow.com/a/10735763/8524395

type ThreadsafeVariable struct {
	mu       sync.Mutex
	variable interface{}
}

func (c *ThreadsafeVariable) Lock() {
	c.mu.Lock()
}

func (c *ThreadsafeVariable) Unlock() {
	c.mu.Unlock()
}

func (c *ThreadsafeVariable) Set(val interface{}) {
	c.variable = val
}

func (c *ThreadsafeVariable) Get() interface{} {
	return c.variable
}
