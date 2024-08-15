package orderedmap

import (
	"bytes"
	"encoding/json"
	"sort"
)

type Element[K comparable, V any] struct {
	Key   K
	Value V
}

// Map is wrapper for map that keeps it's order when deserializing from JSON
// Warning: it's highly inneficient and should only be used for configuration file or similar use cases
type Map[K comparable, V any] struct {
	elems []Element[K, V]
	data  map[K]V
}

func (m *Map[K, V]) UnmarshalJSON(data []byte) error {
	json.Unmarshal(data, &m.data)

	m.elems = make([]Element[K, V], 0, len(m.data))

	index := make(map[K]int)
	for key, value := range m.data {
		m.elems = append(m.elems, Element[K, V]{Key: key, Value: value})
		esc, _ := json.Marshal(key) //Escape the key
		index[key] = bytes.Index(data, esc)
	}

	sort.Slice(m.elems, func(i, j int) bool { return index[m.elems[i].Key] < index[m.elems[j].Key] })
	return nil
}

func (m *Map[K, V]) Elems() []Element[K, V] {
	return m.elems
}
