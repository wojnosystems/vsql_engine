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
	"github.com/stretchr/testify/assert"
	"github.com/wojnosystems/vsql/vparam"
	"github.com/wojnosystems/vsql_engine/engine_context"
	"testing"
)

// TestEngineQuery_Group ensures that grouping creates separate middleware chains
func TestEngineQuery_AllExecuted(t *testing.T) {
	e1 := NewSingle()
	mark1 := make(map[string]int)
	appendToAllNoNest(e1, func(named string) {
		mark1[named] = mark1[named] + 1
	})

	executeAllCallbacksNoNest(e1)
	// Asserts
	for expectedKey, expectedCount := range allMWNoNestCounts {
		v := mark1[expectedKey]
		assert.Equal(t, expectedCount, v, expectedKey, "MK1")
	}
}

// TestEngineQuery_Group ensures that grouping creates separate middleware chains
func TestEngineQuery_Group(t *testing.T) {
	e1 := NewSingle()
	mark1 := make(map[string]int)
	mark2 := make(map[string]int)
	appendToAllNoNest(e1, func(named string) {
		mark1[named] = mark1[named] + 1
	})

	// Create a group. This should copy all middleware to e2 and leave e1 alone.
	e2 := e1.Group()

	// Simulate adding more middleware. It should be added to e2, but not e1
	appendToAllNoNest(e2, func(named string) {
		mark2[named] = mark2[named] + 1
	})

	executeAllCallbacksNoNest(e2)
	// Asserts
	for expectedKey, expectedCount := range allMWNoNestCounts {
		{
			v := mark1[expectedKey]
			assert.Equal(t, expectedCount, v, expectedKey, "MK1")
		}
		{
			v := mark2[expectedKey]
			assert.Equal(t, expectedCount, v, expectedKey, "MK2")
		}
	}
}

// TestEngineQuery_GroupAppended ensures that things are appended
func TestEngineQuery_GroupAppended(t *testing.T) {
	e1 := NewSingle()
	mark1 := make(map[string]int)
	appendToAllNoNest(e1, func(named string) {
		mark1[named] = 1
	})

	// Simulate adding more middleware. It should be added to e2, but not e1
	appendToAllNoNest(e1, func(named string) {
		if _, ok := mark1[named]; !ok {
			assert.Fail(t, "value for ", named, " not already set")
		}
	})

	executeAllCallbacksNoNest(e1)
}

// TestEngineQuery_GroupAppended ensures that things are appended
func TestEngineQuery_GroupPrepended(t *testing.T) {
	e1 := NewSingle()
	mark1 := make(map[string]int)
	appendToAllNoNest(e1, func(named string) {
		if _, ok := mark1[named]; !ok {
			assert.Fail(t, "value for ", named, " already set")
		}
	})

	// Simulate prepending middleware: this should run before the first appendToAllNoNest above
	prependToAllNoNest(e1, func(named string) {
		mark1[named] = 1
	})

	executeAllCallbacksNoNest(e1)
}

// TestEngineQuery_GroupAppended ensures that things are appended
func TestEngineQueryNoNest_NoMiddleware(t *testing.T) {
	e1 := NewSingle()
	executeAllCallbacksNoNest(e1)
}

func appendToAllNoNest(e SingleTXer, f func(named string)) {
	e.BeginMW().Append(func(ctx context.Context, c engine_context.Beginner) {
		f("BeginMW")
		c.Next(ctx)
	})
	e.StatementPrepareMW().Append(func(ctx context.Context, c engine_context.Preparer) {
		f("StatementPrepareMW")
		c.Next(ctx)
	})
	e.QueryMW().Append(func(ctx context.Context, c engine_context.Queryer) {
		f("QueryMW")
		c.Next(ctx)
	})
	e.InsertQueryMW().Append(func(ctx context.Context, c engine_context.Inserter) {
		f("InsertQueryMW")
		c.Next(ctx)
	})
	e.ExecQueryMW().Append(func(ctx context.Context, c engine_context.Execer) {
		f("ExecQueryMW")
		c.Next(ctx)
	})
	e.PingMW().Append(func(ctx context.Context, c engine_context.Er) {
		f("PingMW")
		c.Next(ctx)
	})
	e.StatementCloseMW().Append(func(ctx context.Context, c engine_context.StatementCloser) {
		f("StatementCloseMW")
		c.Next(ctx)
	})
	e.StatementQueryMW().Append(func(ctx context.Context, c engine_context.StatementQueryer) {
		f("StatementQueryMW")
		c.Next(ctx)
	})
	e.StatementInsertQueryMW().Append(func(ctx context.Context, c engine_context.StatementInsertQueryer) {
		f("StatementInsertQueryMW")
		c.Next(ctx)
	})
	e.StatementExecQueryMW().Append(func(ctx context.Context, c engine_context.StatementExecQueryer) {
		f("StatementExecQueryMW")
		c.Next(ctx)
	})
	e.CommitMW().Append(func(ctx context.Context, c engine_context.Beginner) {
		f("CommitMW")
		c.Next(ctx)
	})
	e.RollbackMW().Append(func(ctx context.Context, c engine_context.Beginner) {
		f("RollbackMW")
		c.Next(ctx)
	})
	e.RowsNextMW().Append(func(ctx context.Context, c engine_context.Rowser) {
		f("RowsNextMW")
		c.Next(ctx)
	})
	e.RowsCloseMW().Append(func(ctx context.Context, c engine_context.Rowser) {
		f("RowsCloseMW")
		c.Next(ctx)
	})
	e.ConnCloseMW().Append(func(ctx context.Context, c engine_context.Er) {
		f("ConnCloseMW")
		c.Next(ctx)
	})
}

