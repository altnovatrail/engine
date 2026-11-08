package coderules

import (
	. "math"
	. "strconv"
)

func dotted(x float64) string {
	return Quote(FormatFloat(Exp(x)+Sqrt(x), 'g', -1, 64)) // want math unicode
}
