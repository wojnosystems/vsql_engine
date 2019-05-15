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
	"github.com/wojnosystems/vsql"
	"github.com/wojnosystems/vsql/vresult"
	"github.com/wojnosystems/vsql_engine/context"
)

type InsertResultWare func(c context.Contexter, b vresult.InsertResulter) vresult.InsertResulter

// Middleware for begin
type InsertResultAdder interface {
	Add(w InsertResultWare)
}
type InsertResultApplyer interface {
	Apply(context.Contexter, vsql.QueryExecTransactioner) vsql.QueryExecTransactioner
}

type InsertResulter interface {
	InsertResultMiddleware() InsertResultAdder
}

type InsertResult struct {
	InsertResultAdder
	InsertResultApplyer
	base
}

func NewInsertResult() *InsertResult {
	return &InsertResult{
		base: *newBase(),
	}
}

func (b *InsertResult) Add(w InsertResultWare) {
	b.base.Add(w)
}

func (b *InsertResult) Apply(in vresult.InsertResulter, ctx context.Contexter) vresult.InsertResulter {
	return b.ApplyBase(ctx, in, func(ctx context.Contexter, theMiddleware interface{}, sqlObject interface{}) (sqlObjectOut interface{}) {
		return theMiddleware.(InsertResultWare)(ctx, sqlObject.(vresult.InsertResulter))
	}).(vresult.InsertResulter)
}
