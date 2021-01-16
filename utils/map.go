package utils

import "sync"

// Pair is an element in the map.
type Pair struct {
	Key  string
	Vaue interface{}
}

// Map type that can be safely shared between
// goroutines that require read/write access to a map
type Map struct {
	sync.RWMutex
	items map[string]interface{}
}

// Set sets the given (key, value) pair in the concurrent map.
func (m *Map) Set(key string, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.items[key] = value
}

// Get retrieves a key from a concurrent map.
func (m *Map) Get(key string) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	value, ok := m.items[key]
	return value, ok
}

// Clear removes all entries from the Map.
func (m *Map) Clear() {
	m.Lock()
	defer m.Unlock()
	m.items = make(map[string]interface{})
}

// Iter iterates over the items in a concurrent map and sends
// each of them over a channel, so that callers can iterate over
// the map contents using the built-in range operator.
func (m *Map) Iter() <-chan Pair {
	c := make(chan Pair)
	go func() {
		m.RLock()
		defer m.RUnlock()
		for k, v := range m.items {
			c <- Pair{k, v}
		}
		close(c)
	}()
	return c
}
