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
	"github.com/wojnosystems/vsql_engine/engine_context"
	"github.com/wojnosystems/vsql_engine/engine_ware"
)

// mysqlStatementTx is a prepared statement within the engine_context of a transaction
// sadly, based on the way the database/sql package works, most of this code has to be duplicated with the non-transaction version.
// however, this should save people implementing this code downstream much typing and testing. The amount of code written here is also fairly small due to the way the vsql package is composed so few changes, if any, will be required.
type txStatement struct {
	preparer           engine_context.Preparer
	queryEngineFactory *engineNest
	beginNestedMW      *engine_ware.BeginNestedMW
}

// Query see github.com/wojnosystems/vsql/vstmt/statements.go#Statementer
func (m *txStatement) Query(ctx context.Context, query vparam.Parameterer) (rRows vrows.Rowser, err error) {
	c := engine_context.NewStatementQuery()
	c.SetStatement(m.preparer.Statement())
	c.SetQuery(query)
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.statementQueryMW.PerformMiddleware(ctx, c)
	r := &rows{
		rows:               c.Rows(),
		queryEngineFactory: m.queryEngineFactory.engineQuery,
	}
	return r, c.Error()
}

// Insert see github.com/wojnosystems/vsql/vstmt/statements.go#Statementer
func (m *txStatement) Insert(ctx context.Context, query vparam.Parameterer) (res vresult.InsertResulter, err error) {
	c := engine_context.NewStatementInsertQuery()
	c.SetStatement(m.preparer.Statement())
	c.SetQuery(query)
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.statementInsertQueryMW.PerformMiddleware(ctx, c)
	return c.InsertResult(), c.Error()
}

// Exec see github.com/wojnosystems/vsql/vstmt/statements.go#Statementer
func (m *txStatement) Exec(ctx context.Context, query vparam.Parameterer) (res vresult.Resulter, err error) {
	c := engine_context.NewStatementExecQuery()
	c.SetStatement(m.preparer.Statement())
	c.SetQuery(query)
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	m.queryEngineFactory.statementExecQueryMW.PerformMiddleware(ctx, c)
	return c.Result(), c.Error()
}

// Close see github.com/wojnosystems/vsql/vstmt/statements.go#Statementer
func (m *txStatement) Close() error {
	c := engine_context.NewStatementClose()
	c.SetStatement(m.preparer.Statement())
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	m.queryEngineFactory.statementCloseMW.PerformMiddleware(nil, c)
	return c.Error()
}
