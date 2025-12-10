package board

import (
	"github.com/ghjm/advent_utils"
	"golang.org/x/exp/constraints"
)

type CopyOnWriteStorage[KT constraints.Integer, VT any] struct {
	underlying BoardStorage[KT, VT]
	overlay    Map2D[KT, VT]
	emptyVal   VT
}

// NewCopyOnWriteStorage creates a new CopyOnWriteStorage from an underlying BoardStorage
func NewCopyOnWriteStorage[KT constraints.Integer, VT any](base BoardStorage[KT, VT], emptyVal VT) *CopyOnWriteStorage[KT, VT] {
	return &CopyOnWriteStorage[KT, VT]{
		underlying: base,
		overlay:    Map2D[KT, VT]{},
		emptyVal:   emptyVal,
	}
}

// Allocate is not implemented for CopyOnWriteStorage since there must always be an underlying BoardStorage
func (s *CopyOnWriteStorage[KT, VT]) Allocate(width, height KT, emptyVal VT) {
	panic("CopyOnWriteStorage does not implement Allocate")
}

// Set sets the value at a point
func (s *CopyOnWriteStorage[KT, VT]) Set(p utils.Point[KT], v VT) {
	s.overlay.Set(p, v)
}

// Get gets a value at a point
func (s *CopyOnWriteStorage[KT, VT]) Get(p utils.Point[KT]) (VT, bool) {
	v, ok := s.overlay.Get(p)
	if ok {
		return v, ok
	}
	return s.underlying.Get(p)
}

// Delete sets the value of a point to the provided emptyVal.  Note that this may be different behavior than
// the underlying BoardStorage - this point will continue to appear in Iterate, etc.
func (s *CopyOnWriteStorage[KT, VT]) Delete(p utils.Point[KT]) {
	s.Set(p, s.emptyVal)
}

// GetOrDefault returns a point, or if that point doesn't exist, a default value.  Note that points deleted
// with Delete() will still return emptyVal, not the default value.
func (s *CopyOnWriteStorage[KT, VT]) GetOrDefault(p utils.Point[KT], def VT) VT {
	v, ok := s.Get(p)
	if ok {
		return v
	} else {
		return def
	}
}

// Iterate iterates through the points with defined values, including points that have been deleted
// by Delete().
func (s *CopyOnWriteStorage[KT, VT]) Iterate(iterFunc func(p utils.Point[KT], v VT) bool) {
	overCopy := s.overlay.Copy()
	s.underlying.Iterate(func(p utils.Point[KT], v VT) bool {
		vo, ok := s.overlay.Get(p)
		var cont bool
		if ok {
			cont = iterFunc(p, vo)
			overCopy.Delete(p)
		} else {
			cont = iterFunc(p, v)
		}
		return cont
	})
	overCopy.Iterate(iterFunc)
}

// IterateOrdered iterates through known points in a deterministic order.  Note that this is more expensive
// for CopyOnWriteStorage than for other storage types.
func (s *CopyOnWriteStorage[KT, VT]) IterateOrdered(iterFunc func(p utils.Point[KT], v VT) bool) {
	underCopy := s.underlying.CopyToBoardStorage()
	s.underlying.Iterate(func(p utils.Point[KT], v VT) bool {
		underCopy.Set(p, v)
		return true
	})
	underCopy.IterateOrdered(iterFunc)
}

// CopyToBoardStorage creates a copy of this object's data
func (s *CopyOnWriteStorage[KT, VT]) CopyToBoardStorage() BoardStorage[KT, VT] {
	return &CopyOnWriteStorage[KT, VT]{
		underlying: s.underlying.CopyToBoardStorage(),
		overlay:    s.overlay.Copy(),
		emptyVal:   s.emptyVal,
	}
}
