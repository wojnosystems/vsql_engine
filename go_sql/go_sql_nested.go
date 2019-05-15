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
	"github.com/wojnosystems/vsql/vtxn"
	"github.com/wojnosystems/vsql_engine/strategy"
	"github.com/wojnosystems/vsql_engine/wares"
)

// SQLNestedEnginer
type goSqlNested struct {
	*engineQuery
	beginEngineFactory *wares.BeginNested
}

func NewNested(interpolationFactory strategy.InterpolationFactory, driverFactory func() (db *sql.DB)) SQLNestedEnginer {
	db := driverFactory()
	if db == nil {
		return nil
	}
	m := &goSqlNested{
		engineQuery:        newEngineQuery(db, interpolationFactory),
		beginEngineFactory: wares.NewBeginNested(),
	}
	return m
}

// BeginMiddleware provides a way to add items to the BeginMiddleware
func (m *goSqlNested) BeginMiddleware() wares.BeginNestedAdder {
	return m.beginEngineFactory
}

// Begin see github.com/wojnosystems/vsql/transactions.go#TransactionStarter
func (m *goSqlNested) Begin(ctx context.Context, txOp vtxn.TxOptioner) (n vsql.QueryExecNestedTransactioner, err error) {
	var txo *sql.TxOptions
	if txOp != nil {
		txo = txOp.ToTxOptions()
	}
	r := &goSqlNestedTx{
		goSqlTx: &goSqlTx{
			queryEngineFactory: m.engineQuery.CopyContext(),
		},
		beginEngineFactory: m.beginEngineFactory,
	}
	r.tx, err = m.db.BeginTx(ctx, txo)
	if err != nil {
		m.beginEngineFactory.Apply(r, r.queryEngineFactory.ctx)
	}
	return r, nil
}
