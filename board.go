package utils

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
)

// Board is an abstraction of a 2D map of discrete map points
type Board[KT constraints.Integer, VT any] struct {
	contents Map2D[KT, VT]
	bounds   *Rectangle[KT]
	emptyVal VT
	convFunc func(uint8) VT
	compFunc func(VT, VT) bool
}

// NewBoard allocates and initializes a new Board
func NewBoard[KT constraints.Integer, VT any](emptyVal VT, convFunc func(uint8) VT, compFunc func(VT, VT) bool) *Board[KT, VT] {
	b := Board[KT, VT]{}
	b.emptyVal = emptyVal
	b.convFunc = convFunc
	b.compFunc = compFunc
	return &b
}

// RuneBoard is a Board where all the points are represented by a single rune
type RuneBoard[KT constraints.Integer] struct {
	Board[KT, rune]
}

// NewRuneBoard allocates and initializes a new RuneBoard
func NewRuneBoard[KT constraints.Integer](emptyVal rune) *RuneBoard[KT] {
	return &RuneBoard[KT]{
		Board: *NewBoard[KT, rune](emptyVal,
			func(v uint8) rune {
				return rune(v)
			},
			func(v1 rune, v2 rune) bool {
				return v1 == v2
			},
		),
	}
}

// StdBoard is a RuneBoard whose addresses are of type Point[int] (aka StdPoint) and whose empty value is '.'
type StdBoard struct {
	RuneBoard[int]
}

// NewStandardBoard allocates and initializes a new StandardBoard
func NewStdBoard() *StdBoard {
	return &StdBoard{
		RuneBoard: *NewRuneBoard[int]('.'),
	}
}

// RunePlusData is the data type for elements of a RunePlusBoard
type RunePlusData[ET any] struct {
	Value rune
	Extra ET
}

// RunePlusBoard is a Board that stores a rune plus arbitrary extra data at each location
type RunePlusBoard[KT constraints.Integer, ET any] struct {
	Board[KT, RunePlusData[ET]]
}

// NewRunePlusBoard allocates and initializes a new RunePlusBoard
func NewRunePlusBoard[KT constraints.Integer, ET any](emptyVal rune) *RunePlusBoard[KT, ET] {
	var zve ET
	return &RunePlusBoard[KT, ET]{
		Board: *NewBoard[KT, RunePlusData[ET]](RunePlusData[ET]{emptyVal, zve},
			func(u uint8) RunePlusData[ET] {
				return RunePlusData[ET]{
					Value: rune(u),
					Extra: zve,
				}
			},
			func(v1 RunePlusData[ET], v2 RunePlusData[ET]) bool {
				return v1.Value == v2.Value
			}),
	}
}

// SetBounds sets the boundary rectangle
func (b *Board[KT, VT]) SetBounds(x1, x2, y1, y2 KT) {
	b.bounds = &Rectangle[KT]{
		Point[KT]{x1, y1},
		Point[KT]{x2, y2},
	}
}

// ClearBounds clears the boundary rectangle
func (b *Board[KT, VT]) ClearBounds() {
	b.bounds = nil
}

// Transform iterates through each point of a Board, allowing each to be changed.  The changes are batched till the end.
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

// FromStrings reads a Board from a slice of strings
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

// MustFromStrings reads a Board from a slice of strings, and panics on any error
func (b *Board[KT, VT]) MustFromStrings(s []string) {
	err := b.FromStrings(s)
	if err != nil {
		panic(err)
	}
}

// FromFile reads a Board from a file on disk
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

// MustFromFile reads a Board from a file on disk, and panics on any error
func (b *Board[KT, VT]) MustFromFile(name string) {
	err := b.FromFile(name)
	if err != nil {
		panic(err)
	}
}

// orderBounds ensures the boundary is in the correct order
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

// Bounds returns the boundary rectangle, or the zero value rectangle if no bounds are set
func (b *Board[KT, VT]) Bounds() Rectangle[KT] {
	b.orderBounds()
	if b.bounds == nil {
		return Rectangle[KT]{}
	} else {
		return *b.bounds
	}
}

// ExpandBounds expands the boundary rectangle to include an arbitrary point
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

