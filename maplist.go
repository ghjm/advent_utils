package utils

import "golang.org/x/exp/maps"

type MapList[KT comparable, VT any] struct {
	data map[KT][]VT
}

// Add appends an item to the list for a given map key
func (ml *MapList[KT, VT]) Add(k KT, v VT) {
	if ml.data == nil {
		ml.data = make(map[KT][]VT)
	}
	_, ok := ml.data[k]
	if !ok {
		ml.data[k] = []VT{}
	}
	ml.data[k] = append(ml.data[k], v)
}

// Get returns the array at a given key, or nil if the key doesn't exist
func (ml *MapList[KT, VT]) Get(k KT) []VT {
	if ml.data == nil {
		return nil
	}
	v, ok := ml.data[k]
	if !ok {
		return nil
	}
	return v
}

// Remove removes a key from the map, removing its whole list
func (ml *MapList[KT, VT]) Remove(k KT) {
	if ml.data == nil {
		return
	}
	delete(ml.data, k)
}

// Clear empties the whole map
func (ml *MapList[KT, VT]) Clear() {
	ml.data = nil
}

// Keys returns the keys of the map
func (ml *MapList[KT, VT]) Keys() []KT {
	return maps.Keys(ml.data)
}

// Contains returns true if the key is in the map
func (ml *MapList[KT, VT]) Contains(k KT) bool {
	if ml.data == nil {
		return false
	}
	_, ok := ml.data[k]
	return ok
}

// Len returns the number of keys in the map
func (ml *MapList[KT, VT]) Len() int {
	if ml.data == nil {
		return 0
	}
	return len(ml.data)
}

// Count returns the number of data values in the map
func (ml *MapList[KT, VT]) Count() int {
	if ml.data == nil {
		return 0
	}
	count := 0
	for _, v := range ml.data {
		count += len(v)
	}
	return count
}
