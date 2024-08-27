package orderedmap

import (
	"bytes"
	"encoding/json"
	"iter"
	"sort"

	"github.com/pingooio/stdx/yaml"
)

type Item[K comparable, V any] struct {
	Key   K
	Value V
}

// Map is wrapper for map that keeps it's order when deserializing from JSON
// Warning: it's highly inneficient and should only be used for configuration file or similar use cases
type Map[K comparable, V any] struct {
	items []Item[K, V]
	data  map[K]V
}

func (m *Map[K, V]) UnmarshalJSON(data []byte) (err error) {
	err = json.Unmarshal(data, &m.data)
	if err != nil {
		return
	}

	m.items = make([]Item[K, V], 0, len(m.data))

	index := make(map[K]int)
	for key, value := range m.data {
		m.items = append(m.items, Item[K, V]{Key: key, Value: value})
		esc, _ := json.Marshal(key) //Escape the key
		index[key] = bytes.Index(data, esc)
	}

	sort.Slice(m.items, func(i, j int) bool { return index[m.items[i].Key] < index[m.items[j].Key] })
	return nil
}

func (m *Map[K, V]) UnmarshalYAML(yamlNode *yaml.Node) (err error) {
	err = yamlNode.Decode(&m.data)
	if err != nil {
		return
	}

	m.items = make([]Item[K, V], 0, len(m.data))

	index := make(map[K]int)
	for key, value := range m.data {
		m.items = append(m.items, Item[K, V]{Key: key, Value: value})
		esc, _ := yaml.Marshal(key) //Escape the key
		index[key] = bytes.Index([]byte(yamlNode.Value), esc)
	}

	sort.Slice(m.items, func(i, j int) bool { return index[m.items[i].Key] < index[m.items[j].Key] })
	return nil
}

func (m *Map[K, V]) Items() []Item[K, V] {
	return m.items
}

func (m *Map[K, V]) Len() int {
	return len(m.items)
}

func (m *Map[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, item := range m.items {
			if !yield(item.Key, item.Value) {
				return
			}
		}
	}
}
