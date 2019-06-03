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

type StatementExecQueryHandler func(ctx context.Context, c vsql_context.StatementExecQueryer)

// Middleware for begin
type StatementExecQueryAdder interface {
	Append(w StatementExecQueryHandler)
	Prepend(w StatementExecQueryHandler)
}

type StatementExecQueryWare interface {
	StatementExecQueryMW() StatementExecQueryAdder
}

type StatementExecQueryMW struct {
	StatementExecQueryAdder
	middlewares *list.List // StatementExecQueryHandler
}

func NewStatementExecQueryMW() *StatementExecQueryMW {
	return &StatementExecQueryMW{
		middlewares: list.New(),
	}
}

func (b *StatementExecQueryMW) Append(w StatementExecQueryHandler) {
	b.middlewares.PushBack(statementExecQueryPackageFunc(w))
}

func (b *StatementExecQueryMW) Prepend(w StatementExecQueryHandler) {
	b.middlewares.PushFront(statementExecQueryPackageFunc(w))
}

// PerformMiddleware executes the middleware after injecting vsql_context (if any)
func (b *StatementExecQueryMW) PerformMiddleware(ctx context.Context, c vsql_context.StatementExecQueryer) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(vsql_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b StatementExecQueryMW) Copy() *StatementExecQueryMW {
	r := NewStatementExecQueryMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}

func statementExecQueryPackageFunc(w StatementExecQueryHandler) vsql_context.MiddlewareFunc {
	return func(ctx context.Context, er vsql_context.Er) {
		w(ctx, er.(vsql_context.StatementExecQueryer))
	}
}
