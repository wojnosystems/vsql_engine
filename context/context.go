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

package context

import "github.com/wojnosystems/vsql_engine/key_value"

type Contexter interface {
	KeyValues() key_value.KeyValuer
	IsAborted() bool
	Abort()
	Next()
	Copy() Contexter
}

type context struct {
	kvo       key_value.KeyValuer
	isAborted bool
}

func New() *context {
	return &context{
		kvo: key_value.New(),
	}
}

func (c *context) KeyValues() key_value.KeyValuer {
	return c.kvo
}

func (c context) IsAborted() bool {
	return c.isAborted
}

func (c *context) Abort() {
	c.isAborted = true
}

func (c *context) Next() {
}

func (c *context) Copy() Contexter {
	rc := New()
	rc.kvo = c.kvo.Copy()
	return rc
}
