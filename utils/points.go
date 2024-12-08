package utils

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

type Point[T constraints.Integer] struct{ X, Y T }
type Rectangle[T constraints.Integer] struct {
	P1 Point[T]
	P2 Point[T]
}

func (p Point[T]) String() string {
	return fmt.Sprintf("(X=%d, Y=%d)", p.X, p.Y)
}

func (p Point[T]) Add(q Point[T]) Point[T] {
	return Point[T]{p.X + q.X, p.Y + q.Y}
}

func (p Point[T]) Within(r Rectangle[T]) bool {
	for _, c := range []struct {
		vLo T
		vHi T
		p   T
	}{
		{r.P1.X, r.P2.X, p.X},
		{r.P1.Y, r.P2.Y, p.Y},
	} {
		if c.vLo > c.vHi {
			c.vLo, c.vHi = c.vHi, c.vLo
		}
		if c.p < c.vLo || c.p > c.vHi {
			return false
		}
	}
	return true
}
