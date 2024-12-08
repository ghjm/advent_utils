package utils

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

type Board[T constraints.Integer] struct {
	contents Map2D[T, rune]
	bounds   *Rectangle[T]
	emptyVal rune
}

func (b *Board[T]) SetEmptyVal(ev rune) {
	b.emptyVal = ev
}

func (b *Board[T]) SetBounds(x1, x2, y1, y2 T) {
	b.bounds = &Rectangle[T]{
		Point[T]{x1, y1},
		Point[T]{x2, y2},
	}
}

func (b *Board[T]) ClearBounds() {
	b.bounds = nil
}

func (b *Board[T]) Transform(tFunc func(p Point[T], v rune) rune) {
	type change[T constraints.Integer] struct {
		p Point[T]
		v rune
	}
	var changes []change[T]
	b.contents.Iterate(func(p Point[T], v rune) bool {
		ch := tFunc(p, v)
		if ch != v {
			changes = append(changes, change[T]{p, ch})
		}
		return true
	})
	for _, c := range changes {
		b.contents.Set(c.p, c.v)
	}
}

func (b *Board[T]) FromStrings(s []string, emptyVal rune) error {
	var x, y T
	b.emptyVal = emptyVal
	for y = 0; y < T(len(s)); y++ {
		if len(s[y]) != len(s[0]) {
			return fmt.Errorf("line lengths not uniform")
		}
		for x = 0; x < T(len(s[y])); x++ {
			v := rune(s[y][x])
			if v != emptyVal {
				b.contents.Set(Point[T]{x, y}, v)
			}
		}
	}
	b.bounds = &Rectangle[T]{
		P1: Point[T]{0, 0},
		P2: Point[T]{T(len(s[0]) - 1), T(len(s) - 1)},
	}
	return nil
}

func (b *Board[T]) MustFromStrings(s []string, emptyVal rune) {
	err := b.FromStrings(s, emptyVal)
	if err != nil {
		panic(err)
	}
}

func (b *Board[T]) FromFile(name string, emptyVal rune) error {
	var y, sizeX int
	b.emptyVal = emptyVal
	err := OpenAndReadLines(name, func(line string) error {
		if sizeX == 0 {
			sizeX = len(line)
		} else {
			if sizeX != len(line) {
				return fmt.Errorf("line lengths not uniform")
			}
		}
		for x := range line {
			p := rune(line[x])
			if p != emptyVal {
				b.contents.Set(Point[T]{T(x), T(y)}, p)
			}
		}
		y++
		return nil
	})
	if err != nil {
		return err
	}
	b.bounds = &Rectangle[T]{
		P1: Point[T]{0, 0},
		P2: Point[T]{T(sizeX - 1), T(y - 1)},
	}
	return nil
}

func (b *Board[T]) MustFromFile(name string, emptyVal rune) {
	err := b.FromFile(name, emptyVal)
	if err != nil {
		panic(err)
	}
}

func (b *Board[T]) Bounds() Rectangle[T] {
	if b.bounds == nil {
		return Rectangle[T]{}
	} else {
		return *b.bounds
	}
}

func (b *Board[T]) orderBounds() {
	if b.bounds == nil {
		return
	}
	if b.bounds.P1.X > b.bounds.P2.X {
		b.bounds.P1.X, b.bounds.P2.X = b.bounds.P2.X, b.bounds.P1.X
	}
	if b.bounds.P1.Y > b.bounds.P2.Y {
		b.bounds.P1.Y, b.bounds.P2.Y = b.bounds.P2.Y, b.bounds.P1.Y
	}
}

func (b *Board[T]) ExpandBounds(p Point[T]) {
	if b.bounds == nil {
		b.bounds = &Rectangle[T]{
			P1: Point[T]{
				X: p.X,
				Y: p.Y,
			},
			P2: Point[T]{
				X: p.X,
				Y: p.Y,
			},
		}
	} else {
		b.orderBounds()
		if b.bounds.P1.X > p.X {
			b.bounds.P1.X = p.X
		}
		if b.bounds.P1.Y > p.Y {
			b.bounds.P1.Y = p.Y
		}
		if b.bounds.P2.X < p.X {
			b.bounds.P2.X = p.X
		}
		if b.bounds.P2.Y < p.Y {
			b.bounds.P2.Y = p.Y
		}
	}
}

func (b *Board[T]) Contains(p Point[T]) bool {
	if b.bounds == nil {
		return true
	}
	return p.Within(*b.bounds)
}

func (b *Board[T]) Get(p Point[T]) rune {
	return b.contents.GetOrDefault(p, b.emptyVal)
}

func (b *Board[T]) Set(p Point[T], v rune) {
	b.contents.Set(p, v)
}

func (b *Board[T]) Clear(p Point[T]) {
	b.contents.Delete(p)
}

func (b *Board[T]) SetAndExpandBounds(p Point[T], v rune) {
	b.contents.Set(p, v)
	b.ExpandBounds(p)
}

func (b *Board[T]) IterateBounds(pFunc func(Point[T]) bool) {
	if b.bounds == nil {
		return
	}
	b.orderBounds()
	for y := b.bounds.P1.Y; y <= b.bounds.P2.Y; y++ {
		for x := b.bounds.P1.X; x <= b.bounds.P2.X; x++ {
			if !pFunc(Point[T]{x, y}) {
				return
			}
		}
	}
}

func (b *Board[T]) Copy() Board[T] {
	var nb Board[T]
	nb.contents = b.contents.Copy()
	nb.emptyVal = b.emptyVal
	if b.bounds != nil {
		nb.bounds = &Rectangle[T]{
			P1: Point[T]{
				X: b.bounds.P1.X,
				Y: b.bounds.P1.Y,
			},
			P2: Point[T]{
				X: b.bounds.P2.X,
				Y: b.bounds.P2.Y,
			},
		}
	}
	return nb
}
