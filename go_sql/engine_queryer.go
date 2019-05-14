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

package go_sql

import (
	"context"
	"database/sql"
	"github.com/wojnosystems/vsql/param"
	"github.com/wojnosystems/vsql/vresult"
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine/strategy"
	"github.com/wojnosystems/vsql_engine/wares"
)

type engineQuery struct {
	db                   *sql.DB
	StatementWares       *wares.Statement
	RowsWares            *wares.Rows
	RowWares             *wares.Row
	ResultWares          *wares.Result
	InsertResultWares    *wares.InsertResult
	interpolationFactory strategy.InterpolationFactory
}

func newEngineQuery(db *sql.DB, interpolationFactory strategy.InterpolationFactory) *engineQuery {
	return &engineQuery{
		db:                   db,
		StatementWares:       wares.NewStatement(),
		RowsWares:            wares.NewRows(),
		RowWares:             wares.NewRow(),
		ResultWares:          wares.NewResult(),
		InsertResultWares:    wares.NewInsertResult(),
		interpolationFactory: interpolationFactory,
	}
}

// StatementMiddleware provides a way to add items to the Statement Middleware
func (m *engineQuery) StatementMiddleware() wares.StatementAdder {
	return m.StatementWares
}
// RowsMiddleware provides a way to add items to the RowsMiddleware
func (m *engineQuery) RowsMiddleware() wares.RowsAdder {
	return m.RowsWares
}
// RowMiddleware provides a way to add items to the RowMiddleware
func (m *engineQuery) RowMiddleware() wares.RowAdder {
	return m.RowWares
}
// ResultMiddleware provides a way to add items to the ResultMiddleware
func (m *engineQuery) ResultMiddleware() wares.ResultAdder {
	return m.ResultWares
}
// InsertResultMiddleware provides a way to add items to the InsertResultMiddleware
func (m *engineQuery) InsertResultMiddleware() wares.InsertResultAdder {
	return m.InsertResultWares
}

// Ping see github.com/wojnosystems/vsql/pinger/pinger.go#Pinger
func (m *engineQuery) Ping(ctx context.Context) error {
	return m.db.PingContext(ctx)
}

// Query see github.com/wojnosystems/vsql/vquery/query.go#Queryer
func (m *engineQuery) Query(ctx context.Context, query param.Queryer) (rRows vrows.Rowser, err error) {
	q, ps, err := query.Interpolate(m.interpolationFactory())
	if err != nil {
		return nil, err
	}
	r := &goSqlRows{
		queryEngineFactory: m,
	}
	r.SqlRows, err = m.db.QueryContext(ctx, q, ps...)
	if err != nil {
		m.RowsWares.Apply(r)
	}
	return r, err
}

// Insert see github.com/wojnosystems/vsql/vquery/query.go#Inserter
func (m *engineQuery) Insert(ctx context.Context, query param.Queryer) (res vresult.InsertResulter, err error) {
	q, ps, err := query.Interpolate(m.interpolationFactory())
	if err != nil {
		return nil, err
	}
	r := &goSqlQueryInsertResult{}
	r.sqlResult, err = m.db.ExecContext(ctx, q, ps...)
	if err != nil {
		m.InsertResultWares.Apply(r)
	}
	return r, err
}

// Exec see github.com/wojnosystems/vsql/vquery/query.go#Execer
func (m *engineQuery) Exec(ctx context.Context, query param.Queryer) (res vresult.Resulter, err error) {
	q, ps, err := query.Interpolate(m.interpolationFactory())
	if err != nil {
		return nil, err
	}
	r := &goSqlQueryResult{}
	r.sqlResult, err = m.db.ExecContext(ctx, q, ps...)
	if err != nil {
		m.ResultWares.Apply(r)
	}
	return r, err
}

// Prepare see github.com/wojnosystems/vsql/vstmt/statement.go#Preparer
func (m *engineQuery) Prepare(ctx context.Context, query param.Queryer) (stmtr vstmt.Statementer, err error) {
	q := query.SQLQuery(m.interpolationFactory())
	r := &goSqlStatement{
		queryEngineFactory: m,
	}
	r.stmt, err = m.db.PrepareContext(ctx, q)
	if err != nil {
		m.StatementWares.Apply(r)
	}
	return r, err
}

// Close returns the database connection back to the pool
func (m *engineQuery) Close() (err error) {
	return m.db.Close()
}
