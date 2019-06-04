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
	"github.com/wojnosystems/vsql/param"
	"github.com/wojnosystems/vsql/vresult"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine/vsql_context"
	"testing"
)

func TestEngine_StatementExec(t *testing.T) {
	expectedResult := &vresult.ResulterMock{}

	expectedStatement := &vstmt.StatementerMock{}
	expectedParams := param.New("SELECT * FROM puppies")
	engine := New()
	engine.StatementPrepareMW().Append(func(ctx context.Context, c vsql_context.Preparer) {
		c.SetStatement(expectedStatement)
		c.Next(ctx)
	})
	statement, _ := engine.Prepare(context.Background(), expectedParams)

	engine.StatementExecQueryMW().Append(func(ctx context.Context, c vsql_context.StatementExecQueryer) {
		c.SetResult(expectedResult)
		c.Next(ctx)
	})
	result, _ := statement.Exec(context.Background(), param.NewAppendData([]interface{}{5, "true"}))

	if result != expectedResult {
		t.Error("expected the results to be identical")
	}
}
