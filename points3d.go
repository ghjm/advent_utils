package utils

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

// Point3D is an X, Y, Z coordinate of a given numeric type
type Point3D[T constraints.Integer | constraints.Float] struct{ X, Y, Z T }

// StdPoint3D is a "standard" (i.e. regular int) Point3d
type StdPoint3D = Point3D[int]

// Cuboid is a volume defined by two X, Y, Z coordinates
type Cuboid[T constraints.Integer | constraints.Float] struct {
	P1 Point3D[T]
	P2 Point3D[T]
}

// StdCuboid is a cuboid of type int
type StdCuboid = Cuboid[int]

// String returns a string value of the point
func (p Point3D[T]) String() string {
	return fmt.Sprintf("(X=%v, Y=%v, Z=%v)", p.X, p.Y, p.Z)
}

// Add adds the coordinates of another point to this point
func (p Point3D[T]) Add(q Point3D[T]) Point3D[T] {
	return Point3D[T]{p.X + q.X, p.Y + q.Y, p.Z + q.Z}
}

// Delta returns the difference between the coordinates of this point and another point
func (p Point3D[T]) Delta(q Point3D[T]) Point3D[T] {
	return Point3D[T]{p.X - q.X, p.Y - q.Y, p.Z - q.Z}
}

// Negate returns a point with the negation of the coordinates of this point
func (p Point3D[T]) Negate() Point3D[T] {
	return Point3D[T]{-p.X, -p.Y, -p.Z}
}

// Within returns true if this point is within the bounds of a given cuboid
func (p Point3D[T]) Within(r Cuboid[T]) bool {
	for _, c := range []struct {
		vLo T
		vHi T
		p   T
	}{
		{r.P1.X, r.P2.X, p.X},
		{r.P1.Y, r.P2.Y, p.Y},
		{r.P1.Z, r.P2.Z, p.Z},
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
func (p Point3D[T]) Equal(v Point3D[T]) bool {
	return p.X == v.X && p.Y == v.Y && p.Z == v.Z
}

// OrderCoords ensures that the coordinates are in the correct order
func (c Cuboid[T]) OrderCoords() {
	if c.P1.X > c.P2.X {
		c.P1.X, c.P2.X = c.P2.X, c.P1.X
	}
	if c.P1.Y > c.P2.Y {
		c.P1.Y, c.P2.Y = c.P2.Y, c.P1.Y
	}
	if c.P1.Z > c.P2.Z {
		c.P1.Z, c.P2.Z = c.P2.Z, c.P1.Z
	}
}

// Empty returns true if this cuboid is empty (i.e. has an area of zero)
func (c Cuboid[T]) Empty() bool {
	return c.P1.X == c.P2.X && c.P1.Y == c.P2.Y && c.P1.Z == c.P2.Z
}

// Intersection returns the largest cuboid contained by this and another cuboid.  If they do not overlap
// then the zero cuboid is returned.
func (c Cuboid[T]) Intersection(v Cuboid[T]) Cuboid[T] {
	c.OrderCoords()
	v.OrderCoords()
	if c.P1.X < v.P1.X {
		c.P1.X = v.P1.X
	}
	if c.P1.Y < v.P1.Y {
		c.P1.Y = v.P1.Y
	}
	if c.P1.Z < v.P1.Z {
		c.P1.Z = v.P1.Z
	}
	if c.P2.X > v.P2.X {
		c.P2.X = v.P2.X
	}
	if c.P2.Y > v.P2.Y {
		c.P2.Y = v.P2.Y
	}
	if c.P2.Z > v.P2.Z {
		c.P2.Z = v.P2.Z
	}
	if c.P1.X > c.P2.X || c.P1.Y > c.P2.Y || c.P1.Z > c.P2.Z {
		var zc Cuboid[T]
		return zc
	}
	return c
}

// Overlaps returns true if the given cuboid has an overlap with this cuboid.
func (c Cuboid[T]) Overlaps(v Cuboid[T]) bool {
	c.OrderCoords()
	v.OrderCoords()
	return !c.Empty() && !v.Empty() &&
		c.P1.X < v.P2.X && v.P1.X < c.P2.X &&
		c.P1.Y < v.P2.Y && v.P1.Y < c.P2.Y &&
		c.P1.Z < v.P2.Z && v.P1.Z < c.P2.Z
}

// Union returns the smallest cuboid that contains this and another cuboid.
func (c Cuboid[T]) Union(v Cuboid[T]) Cuboid[T] {
	if c.Empty() {
		return v
	}
	if v.Empty() {
		return c
	}
	c.OrderCoords()
	v.OrderCoords()
	if c.P1.X > v.P1.X {
		c.P1.X = v.P1.X
	}
	if c.P1.Y > v.P1.Y {
		c.P1.Y = v.P1.Y
	}
	if c.P1.Z > v.P1.Z {
		c.P1.Z = v.P1.Z
	}
	if c.P2.X < v.P2.X {
		c.P2.X = v.P2.X
	}
	if c.P2.Y < v.P2.Y {
		c.P2.Y = v.P2.Y
	}
	if c.P2.Z < v.P2.Z {
		c.P2.Z = v.P2.Z
	}
	return c
}

// Contains returns true if this cuboid entirely contains another cuboid.
func (c Cuboid[T]) Contains(v Cuboid[T]) bool {
	c.OrderCoords()
	v.OrderCoords()
	if v.Empty() {
		return true
	}
	return c.P1.X <= v.P1.X && v.P2.X <= c.P2.X &&
		c.P1.Y <= v.P1.Y && v.P2.Y <= c.P2.Y &&
		c.P1.Z <= v.P1.Z && v.P2.Z <= c.P2.Z
}

// Equal returns true if the rectangles are equal
func (c Cuboid[T]) Equal(v Cuboid[T]) bool {
	c.OrderCoords()
	v.OrderCoords()
	return c.P1.Equal(v.P1) && c.P2.Equal(v.P2)
}

// Width returns the width of the cuboid
func (c Cuboid[T]) Width() T {
	a, b := c.P1.X, c.P2.X
	if a > b {
		a, b = b, a
	}
	return b - a + 1
}

// Depth returns the depth of the cuboid
func (c Cuboid[T]) Depth() T {
	a, b := c.P1.Y, c.P2.Y
	if a > b {
		a, b = b, a
	}
	return b - a + 1
}

// Height returns the height of the cuboid
func (c Cuboid[T]) Height() T {
	a, b := c.P1.Z, c.P2.Z
	if a > b {
		a, b = b, a
	}
	return b - a + 1
}

// Volume returns the volume of the cuboid
func (c Cuboid[T]) Volume() T {
	return c.Width() * c.Depth() * c.Height()
}
