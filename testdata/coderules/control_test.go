package coderules

import (
	"encoding/xml" // want import
	"math"
	"testing"
)

func TestControl(t *testing.T) {
	a, b, c := 1.5, 2.5, 3.5
	m := map[string]int{"a": 1}
	for k := range m {
		_ = k
	}
	go func() {}()
	_, _ = xml.Marshal(m)
	if math.Exp(a)+math.Log(b) < a*b+c { // want fused math
		t.Log(math.Pow(a, b)) // want math
	}
	if math.Sqrt(a) < float64(b*c) {
		t.Log(a)
	}
}
