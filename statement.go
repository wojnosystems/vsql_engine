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
	"github.com/wojnosystems/vsql/param"
	"github.com/wojnosystems/vsql/vresult"
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine/engine_context"
)

// mysqlStatement is a representation of a prepared statement that is NOT running in the engine_context of a transaction
type statement struct {
	stmt               vstmt.Statementer
	queryEngineFactory *engineQuery
}

// Query see github.com/wojnosystems/vsql/vstmt/statements.go#Statementer
func (m *statement) Query(ctx context.Context, query param.Parameterer) (rRows vrows.Rowser, err error) {
	c := engine_context.NewStatementQuery()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	c.SetStatement(m.stmt)
	m.queryEngineFactory.statementQueryMW.PerformMiddleware(ctx, c)
	r := &rows{
		rows:               c.Rows(),
		queryEngineFactory: m.queryEngineFactory,
	}
	return r, c.Error()
}

// Insert see github.com/wojnosystems/vsql/vstmt/statements.go#Statementer
func (m *statement) Insert(ctx context.Context, query param.Parameterer) (res vresult.InsertResulter, err error) {
	c := engine_context.NewStatementInsertQuery()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	c.SetStatement(m.stmt)
	m.queryEngineFactory.statementInsertQueryMW.PerformMiddleware(ctx, c)
	return c.InsertResult(), c.Error()
}

// Exec see github.com/wojnosystems/vsql/vstmt/statements.go#Statementer
func (m *statement) Exec(ctx context.Context, query param.Parameterer) (res vresult.Resulter, err error) {
	c := engine_context.NewStatementExecQuery()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetQuery(query)
	c.SetStatement(m.stmt)
	m.queryEngineFactory.statementExecQueryMW.PerformMiddleware(ctx, c)
	return c.Result(), c.Error()
}

// Close see github.com/wojnosystems/vsql/vstmt/statements.go#Statementer
func (m *statement) Close() error {
	c := engine_context.NewStatementClose()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetStatement(m.stmt)
	m.queryEngineFactory.statementCloseMW.PerformMiddleware(context.Background(), c)
	return c.Error()
}
