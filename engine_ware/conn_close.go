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

package engine_ware

import (
	"container/list"
	"context"
	"github.com/wojnosystems/vsql_engine/engine_context"
)

// Middleware for begin
type ConnCloseAdder interface {
	Append(w engine_context.MiddlewareFunc)
	Prepend(w engine_context.MiddlewareFunc)
}

type ConnCloseWare interface {
	ConnCloseMW() ConnCloseAdder
}

type ConnCloseMW struct {
	middlewares *list.List // engine_context.MiddlewareFunc
}

func NewConnCloseMW() *ConnCloseMW {
	return &ConnCloseMW{
		middlewares: list.New(),
	}
}

func (b *ConnCloseMW) Append(w engine_context.MiddlewareFunc) {
	b.middlewares.PushBack(w)
}

func (b *ConnCloseMW) Prepend(w engine_context.MiddlewareFunc) {
	b.middlewares.PushFront(w)
}

// PerformMiddleware executes the middleware after injecting engine_context (if any)
func (b *ConnCloseMW) PerformMiddleware(ctx context.Context, c engine_context.Er) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(engine_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b ConnCloseMW) Copy() *ConnCloseMW {
	r := NewConnCloseMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}
