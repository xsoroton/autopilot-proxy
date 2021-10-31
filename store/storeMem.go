package store

import "sync"

// MEM implementation of store, good for local test running
type MemStore struct {
	mutex sync.Mutex
	// map[<key>]data
	store map[string][]byte
}

func NewMemStore() Store {
	return &MemStore{
		store: make(map[string][]byte),
	}
}

func (m *MemStore) Set(key string, value []byte) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.store[key] = value
	return
}

func (m *MemStore) Get(key string) (value []byte, err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	value, ok := m.store[key]
	if !ok {
		err = ErrNotFound
	}
	return
}

func (m *MemStore) Remove(key string) (err error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.store, key)
	return
}
