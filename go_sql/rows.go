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

package go_sql

import (
	"database/sql"
	"github.com/wojnosystems/vsql/vrows"
)

type goSqlRows struct {
	vrows.Rowser
	queryEngineFactory *engineQuery
	SqlRows            *sql.Rows
}

// Next calls Next() on the sql.Rows object
func (m *goSqlRows) Next() vrows.Rower {
	if m.SqlRows.Next() {
		r := &goSqlRow{SqlRows: m.SqlRows}
		m.queryEngineFactory.RowWares.Apply(r)
		return r
	}
	return nil
}

// Close cleans up the Rows object, releasing it's object back to the pool. Call this when you're done with your vquery results
func (m *goSqlRows) Close() error {
	return m.SqlRows.Close()
}

type goSqlRow struct {
	vrows.Rower
	SqlRows *sql.Rows
}

// Scan calls Scan() on the sql.Rows object and is the primary way to convert values from SQL into Go-native objects
func (m *goSqlRow) Scan(dest ...interface{}) error {
	return m.SqlRows.Scan(dest...)
}

// Columns returns a list of columns available for this vrow
func (m *goSqlRow) Columns() (columnNames []string) {
	// #NOERR: Will never error as this cannot be used while the vrow is closed
	columnNames, _ = m.SqlRows.Columns()
	return
}