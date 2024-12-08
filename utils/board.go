package utils

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

type Board[KT constraints.Integer, VT any] struct {
	contents Map2D[KT, VT]
	bounds   *Rectangle[KT]
	emptyVal VT
	convFunc func(uint8) VT
	compFunc func(VT, VT) bool
}

type StandardBoard struct {
	Board[int, rune]
}

func NewBoard[KT constraints.Integer, VT any](emptyVal VT, convFunc func(uint8) VT, compFunc func(VT, VT) bool) *Board[KT, VT] {
	b := Board[KT, VT]{}
	b.emptyVal = emptyVal
	b.convFunc = convFunc
	b.compFunc = compFunc
	return &b
}

func NewRuneBoard[KT constraints.Integer](emptyVal rune) *Board[KT, rune] {
	return NewBoard[KT, rune](emptyVal,
		func(v uint8) rune {
			return rune(v)
		},
		func(v1 rune, v2 rune) bool {
			return v1 == v2
		})
}

func NewStandardBoard() *StandardBoard {
	return &StandardBoard{
		Board: *NewRuneBoard[int]('.'),
	}
}

func (b *Board[KT, VT]) SetBounds(x1, x2, y1, y2 KT) {
	b.bounds = &Rectangle[KT]{
		Point[KT]{x1, y1},
		Point[KT]{x2, y2},
	}
}

func (b *Board[KT, VT]) ClearBounds() {
	b.bounds = nil
}

func (b *Board[KT, VT]) Transform(tFunc func(p Point[KT], v VT) VT) {
	type change[KT constraints.Integer] struct {
		p Point[KT]
		v VT
	}
	var changes []change[KT]
	b.contents.Iterate(func(p Point[KT], v VT) bool {
		ch := tFunc(p, v)
		if !b.compFunc(ch, v) {
			changes = append(changes, change[KT]{p, ch})
		}
		return true
	})
	for _, c := range changes {
		b.contents.Set(c.p, c.v)
	}
}

func (b *Board[KT, VT]) FromStrings(s []string) error {
	if b.convFunc == nil {
		return fmt.Errorf("board conversion function not initialized")
	}
	var x, y KT
	for y = 0; y < KT(len(s)); y++ {
		if len(s[y]) != len(s[0]) {
			return fmt.Errorf("line lengths not uniform")
		}
		for x = 0; x < KT(len(s[y])); x++ {
			v := b.convFunc(s[y][x])
			if !b.compFunc(v, b.emptyVal) {
				b.contents.Set(Point[KT]{x, y}, v)
			}
		}
	}
	b.bounds = &Rectangle[KT]{
		P1: Point[KT]{0, 0},
		P2: Point[KT]{KT(len(s[0]) - 1), KT(len(s) - 1)},
	}
	return nil
}

func (b *Board[KT, VT]) MustFromStrings(s []string) {
	err := b.FromStrings(s)
	if err != nil {
		panic(err)
	}
}

func (b *Board[KT, VT]) FromFile(name string) error {
	if b.convFunc == nil {
		return fmt.Errorf("board conversion function not initialized")
	}
	var y, sizeX int
	err := OpenAndReadLines(name, func(line string) error {
		if sizeX == 0 {
			sizeX = len(line)
		} else {
			if sizeX != len(line) {
				return fmt.Errorf("line lengths not uniform")
			}
		}
		for x := range line {
			p := b.convFunc(line[x])
			if !b.compFunc(p, b.emptyVal) {
				b.contents.Set(Point[KT]{KT(x), KT(y)}, p)
			}
		}
		y++
		return nil
	})
	if err != nil {
		return err
	}
	b.bounds = &Rectangle[KT]{
		P1: Point[KT]{0, 0},
		P2: Point[KT]{KT(sizeX - 1), KT(y - 1)},
	}
	return nil
}

func (b *Board[KT, VT]) MustFromFile(name string) {
	err := b.FromFile(name)
	if err != nil {
		panic(err)
	}
}

func (b *Board[KT, VT]) Bounds() Rectangle[KT] {
	if b.bounds == nil {
		return Rectangle[KT]{}
	} else {
		return *b.bounds
	}
}

func (b *Board[KT, VT]) orderBounds() {
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

func (b *Board[KT, VT]) ExpandBounds(p Point[KT]) {
	if b.bounds == nil {
		b.bounds = &Rectangle[KT]{
			P1: Point[KT]{
				X: p.X,
				Y: p.Y,
			},
			P2: Point[KT]{
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

func (b *Board[KT, VT]) Contains(p Point[KT]) bool {
	if b.bounds == nil {
		return true
	}
	return p.Within(*b.bounds)
}

func (b *Board[KT, VT]) Get(p Point[KT]) VT {
	return b.contents.GetOrDefault(p, b.emptyVal)
}

func (b *Board[KT, VT]) Set(p Point[KT], v VT) {
	b.contents.Set(p, v)
}

func (b *Board[KT, VT]) Clear(p Point[KT]) {
	b.contents.Delete(p)
}

func (b *Board[KT, VT]) SetAndExpandBounds(p Point[KT], v VT) {
	b.contents.Set(p, v)
	b.ExpandBounds(p)
}

func (b *Board[KT, VT]) IterateBounds(pFunc func(Point[KT]) bool) {
	if b.bounds == nil {
		return
	}
	b.orderBounds()
	for y := b.bounds.P1.Y; y <= b.bounds.P2.Y; y++ {
		for x := b.bounds.P1.X; x <= b.bounds.P2.X; x++ {
			if !pFunc(Point[KT]{x, y}) {
				return
			}
		}
	}
}

func (b *Board[KT, VT]) Copy() Board[KT, VT] {
	var nb Board[KT, VT]
	nb.contents = b.contents.Copy()
	nb.emptyVal = b.emptyVal
	if b.bounds != nil {
		nb.bounds = &Rectangle[KT]{
			P1: Point[KT]{
				X: b.bounds.P1.X,
				Y: b.bounds.P1.Y,
			},
			P2: Point[KT]{
				X: b.bounds.P2.X,
				Y: b.bounds.P2.Y,
			},
		}
	}
	return nb
}
