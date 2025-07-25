package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const filename = "storage.json"

type Storage struct {
	kvStore map[string]string
	mu      sync.RWMutex
}

func NewStorage() *Storage {
	s := &Storage{
		kvStore: make(map[string]string),
		mu:      sync.RWMutex{},
	}

	err := s.load()
	if err != nil {
		fmt.Printf("Error loading storage: %v\n", err)
	}

	return s
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
	s.persist()
}

func (s *Storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.kvStore[key]; !ok {
		return fmt.Errorf("key %s does not exist", key)
	}
	delete(s.kvStore, key)
	s.persist()

	return nil
}

func (s *Storage) persist() error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(s.kvStore)
	if err != nil {
		return fmt.Errorf("failed to encode JSON to file %s: %v", filename, err)
	}

	return nil
}

func (s *Storage) load() error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File does not exist, nothing to load
		}
		return fmt.Errorf("failed to open file %s: %v", filename, err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&s.kvStore)
	if err != nil {
		return fmt.Errorf("failed to decode JSON from file %s: %v", filename, err)
	}

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
	fmt.Println()
}