func prependToAllNoNest(e SingleTXer, f func(named string)) {
	e.BeginMW().Prepend(func(ctx context.Context, c engine_context.Beginner) {
		f("BeginMW")
		c.Next(ctx)
	})
	e.StatementPrepareMW().Prepend(func(ctx context.Context, c engine_context.Preparer) {
		f("StatementPrepareMW")
		c.Next(ctx)
	})
	e.QueryMW().Prepend(func(ctx context.Context, c engine_context.Queryer) {
		f("QueryMW")
		c.Next(ctx)
	})
	e.InsertQueryMW().Prepend(func(ctx context.Context, c engine_context.Inserter) {
		f("InsertQueryMW")
		c.Next(ctx)
	})
	e.ExecQueryMW().Prepend(func(ctx context.Context, c engine_context.Execer) {
		f("ExecQueryMW")
		c.Next(ctx)
	})
	e.PingMW().Prepend(func(ctx context.Context, c engine_context.Er) {
		f("PingMW")
		c.Next(ctx)
	})
	e.StatementCloseMW().Prepend(func(ctx context.Context, c engine_context.StatementCloser) {
		f("StatementCloseMW")
		c.Next(ctx)
	})
	e.StatementQueryMW().Prepend(func(ctx context.Context, c engine_context.StatementQueryer) {
		f("StatementQueryMW")
		c.Next(ctx)
	})
	e.StatementInsertQueryMW().Prepend(func(ctx context.Context, c engine_context.StatementInsertQueryer) {
		f("StatementInsertQueryMW")
		c.Next(ctx)
	})
	e.StatementExecQueryMW().Prepend(func(ctx context.Context, c engine_context.StatementExecQueryer) {
		f("StatementExecQueryMW")
		c.Next(ctx)
	})
	e.CommitMW().Prepend(func(ctx context.Context, c engine_context.Beginner) {
		f("CommitMW")
		c.Next(ctx)
	})
	e.RollbackMW().Prepend(func(ctx context.Context, c engine_context.Beginner) {
		f("RollbackMW")
		c.Next(ctx)
	})
	e.RowsNextMW().Prepend(func(ctx context.Context, c engine_context.Rowser) {
		f("RowsNextMW")
		c.Next(ctx)
	})
	e.RowsCloseMW().Prepend(func(ctx context.Context, c engine_context.Rowser) {
		f("RowsCloseMW")
		c.Next(ctx)
	})
	e.ConnCloseMW().Prepend(func(ctx context.Context, c engine_context.Er) {
		f("ConnCloseMW")
		c.Next(ctx)
	})
}

// Very basic fire or no fire test of all callbacks at every level (assuming no nesting)
var allMWNoNestCounts = map[string]int{
	"BeginMW":                1,
	"StatementPrepareMW":     2,
	"QueryMW":                2,
	"InsertQueryMW":          2,
	"ExecQueryMW":            2,
	"PingMW":                 1,
	"StatementCloseMW":       2,
	"StatementQueryMW":       2,
	"StatementInsertQueryMW": 2,
	"StatementExecQueryMW":   2,
	"CommitMW":               1,
	"RollbackMW":             1,
	"RowsNextMW":             4,
	"RowsCloseMW":            4,
	"ConnCloseMW":            1,
}

// executeAllCallbacksNoNest triggers every single middleware. Should generate allMWNoNestCounts
func executeAllCallbacksNoNest(e SingleTXer) {
	p := vparam.New("Some Query")
	ctx := context.Background()
	_ = e.Ping(ctx)
	{
		r, _ := e.Query(ctx, p)
		r.Next()
		_ = r.Close()
	}
	_, _ = e.Exec(ctx, p)
	_, _ = e.Insert(ctx, p)
	{ // Statements
		s, _ := e.Prepare(ctx, p)
		{
			r, _ := s.Query(ctx, p)
			r.Next()
			_ = r.Close()
		}
		_, _ = s.Insert(ctx, p)
		_, _ = s.Exec(ctx, p)
		_ = s.Close()
	}
	{ // Transaction
		tx, _ := e.Begin(ctx, nil)
		{
			r, _ := tx.Query(ctx, p)
			r.Next()
			_ = r.Close()
		}
		_, _ = tx.Insert(ctx, p)
		_, _ = tx.Exec(ctx, p)

		{ // Statements
			s, _ := tx.Prepare(ctx, p)
			{
				r, _ := s.Query(ctx, p)
				r.Next()
				_ = r.Close()
			}
			_, _ = s.Insert(ctx, p)
			_, _ = s.Exec(ctx, p)
			_ = s.Close()
		}
		_ = tx.Rollback()
		_ = tx.Commit()
	}
	_ = e.Close()
}
