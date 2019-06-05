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
	"errors"
	"github.com/wojnosystems/vsql/param"
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql_engine/engine_context"
	"testing"
)

func TestEngine_Query(t *testing.T) {
	expectedRows := &vrows.RowserMock{}
	var actualRows vrows.Rowser
	expectedParams := param.New("SELECT * FROM puppies")
	var actualParams param.Queryer
	engine := NewSingle()
	engine.QueryMW().Append(func(ctx context.Context, c engine_context.Queryer) {
		actualParams = c.Query()
		c.SetRows(expectedRows)
		c.Next(ctx)
	})
	rows, _ := engine.Query(context.Background(), expectedParams)
	if actualParams != expectedParams {
		t.Error("expected parameters to be passed")
	}

	engine.RowsNextMW().Append(func(ctx context.Context, c engine_context.Rowser) {
		actualRows = c.Rows()
		c.Next(ctx)
	})
	rows.Next()
	if actualRows != expectedRows {
		t.Error("expected actual rows to be set")
	}

	expectedErr := errors.New("boom")
	didRun := false
	engine.RowsCloseMW().Append(func(ctx context.Context, c engine_context.Rowser) {
		didRun = true
		if actualRows != c.Rows() {
			t.Error("expected the same rows object to be returned")
		}
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

}
