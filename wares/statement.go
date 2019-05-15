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
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine/context"
)

type StatementWare func(c context.Contexter, b vstmt.Statementer) vstmt.Statementer

type StatementAdder interface {
	Add(w StatementWare)
}
type StatementApplyer interface {
	Apply(ctx context.Contexter, s vstmt.Statementer) vstmt.Statementer
}

type Statementer interface {
	StatementMiddleware() StatementAdder
}

type Statement struct {
	StatementAdder
	StatementApplyer
	base
}

func NewStatement() *Statement {
	return &Statement{
		base: *newBase(),
	}
}

func (b *Statement) Add(w StatementWare) {
	b.base.Add(w)
}

func (b *Statement) Apply(in vstmt.Statementer, ctx context.Contexter) vstmt.Statementer {
	return b.ApplyBase(ctx, in, func(ctx context.Contexter, theMiddleware interface{}, sqlObject interface{}) (sqlObjectOut interface{}) {
		return theMiddleware.(StatementWare)(ctx, sqlObject.(vstmt.Statementer))
	}).(vstmt.Statementer)
}