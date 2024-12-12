package utils

type DefaultMap[K comparable, D any] interface {
	Get(K) D
	Set(K, D)
	Delete(K)
}

type defaultMap[K comparable, D any] struct {
	data         map[K]D
	defaultValue D
}

func NewDefaultMap[K comparable, D any](defaultValue D) DefaultMap[K, D] {
	return &defaultMap[K, D]{
		data:         make(map[K]D),
		defaultValue: defaultValue,
	}
}

func (m *defaultMap[K, D]) Get(key K) D {
	if m.data == nil {
		m.data = make(map[K]D)
	}
	v, ok := m.data[key]
	if ok {
		return v
	}
	return m.defaultValue
}

func (m *defaultMap[K, D]) Set(key K, value D) {
	if m.data == nil {
		m.data = make(map[K]D)
	}
	m.data[key] = value
}

func (m *defaultMap[K, D]) Delete(key K) {
	if m.data == nil {
		m.data = make(map[K]D)
	}
	delete(m.data, key)
}

// GetOrDefault provides default map functionality to existing maps that aren't a DefaultMap
func GetOrDefault[K comparable, D any](d map[K]D, key K, defaultValue D) D {
	if v, ok := d[key]; ok {
		return v
	}
	return defaultValue
}
