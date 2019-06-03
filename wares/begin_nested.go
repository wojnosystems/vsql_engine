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

type BeginNestedHandler func(ctx context.Context, c vsql_context.NestedBeginner)

// Middleware for begin
type BeginNestedAdder interface {
	Append(w BeginNestedHandler)
	Prepend(w BeginNestedHandler)
}

type BeginNestedWare interface {
	BeginNestedMW() BeginNestedAdder
}

type BeginNestedMW struct {
	BeginNestedAdder
	middlewares *list.List // vsql_context.BeginWareType
}

func NewBeginNestedMW() *BeginNestedMW {
	return &BeginNestedMW{
		middlewares: list.New(),
	}
}

func (b *BeginNestedMW) Append(w BeginNestedHandler) {
	b.middlewares.PushBack(beginNestedPackageFunc(w))
}

func (b *BeginNestedMW) Prepend(w BeginNestedHandler) {
	b.middlewares.PushFront(beginNestedPackageFunc(w))
}

// PerformMiddleware executes the middleware after injecting vsql_context (if any)
func (b *BeginNestedMW) PerformMiddleware(ctx context.Context, c vsql_context.NestedBeginner) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(vsql_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b BeginNestedMW) Copy() *BeginNestedMW {
	r := NewBeginNestedMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}

func beginNestedPackageFunc(w BeginNestedHandler) vsql_context.MiddlewareFunc {
	return func(ctx context.Context, er vsql_context.Er) {
		w(ctx, er.(vsql_context.NestedBeginner))
	}
}
