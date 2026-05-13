package utility

import (
	"fmt"
	"sync"
)

type UpdateFunc[T any] func(*T)

type SafeMap[K comparable, V any] struct {
	sync.RWMutex
	Data map[K]*V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		Data: make(map[K]*V),
	}
}

func (m *SafeMap[K, V]) SafeUpdate(key K, updates ...UpdateFunc[V]) *SafeMap[K, V] {
	m.Lock()
	defer m.Unlock()

	item, exists := m.Data[key]
	if !exists {
		var zero V
		item = &zero
		m.Data[key] = item
	}

	for _, update := range updates {
		update(item)
	}
	return m
}

// 返回副本（返回副本，避免外部修改）
func (m *SafeMap[K, V]) GetCopy(key K) (V, bool) {
	m.RLock()
	defer m.RUnlock()

	if item, exists := m.Data[key]; exists {
		return *item, true
	}
	var zero V
	return zero, false
}

func (m *SafeMap[K, V]) Delete(key K) *SafeMap[K, V] {
	m.Lock()
	defer m.Unlock()
	delete(m.Data, key)
	return m
}

func (m *SafeMap[K, V]) Keys() []K {
	m.RLock()
	defer m.RUnlock()

	keys := make([]K, 0, len(m.Data))
	for k := range m.Data {
		keys = append(keys, k)
	}
	return keys
}

func (m *SafeMap[K, V]) String() string {
	m.RLock()
	defer m.RUnlock()

	if len(m.Data) == 0 {
		return "SafeMap{}"
	}

	result := "SafeMap{\n"
	for k, v := range m.Data {
		result += fmt.Sprintf("  %v: %+v\n", k, *v)
	}
	result += "}"
	return result
}
