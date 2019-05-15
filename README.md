# Getting started

## Who this module is for

This module is used in the various database implementations for vsql, or rather, it can be. This is a factory-model that will allow middleware to build and execute requests through the vsql interfaces. This module is not intended to be used by "normal" developers who just want to make calls to databases. This is intended for people who wish to make drivers or middleware for those drivers.

# Purpose

While proving out the vsql interfaces with a co-worker, he indicated that there was a need to track certain calls and states from behind the scenes with databases in such a way that the implementing code is not aware of these calls and, indeed, have no use of this information, but system builders do. He wanted to know when a statement had been prepared but had not been closed when the connection was released back to the pool.

Sure, we could just solve that problem in the system and mark it as a feature of the library, but that's not extensible. Despite my desire to build a sensible architecture, I did implement a row-builder for expediency's sake and I can see now why that was a mistake. It really does need to be some sort of extensible middleware. So that's what this is.

If, one day, you decide that Go's database/sql package sucks, you can still use the vsql interfaces without changing code.