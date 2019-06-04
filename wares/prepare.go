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

type PrepareHandler func(ctx context.Context, c vsql_context.Preparer)

// Middleware for begin
type StatementPrepareAdder interface {
	Append(w PrepareHandler)
	Prepend(w PrepareHandler)
}

type StatementPrepareWare interface {
	StatementPrepareMW() StatementPrepareAdder
}

type StatementPrepareMW struct {
	middlewares *list.List // vsql_context.MiddlewareFunc
}

func NewStatementPrepareMW() *StatementPrepareMW {
	return &StatementPrepareMW{
		middlewares: list.New(),
	}
}

func (b *StatementPrepareMW) Append(w PrepareHandler) {
	b.middlewares.PushBack(preparePackageFunc(w))
}

func (b *StatementPrepareMW) Prepend(w PrepareHandler) {
	b.middlewares.PushFront(preparePackageFunc(w))
}

// PerformMiddleware executes the middleware after injecting vsql_context (if any)
func (b *StatementPrepareMW) PerformMiddleware(ctx context.Context, c vsql_context.Preparer) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(vsql_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b StatementPrepareMW) Copy() *StatementPrepareMW {
	r := NewStatementPrepareMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}

func preparePackageFunc(w PrepareHandler) vsql_context.MiddlewareFunc {
	return func(ctx context.Context, er vsql_context.Er) {
		w(ctx, er.(vsql_context.Preparer))
	}
}
