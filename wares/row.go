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

type Row struct {
	RowAdder
	RowApplyer
	middlewares []RowWare
}

func NewRow() *Row {
	return &Row {
		middlewares: make([]RowWare, 0),
	}
}

func (b *Row) Add(w RowWare) {
	if b.middlewares == nil {
		b.middlewares = make([]RowWare, 0)
	}
	b.middlewares = append(b.middlewares, w)
}

func (b *Row) Apply(in vrows.Rower) vrows.Rower {
	if len(b.middlewares) == 0 {
		return in
	}
	ctx := context.New()
	for _, m := range b.middlewares {
		in = m(ctx, in)
	}
	return in
}
