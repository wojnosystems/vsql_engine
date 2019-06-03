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
	"github.com/wojnosystems/vsql"
	"github.com/wojnosystems/vsql_engine/wares"
)

// SQLEnginer is the base engine with non-nested transactions
type SQLEnginer interface {
	SQLEngineQueryer
	// This Engine conforms to the SQLNester interface
	vsql.SQLer
	// Enables the transactions to be used, but not nested transactions
	wares.BeginWare
	// Group creates a copy of the middleware and context created thus far so you can have customized middleware for parts of your application
	Group() SQLEnginer
}

// SQLEnginer is the version of the engine with nested transactions
type SQLNestedEnginer interface {
	SQLEngineQueryer
	// This Engine conforms to the SQLNester interface
	vsql.SQLNester
	// Enables nested transactions
	wares.BeginNestedWare
	// Group creates a copy of the middleware and context created thus far so you can have customized middleware for parts of your application
	Group() SQLNestedEnginer
}

// SQLEngineQueryer is the part of the engine that only supports the operations common to both the nested and non-nested engine
// this includes regular queries (SELECT) and other operations on the database (DELETE/UPDATE/INSERT/ETC)
type SQLEngineQueryer interface {
	// Enables transactions to be committed
	wares.CommitWare
	// Enables transactions to be rolled back
	wares.RollbackWare
	// Enables (result-returning) queries to be run
	wares.QueryWare
	// Enables INSERT INTO to be performed
	wares.InsertQueryWare
	// Enables Exec to be called: DELETE, UPDATE, INSERT INTO
	wares.ExecWare
	// Enables the ability to create statements
	wares.StatementPrepareWare
	// Enables the ability to close statements
	wares.StatementCloseWare
	// Enables Statement-based Queries
	wares.StatementQueryWare
	// Enables Statement-based Inserts
	wares.StatementInsertQueryWare
	// Enables Statement-based Execs
	wares.StatementExecQueryWare
	// Enables result Rows records to be closed
	wares.RowsCloseWare
	// Enables the next row to be fetched
	wares.RowsNextWare
	// Enables the ability to ping the database server to check for connectivity
	wares.PingWare
	// Enables the connection to be closed
	wares.ConnCloseWare
}
