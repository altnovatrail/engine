// Package fma is the control of fma-diff in the justfile: three shapes the compiler fuses, one of
// which at least it fuses on every target the check compiles for, so that the check must see a
// difference here (docs/architecture.md, section 3.7). Nothing else builds it, and no rule of code
// holds it.
package fma

// MulAdd is a*b + c: a fused multiply-add on every target that fuses, amd64 with v3 too.
func MulAdd(a, b, c float64) float64 { return a*b + c }

// MulSub is c - a*b: a fused multiply-subtract where the target has one; amd64 has none.
func MulSub(a, b, c float64) float64 { return c - a*b }

// Middle is (p1+p2)/2 - 0.775, the engine's own shape: the mean of two P against the target of
// the corridor. The division by 2, whose reciprocal is exact, becomes a multiplication, and fuses
// with the subtraction on arm64, ppc64le, s390x, riscv64 and loong64; float64() round the quotient
// keeps it apart.
func Middle(p1, p2 float64) float64 { return (p1+p2)/2 - 0.775 }
