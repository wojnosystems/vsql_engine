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

type StatementCloseHandler func(ctx context.Context, c engine_context.StatementCloser)

// Middleware for begin
type StatementCloseAdder interface {
	Append(w StatementCloseHandler)
	Prepend(w StatementCloseHandler)
}

type StatementCloseWare interface {
	StatementCloseMW() StatementCloseAdder
}

type StatementCloseMW struct {
	middlewares *list.List // StatementCloseHandler
}

func NewStatementCloseMW() *StatementCloseMW {
	return &StatementCloseMW{
		middlewares: list.New(),
	}
}

func (b *StatementCloseMW) Append(w StatementCloseHandler) {
	b.middlewares.PushBack(statementClosePackageFunc(w))
}

func (b *StatementCloseMW) Prepend(w StatementCloseHandler) {
	b.middlewares.PushFront(statementClosePackageFunc(w))
}

// PerformMiddleware executes the middleware after injecting engine_context (if any)
func (b *StatementCloseMW) PerformMiddleware(ctx context.Context, c engine_context.StatementCloser) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(engine_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b StatementCloseMW) Copy() *StatementCloseMW {
	r := NewStatementCloseMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}

func statementClosePackageFunc(w StatementCloseHandler) engine_context.MiddlewareFunc {
	return func(ctx context.Context, er engine_context.Er) {
		w(ctx, er.(engine_context.StatementCloser))
	}
}
