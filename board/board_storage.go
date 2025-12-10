package board

import (
	utils "github.com/ghjm/advent_utils"
	"golang.org/x/exp/constraints"
)

// BoardStorage is an interface to pluggable back-end storage for a Board
type BoardStorage[KT constraints.Integer, VT any] interface {
	Allocate(width, height KT, emptyVal VT)
	Set(p utils.Point[KT], v VT)
	Get(p utils.Point[KT]) (VT, bool)
	Delete(p utils.Point[KT])
	GetOrDefault(p utils.Point[KT], def VT) VT
	Iterate(iterFunc func(p utils.Point[KT], v VT) bool)
	IterateOrdered(iterFunc func(p utils.Point[KT], v VT) bool)
	CopyToBoardStorage() BoardStorage[KT, VT]
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

func (fb *FlatBoard) IterateOrdered(iterFunc func(p utils.StdPoint, v rune) bool) {
	fb.Iterate(iterFunc)
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
