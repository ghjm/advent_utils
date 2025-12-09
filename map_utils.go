package utils

// GetArbitraryKey returns a single arbitrary key from the given map.  The bool is returned false if the map is empty.
func GetArbitraryKey[K comparable, V any](m map[K]V) (K, bool) {
	for k := range m {
		return k, true
	}
	var k K
	return k, false
}

// MustGetArbitraryKey returns a single arbitrary key from the given map, and panics if the map is empty.
func MustGetArbitraryKey[K comparable, V any](m map[K]V) K {
	for k := range m {
		return k
	}
	panic("cannot get key from empty map")
}

// PopArbitraryKey returns and deletes a single arbitrary key from the given map.  The bool is returned false if the map is empty.
func PopArbitraryKey[K comparable, V any](m map[K]V) (K, bool) {
	k, ok := GetArbitraryKey(m)
	if !ok {
		return k, false
	}
	delete(m, k)
	return k, true
}

// MustPopArbitraryKey returns and deletes a single arbitrary key from the given map, and panics if the map is empty.
func MustPopArbitraryKey[K comparable, V any](m map[K]V) K {
	k := MustGetArbitraryKey[K, V](m)
	delete(m, k)
	return k
}
