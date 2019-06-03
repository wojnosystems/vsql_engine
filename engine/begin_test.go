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
	"errors"
	"github.com/wojnosystems/vsql"
	"github.com/wojnosystems/vsql_engine/vsql_context"
	"testing"
)

func TestNoNest_Begin_NilErrorPassed(t *testing.T) {
	didRun := false
	engine := New()
	engine.BeginMW().Append(func(ctx context.Context, c vsql_context.Beginner) {
		didRun = true
		c.Next(ctx)
	})
	_, err := engine.Begin(context.Background(), nil)
	if err != nil {
		t.Error("expected nil error")
	}
	if didRun == false {
		t.Error("expected middleware to run")
	}
}

func TestNoNest_Begin_ErrorPassed(t *testing.T) {
	expectedErr := errors.New("test")
	engine := New()
	engine.BeginMW().Append(func(ctx context.Context, c vsql_context.Beginner) {
		c.SetError(expectedErr)
		c.Next(ctx)
	})
	_, err := engine.Begin(context.Background(), nil)
	if err != expectedErr {
		t.Error("expected an error to be set")
	}
}

func TestNoNest_Begin_QETPassedCommit(t *testing.T) {
	engine := New()
	expectedQET := &vsql.QueryExecTransactionerMock{}
	var actualQET vsql.QueryExecTransactioner
	engine.BeginMW().Append(func(ctx context.Context, c vsql_context.Beginner) {
		c.SetQueryExecTransactioner(expectedQET)
		c.Next(ctx)
	})
	engine.CommitMW().Append(func(ctx context.Context, c vsql_context.Beginner) {
		actualQET = c.QueryExecTransactioner()
		c.Next(ctx)
	})
	r, _ := engine.Begin(context.Background(), nil)
	_ = r.Commit()
	if expectedQET != actualQET {
		t.Error("expected a query Exec transactioner to be returned")
	}
}

func TestNoNest_Begin_QETPassedRollback(t *testing.T) {
	engine := New()
	expectedQET := &vsql.QueryExecTransactionerMock{}
	var actualQET vsql.QueryExecTransactioner
	engine.BeginMW().Append(func(ctx context.Context, c vsql_context.Beginner) {
		c.SetQueryExecTransactioner(expectedQET)
		c.Next(ctx)
	})
	engine.RollbackMW().Append(func(ctx context.Context, c vsql_context.Beginner) {
		actualQET = c.QueryExecTransactioner()
		c.Next(ctx)
	})
	r, _ := engine.Begin(context.Background(), nil)
	_ = r.Rollback()
	if expectedQET != actualQET {
		t.Error("expected a query Exec transactioner to be returned")
	}
}
