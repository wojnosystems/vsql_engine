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

type RowsCloseHandler func(ctx context.Context, c engine_context.Rowser)

// Middleware for begin
type RowsCloseAdder interface {
	Append(w RowsCloseHandler)
	Prepend(w RowsCloseHandler)
}

type RowsCloseWare interface {
	RowsCloseMW() RowsCloseAdder
}

type RowsCloseMW struct {
	middlewares *list.List // engine_context.MiddlewareFunc
}

func NewRowsCloseMW() *RowsCloseMW {
	return &RowsCloseMW{
		middlewares: list.New(),
	}
}

func (b *RowsCloseMW) Append(w RowsCloseHandler) {
	b.middlewares.PushBack(rowsClosePackageFunc(w))
}

func (b *RowsCloseMW) Prepend(w RowsCloseHandler) {
	b.middlewares.PushFront(rowsClosePackageFunc(w))
}

// PerformMiddleware executes the middleware after injecting engine_context (if any)
func (b *RowsCloseMW) PerformMiddleware(ctx context.Context, c engine_context.Rowser) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(engine_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b RowsCloseMW) Copy() *RowsCloseMW {
	r := NewRowsCloseMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}

func rowsClosePackageFunc(w RowsCloseHandler) engine_context.MiddlewareFunc {
	return func(ctx context.Context, er engine_context.Er) {
		w(ctx, er.(engine_context.Rowser))
	}
}
