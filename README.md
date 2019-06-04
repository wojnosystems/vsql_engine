# Getting started

## Who this module is for

This module is for any developer who said: "Middleware is great for web servers, I sure wish I had that for my database!" This module includes a framework that conforms to the interface as defined in: https://github.com/wojnosystems/vsql and allows you to inject middleware whenever calls are made. You can alter the query before they are run, or after you get the results.

You might be just someone who's using this module because you need a way of injecting data into your SQL queries, or someone else used this and you're trying to figure out what's going on. For the most developers, they really shouldn't need to care about this module. All of the interfaces used in the code should be limited to the [VSQL](https://github.com/wojnosystems/vsql) interfaces.

Before playing with this, you should be intimately familiar with Go's [database/sql](https://golang.org/pkg/database/sql/) package, as the VSQL and this module is based upon that interface.

# How does it work?

It works very similarly to Go-Gin's contexts, however, since each call is potentially different, nearly every Middleware (MW) stack has it's own, special context.

You inject callbacks into specific events in the VSQL interface. This allows me to test the interface without actually hooking up a database. This also means that database drivers are ALSO middleware. This means the engine is database-agnostic and will run regardless of which database is hooked into it. This abstraction is extremely powerful if you need to swap out your data-store without re-writing all of your code.

## Cloning Middleware Stacks

You can create an engine, add callbacks, then clone the set and add more callbacks using the Group method. This allows you to mix and match quite easily. You can prepend or append callbacks to the chain.

Potential uses (or at least the ones I needed/could come up with):

 * Query Logging
 * Begin + Commit/Rollback tracking (be sure transactions are ended to avoid leaking resources)
 * Prepare + Close tracking (be sure transactions are ended to avoid leaking resources)
 * Global caching with Redis/Memcache

# What does this module do?

This module implements the https://github.com/wojnosystems/vsql interfaces, but instead of actually making database calls, it simply implements the interface and provides a way to plug into the interfaces calls. If you wish to use a database implementation, another module provides that for you using Go's standard database/sql package.

You can hook into the following callback chains for the following events:

 * Ping (connection)
 * Query (connection)
 * Insert (connection)
 * Exec (connection)
 * Close (connection)
 * Begin (transaction)
 * Begin (nested transaction)
 * Commit (transaction)
 * Rollback (transaction)
 * Prepare (statement)
 * Query (statement)
 * Insert (statement)
 * Exec (statement)
 * Close (statement)
 * Prepare (transaction statement)
 * Query (transaction statement)
 * Insert (transaction statement)
 * Exec (transaction statement)
 * Close (transaction statement)
 * Next (rows (result of query))
 * Close (rows (result of query))
 
These callbacks work like go-Gin, in that you pass in a function that is executed with a way to persist state. This context state is shared using a RWLock, so it's safe to set and get values, but not safe to override values, based on a conditional or if missing. Keep that in mind.

This context is shared with ALL of the calls. So if you add a value to the context of a Begin, you can look that up in Rollback or Result callbacks.

Adding values to the key-value store of the context is thread-safe.

This was done to provide easier scoping, as transactions should not be communicating and this allows you to use the same key in transactions for values.

Any values set in the context BEFORE begin/prepare is called, will be CLONED into the resulting nested context. Nested transactions will CLONE the current context as well, so every nested transaction gets it's own local namespace.

## Using it

You can inject middleware by calling the appropriate middleware callback end-point.

### Calling Next!

You MUST call c.Next(ctx) to allow appended middleware to run. If you do not, no further callbacks will be executed.

Anything running before the call to c.Next(ctx) should occur "before" the database middleware and anything after the c.Next(ctx) call should occur "after" the database middleware has run. Yes, even the database is middleware!

## Passing data

Every vsql_context.* object has a [KeyValuer](https://github.com/wojnosystems/go_keyvaluer) object. You can store arbitrary data here in a thread-safe way. If you need to store data that is transaction-specific, you can create your own substructure and key off of that transaction object. It's guaranteed to be unique (if you clean it up after closing transactions) and can identify the transaction. This is not directly supported by KeyValuer, but it's possible with a little leg-work on your end.

# Examples

```go
package main

import(
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wojnosystems/vsql/vstmt"
	"github.com/wojnosystems/vsql_engine"
	"github.com/wojnosystems/vsql_engine/engine"
)

func main() {
    // The engine is the magic part. It has all of the middleware
    myEngine := engine.New()
    // Install MySQL using Go's database/sql package
    vsql_go.Install(myEngine, sql.Open("mysql", createMySQLConfig().FormatDSN()) )

    // Install your own statement close check middleware
    statementCloseCheck(myEngine)
    
    stmt, err := myEngine.Prepare( context.Background(), param.New("SELECT * FROM users") )
    // Log has message: "statement prepared:w00t"
    stmt.Close()
    // Log has message: "statement closed:hawt"
}

// statementCloseCheck is custom middleware that installs itself into the SQL Engine. When a statement is prepared, it logs it, when a statement is closed, it logs it
func statementCloseCheck( e vsql_engine.SQLEnginer ) {
	// Prepend is used as the MySQL engine is already installed and we want to run BEFORE the database gets a hold of things. We don't have to for this example, but it's generally what you want.
	e.StatementPrepareMW().Prepend(func(ctx context.Context, c vsql_context.Preparer) {
        log.Println("statement prepared:w00t")
        c.Next(ctx)
	} )
	e.StatementCloseMW().Prepend(func(ctx context.Context, c vsql_context.StatementCloser) {
        log.Println("statement closed:hawt")
        c.Next(ctx)
	} )
}

func createMySQLConfig() (cfg mysql.Config) {
    cfg = mysql.Config{
        User:                 os.Getenv("MYSQL_USER"),
        Passwd:               os.Getenv("MYSQL_PASSWORD"),
        Addr:                 os.Getenv("MYSQL_ADDR"),
        DBName:               os.Getenv("MYSQL_DBNAME"),
        AllowNativePasswords: true,
        AllowOldPasswords:    true,
    }
    if strings.HasPrefix(cfg.Addr, "unix") {
        cfg.Net = "unix"
    } else {
        cfg.Net = "tcp"
    }
    return
}
```

# Creating your own middleware

You can create your own middleware and store data in the vsql_context.Er object by using the KeyValue() object.

## Naming your keys

To avoid name collisions, you should name your keys for any data stored in vsql_context.Er.KeyValue() using the full name of your module, including the github.com part or where ever it's hosted. This should guarantee no collisions.

# Purpose

While proving out the vsql interfaces with a co-worker, he indicated that there was a need to track certain calls and states from behind the scenes with databases in such a way that the implementing code is not aware of these calls and, indeed, have no use of this information, but system builders do. He wanted to know when a statement had been prepared but had not been closed when the connection was released back to the pool.

Sure, we could just solve that problem in the system and mark it as a feature of the library, but that's not extensible. Despite my desire to build a sensible architecture, I did implement a row-builder for expediency's sake and I can see now why that was a mistake. It really does need to be some sort of extensible middleware. So that's what this is.

If, one day, you decide that Go's database/sql package sucks, you can still use the vsql interfaces without changing code that uses the interfaces.

# License 

Copyright 2019 Chris Wojno

Permission is hereby granted, free of charge, to any person obtaining a copy of this software 
and associated documentation files (the "Software"), to deal in the Software without restriction, 
including without limitation the rights to use, copy, modify, merge, publish, distribute, 
sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is 
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or 
substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING 
BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND 
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, 
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, 
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
