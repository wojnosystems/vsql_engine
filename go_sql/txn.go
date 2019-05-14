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
	"github.com/wojnosystems/vsql"
	"github.com/wojnosystems/vsql/param"
	"github.com/wojnosystems/vsql/vresult"
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql/vstmt"
)

type goSqlTx struct {
	vsql.QueryExecTransactioner
	queryEngineFactory *engineQuery
	tx                 *sql.Tx
}

// Commit see github.com/wojnosystems/vsql/transactions.go#Transactioner
func (m *goSqlTx) Commit() error {
	return m.tx.Commit()
}

// Rollback see github.com/wojnosystems/vsql/transactions.go#Transactioner
func (m *goSqlTx) Rollback() error {
	return m.tx.Rollback()
}

// Query goSqlTx github.com/wojnosystems/vsql/strategy.go#QueryExecer
func (m *goSqlTx) Query(ctx context.Context, query param.Queryer) (rRows vrows.Rowser, err error) {
	var q string
	var ps []interface{}
	r := &goSqlRows{
		queryEngineFactory: m.queryEngineFactory,
	}
	q, ps, err = query.Interpolate(m.queryEngineFactory.interpolationFactory())
	if err != nil {
		return r, err
	}
	r.SqlRows, err = m.tx.QueryContext(ctx, q, ps...)
	if err != nil {
		m.queryEngineFactory.RowsWares.Apply(r)
	}
	return r, err
}

// Insert see github.com/wojnosystems/vsql/strategy.go#QueryExecer
func (m *goSqlTx) Insert(ctx context.Context, query param.Queryer) (res vresult.InsertResulter, err error) {
	var q string
	var ps []interface{}
	r := &goSqlQueryInsertResult{}
	q, ps, err = query.Interpolate(m.queryEngineFactory.interpolationFactory())
	if err != nil {
		return r, err
	}
	r.sqlResult, err = m.tx.ExecContext(ctx, q, ps...)
	if err != nil {
		m.queryEngineFactory.InsertResultWares.Apply(r)
	}
	return r, err
}

// Exec see github.com/wojnosystems/vsql/strategy.go#QueryExecer
func (m *goSqlTx) Exec(ctx context.Context, query param.Queryer) (res vresult.Resulter, err error) {
	var q string
	var ps []interface{}
	r := &goSqlQueryResult{}
	q, ps, err = query.Interpolate(m.queryEngineFactory.interpolationFactory())
	if err != nil {
		return r, err
	}
	r.sqlResult, err = m.tx.ExecContext(ctx, q, ps...)
	if err != nil {
		m.queryEngineFactory.ResultWares.Apply(r)
	}
	return r, err
}

// Prepare see github.com/wojnosystems/vsql/vstmt/statements.go#Preparer
func (m *goSqlTx) Prepare(ctx context.Context, query param.Queryer) (stmtr vstmt.Statementer, err error) {
	q := query.SQLQuery(m.queryEngineFactory.interpolationFactory())
	r := &goSqlTxStatement{
		queryEngineFactory: m.queryEngineFactory,
		tx: m.tx,
	}
	r.stmt, err = m.tx.PrepareContext(ctx, q)
	if err != nil {
		m.queryEngineFactory.StatementWares.Apply(r)
	}
	return r, err
}