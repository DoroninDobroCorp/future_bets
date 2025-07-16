package bikeymap

import "sync"

type key interface {
	int | string
}

type BiKeyMapInterface[K key, V any] interface {
	WriteBothKeys(key1 K, key2 K, value V)
	ReadAll() (map[K]K, map[K]V)
	// Read(key K) (key1, key2 *K, value *V)
	ReadKey(key K) (K, bool)
	Delete(key K)
	Len() (keysLen, valuesLen int)
}

type BiKeyMap[K key, V any] struct {
	sync.RWMutex
	keys   map[K]K
	values map[K]V
}

func NewBiKeyMap[K key, V any]() BiKeyMapInterface[K, V] {
	return &BiKeyMap[K, V]{
		keys:   make(map[K]K),
		values: make(map[K]V),
	}
}

func (m *BiKeyMap[K, V]) WriteBothKeys(key1 K, key2 K, value V) {
	m.Lock()
	defer m.Unlock()

	m.keys[key1] = key2
	m.keys[key2] = key1

	m.values[key1] = value
	m.values[key2] = value
}

func (m *BiKeyMap[K, V]) ReadKey(key K) (K, bool) {
	m.RLock()
	defer m.RUnlock()

	key2, ok := m.keys[key]

	return key2, ok
}

// func (m *BiKeyMap[K, V]) Read(key K) (key1, key2 *K, value *V) {
// 	m.RLock()
// 	defer m.RUnlock()

// 	key2, ok := m.keys[key]
// 	if !ok {
// 		return nil, nil, nil
// 	}

// 	value, ok = m.values[key]
// 	if !ok {
// 		return &key, key2, nil
// 	}

// 	return &key, key2, value
// }

func (m *BiKeyMap[K, V]) Len() (keysLen, valuesLen int) {
	m.RLock()
	defer m.RUnlock()

	return len(m.keys), len(m.values)
}

func (m *BiKeyMap[K, V]) Delete(key K) {
	m.Lock()
	defer m.Unlock()

	key2, ok := m.keys[key]
	if ok {
		delete(m.keys, key2)
		delete(m.values, key2)
	}

	delete(m.keys, key)
	delete(m.values, key)
}

func (m *BiKeyMap[K, V]) ReadAll() (map[K]K, map[K]V) {
	m.RLock()
	defer m.RUnlock()

	newMapKeys := make(map[K]K, len(m.keys))
	for k, v := range m.keys {
		newMapKeys[k] = v
	}

	newMapValues := make(map[K]V, len(m.values))
	for k, v := range m.values {
		newMapValues[k] = v
	}

	return newMapKeys, newMapValues
}
