package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"golang.org/x/exp/constraints"
	"sort"
)

type Map2D[KT constraints.Integer, VT any] struct {
	data       map[Point[KT]]VT
	boundsSet  bool
	boundsLow  Point[KT]
	boundsHigh Point[KT]
}

type Hashable interface {
	HashString() string
}

type Map2DHashable[KT constraints.Integer, VT Hashable] struct {
	Map2D[KT, VT]
}

func (m2 *Map2D[KT, VT]) Set(p Point[KT], v VT) {
	if m2.data == nil {
		m2.data = make(map[Point[KT]]VT)
	}
	m2.data[p] = v
}

func (m2 *Map2D[KT, VT]) Get(p Point[KT]) (VT, bool) {
	if m2.data == nil {
		var zv VT
		return zv, false
	}
	v, ok := m2.data[p]
	if !ok {
		var zv VT
		return zv, false
	}
	return v, true
}

func (m2 *Map2D[KT, VT]) Delete(p Point[KT]) {
	if m2.data == nil {
		return
	}
	delete(m2.data, p)
}

func (m2 *Map2D[KT, VT]) Contains(p Point[KT]) bool {
	_, ok := m2.Get(p)
	return ok
}

func (m2 *Map2D[KT, VT]) GetOrDefault(p Point[KT], def VT) VT {
	v, ok := m2.Get(p)
	if ok {
		return v
	} else {
		return def
	}
}

func (m2 *Map2D[KT, VT]) MustGet(p Point[KT]) VT {
	v, ok := m2.Get(p)
	if ok {
		return v
	} else {
		panic("missing map value")
	}
}

func (m2 *Map2D[KT, VT]) Len() int {
	return len(m2.data)
}

func (m2 *Map2D[KT, VT]) Iterate(iterFunc func(p Point[KT], v VT) bool) {
	for k, v := range m2.data {
		if !iterFunc(k, v) {
			return
		}
	}
}

func (m2 *Map2D[KT, VT]) IterateOrdered(iterFunc func(p Point[KT], v VT) bool) {
	type tuple = struct {
		k Point[KT]
		v VT
	}
	var data []tuple
	for k, v := range m2.data {
		data = append(data, tuple{k, v})
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].k.Y < data[j].k.Y || (data[i].k.Y == data[j].k.Y && data[i].k.X < data[j].k.X)
	})
	for _, t := range data {
		if !iterFunc(t.k, t.v) {
			return
		}
	}
}

func (m2 *Map2D[KT, VT]) Copy() Map2D[KT, VT] {
	c := Map2D[KT, VT]{}
	m2.Iterate(func(p Point[KT], v VT) bool {
		c.Set(p, v)
		return true
	})
	return c
}

func (m2 *Map2DHashable[KT, VT]) Hash() uint64 {
	s := sha256.New()
	m2.IterateOrdered(func(p Point[KT], v VT) bool {
		s.Write([]byte(fmt.Sprintf("(%d,%d)", p.X, p.Y)))
		s.Write([]byte{0})
		s.Write([]byte(v.HashString()))
		s.Write([]byte{0})
		return true
	})
	return binary.BigEndian.Uint64(s.Sum(nil))
}
