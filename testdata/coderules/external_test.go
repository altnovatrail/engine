package coderules_test

import (
	"testing"

	"example.com/coderules"
	"example.com/coderules/sub"
)

func TestExternal(t *testing.T) {
	if coderules.Exported(2) != 2 {
		t.Fail()
	}
	// coderules has tests in the package, and sub, which it imports, comes here through it and
	// directly: the check must see one sub.
	if sub.Take(coderules.Make()) != 1 {
		t.Fail()
	}
	_ = coderules.Exported(2) * 3.0 // want fused
}
