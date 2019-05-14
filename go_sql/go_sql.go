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

// This is the implementation of Go's native database/sql queryEngineFactory. It is intended to be used in other drivers, but can be thrown away and re-written if necessary. As long as it supports the queryEngineFactory interface, there is no need to re-write the code for vsql drivers.
// Some of the original work done in the vsql was done using Go's database/sql library and this is wrong. The vsql library should have always been a vanilla interface package (with some helper functions operating on-top of those abstractions). This package is an attempt to fix that by de-coupling the interface from any specific implementation.

// SQLEnginer
type goSql struct {
	*engineQuery
	beginEngineFactory *wares.Begin
}

func New(interpolationFactory strategy.InterpolationFactory, driverFactory func() (db *sql.DB)) SQLEnginer {
	db := driverFactory()
	if db == nil {
		return nil
	}
	m := &goSql{
		engineQuery:        newEngineQuery(db, interpolationFactory),
		beginEngineFactory: wares.NewBegin(),
	}
	return m
}

// BeginMiddleware provides a way to add items to the BeginMiddleware
func (m *goSql) BeginMiddleware() wares.BeginAdder {
	return m.beginEngineFactory
}

// Begin see github.com/wojnosystems/vsql/transactions.go#TransactionStarter
func (m *goSql) Begin(ctx context.Context, txOp vtxn.TxOptioner) (n vsql.QueryExecTransactioner, err error) {
	var txo *sql.TxOptions
	if txOp != nil {
		txo = txOp.ToTxOptions()
	}
	r := &goSqlTx{
		queryEngineFactory: m.engineQuery,
	}
	r.tx, err = m.db.BeginTx(ctx, txo)
	if err != nil {
		m.beginEngineFactory.Apply(r)
	}
	return r, nil
}
