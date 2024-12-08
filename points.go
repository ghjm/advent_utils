package utils

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

// Point is an X, Y coordinate of a given integer type
type Point[T constraints.Integer] struct{ X, Y T }

// StdPoint is a "standard" (i.e. regular int) point
type StdPoint = Point[int]

// PointPlusData is an X, Y coordinate, and some associated data
type PointPlusData[T constraints.Integer, ET any] struct {
	Point Point[T]
	Data  ET
}

// Rectangle is an area defined by two X, Y coordinates at the top left and bottom right corners
type Rectangle[T constraints.Integer] struct {
	P1 Point[T]
	P2 Point[T]
}

// String returns a string value of the point
func (p Point[T]) String() string {
	return fmt.Sprintf("(X=%d, Y=%d)", p.X, p.Y)
}

// Add adds the X and Y values of another point to this point
func (p Point[T]) Add(q Point[T]) Point[T] {
	return Point[T]{p.X + q.X, p.Y + q.Y}
}

// Delta returns the difference between the X and Y values of this point and another point
func (p Point[T]) Delta(q Point[T]) Point[T] {
	return Point[T]{p.X - q.X, p.Y - q.Y}
}

// Negate returns a point with the negation of the X and Y values of this point
func (p Point[T]) Negate() Point[T] {
	return Point[T]{-p.X, -p.Y}
}

// Within returns true if this point is within the bounds of a given rectangle
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

// Width returns the width of the rectangle
func (r Rectangle[T]) Width() T {
	a, b := r.P1.X, r.P2.X
	if a > b {
		a, b = b, a
	}
	return b - a + 1
}

// Height returns the width of the rectangle
func (r Rectangle[T]) Height() T {
	a, b := r.P1.Y, r.P2.Y
	if a > b {
		a, b = b, a
	}
	return b - a + 1
}
