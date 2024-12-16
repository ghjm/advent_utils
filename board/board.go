package board

import (
	"fmt"
	"github.com/ghjm/advent_utils"
	"golang.org/x/exp/constraints"
	"strings"
)

// BoardStorage is an interface to pluggable back-end storage for a Board
type BoardStorage[KT constraints.Integer, VT any] interface {
	Allocate(width, height KT, emptyVal VT)
	Set(p utils.Point[KT], v VT)
	Get(p utils.Point[KT]) (VT, bool)
	Delete(p utils.Point[KT])
	GetOrDefault(p utils.Point[KT], def VT) VT
	Iterate(iterFunc func(p utils.Point[KT], v VT) bool)
	CopyToBoardStorage() BoardStorage[KT, VT]
}

// Board is an abstraction of a 2D map of discrete map points
type Board[KT constraints.Integer, VT any] struct {
	BoardOptions[KT, VT]
}

// BoardOptions collects extra options when initializing a Board
type BoardOptions[KT constraints.Integer, VT any] struct {
	storage  BoardStorage[KT, VT]
	bounds   *utils.Rectangle[KT]
	emptyVal VT
	convFunc func(uint8) VT
	compFunc func(VT, VT) bool
}

// WithStorage provides a storage backend to a Board
func WithStorage[KT constraints.Integer, VT any](storage BoardStorage[KT, VT]) func(*BoardOptions[KT, VT]) {
	return func(options *BoardOptions[KT, VT]) {
		options.storage = storage
	}
}

// WithBounds provides initial bounds to a Board
func WithBounds[KT constraints.Integer, VT any](bounds utils.Rectangle[KT]) func(*BoardOptions[KT, VT]) {
	return func(options *BoardOptions[KT, VT]) {
		options.bounds = &bounds
	}
}

// WithEmptyVal provides an empty value
func WithEmptyVal[KT constraints.Integer, VT any](emptyVal VT) func(*BoardOptions[KT, VT]) {
	return func(options *BoardOptions[KT, VT]) {
		options.emptyVal = emptyVal
	}
}

// WithConvFunc provides a conversion function, needed for loading from strings/files
func WithConvFunc[KT constraints.Integer, VT any](convFunc func(uint8) VT) func(*BoardOptions[KT, VT]) {
	return func(options *BoardOptions[KT, VT]) {
		options.convFunc = convFunc
	}
}

// WithCompareFunc provides a conversion function, needed for loading from strings/files
func WithCompareFunc[KT constraints.Integer, VT any](compFunc func(VT, VT) bool) func(*BoardOptions[KT, VT]) {
	return func(options *BoardOptions[KT, VT]) {
		options.compFunc = compFunc
	}
}

// NewBoard allocates and initializes a new Board
func NewBoard[KT constraints.Integer, VT any](options ...func(board *BoardOptions[KT, VT])) *Board[KT, VT] {
	b := Board[KT, VT]{}
	for _, opt := range options {
		opt(&b.BoardOptions)
	}
	if b.storage == nil {
		b.storage = &Map2D[KT, VT]{}
	}
	return &b
}

// RuneBoard is a Board where all the points are represented by a single rune
type RuneBoard[KT constraints.Integer] struct {
	Board[KT, rune]
}

// NewRuneBoard allocates and initializes a new RuneBoard
func NewRuneBoard[KT constraints.Integer](options ...func(board *BoardOptions[KT, rune])) *RuneBoard[KT] {
	b := &RuneBoard[KT]{
		Board: *NewBoard[KT, rune](options...),
	}
	if b.emptyVal == 0 {
		b.emptyVal = '.'
	}
	if b.convFunc == nil {
		b.convFunc = func(v uint8) rune {
			return rune(v)
		}
	}
	if b.compFunc == nil {
		b.compFunc = func(v1 rune, v2 rune) bool {
			return v1 == v2
		}
	}
	return b
}

// StdBoard is a convenience name for RuneBoard[int]
type StdBoard struct {
	RuneBoard[int]
}

