package utils

type Map2D[K1 comparable, K2 comparable, V any] struct {
	data map[K1]map[K2]V
}

func (m2 *Map2D[K1, K2, V]) Set(k1 K1, k2 K2, v V) {
	if m2.data == nil {
		m2.data = make(map[K1]map[K2]V)
	}
	if m2.data[k1] == nil {
		m2.data[k1] = make(map[K2]V)
	}
	m2.data[k1][k2] = v
}

func (m2 *Map2D[K1, K2, V]) Get(k1 K1, k2 K2) (V, bool) {
	if m2.data == nil {
		var zv V
		return zv, false
	}
	m, ok := m2.data[k1]
	if !ok {
		var zv V
		return zv, false
	}
	v, ok := m[k2]
	if !ok {
		var zv V
		return zv, false
	}
	return v, true
}

func (m2 *Map2D[K1, K2, V]) Contains(k1 K1, k2 K2) bool {
	_, ok := m2.Get(k1, k2)
	return ok
}

func (m2 *Map2D[K1, K2, V]) GetOrDefault(k1 K1, k2 K2, def V) V {
	v, ok := m2.Get(k1, k2)
	if ok {
		return v
	} else {
		return def
	}
}

func (m2 *Map2D[K1, K2, V]) MustGet(k1 K1, k2 K2) V {
	v, ok := m2.Get(k1, k2)
	if ok {
		return v
	} else {
		panic("missing map value")
	}
}
