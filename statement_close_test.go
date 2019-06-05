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
	"github.com/wojnosystems/vsql/param"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine/engine_context"
	"testing"
)

func TestEngine_StatementClose(t *testing.T) {

	expectedStatement := &vstmt.StatementerMock{}
	expectedParams := param.New("SELECT * FROM puppies")
	engine := NewSingle()
	engine.StatementPrepareMW().Append(func(ctx context.Context, c engine_context.Preparer) {
		c.SetStatement(expectedStatement)
		c.Next(ctx)
	})
	statement, _ := engine.Prepare(context.Background(), expectedParams)

	didRun := false
	engine.StatementCloseMW().Append(func(ctx context.Context, c engine_context.StatementCloser) {
		didRun = true
		if expectedStatement != c.Statement() {
			t.Error("expected to get the same statement provided")
		}
		c.Next(ctx)
	})
	_ = statement.Close()

	if !didRun {
		t.Error("expected to run")
	}
}