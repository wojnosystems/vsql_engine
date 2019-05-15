package context

import "testing"

func TestContext_IsAborted(t *testing.T) {
	c := New()
	if c.IsAborted() {
		t.Error("should not default to aborted")
	}

	c.Abort()
	if !c.IsAborted() {
		t.Error("should be aborted")
	}
}
