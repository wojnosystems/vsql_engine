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

package engine_context

import (
	"context"
	"github.com/wojnosystems/vsql/vresult"
)

type StatementExecQueryer interface {
	statementCommoner
	statementQueryCommoner

	SetResult(resulter vresult.Resulter)
	Result() vresult.Resulter
}

func NewStatementExecQuery() StatementExecQueryer {
	return &StatementExecQuery{
		statementQueryCommon: newStatementQueryCommon(),
	}
}

type StatementExecQuery struct {
	*statementQueryCommon
	result vresult.Resulter
}

func (c *StatementExecQuery) SetResult(resulter vresult.Resulter) {
	c.result = resulter
}

func (c StatementExecQuery) Result() vresult.Resulter {
	return c.result
}

// Next runs the middleware, if any is available, null op if not. Next is only intended to be run once each middleware layer.
func (c *StatementExecQuery) Next(ctx context.Context) {
	if c.currentMiddleware != nil {
		cm := c.currentMiddleware
		c.moveToNextMiddleware()
		cm.Value.(MiddlewareFunc)(ctx, c)
	}
}
