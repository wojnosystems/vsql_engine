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

type StatementQueryHandler func(ctx context.Context, c engine_context.StatementQueryer)

// Middleware for begin
type StatementQueryAdder interface {
	Append(w StatementQueryHandler)
	Prepend(w StatementQueryHandler)
}

type StatementQueryWare interface {
	StatementQueryMW() StatementQueryAdder
}

type StatementQueryMW struct {
	middlewares *list.List // StatementQueryHandler
}

func NewStatementQueryMW() *StatementQueryMW {
	return &StatementQueryMW{
		middlewares: list.New(),
	}
}

func (b *StatementQueryMW) Append(w StatementQueryHandler) {
	b.middlewares.PushBack(statementQueryPackageFunc(w))
}

func (b *StatementQueryMW) Prepend(w StatementQueryHandler) {
	b.middlewares.PushFront(statementQueryPackageFunc(w))
}

// PerformMiddleware executes the middleware after injecting engine_context (if any)
func (b *StatementQueryMW) PerformMiddleware(ctx context.Context, c engine_context.StatementQueryer) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(engine_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b StatementQueryMW) Copy() *StatementQueryMW {
	r := NewStatementQueryMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}

func statementQueryPackageFunc(w StatementQueryHandler) engine_context.MiddlewareFunc {
	return func(ctx context.Context, er engine_context.Er) {
		w(ctx, er.(engine_context.StatementQueryer))
	}
}
