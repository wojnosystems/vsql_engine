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
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine/context"
)

type BeginWare func(c context.Contexter, b vsql.QueryExecTransactioner) vsql.QueryExecTransactioner

// Middleware for begin
type BeginAdder interface {
	Add(w BeginWare)
}
type BeginApplyer interface {
	Apply(vsql.QueryExecTransactioner) vsql.QueryExecTransactioner
}

type Beginner interface {
	BeginMiddleware() BeginAdder
}


type BeginNestedWare func(c context.Contexter, b vsql.QueryExecNestedTransactioner) vsql.QueryExecNestedTransactioner

// Middleware for begin
type BeginNestedAdder interface {
	Add(w BeginNestedWare)
}
type BeginNestedApplyer interface {
	Apply(vsql.QueryExecNestedTransactioner) vsql.QueryExecNestedTransactioner
}

type BeginNester interface {
	BeginMiddleware() BeginNestedAdder
}


type StatementWare func(c context.Contexter, b vstmt.Statementer) vstmt.Statementer

type StatementAdder interface {
	Add(w StatementWare)
}
type StatementApplyer interface {
	Apply(vstmt.Statementer) vstmt.Statementer
}

type Statementer interface {
	StatementMiddleware() StatementAdder
}


type RowsWare func(c context.Contexter, b vrows.Rowser) vrows.Rowser

type RowsAdder interface {
	Add(w RowsWare)
}
type RowsApplyer interface {
	Apply(vrows.Rowser) vrows.Rowser
}

type Rowser interface {
	RowsMiddleware() RowsAdder
}


type RowWare func(c context.Contexter, b vrows.Rower) vrows.Rower

type RowAdder interface {
	Add(w RowWare)
}
type RowApplyer interface {
	Apply(vrows.Rower) vrows.Rower
}

type Rower interface {
	RowMiddleware() RowAdder
}


type ResultWare func(c context.Contexter, b vresult.Resulter) vresult.Resulter

type ResultAdder interface {
	Add(w ResultWare)
}
type ResultApplyer interface {
	Apply(vresult.Resulter) vresult.Resulter
}

type Resulter interface {
	ResultMiddleware() ResultAdder
}


type InsertResultWare func(c context.Contexter, b vresult.InsertResulter) vresult.InsertResulter

type InsertResultAdder interface {
	Add(w InsertResultWare)
}
type InsertResultApplyer interface {
	Apply(vresult.InsertResulter) vresult.InsertResulter
}

type InsertResulter interface {
	InsertResultMiddleware() InsertResultAdder
}
