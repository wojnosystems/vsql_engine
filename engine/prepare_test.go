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
	"github.com/wojnosystems/vsql/param"
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine/vsql_context"
	"testing"
)

func TestEngine_Prepare(t *testing.T) {
	expectedRows := &vrows.RowserMock{}
	var actualRows vrows.Rowser

	expectedStatement := &vstmt.StatementerMock{}
	var actualStatement vstmt.Statementer
	expectedParams := param.New("SELECT * FROM puppies")
	var actualParams param.Queryer
	engine := New()
	engine.PrepareMW().Append(func(ctx context.Context, c vsql_context.Preparer) {
		actualParams = c.Query()
		c.SetStatement(expectedStatement)
		c.Next(ctx)
	})
	statement, _ := engine.Prepare(context.Background(), expectedParams)
	if actualParams != expectedParams {
		t.Error("expected parameters to be passed")
	}

	engine.StatementQueryMW().Append(func(ctx context.Context, c vsql_context.StatementQueryer) {
		actualStatement = c.Statement()
		c.SetRows(expectedRows)
		c.Next(ctx)
	})
	rows, _ := statement.Query(context.Background(), param.NewAppendData([]interface{}{5, "true"}))
	if actualStatement != expectedStatement {
		t.Error("expected actual statement to be set")
	}

	engine.RowsNextMW().Append(func(ctx context.Context, c vsql_context.Rowser) {
		actualRows = c.Rows()
		c.Next(ctx)
	})
	rows.Next()

	if expectedRows != actualRows {
		t.Error("expected the same rows object to be returned")
	}

	expectedErr := errors.New("boom")
	didRun := false
	engine.RowsCloseMW().Append(func(ctx context.Context, c vsql_context.Rowser) {
		didRun = true
		c.SetError(expectedErr)
		c.Next(ctx)
	})
	err := rows.Close()
	if !didRun {
		t.Error("expected RowsClose to run")
	}
	if err != expectedErr {
		t.Error("expected error to be returned")
	}

	expectedStatement.AssertExpectations(t)
	expectedRows.AssertExpectations(t)
}
