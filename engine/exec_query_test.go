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
	"github.com/wojnosystems/vsql_engine/vsql_context"
	"testing"
)

func TestEngine_ExecQuery(t *testing.T) {
	expectedResult := &vresult.ResulterMock{}
	expectedParams := param.New("SELECT * FROM puppies")
	var actualParams param.Queryer
	engine := New()
	engine.ExecQueryMW().Append(func(ctx context.Context, c vsql_context.Execer) {
		actualParams = c.Query()
		c.SetResult(expectedResult)
		c.Next(ctx)
	})
	actualResult, _ := engine.Exec(context.Background(), expectedParams)
	if actualParams != expectedParams {
		t.Error("expected parameters to be passed")
	}

	if actualResult != expectedResult {
		t.Error("expected the actual result to be the one we set")
	}
}
