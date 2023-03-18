package syncg

import "sync"

// Map is a generic wrapper around sync.Map.
type Map[K comparable, V any] sync.Map

// Load returns the value stored in the Map for the given key.
func (m *Map[K, V]) Load(key K) (V, bool) {
	v, ok := (*sync.Map)(m).Load(key)
	if ok {
		return v.(V), true
	}
	var z V
	return z, false
}

// Store stores the given value in the Map for the given key.
func (m *Map[K, V]) Store(key K, val V) {
	(*sync.Map)(m).Store(key, val)
}

// Delete deletes the value stored in the Map for the given key.
func (m *Map[K, V]) Delete(key K) {
	(*sync.Map)(m).Delete(key)
}

// Range calls f sequentially for each key and value present in the Map. If f
// returns false, range stops the iteration.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	(*sync.Map)(m).Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

// LoadOrStore returns the existing value for the key if present. Otherwise, it
// stores and returns the given value. The loaded result is true if the value was
// loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, val V) (actual V, loaded bool) {
	v, loaded := (*sync.Map)(m).LoadOrStore(key, val)
	return v.(V), loaded
}

// LoadAndDelete deletes the value for the key if present. The loaded result is
// true if the value was loaded, false if the Map contains no value for the key.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := (*sync.Map)(m).LoadAndDelete(key)
	if loaded {
		return v.(V), true
	}
	var z V
	return z, false
}
