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
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql_engine/context"
)

type RowsWare func(c context.Contexter, b vrows.Rowser) vrows.Rowser

type RowsAdder interface {
	Add(w RowsWare)
}
type RowsApplyer interface {
	Apply(contexter context.Contexter, r vrows.Rowser) vrows.Rowser
}

type Rowser interface {
	RowsMiddleware() RowsAdder
}

type Rows struct {
	RowsAdder
	RowsApplyer
	base
}

func NewRows() *Rows {
	return &Rows{
		base: *newBase(),
	}
}

func (b *Rows) Add(w RowsWare) {
	b.base.Add(w)
}

func (b *Rows) Apply(in vrows.Rowser, ctx context.Contexter) vrows.Rowser {
	return b.ApplyBase(ctx, in, func(ctx context.Contexter, theMiddleware interface{}, sqlObject interface{}) (sqlObjectOut interface{}) {
		return theMiddleware.(RowsWare)(ctx, sqlObject.(vrows.Rowser))
	}).(vrows.Rowser)
}
