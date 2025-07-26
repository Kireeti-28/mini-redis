package storage

import (
	"os"
	"testing"
)

func setupTest(t *testing.T) (string, func()) {
	file, err := os.CreateTemp("", "test_storage_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	cleanup := func() {
		file.Close()
		os.Remove(file.Name())
	}

	return file.Name(), cleanup
}

func TestStorageSetAndGet(t *testing.T) {
	filename, cleanUpfn := setupTest(t)
	defer cleanUpfn()

	storage, err := NewStorage(filename)
	if err != nil {
		t.Fatalf("Failed to create storage")
	}

	key := "test_name"
	value := "set_and_get"

	storage.Set(key, value)
	savedValue, err := storage.Get(key)

	if err != nil {
		t.Fatalf("error getting key: %s", key)
	}

	if savedValue != value {
		t.Errorf("expected value to be %s but recieved %s", value, savedValue)
	}
}

func TestStorageGetNonExist(t *testing.T) {
	filename, deferFn := setupTest(t)
	defer deferFn()

	storage, err := NewStorage(filename)
	if err != nil {
		t.Fatalf("error creating storage")
	}

	key := "name"
	_, err = storage.Get("name") //name doesn't exist

	if err.Error() != "key "+key+" does not exist" {
		t.Errorf("expected error: %s\n but got: %s\n", "key "+key+" does not exist", err.Error())
	}
}

func TestStorageSetOverWrite(t *testing.T) {
	filename, deferFn := setupTest(t)
	defer deferFn()

	storage, err := NewStorage(filename)
	if err != nil {
		t.Fatalf("error creating storage")
	}

	key := "name"
	value := "value"
	storage.Set(key, value)

	overrideValue := "overridevalue"
	storage.Set(key, overrideValue)

	savedValue, err := storage.Get(key)
	if err != nil {
		t.Fatalf("failed to retrive key %s", key)
	}

	if savedValue != overrideValue {
		t.Errorf("expected: %s\n but got: %s\n", overrideValue, savedValue)
	}
}

func TestStorageDelete(t *testing.T) {
	filename, deferFn := setupTest(t)
	defer deferFn()

	storage, err := NewStorage(filename)
	if err != nil {
		t.Fatalf("error creating storage")
	}

	key := "key"
	value := "value"

	err = storage.Set(key, value)
	if err != nil {
		t.Fatalf("failed to set key %s, value %s", key, value)
	}

	err = storage.Delete(key)
	if err != nil {
		t.Fatalf("failed to delete key %s", key)
	}

	_, err = storage.Get(key)
	if err.Error() != "key "+key+" does not exist" {
		t.Errorf("expected error: %s\n but got: %s\n", "key "+key+" does not exist", err.Error())
	}
}

func TestStoragePersistence(t *testing.T) {
	filename, deferFn := setupTest(t)
	defer deferFn()

	storage1, err := NewStorage(filename)
	if err != nil {
		t.Fatalf("error creating storage1")
	}

	nameKey := "name"
	nameValue := "alpha"
	fatherKey := "father"
	fatherValue := "beta"
	motherKey := "mother"
	motherValue := "omega"

	err = storage1.Set(nameKey, nameValue)
	if err != nil {
		t.Fatalf("failed to set key %s value %s", nameKey, nameValue)
	}
	err = storage1.Set(fatherKey, fatherValue)
	if err != nil {
		t.Fatalf("failed to set key %s value %s", fatherKey, fatherValue)
	}
	err = storage1.Set(motherKey, motherValue)
	if err != nil {
		t.Fatalf("failed to set key %s value %s", motherKey, motherValue)
	}

	storage1.Close()

	// this instance should have all key-vals of prev instance
	storage2, err := NewStorage(filename)
	if err != nil {
		t.Fatalf("failed to create storage2")
	}

	nameSValue, err := storage2.Get(nameKey)
	if err != nil {
		t.Fatalf("error getting key %s", nameKey)
	}

	if nameSValue != nameValue {
		t.Errorf("expected %s as value but got %s", nameValue, nameSValue)
	}

	fatherSValue, err := storage2.Get(fatherKey)
	if err != nil {
		t.Fatalf("error getting key %s", fatherKey)
	}

	if fatherSValue != fatherValue {
		t.Errorf("expected %s got %s", fatherValue, fatherSValue)
	}

	motherSValue, err := storage2.Get(motherKey)
	if err != nil {
		t.Fatalf("error getting key %s", motherKey)
	}

	if motherSValue != motherValue {
		t.Errorf("expected %s got %s", motherValue, motherSValue)
	}

	storage2.Close()
}
