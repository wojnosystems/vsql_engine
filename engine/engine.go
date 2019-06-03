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
	"github.com/wojnosystems/vsql/param"
	"github.com/wojnosystems/vsql/vresult"
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine/vsql_context"
	"github.com/wojnosystems/vsql_engine/wares"
)

type engineQuery struct {
	queryMW                *wares.QueryMW
	insertQueryMW          *wares.InsertQueryMW
	execQueryMW            *wares.ExecMW
	pingMW                 *wares.PingMW
	rowsCloseMW            *wares.RowsCloseMW
	rowsNextMW             *wares.RowsNextMW
	connCloseMW            *wares.ConnCloseMW
	commitMW               *wares.CommitMW
	rollbackMW             *wares.RollbackMW
	prepareMW              *wares.PrepareMW
	statementCloseMW       *wares.StatementCloseMW
	statementQueryMW       *wares.StatementQueryMW
	statementInsertQueryMW *wares.StatementInsertQueryMW
	statementExecQueryMW   *wares.StatementExecQueryMW

	middlewareContext vsql_context.WithMiddlewarer
}

func newEngineQuery() *engineQuery {
	return &engineQuery{
		queryMW:                wares.NewQueryMW(),
		insertQueryMW:          wares.NewInsertQueryMW(),
		execQueryMW:            wares.NewExecMW(),
		pingMW:                 wares.NewPingMW(),
		rowsNextMW:             wares.NewRowsNextMW(),
		rowsCloseMW:            wares.NewRowsCloseMW(),
		connCloseMW:            wares.NewConnCloseMW(),
		commitMW:               wares.NewCommitMW(),
		rollbackMW:             wares.NewRollbackMW(),
		prepareMW:              wares.NewPrepareMW(),
		statementCloseMW:       wares.NewStatementCloseMW(),
		statementQueryMW:       wares.NewStatementQueryMW(),
		statementInsertQueryMW: wares.NewStatementInsertQueryMW(),
		statementExecQueryMW:   wares.NewStatementExecQueryMW(),

		middlewareContext: vsql_context.New(),
	}
}

// Group returns a copy of engineQuery with the same middlewares, but a cloned version of the vsql_context
// This allows you to mix and match middlewares using the same database connection pool, much like Go-Gin allows you to group routes and middleware.
func (m *engineQuery) Group() *engineQuery {
	rc := newEngineQuery()
	rc.queryMW = m.queryMW
	rc.insertQueryMW = m.insertQueryMW
	rc.execQueryMW = m.execQueryMW
	rc.pingMW = m.pingMW
	rc.rowsNextMW = m.rowsNextMW
	rc.rowsCloseMW = m.rowsCloseMW
	rc.connCloseMW = m.connCloseMW
	rc.commitMW = m.commitMW
	rc.rollbackMW = m.rollbackMW
	rc.statementCloseMW = m.statementCloseMW
	rc.statementQueryMW = m.statementQueryMW
	rc.statementInsertQueryMW = m.statementInsertQueryMW
	rc.statementExecQueryMW = m.statementExecQueryMW

	// OK to cast this as we KNOW it will be a context.WithMiddlewarer
	rc.middlewareContext = m.middlewareContext.Copy().(vsql_context.WithMiddlewarer)
	return rc
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) PrepareMW() wares.StatementPrepareAdder {
	return m.prepareMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) QueryMW() wares.QueryAdder {
	return m.queryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) InsertQueryMW() wares.InsertQueryAdder {
	return m.insertQueryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) ExecQueryMW() wares.ExecAdder {
	return m.execQueryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) PingMW() wares.PingAdder {
	return m.pingMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) StatementCloseMW() wares.StatementCloseAdder {
	return m.statementCloseMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) StatementQueryMW() wares.StatementQueryAdder {
	return m.statementQueryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) StatementInsertQueryMW() wares.StatementInsertQueryAdder {
	return m.statementInsertQueryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) StatementExecQueryMW() wares.StatementExecQueryAdder {
	return m.statementExecQueryMW
}

// CommitMiddleware provides a way to add items to the commitMW
func (m *engineQuery) CommitMW() wares.CommitAdder {
	return m.commitMW
}

// RollbackMiddleware provides a way to add items to the rollbackMW
func (m *engineQuery) RollbackMW() wares.RollbackAdder {
	return m.rollbackMW
}

// RowsNextMW provides a way to add items to the Connection NoNestCloser
func (m *engineQuery) RowsNextMW() wares.RowsNextAdder {
	return m.rowsNextMW
}

// ConnCloseMW provides a way to add items to the Connection NoNestCloser
func (m *engineQuery) RowsCloseMW() wares.RowsCloseAdder {
	return m.rowsCloseMW
}

// ConnCloseMW provides a way to add items to the Connection NoNestCloser
func (m *engineQuery) ConnCloseMW() wares.ConnCloseAdder {
	return m.connCloseMW
}

// Ping see github.com/wojnosystems/vsql/pinger/pinger.go#Pinger
func (m *engineQuery) Ping(ctx context.Context) error {
	c := m.middlewareContext.Copy().(vsql_context.WithMiddlewarer)
	m.pingMW.PerformMiddleware(ctx, c)
	return c.Error()
}

// Query see github.com/wojnosystems/vsql/vquery/queryer.go#Queryer
func (m *engineQuery) Query(ctx context.Context, query param.Queryer) (rRows vrows.Rowser, err error) {
	c := vsql_context.NewQuery()
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetQuery(query)
	m.queryMW.PerformMiddleware(ctx, c)
	r := &rows{
		rows:               c.Rows(),
		queryEngineFactory: m,
	}
	return r, c.Error()
}

// Insert see github.com/wojnosystems/vsql/vquery/queryer.go#Inserter
func (m *engineQuery) Insert(ctx context.Context, query param.Queryer) (res vresult.InsertResulter, err error) {
	c := vsql_context.NewInsertQuery()
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetQuery(query)
	m.insertQueryMW.PerformMiddleware(ctx, c)
	return c.InsertResult(), c.Error()
}

// Exec see github.com/wojnosystems/vsql/vquery/queryer.go#Execer
func (m *engineQuery) Exec(ctx context.Context, query param.Queryer) (res vresult.Resulter, err error) {
	c := vsql_context.NewExecQuery()
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetQuery(query)
	m.execQueryMW.PerformMiddleware(ctx, c)
	return c.Result(), c.Error()
}

// Prepare see github.com/wojnosystems/vsql/vstmt/statement.go#Preparer
func (m *engineQuery) Prepare(ctx context.Context, query param.Queryer) (stmtr vstmt.Statementer, err error) {
	c := vsql_context.NewPreparer()
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetQuery(query)
	m.prepareMW.PerformMiddleware(ctx, c)
	s := &statement{
		stmt:               c.Statement(),
		queryEngineFactory: m,
	}
	return s, c.Error()
}

// Close returns the database connection back to the pool
func (m *engineQuery) Close() (err error) {
	c := vsql_context.New()
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	m.connCloseMW.PerformMiddleware(nil, c)
	return c.Error()
}