// NewStdBoard allocates and initializes a new StandardBoard
func NewStdBoard(options ...func(board *BoardOptions[int, rune])) *StdBoard {
	return &StdBoard{
		RuneBoard: *NewRuneBoard[int](options...),
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

// NewRunePlusBoard allocates and initializes a new RunePlusBoard.
func NewRunePlusBoard[KT constraints.Integer, ET any](options ...func(board *BoardOptions[KT, RunePlusData[ET]])) *RunePlusBoard[KT, ET] {
	b := &RunePlusBoard[KT, ET]{
		Board: *NewBoard[KT, RunePlusData[ET]](options...),
	}
	if b.emptyVal.Value == 0 {
		b.emptyVal.Value = '.'
	}
	var zve ET
	if b.convFunc == nil {
		b.convFunc = func(u uint8) RunePlusData[ET] {
			return RunePlusData[ET]{
				Value: rune(u),
				Extra: zve,
			}
		}
	}
	if b.compFunc == nil {
		b.compFunc = func(v1 RunePlusData[ET], v2 RunePlusData[ET]) bool {
			return v1.Value == v2.Value
		}
	}
	return b
}

// Transform iterates through each point of a Board, allowing each to be changed.  The changes are batched till the end.
func (b *Board[KT, VT]) Transform(tFunc func(p utils.Point[KT], v VT) VT) {
	type change[KT constraints.Integer] struct {
		p utils.Point[KT]
		v VT
	}
	var changes []change[KT]
	b.storage.Iterate(func(p utils.Point[KT], v VT) bool {
		ch := tFunc(p, v)
		if !b.compFunc(ch, v) {
			changes = append(changes, change[KT]{p, ch})
		}
		return true
	})
	for _, c := range changes {
		b.storage.Set(c.p, c.v)
	}
}

// FromStrings reads a Board from a slice of strings
func (b *Board[KT, VT]) FromStrings(s []string) error {
	if b.convFunc == nil {
		return fmt.Errorf("board conversion function not initialized")
	}
	b.storage.Allocate(KT(len(s[0])), KT(len(s)), b.emptyVal)
	var x, y KT
	for y = 0; y < KT(len(s)); y++ {
		if len(s[y]) != len(s[0]) {
			return fmt.Errorf("line lengths not uniform")
		}
		for x = 0; x < KT(len(s[y])); x++ {
			v := b.convFunc(s[y][x])
			if !b.compFunc(v, b.emptyVal) {
				b.storage.Set(utils.Point[KT]{x, y}, v)
			}
		}
	}
	b.bounds = &utils.Rectangle[KT]{
		P1: utils.Point[KT]{0, 0},
		P2: utils.Point[KT]{KT(len(s[0]) - 1), KT(len(s) - 1)},
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
	var lines []string
	err := utils.OpenAndReadLines(name, func(line string) error {
		lines = append(lines, line)
		return nil
	})
	if err != nil {
		return err
	}
	b.bounds = &utils.Rectangle[KT]{
		P1: utils.Point[KT]{0, 0},
		P2: utils.Point[KT]{KT(len(lines[0]) - 1), KT(len(lines) - 1)},
	}
	b.storage.Allocate(KT(len(lines[0])), KT(len(lines)), b.emptyVal)
	for y, line := range lines {
		if len(line) != len(lines[0]) {
			return fmt.Errorf("line lengths not uniform")
		}
		for x := range line {
			p := b.convFunc(line[x])
			if !b.compFunc(p, b.emptyVal) {
				b.storage.Set(utils.Point[KT]{KT(x), KT(y)}, p)
			}
		}
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
func (b *Board[KT, VT]) Bounds() utils.Rectangle[KT] {
	b.orderBounds()
	if b.bounds == nil {
		return utils.Rectangle[KT]{}
	} else {
		return *b.bounds
	}
}

// ExpandBounds expands the boundary rectangle to include an arbitrary point
func (b *Board[KT, VT]) ExpandBounds(p utils.Point[KT]) {
	if b.bounds == nil {
		b.bounds = &utils.Rectangle[KT]{
			P1: utils.Point[KT]{
				X: p.X,
				Y: p.Y,
			},
			P2: utils.Point[KT]{
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
func (b *Board[KT, VT]) Contains(p utils.Point[KT]) bool {
	if b.bounds == nil {
		return true
	}
	return p.Within(*b.bounds)
}

// Get returns the value of a location on the board
func (b *Board[KT, VT]) Get(p utils.Point[KT]) VT {
	return b.storage.GetOrDefault(p, b.emptyVal)
}

// GetRune returns only the rune from a location on the board
func (b *RunePlusBoard[KT, ET]) GetRune(p utils.Point[KT]) rune {
	return b.storage.GetOrDefault(p, b.emptyVal).Value
}

// GetExtra returns only the extra value from a location on the board
func (b *RunePlusBoard[KT, ET]) GetExtra(p utils.Point[KT]) ET {
	return b.storage.GetOrDefault(p, b.emptyVal).Extra
}

// Set sets the value of a location on the board
func (b *Board[KT, VT]) Set(p utils.Point[KT], v VT) {
	b.storage.Set(p, v)
}

// SetRuneOnly sets the rune and clears any extra data
func (b *RunePlusBoard[KT, ET]) SetRuneOnly(p utils.Point[KT], v rune) {
	b.storage.Set(p, RunePlusData[ET]{Value: v})
}

// SetRune sets the rune, preserving extra data if present
func (b *RunePlusBoard[KT, ET]) SetRune(p utils.Point[KT], v rune) {
	c, ok := b.storage.Get(p)
	var ev ET
	if ok {
		ev = c.Extra
	}
	b.storage.Set(p, RunePlusData[ET]{Value: v, Extra: ev})
}

// SetExtra sets the extra data, preserving the rune value.  If the rune had no value, the empty value is added.
func (b *RunePlusBoard[KT, ET]) SetExtra(p utils.Point[KT], v ET) {
	c := b.storage.GetOrDefault(p, b.emptyVal)
	b.storage.Set(p, RunePlusData[ET]{Value: c.Value, Extra: v})
}

// Clear clears the value of a location on the board
func (b *Board[KT, VT]) Clear(p utils.Point[KT]) {
	b.storage.Delete(p)
}

// SetAndExpandBounds sets a point and also ensures that this point is within the boundary rectangle
func (b *Board[KT, VT]) SetAndExpandBounds(p utils.Point[KT], v VT) {
	b.storage.Set(p, v)
	b.ExpandBounds(p)
}

// Iterate calls a function for every populated location on the board
func (b *Board[KT, VT]) Iterate(iterFunc func(p utils.Point[KT], v VT) bool) {
	b.storage.Iterate(iterFunc)
}

// IterateRunes calls a function for every populated location on the board, returning only the rune
func (b *RunePlusBoard[KT, VT]) IterateRunes(iterFunc func(p utils.Point[KT], v rune) bool) {
	b.storage.Iterate(func(p utils.Point[KT], v RunePlusData[VT]) bool {
		return iterFunc(p, v.Value)
	})
}

// IterateBounds calls a function for every point within the boundary rectangle, whether or not it is populated
func (b *Board[KT, VT]) IterateBounds(pFunc func(utils.Point[KT]) bool) {
	if b.bounds == nil {
		return
	}
	b.orderBounds()
	for y := b.bounds.P1.Y; y <= b.bounds.P2.Y; y++ {
		for x := b.bounds.P1.X; x <= b.bounds.P2.X; x++ {
			if !pFunc(utils.Point[KT]{x, y}) {
				return
			}
		}
	}
}

// Copy returns a new copy of the board
func (b *Board[KT, VT]) Copy() *Board[KT, VT] {
	var nb Board[KT, VT]
	nb.storage = b.storage.CopyToBoardStorage()
	nb.emptyVal = b.emptyVal
	if b.bounds != nil {
		nb.bounds = &utils.Rectangle[KT]{
			P1: utils.Point[KT]{
				X: b.bounds.P1.X,
				Y: b.bounds.P1.Y,
			},
			P2: utils.Point[KT]{
				X: b.bounds.P2.X,
				Y: b.bounds.P2.Y,
			},
		}
	}
	return &nb
}

// Copy returns a new copy of the board
func (b *RuneBoard[KT]) Copy() *RuneBoard[KT] {
	var nb RuneBoard[KT]
	nb.Board = *b.Board.Copy()
	return &nb
}

// Copy returns a new copy of the board
func (b *StdBoard) Copy() *StdBoard {
	var nb StdBoard
	nb.Board = *b.Board.Copy()
	return &nb
}

// Serial returns a unique serial number for this point, equal to y*width+x
func (b *Board[KT, VT]) Serial(p utils.Point[KT]) KT {
	return p.Y*b.Bounds().Width() + p.X
}

// Format returns a string representation of the board, suitable for printing.  The user must supply a conversion function.
func (b *Board[KT, VT]) Format(fFunc func(VT) rune) []string {
	if b.bounds == nil {
		return nil
	}
	b.orderBounds()
	var results []string
	for y := b.bounds.P1.Y; y <= b.bounds.P2.Y; y++ {
		var builder strings.Builder
		for x := b.bounds.P1.X; x <= b.bounds.P2.X; x++ {
			builder.WriteRune(fFunc(b.Get(utils.Point[KT]{X: x, Y: y})))
		}
		results = append(results, builder.String())
	}
	return results
}

// Print prints a board to stdout.  The user must supply a conversion function.
func (b *Board[KT, VT]) Print(fFunc func(VT) rune) {
	for _, line := range b.Format(fFunc) {
		fmt.Printf("%s\n", line)
	}
}

// Cardinals returns the four cardinal points adjacent to a given point.
func (b *Board[KT, VT]) Cardinals(p utils.Point[KT], includeOffBoard bool) []utils.Point[KT] {
	var results []utils.Point[KT]
	for _, d := range []utils.StdPoint{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		np := utils.Point[KT]{
			X: p.X + KT(d.X),
			Y: p.Y + KT(d.Y),
		}
		if includeOffBoard || b.Contains(np) {
			results = append(results, np)
		}
	}
	return results
}

// Diagonals returns the eight diagonal (including cardinal) points adjacent to a given point.
func (b *Board[KT, VT]) Diagonals(p utils.Point[KT], includeOffBoard bool) []utils.Point[KT] {
	var results []utils.Point[KT]
	for _, d := range []utils.StdPoint{{-1, -1}, {0, -1}, {1, -1}, {-1, 0}, {1, 0}, {-1, 1}, {0, 1}, {1, 1}} {
		np := utils.Point[KT]{
			X: p.X + KT(d.X),
			Y: p.Y + KT(d.Y),
		}
		if includeOffBoard || b.Contains(np) {
			results = append(results, np)
		}
	}
	return results
}

// Search performs a flood fill type search of a board from a given start point and with a given neighbors function.
func (b *Board[KT, VT]) Search(start utils.Point[KT], neighbors func(p utils.Point[KT]) []utils.Point[KT]) map[utils.Point[KT]]struct{} {
	open := []utils.Point[KT]{start}
	visited := make(map[utils.Point[KT]]struct{})
	for len(open) > 0 {
		cur := open[0]
		open = open[1:]
		if _, ok := visited[cur]; ok {
			continue
		}
		visited[cur] = struct{}{}
		for _, p := range neighbors(cur) {
			if _, ok := visited[p]; !ok {
				open = append(open, p)
			}
		}
	}
	return visited
}

// FindRegions groups a board into same-valued regions based on cardinal neighbors.
func (b *Board[KT, VT]) FindRegions(includeEmptyVal bool) []map[utils.Point[KT]]struct{} {
	if b.compFunc == nil {
		panic("compFunc not defined")
	}
	var regions []map[utils.Point[KT]]struct{}
	visited := make(map[utils.Point[KT]]struct{})
	b.IterateBounds(func(p utils.Point[KT]) bool {
		if _, ok := visited[p]; ok {
			return true
		}
		if (!includeEmptyVal) && b.compFunc(b.Get(p), b.emptyVal) {
			return true
		}
		reg := b.Search(p, func(pn utils.Point[KT]) []utils.Point[KT] {
			var results []utils.Point[KT]
			for _, cpn := range b.Cardinals(pn, false) {
				if b.compFunc(b.Get(pn), b.Get(cpn)) {
					results = append(results, cpn)
				}
			}
			return results
		})
		for rp := range reg {
			visited[rp] = struct{}{}
		}
		regions = append(regions, reg)
		return true
	})
	return regions
}

// Format returns a string representation of the board, suitable for printing.
func (b *RuneBoard[KT]) Format() []string {
	return b.Board.Format(func(r rune) rune {
		return r
	})
}

// Print prints a board to stdout.
func (b *RuneBoard[KT]) Print() {
	b.Board.Print(func(r rune) rune {
		return r
	})
}

// Format returns a string representation of the board, suitable for printing.  Extra data is not shown.
func (b *RunePlusBoard[KT, ET]) Format() []string {
	return b.Board.Format(func(r RunePlusData[ET]) rune {
		return r.Value
	})
}

// Print prints a board to stdout.  Extra data is not printed.
func (b *RunePlusBoard[KT, ET]) Print() {
	b.Board.Print(func(r RunePlusData[ET]) rune {
		return r.Value
	})
}

type FlatBoard struct {
	board    [][]rune
	emptyVal rune
}

func (fb *FlatBoard) Allocate(width, height int, emptyVal rune) {
	fb.board = make([][]rune, 0, height)
	for y := 0; y < height; y++ {
		line := make([]rune, 0, width)
		for x := 0; x < width; x++ {
			line = append(line, emptyVal)
		}
		fb.board = append(fb.board, line)
	}
	fb.emptyVal = emptyVal
}

func (fb *FlatBoard) GetBounds() utils.StdRectangle {
	return utils.StdRectangle{
		P1: utils.Point[int]{},
		P2: utils.Point[int]{
			X: len(fb.board[0]) - 1,
			Y: len(fb.board) - 1,
		},
	}
}

func (fb *FlatBoard) Set(p utils.StdPoint, v rune) {
	fb.board[p.Y][p.X] = v
}

func (fb *FlatBoard) Get(p utils.StdPoint) (rune, bool) {
	if p.X >= 0 && p.X < len(fb.board[0]) && p.Y >= 0 && p.Y < len(fb.board) {
		return fb.board[p.Y][p.X], true
	}
	return 0, false
}

func (fb *FlatBoard) Delete(p utils.StdPoint) {
	fb.Set(p, fb.emptyVal)
}

func (fb *FlatBoard) GetOrDefault(p utils.StdPoint, def rune) rune {
	if p.X >= 0 && p.X < len(fb.board[0]) && p.Y >= 0 && p.Y < len(fb.board) {
		return fb.board[p.Y][p.X]
	}
	return def
}

func (fb *FlatBoard) Iterate(iterFunc func(p utils.StdPoint, v rune) bool) {
	for y := 0; y < len(fb.board); y++ {
		for x := 0; x < len(fb.board[0]); x++ {
			if !iterFunc(utils.StdPoint{x, y}, fb.board[y][x]) {
				return
			}
		}
	}
}

func (fb *FlatBoard) CopyToBoardStorage() BoardStorage[int, rune] {
	nb := new(FlatBoard)
	nb.emptyVal = fb.emptyVal
	for y := 0; y < len(fb.board); y++ {
		var line []rune
		for x := 0; x < len(fb.board[0]); x++ {
			line = append(line, fb.board[y][x])
		}
		nb.board = append(nb.board, line)
	}
	return nb
}
