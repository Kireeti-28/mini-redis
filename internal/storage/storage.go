package storage

import (
	"fmt"
	"sync"
)

type Storage struct {
	kvStore map[string]string
	mu      sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		kvStore: map[string]string{},
		mu:      sync.RWMutex{},
	}
}

func (s *Storage) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.kvStore[key]
	if !ok {
		return "", fmt.Errorf("key %s does not exist", key)
	}

	return value, nil
}

func (s *Storage) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.kvStore[key] = value
}

func (s *Storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.kvStore[key]; !ok {
		return fmt.Errorf("key %s does not exist", key)
	}
	delete(s.kvStore, key)

	return nil
}

func (s *Storage) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.kvStore)
}

func (s *Storage) View() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.kvStore) == 0 {
		fmt.Printf("Storage is empty")
		return
	}

	for key, value := range s.kvStore {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
}
