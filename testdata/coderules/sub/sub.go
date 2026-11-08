// Package sub is imported by the control package and has tests of both kinds, so that its
// external tests import it through another package of the module.
package sub

type Value struct{ X int }

func Take(v Value) int { return v.X }
