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

package engine_context

import (
	"container/list"
	"context"
	go_keyvaluer "github.com/wojnosystems/go_keyvaluer"
)

type Er interface {
	KeyValues() go_keyvaluer.KeyValuer
	// Next executes the next middleware in the chain until none are left
	SetError(err error)
	Error() error

	Next(ctx context.Context)
	Copy() Er
	middlewares() *list.List
}

type MiddlewareFunc func(ctx context.Context, er Er)

type WithMiddlewarer interface {
	Er
	SetMiddlewares(*list.List)
	// ShallowCopyFrom only copies the parts known to WithMiddlewarer, the rest of the configuration is up to the inheriting object
	// This copies a reference to kvo and middlewareV from the object passed to the argument and into the receiver.
	ShallowCopyFrom(WithMiddlewarer)
}

type contextBase struct {
	kvo go_keyvaluer.KeyValuer
	err error
	// middlewareV is always of type: MiddlewareFunc
	middlewareV       *list.List
	currentMiddleware *list.Element
}

func New() WithMiddlewarer {
	return newContextBase()
}

func newContextBase() *contextBase {
	return &contextBase{
		kvo: go_keyvaluer.New(),
	}
}

func (c *contextBase) KeyValues() go_keyvaluer.KeyValuer {
	return c.kvo
}

func (c *contextBase) Copy() Er {
	rc := &contextBase{}
	// kvo is thread-safe. This copy is just a reference copy to ensure that the new context can reference any values in that KVO
	rc.kvo = c.kvo
	rc.SetMiddlewares(cloneMiddlewareList(c.middlewareV))
	rc.err = nil
	return rc
}

func cloneMiddlewareList(l *list.List) (o *list.List) {
	o = list.New()
	if l != nil {
		o.PushBackList(l)
	}
	return
}

func (c *contextBase) ShallowCopyFrom(o WithMiddlewarer) {
	c.kvo = o.KeyValues()
	c.middlewareV = o.middlewares()
}

func (c *contextBase) SetError(err error) {
	c.err = err
}

func (c contextBase) Error() error {
	return c.err
}

func (c *contextBase) SetMiddlewares(m *list.List) {
	if m == nil {
		c.middlewareV = list.New()
	} else {
		c.middlewareV = m
	}
	c.currentMiddleware = c.middlewareV.Front()
}

func (c *contextBase) moveToNextMiddleware() {
	// Just starting a middleware chain, start with the first node
	if c.currentMiddleware != nil {
		c.currentMiddleware = c.currentMiddleware.Next()
	}
}

// Next runs the middleware, if any is available, null op if not. Next is only intended to be run once each middleware layer.
func (c *contextBase) Next(ctx context.Context) {
	if c.currentMiddleware != nil {
		cm := c.currentMiddleware
		c.moveToNextMiddleware()
		cm.Value.(MiddlewareFunc)(ctx, c)
	}
}

func (c contextBase) middlewares() *list.List {
	return c.middlewareV
}
