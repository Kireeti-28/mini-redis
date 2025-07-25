package storage

import "fmt"

type Storage struct {
	kvStore map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		kvStore: map[string]string{},
	}
}

func (s *Storage) Get(key string) (string, error) {
	value, ok := s.kvStore[key]
	if !ok {
		return "", fmt.Errorf("key %s does not exist", key)
	}

	return value, nil
}

func (s *Storage) Set(key, value string) {
	s.kvStore[key] = value
}

func (s *Storage) Delete(key string) error {
	if _, ok := s.kvStore[key]; !ok {
		return fmt.Errorf("key %s does not exist", key)
	}

	delete(s.kvStore, key)

	return nil
}

func (s *Storage) Size() int {
	return len(s.kvStore)
}

func (s *Storage) View() {
	if len(s.kvStore) == 0 {
		fmt.Printf("Storage is empty")
		return
	}

	for key, value := range s.kvStore {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
}
