package storage

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var storageLogger = slog.Default()

type Storage struct {
	kvStore map[string]string
	mu      sync.RWMutex
	logFile *os.File
}

func NewStorage(filename string) (*Storage, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %v", filename, err)
	}

	s := &Storage{
		kvStore: make(map[string]string),
		mu:      sync.RWMutex{},
		logFile: file,
	}

	err = s.replayLogs(filename)
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("failed to replay logs: %v", err)
	}

	return s, nil
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

func (s *Storage) GetAll() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.kvStore
}

func (s *Storage) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.kvStore[key] = value

	_, err := s.logFile.WriteString(fmt.Sprintf("SET,%s,%s\n", key, value))
	if err != nil {
		return fmt.Errorf("failed to write to log file: %v", err)
	}

	return nil
}

func (s *Storage) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.kvStore[key]; !ok {
		return fmt.Errorf("key %s does not exist", key)
	}
	delete(s.kvStore, key)

	_, err := s.logFile.WriteString(fmt.Sprintf("DEL,%s\n", key))
	if err != nil {
		return fmt.Errorf("failed to write to log file: %v", err)
	}

	return nil
}

func (s *Storage) replayLogs(filename string) error {
	storageLogger.Info("Replaying logs from file", "filename", filename)

	s.logFile.Seek(0, 0) // Reset file pointer to the beginning
	scanner := bufio.NewScanner(s.logFile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 2 {
			storageLogger.Info("Invalid log entry", "line", line)
			continue
		}

		switch parts[0] {
		case "SET":
			if len(parts) != 3 {
				storageLogger.Info("Invalid SET entry", "line", line)
				continue
			}
			s.kvStore[parts[1]] = parts[2]
		case "DEL":
			if len(parts) != 2 {
				storageLogger.Info("Invalid DEL entry", "line", line)
				continue
			}
			delete(s.kvStore, parts[1])
		default:
			storageLogger.Info("Unknown command in log", "line", line)
			continue
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	storageLogger.Info("Finished replaying logs", "count", len(s.kvStore))
	return nil
}

func (s *Storage) Close() error {
	return s.logFile.Close()
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