// Contains returns true if the given point is contained within the board's boundary rectangle
func (b *Board[KT, VT]) Contains(p Point[KT]) bool {
	if b.bounds == nil {
		return true
	}
	return p.Within(*b.bounds)
}

// Get returns the value of a location on the board
func (b *Board[KT, VT]) Get(p Point[KT]) VT {
	return b.contents.GetOrDefault(p, b.emptyVal)
}

// GetRune returns only the rune from a location on the board
func (b *RunePlusBoard[KT, ET]) GetRune(p Point[KT]) rune {
	return b.contents.GetOrDefault(p, b.emptyVal).Value
}

// GetExtra returns only the extra value from a location on the board
func (b *RunePlusBoard[KT, ET]) GetExtra(p Point[KT]) ET {
	return b.contents.GetOrDefault(p, b.emptyVal).Extra
}

// Set sets the value of a location on the board
func (b *Board[KT, VT]) Set(p Point[KT], v VT) {
	b.contents.Set(p, v)
}

// SetRuneOnly sets the rune and clears any extra data
func (b *RunePlusBoard[KT, ET]) SetRuneOnly(p Point[KT], v rune) {
	b.contents.Set(p, RunePlusData[ET]{Value: v})
}

// SetRune sets the rune, preserving extra data if present
func (b *RunePlusBoard[KT, ET]) SetRune(p Point[KT], v rune) {
	c, ok := b.contents.Get(p)
	var ev ET
	if ok {
		ev = c.Extra
	}
	b.contents.Set(p, RunePlusData[ET]{Value: v, Extra: ev})
}

// SetExtra sets the extra data, preserving the rune value.  If the rune had no value, the empty value is added.
func (b *RunePlusBoard[KT, ET]) SetExtra(p Point[KT], v ET) {
	c := b.contents.GetOrDefault(p, b.emptyVal)
	b.contents.Set(p, RunePlusData[ET]{Value: c.Value, Extra: v})
}

// Clear clears the value of a location on the board
func (b *Board[KT, VT]) Clear(p Point[KT]) {
	b.contents.Delete(p)
}

// SetAndExpandBounds sets a point and also ensures that this point is within the boundary rectangle
func (b *Board[KT, VT]) SetAndExpandBounds(p Point[KT], v VT) {
	b.contents.Set(p, v)
	b.ExpandBounds(p)
}

// Iterate calls a function for every populated location on the board
func (b *Board[KT, VT]) Iterate(iterFunc func(p Point[KT], v VT) bool) {
	b.contents.Iterate(iterFunc)
}

// IterateRunes calls a function for every populated location on the board, returning only the rune
func (b *RunePlusBoard[KT, VT]) IterateRunes(iterFunc func(p Point[KT], v rune) bool) {
	b.contents.Iterate(func(p Point[KT], v RunePlusData[VT]) bool {
		return iterFunc(p, v.Value)
	})
}

// IterateBounds calls a function for every point within the boundary rectangle, whether or not it is populated
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

// Copy returns a new copy of the board
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

// Serial returns a unique serial number for this point, equal to y*width+x
func (b *Board[KT, VT]) Serial(p Point[KT]) KT {
	return p.Y*b.Bounds().Width() + p.X
}

// Format returns a string representation of the board, suitable for printing.  The user must supply a conversion function.
func (b *Board[KT, VT]) Format(fFunc func(VT) rune) []string {
	var results []string
	var builder strings.Builder
	b.IterateBounds(func(p Point[KT]) bool {
		if p.X == 0 && p.Y != 0 {
			results = append(results, builder.String())
			builder.Reset()
		}
		builder.WriteRune(fFunc(b.Get(p)))
		return true
	})
	results = append(results, builder.String())
	return results
}

// Format returns a string representation of the board, suitable for printing.
func (b *RuneBoard[KT]) Format() []string {
	return b.Board.Format(func(r rune) rune {
		return r
	})
}

// Format returns a string representation of the board, suitable for printing.  Extra data is not shown.
func (b *RunePlusBoard[KT, ET]) Format() []string {
	return b.Board.Format(func(r RunePlusData[ET]) rune {
		return r.Value
	})
}
