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

type CommitHandler func(ctx context.Context, c vsql_context.Beginner)

// Middleware for begin
type CommitAdder interface {
	Append(w CommitHandler)
	Prepend(w CommitHandler)
}

type CommitWare interface {
	CommitMW() CommitAdder
}

type CommitMW struct {
	middlewares *list.List // vsql_context.MiddlewareFunc
}

func NewCommitMW() *CommitMW {
	return &CommitMW{
		middlewares: list.New(),
	}
}

func (b *CommitMW) Append(w CommitHandler) {
	b.middlewares.PushBack(commitPackageFunc(w))
}

func (b *CommitMW) Prepend(w CommitHandler) {
	b.middlewares.PushFront(commitPackageFunc(w))
}

// PerformMiddleware executes the middleware after injecting vsql_context (if any)
func (b *CommitMW) PerformMiddleware(ctx context.Context, c vsql_context.Beginner) {
	if b.middlewares.Len() == 0 {
		return
	}
	c.(vsql_context.WithMiddlewarer).SetMiddlewares(b.middlewares)
	c.Next(ctx)
}

func (b CommitMW) Copy() *CommitMW {
	r := NewCommitMW()
	r.middlewares.PushBackList(b.middlewares)
	return r
}

func commitPackageFunc(w CommitHandler) vsql_context.MiddlewareFunc {
	return func(ctx context.Context, er vsql_context.Er) {
		w(ctx, er.(vsql_context.Beginner))
	}
}
