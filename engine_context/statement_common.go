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
	"github.com/wojnosystems/vsql/vparam"
	"github.com/wojnosystems/vsql/vstmt"
)

type statementCommoner interface {
	Er
	SetStatement(statementer vstmt.Statementer)
	Statement() vstmt.Statementer
	SetParameterer(parameterer vparam.Parameterer)
	Parameterer() vparam.Parameterer
}

func newStatementCommon() *statementCommon {
	return &statementCommon{
		contextBase: newContextBase(),
	}
}

type statementCommon struct {
	*contextBase
	statement   vstmt.Statementer
	parameterer vparam.Parameterer
}

func (c *statementCommon) SetStatement(s vstmt.Statementer) {
	c.statement = s
}

func (c statementCommon) Statement() vstmt.Statementer {
	return c.statement
}

func (c *statementCommon) SetParameterer(parameterer vparam.Parameterer) {
	c.parameterer = parameterer
}

func (c statementCommon) Parameterer() vparam.Parameterer {
	return c.parameterer
}
