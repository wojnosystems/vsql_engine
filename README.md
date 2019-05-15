# Getting started

## Who this module is for

This module is used in the various database implementations for vsql, or rather, it can be. This is a factory-model that will allow middleware to build and execute requests through the vsql interfaces. This module is not intended to be used by "normal" developers who just want to make calls to databases. This is intended for people who wish to make drivers or middleware for those drivers and those who are using those drivers to inject middleware into their database calls.

# What does this module do?

This module implements the https://github.com/wojnosystems/vsql interfaces using Go's database/sql package. It also provides a mechanism to inject callbacks for certain events:

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

Also, when a transaction is begun, or a statement is prepared, the context is cloned at that point. Any context added when the transaction was begun (if set in a Begin callback), is forgotten when the transaction is committed or rolled back. This applies to nested transactions, too.

This was done to provide easier scoping, as transactions should not be communicating and this allows you to use the same key in transactions for values.

Any values set in the context BEFORE begin/prepare is called, will be CLONED into the resulting nested context. Nested transactions will CLONE the current context as well, so every nested transaction gets it's own local namespace.

## Using it

You can inject middleware by calling the appropriate middleware callback end-point. If you wish to inject some middleware when a Row is created, call the 

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
