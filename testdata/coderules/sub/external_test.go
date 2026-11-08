package sub_test

import (
	"testing"

	"example.com/coderules"
	"example.com/coderules/sub"
)

func TestMake(t *testing.T) {
	if sub.Take(coderules.Make()) != 1 {
		t.Fail()
	}
}
