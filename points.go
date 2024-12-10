package utils

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

// Point is an X, Y coordinate of a given numeric type
type Point[T constraints.Integer | constraints.Float] struct{ X, Y T }

// StdPoint is a "standard" (i.e. regular int) point
type StdPoint = Point[int]

// PointPlusData is an X, Y coordinate, and some associated data
type PointPlusData[T constraints.Integer | constraints.Float, ET any] struct {
	Point Point[T]
	Data  ET
}

// Rectangle is an area defined by two X, Y coordinates at the top left and bottom right corners
type Rectangle[T constraints.Integer | constraints.Float] struct {
	P1 Point[T]
	P2 Point[T]
}

// StdRectangle is a rectangle of type int
type StdRectangle = Rectangle[int]

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

// Equal returns true if the points are equal
func (p Point[T]) Equal(q Point[T]) bool {
	return p.X == q.X && p.Y == q.Y
}

// OrderCoords makes sure the coordinates are properly ordered
func (r Rectangle[T]) OrderCoords() {
	if r.P1.X > r.P2.X {
		r.P1.X, r.P2.X = r.P2.X, r.P1.X
	}
	if r.P1.Y > r.P2.Y {
		r.P1.Y, r.P2.Y = r.P2.Y, r.P1.Y
	}
}

// Empty returns true if this rectangle is empty (i.e. has an area of zero)
func (r Rectangle[T]) Empty() bool {
	return r.P1.X == r.P2.X && r.P1.Y == r.P2.Y
}

// Equal returns true if the rectangles are equal
func (r Rectangle[T]) Equal(v Rectangle[T]) bool {
	r.OrderCoords()
	v.OrderCoords()
	return r.P1.Equal(v.P1) && r.P2.Equal(v.P2)
}

// Intersection returns the largest rectangle contained by this and another rectangle.  If they do not overlap
// then the zero rectangle is returned.
func (r Rectangle[T]) Intersection(v Rectangle[T]) Rectangle[T] {
	r.OrderCoords()
	v.OrderCoords()
	if r.P1.X < v.P1.X {
		r.P1.X = v.P1.X
	}
	if r.P1.Y < v.P1.Y {
		r.P1.Y = v.P1.Y
	}
	if r.P2.X > v.P2.X {
		r.P2.X = v.P2.X
	}
	if r.P2.Y > v.P2.Y {
		r.P2.Y = v.P2.Y
	}
	if r.P1.X > r.P2.X || r.P1.Y > r.P2.Y {
		var zr Rectangle[T]
		return zr
	}
	return r
}

// Overlaps returns true if the given rectangle has an overlap with this rectangle.
func (r Rectangle[T]) Overlaps(v Rectangle[T]) bool {
	r.OrderCoords()
	v.OrderCoords()
	return !r.Empty() && !v.Empty() &&
		r.P1.X < v.P2.X && v.P1.X < r.P2.X &&
		r.P1.Y < v.P2.Y && v.P1.Y < r.P2.Y
}

// Union returns the smallest rectangle that contains this and another rectangle.
func (r Rectangle[T]) Union(v Rectangle[T]) Rectangle[T] {
	if r.Empty() {
		return v
	}
	if v.Empty() {
		return r
	}
	r.OrderCoords()
	v.OrderCoords()
	if r.P1.X > v.P1.X {
		r.P1.X = v.P1.X
	}
	if r.P1.Y > v.P1.Y {
		r.P1.Y = v.P1.Y
	}
	if r.P2.X < v.P2.X {
		r.P2.X = v.P2.X
	}
	if r.P2.Y < v.P2.Y {
		r.P2.Y = v.P2.Y
	}
	return r
}

// Contains returns true if this rectangle entirely contains another rectangle.
func (r Rectangle[T]) Contains(v Rectangle[T]) bool {
	r.OrderCoords()
	v.OrderCoords()
	if v.Empty() {
		return true
	}
	return r.P1.X <= v.P1.X && v.P2.X <= r.P2.X &&
		r.P1.Y <= v.P1.Y && v.P2.Y <= r.P2.Y
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

// Area returns the area of the rectangle
func (r Rectangle[T]) Area() T {
	return r.Width() * r.Height()
}
