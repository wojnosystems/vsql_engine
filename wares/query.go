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

package wares

import (
	"container/list"
	"context"
	"github.com/wojnosystems/vsql_engine/vsql_context"
)

type QueryHandler func(ctx context.Context, c vsql_context.Queryer)

// Middleware for begin
type QueryAdder interface {
	Append(w QueryHandler)
	Prepend(w QueryHandler)
}

type QueryWare interface {
	QueryMW() QueryAdder
}

type QueryMW struct {
	QueryAdder
	middlewares *list.List // vsql_context.MiddlewareFunc
}

func NewQueryMW() *QueryMW {
	return &QueryMW{
		middlewares: list.New(),
	}
}

func (b *QueryMW) Append(w QueryHandler) {
	b.middlewares.PushBack(queryPackageFunc(w))
}

func (b *QueryMW) Prepend(w QueryHandler) {
	b.middlewares.PushFront(queryPackageFunc(w))
}

// PerformMiddleware executes the middleware after injecting vsql_context (if any)
func (b *QueryMW) PerformMiddleware(ctx context.Context, c vsql_context.Queryer) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(vsql_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b QueryMW) Copy() *QueryMW {
	r := NewQueryMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}

func queryPackageFunc(w QueryHandler) vsql_context.MiddlewareFunc {
	return func(ctx context.Context, er vsql_context.Er) {
		w(ctx, er.(vsql_context.Queryer))
	}
}
