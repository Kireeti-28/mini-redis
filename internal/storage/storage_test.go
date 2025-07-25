package storage

import (
	"fmt"
	"sync"
	"testing"
)

func TestStorageSet(t *testing.T) {
	s := NewStorage()
	s.Set("foo", "bar")

	val, err := s.Get("foo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if val != "bar" {
		t.Errorf("expected value 'bar', got '%s'", val)
	}
}

func TestStorageSetOverwrite(t *testing.T) {
	s := NewStorage()
	s.Set("foo", "bar")
	s.Set("foo", "baz")

	val, err := s.Get("foo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if val != "baz" {
		t.Errorf("expected value 'baz', got '%s'", val)
	}
}

func TestStorageSetEmptyKey(t *testing.T) {
	s := NewStorage()
	s.Set("", "empty")
	val, err := s.Get("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if val != "empty" {
		t.Errorf("expected value 'empty', got '%s'", val)
	}
}

func TestStorageGet(t *testing.T) {
	s := NewStorage()
	s.Set("foo", "bar")

	val, err := s.Get("foo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if val != "bar" {
		t.Errorf("expected value 'bar', got '%s'", val)
	}
}

func TestStorageGetNonExistentKey(t *testing.T) {
	s := NewStorage()
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent key, got nil")
	}
	expectedErr := "key nonexistent does not exist"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestStorageDelete(t *testing.T) {
	s := NewStorage()
	s.Set("foo", "bar")

	err := s.Delete("foo")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = s.Get("foo")
	if err == nil {
		t.Fatal("expected error for deleted key, got nil")
	}
	expectedErr := "key foo does not exist"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestStorageDeleteNonExistentKey(t *testing.T) {
	s := NewStorage()
	err := s.Delete("nonexistent")
	if err == nil {
		t.Fatal("expected error for non-existent key, got nil")
	}
	expectedErr := "key nonexistent does not exist"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

// test for concurrent write and read operations
func TestStorageConcurrentAccess(t *testing.T) {
	s := NewStorage()
	var wg sync.WaitGroup

	// Start concurrent writes
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			s.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
		}
	}()

	// Start concurrent reads
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			// It's possible some keys aren't written yet, so ignore missing keys
			val, err := s.Get(fmt.Sprintf("key%d", i))
			if err == nil {
				expected := fmt.Sprintf("value%d", i)
				if val != expected {
					t.Errorf("expected %s, got %s", expected, val)
				}
			}
		}
	}()

	wg.Wait()

	if s.Size() != 100 {
		t.Errorf("expected size 100, got %d", s.Size())
	}
}
