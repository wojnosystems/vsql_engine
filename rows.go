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
	"github.com/wojnosystems/vsql/vrows"
	"github.com/wojnosystems/vsql_engine/engine_context"
)

type rows struct {
	rows               vrows.Rowser
	queryEngineFactory *engineQuery
}

// Next calls Next() on the sql.Rows object
func (m *rows) Next() vrows.Rower {
	c := engine_context.NewRowNext()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetRows(m.rows)
	m.queryEngineFactory.rowsNextMW.PerformMiddleware(context.Background(), c)
	return c.Row()
}

// Close cleans up the Rows object, releasing it's object back to the pool. Call this when you're done with your vquery results
func (m *rows) Close() error {
	c := engine_context.NewRows()
	c.(engine_context.WithMiddlewarer).ShallowCopyFrom(m.queryEngineFactory.middlewareContext)
	c.SetRows(m.rows)
	m.queryEngineFactory.rowsCloseMW.PerformMiddleware(context.Background(), c)
	return c.Error()
}
