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
	"github.com/wojnosystems/vsql_engine/engine_ware"
)

// SQLEnginer is the base engine with non-nested transactions
type SingleTXer interface {
	SQLQueryer
	// This Engine conforms to the MultiTXer interface
	vsql.SQLer
	// Enables the transactions to be used, but not nested transactions
	engine_ware.BeginWare
	// Group creates a copy of the middleware and context created thus far so you can have customized middleware for parts of your application
	Group() SingleTXer
}

// SQLEnginer is the version of the engine with nested transactions
type MultiTXer interface {
	SQLQueryer
	// This Engine conforms to the MultiTXer interface
	vsql.SQLNester
	// Enables nested transactions
	engine_ware.BeginNestedWare
	// Group creates a copy of the middleware and context created thus far so you can have customized middleware for parts of your application
	Group() MultiTXer
}

// SQLEngineQueryer is the part of the engine that only supports the operations common to both the nested and non-nested engine
// this includes regular queries (SELECT) and other operations on the database (DELETE/UPDATE/INSERT/ETC)
type SQLQueryer interface {
	// Enables transactions to be committed
	engine_ware.CommitWare
	// Enables transactions to be rolled back
	engine_ware.RollbackWare
	// Enables (result-returning) queries to be run
	engine_ware.QueryWare
	// Enables INSERT INTO to be performed
	engine_ware.InsertQueryWare
	// Enables Exec to be called: DELETE, UPDATE, INSERT INTO
	engine_ware.ExecWare
	// Enables the ability to create statements
	engine_ware.StatementPrepareWare
	// Enables the ability to close statements
	engine_ware.StatementCloseWare
	// Enables Statement-based Queries
	engine_ware.StatementQueryWare
	// Enables Statement-based Inserts
	engine_ware.StatementInsertQueryWare
	// Enables Statement-based Execs
	engine_ware.StatementExecQueryWare
	// Enables result Rows records to be closed
	engine_ware.RowsCloseWare
	// Enables the next row to be fetched
	engine_ware.RowsNextWare
	// Enables the ability to ping the database server to check for connectivity
	engine_ware.PingWare
	// Enables the connection to be closed
	engine_ware.ConnCloseWare
}
