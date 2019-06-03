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

package engine

import (
	"context"
	"github.com/wojnosystems/vsql"
	"github.com/wojnosystems/vsql/param"
	"github.com/wojnosystems/vsql/vresult"
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql/vtxn"
	"github.com/wojnosystems/vsql_engine/vsql_context"
	"github.com/wojnosystems/vsql_engine/wares"
)

//vsql.QueryExecNestedTransactioner
type nestedTx struct {
	beginnerNestedContext vsql_context.NestedBeginner
	queryEngineFactory    *engineNest
	beginNestedMW         *wares.BeginNestedMW
}

// Begin see github.com/wojnosystems/vsql/transactions.go#TransactionStarter
func (m *nestedTx) Begin(ctx context.Context, txOp vtxn.TxOptioner) (n vsql.QueryExecNestedTransactioner, err error) {
	c := vsql_context.NewNestedBeginner()
	c.SetQueryExecNestedTransactioner(m)
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetTxOptions(txOp)
	m.beginNestedMW.PerformMiddleware(ctx, c)
	s := &nestedTx{
		beginnerNestedContext: c,
		queryEngineFactory:    m.queryEngineFactory,
	}
	return s, c.Error()
}

// Commit see github.com/wojnosystems/vsql/transactions.go#Transactioner
func (m *nestedTx) Commit() error {
	c := vsql_context.NewBeginner()
	c.SetQueryExecTransactioner(m.beginnerNestedContext.QueryExecNestedTransactioner())
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	m.queryEngineFactory.commitMW.PerformMiddleware(nil, c)
	return c.Error()
}

// Rollback see github.com/wojnosystems/vsql/transactions.go#Transactioner
func (m *nestedTx) Rollback() error {
	c := vsql_context.NewBeginner()
	c.SetQueryExecTransactioner(m.beginnerNestedContext.QueryExecNestedTransactioner())
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	m.queryEngineFactory.rollbackMW.PerformMiddleware(nil, c)
	return c.Error()
}

// Query nonNestedTx github.com/wojnosystems/vsql/strategy.go#QueryExecer
func (m *nestedTx) Query(ctx context.Context, query param.Queryer) (rRows vrows.Rowser, err error) {
	c := vsql_context.NewQuery()
	c.SetQueryExecTransactioner(m)
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.queryMW.PerformMiddleware(ctx, c)
	return c.Rows(), c.Error()
}

// Insert see github.com/wojnosystems/vsql/strategy.go#QueryExecer
func (m *nestedTx) Insert(ctx context.Context, query param.Queryer) (res vresult.InsertResulter, err error) {
	c := vsql_context.NewInsertQuery()
	c.SetQueryExecTransactioner(m)
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.insertQueryMW.PerformMiddleware(ctx, c)
	return c.InsertResult(), c.Error()
}

// Exec see github.com/wojnosystems/vsql/strategy.go#QueryExecer
func (m *nestedTx) Exec(ctx context.Context, query param.Queryer) (res vresult.Resulter, err error) {
	c := vsql_context.NewExecQuery()
	c.SetQueryExecTransactioner(m)
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.execQueryMW.PerformMiddleware(ctx, c)
	return c.Result(), c.Error()
}

// Prepare see github.com/wojnosystems/vsql/vstmt/statements.go#Preparer
func (m *nestedTx) Prepare(ctx context.Context, query param.Queryer) (stmtr vstmt.Statementer, err error) {
	c := vsql_context.NewPreparer()
	c.SetQueryExecTransactioner(m)
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.prepareMW.PerformMiddleware(ctx, c)
	return c.Statement(), c.Error()
}