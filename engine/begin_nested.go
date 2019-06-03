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

// SQLNestedEnginer
type engineNest struct {
	*engineQuery
	beginMW *wares.BeginNestedMW
}

func NewNested() vsql_engine.SQLNestedEnginer {
	m := &engineNest{
		engineQuery: newEngineQuery(),
		beginMW:     wares.NewBeginNestedMW(),
	}
	return m
}

// BeginMW provides a way to add items to the BeginMW
func (m *engineNest) BeginNestedMW() wares.BeginNestedAdder {
	return m.beginMW
}

// Begin see github.com/wojnosystems/vsql/transactions.go#TransactionStarter
func (m *engineNest) Begin(ctx context.Context, txOp vtxn.TxOptioner) (n vsql.QueryExecNestedTransactioner, err error) {
	c := vsql_context.NewNestedBeginner()
	c.(vsql_context.WithMiddlewarer).ShallowCopyFrom(m.middlewareContext)
	c.SetTxOptions(txOp)
	m.beginMW.PerformMiddleware(ctx, c)
	s := &nestedTx{
		beginNestedMW:         m.beginMW,
		beginnerNestedContext: c,
		queryEngineFactory:    m,
	}
	return s, c.Error()
}

func (m *engineNest) Group() vsql_engine.SQLNestedEnginer {
	r := &engineNest{
		engineQuery: m.engineQuery.Group(),
		beginMW:     m.beginMW.Copy(),
	}
	return r
}
