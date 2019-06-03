# Getting started

## Who this module is for

This module is for any developer who said: "Middleware is great for web servers, I sure wish I had that for my database!" This module includes a framework that conforms to the interface as defined in: https://github.com/wojnosystems/vsql and allows you to inject middleware whenever calls are made. You can alter the query before they are run, or after you get the results.

Potential uses (or at least the ones I needed):

 * Query Logging
 * Begin + Commit/Rollback tracking (be sure transactions are ended to avoid leaking resources)
 * Prepare + Close tracking (be sure transactions are ended to avoid leaking resources)
 * Global caching with Redis/Memcache

# What does this module do?

This module implements the https://github.com/wojnosystems/vsql interfaces, but instead of actually making database calls, it simply implements the interface and provides a way to plug into the interfaces calls. If you wish to use a database implementation, another module provides that for you using Go's standard database/sql package

 * Global (applies to all, is set when the engine started)
 * Begin (transaction)
 * Begin (nested transaction)
 * Commit (transaction)
 * Rollback (transaction)
 * Prepare (statement)
 * Close (statement)
 * Rows are returned (query)
 * Row is created (query)
 * Result is created (query)
 * InsertResult is created (query)
 
These callbacks work like go-Gin, in that you pass in a function that is executed with a way to persist state. This context state is shared using a RWLock, so it's safe to set and get values, but not safe to override values, based on a conditional or if missing. Keep that in mind.

This context is shared with ALL of the calls. So if you add a value to the context of a Begin, you can look that up in Rollback or Result callbacks.

Adding values to the key-value store of the context is thread-safe.

This was done to provide easier scoping, as transactions should not be communicating and this allows you to use the same key in transactions for values.

Any values set in the context BEFORE begin/prepare is called, will be CLONED into the resulting nested context. Nested transactions will CLONE the current context as well, so every nested transaction gets it's own local namespace.

## Using it

You can inject middleware by calling the appropriate middleware callback end-point.


# Examples

```go
func main() {

    engine := vsql_engine.New()
    // Install MySQL for use
    vsql_mysql.Install(engine)
    
    vsql_txnCloseCheck.Install(engine)

    // statement close check
    statementCloseCheck(engine)
    
    stmt, err := engine.Prepare( context.Background(), param.New("SELECT * FROM users") )
    // Log has message: "statement prepared:w00t"
    stmt.Close()
    // Log has message: "statement closed:hawt"
}

// statementCloseCheck is custom middleware that installs itself into the SQL Engine. 
func statementCloseCheck( engine vsql_engine.R ) {
	engine.PrepareMW().Prepend( func()wares.Prepare {
		return func(ctx context.Context, c vsql_context.Er, query vquery.Queryer) {
			log.Println("statement prepared:w00t")
			return c.Next(ctx, c, query)
		}
	} )
	engine.StatementCloseMW().Prepend( func()wares.Prepare {
		return func(ctx context.Context, c vsql_context.Er) {
			log.Println("statement closed:hawt")
			return c.Next(ctx, c)
		}
	} )
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
