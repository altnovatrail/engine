package sub

import "testing"

func TestTake(t *testing.T) {
	if Take(Value{X: 2}) != 2 {
		t.Fail()
	}
}
