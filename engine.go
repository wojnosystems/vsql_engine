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
	"github.com/wojnosystems/vsql_engine/engine_ware"
)

type engineQuery struct {
	queryMW                *engine_ware.QueryMW
	insertQueryMW          *engine_ware.InsertQueryMW
	execQueryMW            *engine_ware.ExecMW
	pingMW                 *engine_ware.PingMW
	rowsCloseMW            *engine_ware.RowsCloseMW
	rowsNextMW             *engine_ware.RowsNextMW
	connCloseMW            *engine_ware.ConnCloseMW
	commitMW               *engine_ware.CommitMW
	rollbackMW             *engine_ware.RollbackMW
	statementPrepareMW     *engine_ware.StatementPrepareMW
	statementCloseMW       *engine_ware.StatementCloseMW
	statementQueryMW       *engine_ware.StatementQueryMW
	statementInsertQueryMW *engine_ware.StatementInsertQueryMW
	statementExecQueryMW   *engine_ware.StatementExecQueryMW

	middlewareContext engine_context.WithMiddlewarer
}

func newEngineQuery() *engineQuery {
	return &engineQuery{
		queryMW:                engine_ware.NewQueryMW(),
		insertQueryMW:          engine_ware.NewInsertQueryMW(),
		execQueryMW:            engine_ware.NewExecMW(),
		pingMW:                 engine_ware.NewPingMW(),
		rowsNextMW:             engine_ware.NewRowsNextMW(),
		rowsCloseMW:            engine_ware.NewRowsCloseMW(),
		connCloseMW:            engine_ware.NewConnCloseMW(),
		commitMW:               engine_ware.NewCommitMW(),
		rollbackMW:             engine_ware.NewRollbackMW(),
		statementPrepareMW:     engine_ware.NewStatementPrepareMW(),
		statementCloseMW:       engine_ware.NewStatementCloseMW(),
		statementQueryMW:       engine_ware.NewStatementQueryMW(),
		statementInsertQueryMW: engine_ware.NewStatementInsertQueryMW(),
		statementExecQueryMW:   engine_ware.NewStatementExecQueryMW(),

		middlewareContext: engine_context.New(),
	}
}

// Group returns a copy of engineQuery with the same middlewares, but a cloned version of the engine_context
// This allows you to mix and match middlewares using the same database connection pool, much like Go-Gin allows you to group routes and middleware.
func (m *engineQuery) Group() *engineQuery {
	rc := newEngineQuery()
	rc.queryMW = m.queryMW.Copy()
	rc.insertQueryMW = m.insertQueryMW.Copy()
	rc.execQueryMW = m.execQueryMW.Copy()
	rc.pingMW = m.pingMW.Copy()
	rc.rowsNextMW = m.rowsNextMW.Copy()
	rc.rowsCloseMW = m.rowsCloseMW.Copy()
	rc.connCloseMW = m.connCloseMW.Copy()
	rc.commitMW = m.commitMW.Copy()
	rc.rollbackMW = m.rollbackMW.Copy()
	rc.statementPrepareMW = m.statementPrepareMW.Copy()
	rc.statementCloseMW = m.statementCloseMW.Copy()
	rc.statementQueryMW = m.statementQueryMW.Copy()
	rc.statementInsertQueryMW = m.statementInsertQueryMW.Copy()
	rc.statementExecQueryMW = m.statementExecQueryMW.Copy()

	// OK to cast this as we KNOW it will be a context.WithMiddlewarer
	rc.middlewareContext = m.middlewareContext.Copy().(engine_context.WithMiddlewarer)
	return rc
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) StatementPrepareMW() engine_ware.StatementPrepareAdder {
	return m.statementPrepareMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) QueryMW() engine_ware.QueryAdder {
	return m.queryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) InsertQueryMW() engine_ware.InsertQueryAdder {
	return m.insertQueryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) ExecQueryMW() engine_ware.ExecAdder {
	return m.execQueryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) PingMW() engine_ware.PingAdder {
	return m.pingMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) StatementCloseMW() engine_ware.StatementCloseAdder {
	return m.statementCloseMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) StatementQueryMW() engine_ware.StatementQueryAdder {
	return m.statementQueryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) StatementInsertQueryMW() engine_ware.StatementInsertQueryAdder {
	return m.statementInsertQueryMW
}

// StatementCloseMiddleware provides a way to add items to the StatementCloseWares
func (m *engineQuery) StatementExecQueryMW() engine_ware.StatementExecQueryAdder {
	return m.statementExecQueryMW
}

// CommitMiddleware provides a way to add items to the commitMW
func (m *engineQuery) CommitMW() engine_ware.CommitAdder {
	return m.commitMW
}

// RollbackMiddleware provides a way to add items to the rollbackMW
func (m *engineQuery) RollbackMW() engine_ware.RollbackAdder {
	return m.rollbackMW
}

// RowsNextMW provides a way to add items to the Connection NoNestCloser
func (m *engineQuery) RowsNextMW() engine_ware.RowsNextAdder {
	return m.rowsNextMW
}

// ConnCloseMW provides a way to add items to the Connection NoNestCloser
func (m *engineQuery) RowsCloseMW() engine_ware.RowsCloseAdder {
	return m.rowsCloseMW
}

// ConnCloseMW provides a way to add items to the Connection NoNestCloser
func (m *engineQuery) ConnCloseMW() engine_ware.ConnCloseAdder {
	return m.connCloseMW
}

// Ping see github.com/wojnosystems/vsql/pinger/pinger.go#Pinger
func (m *engineQuery) Ping(ctx context.Context) error {
	c := m.middlewareContext.Copy().(engine_context.WithMiddlewarer)
	m.pingMW.PerformMiddleware(ctx, c)
	return c.Error()
}

// Query see github.com/wojnosystems/vsql/vquery/queryer.go#Queryer
func (m *engineQuery) Query(ctx context.Context, query vparam.Queryer) (rRows vrows.Rowser, err error) {
	c := engine_context.NewQuery()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetQuery(query)
	m.queryMW.PerformMiddleware(ctx, c)
	r := &rows{
		rows:               c.Rows(),
		queryEngineFactory: m,
	}
	return r, c.Error()
}

// Insert see github.com/wojnosystems/vsql/vquery/queryer.go#Inserter
func (m *engineQuery) Insert(ctx context.Context, query vparam.Queryer) (res vresult.InsertResulter, err error) {
	c := engine_context.NewInsertQuery()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetQuery(query)
	m.insertQueryMW.PerformMiddleware(ctx, c)
	return c.InsertResult(), c.Error()
}

// Exec see github.com/wojnosystems/vsql/vquery/queryer.go#Execer
func (m *engineQuery) Exec(ctx context.Context, query vparam.Queryer) (res vresult.Resulter, err error) {
	c := engine_context.NewExecQuery()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetQuery(query)
	m.execQueryMW.PerformMiddleware(ctx, c)
	return c.Result(), c.Error()
}

// Prepare see github.com/wojnosystems/vsql/vstmt/statement.go#Preparer
func (m *engineQuery) Prepare(ctx context.Context, query vparam.Queryer) (stmtr vstmt.Statementer, err error) {
	c := engine_context.NewPreparer()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetQuery(query)
	m.statementPrepareMW.PerformMiddleware(ctx, c)
	s := &statement{
		stmt:               c.Statement(),
		queryEngineFactory: m,
	}
	return s, c.Error()
}

// Close returns the database connection back to the pool
func (m *engineQuery) Close() (err error) {
	c := engine_context.New()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	m.connCloseMW.PerformMiddleware(nil, c)
	return c.Error()
}
