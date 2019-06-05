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

package vsql_engine

import (
	"context"
	"github.com/wojnosystems/vsql/vparam"
	"github.com/wojnosystems/vsql/vresult"
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine/engine_context"
)

type nonNestedTx struct {
	beginnerContext    engine_context.Beginner
	queryEngineFactory *engineQuery
}

// Commit see github.com/wojnosystems/vsql/transactions.go#Transactioner
func (m *nonNestedTx) Commit() error {
	c := engine_context.NewBeginner()
	c.SetQueryExecTransactioner(m.beginnerContext.QueryExecTransactioner())
	m.queryEngineFactory.commitMW.PerformMiddleware(nil, c)
	return c.Error()
}

// Rollback see github.com/wojnosystems/vsql/transactions.go#Transactioner
func (m *nonNestedTx) Rollback() error {
	c := engine_context.NewBeginner()
	c.SetQueryExecTransactioner(m.beginnerContext.QueryExecTransactioner())
	m.queryEngineFactory.rollbackMW.PerformMiddleware(nil, c)
	return c.Error()
}

// Query nonNestedTx github.com/wojnosystems/vsql/strategy.go#QueryExecer
func (m *nonNestedTx) Query(ctx context.Context, query vparam.Queryer) (rRows vrows.Rowser, err error) {
	c := engine_context.NewQuery()
	c.SetQueryExecTransactioner(m)
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.queryMW.PerformMiddleware(ctx, c)
	r := &rows{
		rows:               c.Rows(),
		queryEngineFactory: m.queryEngineFactory,
	}
	return r, c.Error()
}

// Insert see github.com/wojnosystems/vsql/strategy.go#QueryExecer
func (m *nonNestedTx) Insert(ctx context.Context, query vparam.Queryer) (res vresult.InsertResulter, err error) {
	c := engine_context.NewInsertQuery()
	c.SetQueryExecTransactioner(m)
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.insertQueryMW.PerformMiddleware(ctx, c)
	return c.InsertResult(), c.Error()
}

// Exec see github.com/wojnosystems/vsql/strategy.go#QueryExecer
func (m *nonNestedTx) Exec(ctx context.Context, query vparam.Queryer) (res vresult.Resulter, err error) {
	c := engine_context.NewExecQuery()
	c.SetQueryExecTransactioner(m)
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.execQueryMW.PerformMiddleware(ctx, c)
	return c.Result(), c.Error()
}

// Prepare see github.com/wojnosystems/vsql/vstmt/statements.go#Preparer
func (m *nonNestedTx) Prepare(ctx context.Context, query vparam.Queryer) (stmtr vstmt.Statementer, err error) {
	c := engine_context.NewPreparer()
	c.SetQueryExecTransactioner(m)
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.statementPrepareMW.PerformMiddleware(ctx, c)
	s := &statement{
		stmt:               c.Statement(),
		queryEngineFactory: m.queryEngineFactory,
	}
	return s, c.Error()
}
