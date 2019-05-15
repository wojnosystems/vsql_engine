//Copyright 2019 Chris Wojno
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS
// OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
// OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package key_value

import (
	"fmt"
	"sync"
)

// KeyValuer is thread-safe
type KeyValuer interface {
	Set(key string, value interface{})
	Get(key string) (v interface{}, ok bool)
	Del(key string)
	MustGet(key string) (v interface{})
	// Copy creates a shallow clone of this KeyValuer so that it may be used in nested transactions in isolation
	Copy() KeyValuer
}

// InvocationContext is a key-value store for persisting keyValue throughout a middleware invocation chain
type keyValue struct {
	mu     *sync.RWMutex
	values map[string]interface{}
}

func New() *keyValue {
	return &keyValue{
		values: make(map[string]interface{}),
		mu:     &sync.RWMutex{},
	}
}

func (c *keyValue) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = value
}

func (c *keyValue) Get(key string) (v interface{}, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok = c.values[key]
	return
}

func (c *keyValue) Del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.values, key)
}

func (c *keyValue) MustGet(key string) (v interface{}) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var ok bool
	v, ok = c.values[key]
	if !ok {
		panic(fmt.Errorf(`"%s" was not found`, key))
	}
	return
}

func (c *keyValue) Copy() KeyValuer {
	c.mu.Lock()
	defer c.mu.Unlock()
	rc := New()
	for key, value := range c.values {
		rc.values[key] = value
	}
	return rc
}
