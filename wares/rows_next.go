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

type RowsNextHandler func(ctx context.Context, c vsql_context.Rowser)

// Middleware for begin
type RowsNextAdder interface {
	Append(w RowsNextHandler)
	Prepend(w RowsNextHandler)
}

type RowsNextWare interface {
	RowsNextMW() RowsNextAdder
}

type RowsNextMW struct {
	RowsNextAdder
	middlewares *list.List // vsql_context.MiddlewareFunc
}

func NewRowsNextMW() *RowsNextMW {
	return &RowsNextMW{
		middlewares: list.New(),
	}
}

func (b *RowsNextMW) Append(w RowsNextHandler) {
	b.middlewares.PushBack(rowsNextPackageFunc(w))
}

func (b *RowsNextMW) Prepend(w RowsNextHandler) {
	b.middlewares.PushFront(rowsNextPackageFunc(w))
}

// PerformMiddleware executes the middleware after injecting vsql_context (if any)
func (b *RowsNextMW) PerformMiddleware(ctx context.Context, c vsql_context.Rowser) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(vsql_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b RowsNextMW) Copy() *RowsNextMW {
	r := NewRowsNextMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}

func rowsNextPackageFunc(w RowsNextHandler) vsql_context.MiddlewareFunc {
	return func(ctx context.Context, er vsql_context.Er) {
		w(ctx, er.(vsql_context.Rowser))
	}
}
