package utils

import (
	"sort"
)

// Pair holds a key-value pair from the map
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// GetSortedKeys gets sorted list of keys
func GetSortedKeys[K comparable, V any](m *map[K]V, less func(a, b K) bool) []K {
	keys := make([]K, 0)
	for k := range *m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return less(keys[i], keys[j])
	})

	return keys
}
