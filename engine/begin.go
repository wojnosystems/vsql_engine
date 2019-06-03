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
	"github.com/wojnosystems/vsql/vtxn"
	"github.com/wojnosystems/vsql_engine"
	"github.com/wojnosystems/vsql_engine/vsql_context"
	"github.com/wojnosystems/vsql_engine/wares"
)

// This is the implementation of Go's native database/sql queryEngineFactory. It is intended to be used in other drivers, but can be thrown away and re-written if necessary. As long as it supports the queryEngineFactory interface, there is no need to re-write the code for vsql drivers.
// Some of the original work done in the vsql was done using Go's database/sql library and this is wrong. The vsql library should have always been a vanilla interface package (with some helper functions operating on-top of those abstractions). This package is an attempt to fix that by de-coupling the interface from any specific implementation.

// SQLEnginer
type engineNoNest struct {
	*engineQuery
	beginMW *wares.BeginMW
}

func New() vsql_engine.SQLEnginer {
	m := &engineNoNest{
		engineQuery: newEngineQuery(),
		beginMW:     wares.NewBeginMW(),
	}
	return m
}

// BeginMW provides a way to add items to the BeginMW
func (m *engineNoNest) BeginMW() wares.BeginAdder {
	return m.beginMW
}

// Begin see github.com/wojnosystems/vsql/transactions.go#TransactionStarter
func (m *engineNoNest) Begin(ctx context.Context, txOp vtxn.TxOptioner) (n vsql.QueryExecTransactioner, err error) {
	c := vsql_context.NewBeginner()
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetTxOptions(txOp)
	m.beginMW.PerformMiddleware(ctx, c)
	s := &nonNestedTx{
		beginnerContext:    c,
		queryEngineFactory: m.engineQuery,
	}
	return s, c.Error()
}

func (m *engineNoNest) Group() vsql_engine.SQLEnginer {
	r := &engineNoNest{
		engineQuery: m.engineQuery.Group(),
		beginMW:     m.beginMW.Copy(),
	}
	return r
}
