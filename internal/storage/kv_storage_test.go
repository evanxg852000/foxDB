package storage

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/dgraph-io/badger/v3"
)

func setupTestDB(t *testing.T) (*KvStorage, func()) {
	tmpDir, err := os.MkdirTemp("", "foxdb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	storage, err := NewKvStorage(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create storage: %v", err)
	}

	cleanup := func() {
		storage.Remove()
		os.RemoveAll(tmpDir)
	}

	return storage, cleanup
}

func TestNewKvStorage(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "foxdb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storage, err := NewKvStorage(tmpDir)
	if err != nil {
		t.Fatalf("NewKvStorage failed: %v", err)
	}
	defer storage.Close()

	if storage.db == nil {
		t.Error("Expected db to be initialized")
	}
}

func TestNewKvStorageInvalidPath(t *testing.T) {
	invalidPath := "/invalid/path/that/does/not/exist"
	_, err := NewKvStorage(invalidPath)
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestKvStorageSetAndGet(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	testCases := []struct {
		key   []byte
		value []byte
	}{
		{[]byte("key1"), []byte("value1")},
		{[]byte("key2"), []byte("value2")},
		{[]byte("single_char"), []byte("x")},
		{[]byte("special_chars"), []byte("value with spaces and symbols !@#$%")},
		{[]byte("binary_data"), []byte{0x00, 0x01, 0x02, 0xFF}},
	}

	// Test Set operations
	for _, tc := range testCases {
		err := storage.Set(tc.key, tc.value)
		if err != nil {
			t.Errorf("Set failed for key %q: %v", tc.key, err)
		}
	}

	// Test Get operations
	for _, tc := range testCases {
		value, err := storage.Get(tc.key)
		if err != nil {
			t.Errorf("Get failed for key %q: %v", tc.key, err)
			continue
		}
		if !bytes.Equal(value, tc.value) {
			t.Errorf("Get returned wrong value for key %q: expected %q, got %q", tc.key, tc.value, value)
		}
	}
}

func TestKvStorageGetNonExistentKey(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	_, err := storage.Get([]byte("non_existent_key"))
	if err == nil {
		t.Error("Expected error for non-existent key, got nil")
	}
}

func TestKvStorageEmptyKey(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Test that empty keys are rejected
	err := storage.Set([]byte(""), []byte("value"))
	if err == nil {
		t.Error("Expected error for empty key, got nil")
	}

	_, err = storage.Get([]byte(""))
	if err == nil {
		t.Error("Expected error for empty key get, got nil")
	}
}

func TestKvStorageDelete(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	key := []byte("test_key")
	value := []byte("test_value")

	// Set a key-value pair
	err := storage.Set(key, value)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Verify it exists
	_, err = storage.Get(key)
	if err != nil {
		t.Fatalf("Get failed after Set: %v", err)
	}

	// Delete the key
	err = storage.Delete(key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify it's deleted
	_, err = storage.Get(key)
	if err == nil {
		t.Error("Expected error after delete, got nil")
	}
}

func TestKvStorageDeleteNonExistentKey(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Deleting non-existent key should not return error
	err := storage.Delete([]byte("non_existent_key"))
	if err != nil {
		t.Errorf("Delete non-existent key should not fail: %v", err)
	}
}

func TestKvStorageSync(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	err := storage.Set([]byte("key"), []byte("value"))
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	err = storage.Sync()
	if err != nil {
		t.Errorf("Sync failed: %v", err)
	}
}

func TestKvStorageClose(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "foxdb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storage, err := NewKvStorage(tmpDir)
	if err != nil {
		t.Fatalf("NewKvStorage failed: %v", err)
	}

	err = storage.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Trying to use storage after close should fail
	err = storage.Set([]byte("key"), []byte("value"))
	if err == nil {
		t.Error("Expected error when using storage after close")
	}
}

func TestKvStorageRemove(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "foxdb_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	storage, err := NewKvStorage(tmpDir)
	if err != nil {
		t.Fatalf("NewKvStorage failed: %v", err)
	}

	// Add some data
	err = storage.Set([]byte("key"), []byte("value"))
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Remove the database
	err = storage.Remove()
	if err != nil {
		t.Errorf("Remove failed: %v", err)
	}

	// Verify directory is removed
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		t.Error("Database directory should be removed")
	}
}

func TestKvStorageBatch(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Test batch operations
	err := storage.Batch(func(txn *badger.Txn) error {
		if err := txn.Set([]byte("batch_key1"), []byte("batch_value1")); err != nil {
			return err
		}
		if err := txn.Set([]byte("batch_key2"), []byte("batch_value2")); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Batch operation failed: %v", err)
	}

	// Verify both keys were set
	value1, err := storage.Get([]byte("batch_key1"))
	if err != nil {
		t.Errorf("Get batch_key1 failed: %v", err)
	} else if !bytes.Equal(value1, []byte("batch_value1")) {
		t.Errorf("Wrong value for batch_key1: expected %q, got %q", "batch_value1", value1)
	}

	value2, err := storage.Get([]byte("batch_key2"))
	if err != nil {
		t.Errorf("Get batch_key2 failed: %v", err)
	} else if !bytes.Equal(value2, []byte("batch_value2")) {
		t.Errorf("Wrong value for batch_key2: expected %q, got %q", "batch_value2", value2)
	}
}

func TestKvStorageBatchError(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Test batch with error - transaction should be rolled back
	err := storage.Batch(func(txn *badger.Txn) error {
		if err := txn.Set([]byte("batch_key1"), []byte("batch_value1")); err != nil {
			return err
		}
		// Return an error to simulate failure
		return fmt.Errorf("simulated error")
	})

	if err == nil {
		t.Error("Expected batch to fail with error")
	}

	// Verify key was not set due to rollback
	_, err = storage.Get([]byte("batch_key1"))
	if err == nil {
		t.Error("Key should not exist after failed batch transaction")
	}
}

func TestKvScanBasic(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Set up test data with common prefix
	testData := map[string]string{
		"user:1": "alice",
		"user:2": "bob",
		"user:3": "charlie",
		"item:1": "laptop",
		"item:2": "mouse",
	}

	for key, value := range testData {
		err := storage.Set([]byte(key), []byte(value))
		if err != nil {
			t.Fatalf("Set failed for %s: %v", key, err)
		}
	}

	// Test scanning with "user:" prefix
	scan := storage.Scan([]byte("user:"))
	defer scan.Close()

	foundItems := make(map[string]string)
	for scan.Valid() {
		key, value, err := scan.Item()
		if err != nil {
			t.Fatalf("Item() failed: %v", err)
		}
		foundItems[string(key)] = string(value)
		scan.Next()
	}

	expectedUserItems := map[string]string{
		"user:1": "alice",
		"user:2": "bob",
		"user:3": "charlie",
	}

	if len(foundItems) != len(expectedUserItems) {
		t.Errorf("Expected %d items, got %d", len(expectedUserItems), len(foundItems))
	}

	for key, expectedValue := range expectedUserItems {
		if actualValue, exists := foundItems[key]; !exists {
			t.Errorf("Key %s not found in scan results", key)
		} else if actualValue != expectedValue {
			t.Errorf("Wrong value for key %s: expected %s, got %s", key, expectedValue, actualValue)
		}
	}
}

func TestKvScanEmptyPrefix(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Set up test data
	testKeys := []string{"a", "b", "c"}
	for _, key := range testKeys {
		err := storage.Set([]byte(key), []byte("value_"+key))
		if err != nil {
			t.Fatalf("Set failed for %s: %v", key, err)
		}
	}

	// Scan with empty prefix should return all keys
	scan := storage.Scan([]byte(""))
	defer scan.Close()

	count := 0
	for scan.Valid() {
		_, _, err := scan.Item()
		if err != nil {
			t.Fatalf("Item() failed: %v", err)
		}
		count++
		scan.Next()
	}

	if count < len(testKeys) {
		t.Errorf("Expected at least %d items, got %d", len(testKeys), count)
	}
}

func TestKvScanNoMatches(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Set up test data
	err := storage.Set([]byte("key1"), []byte("value1"))
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}

	// Scan with prefix that doesn't match anything
	scan := storage.Scan([]byte("nonexistent:"))
	defer scan.Close()

	if scan.Valid() {
		t.Error("Expected no items for non-matching prefix")
	}
}

func TestKvScanIteratorLifecycle(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Set up test data
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("test:%d", i)
		value := fmt.Sprintf("value%d", i)
		err := storage.Set([]byte(key), []byte(value))
		if err != nil {
			t.Fatalf("Set failed for %s: %v", key, err)
		}
	}

	scan := storage.Scan([]byte("test:"))

	// Test multiple iterations
	count := 0
	for scan.Valid() {
		key, value, err := scan.Item()
		if err != nil {
			t.Fatalf("Item() failed: %v", err)
		}
		if len(key) == 0 || len(value) == 0 {
			t.Error("Key or value should not be empty")
		}
		count++
		scan.Next()
	}

	if count != 5 {
		t.Errorf("Expected 5 items, got %d", count)
	}

	// After iteration is complete, Valid() should return false
	if scan.Valid() {
		t.Error("Iterator should not be valid after complete iteration")
	}

	// Close should not panic
	scan.Close()
}

func TestKvScanMultipleScans(t *testing.T) {
	storage, cleanup := setupTestDB(t)
	defer cleanup()

	// Set up test data
	testData := map[string]string{
		"prefix1:a": "value1a",
		"prefix1:b": "value1b",
		"prefix2:a": "value2a",
		"prefix2:b": "value2b",
	}

	for key, value := range testData {
		err := storage.Set([]byte(key), []byte(value))
		if err != nil {
			t.Fatalf("Set failed for %s: %v", key, err)
		}
	}

	// Test multiple concurrent scans
	scan1 := storage.Scan([]byte("prefix1:"))
	scan2 := storage.Scan([]byte("prefix2:"))

	count1 := 0
	for scan1.Valid() {
		_, _, err := scan1.Item()
		if err != nil {
			t.Fatalf("Scan1 Item() failed: %v", err)
		}
		count1++
		scan1.Next()
	}

	count2 := 0
	for scan2.Valid() {
		_, _, err := scan2.Item()
		if err != nil {
			t.Fatalf("Scan2 Item() failed: %v", err)
		}
		count2++
		scan2.Next()
	}

	scan1.Close()
	scan2.Close()

	if count1 != 2 {
		t.Errorf("Expected 2 items for prefix1, got %d", count1)
	}
	if count2 != 2 {
		t.Errorf("Expected 2 items for prefix2, got %d", count2)
	}
}

// Benchmark tests
func BenchmarkKvStorageSet(b *testing.B) {
	storage, cleanup := setupTestBench(b)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		value := fmt.Sprintf("bench_value_%d", i)
		err := storage.Set([]byte(key), []byte(value))
		if err != nil {
			b.Fatalf("Set failed: %v", err)
		}
	}
}

func BenchmarkKvStorageGet(b *testing.B) {
	storage, cleanup := setupTestBench(b)
	defer cleanup()

	// Setup data
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		value := fmt.Sprintf("bench_value_%d", i)
		err := storage.Set([]byte(key), []byte(value))
		if err != nil {
			b.Fatalf("Set failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i%1000)
		_, err := storage.Get([]byte(key))
		if err != nil {
			b.Fatalf("Get failed: %v", err)
		}
	}
}

func BenchmarkKvStorageScan(b *testing.B) {
	storage, cleanup := setupTestBench(b)
	defer cleanup()

	// Setup data
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("scan_prefix_%04d", i)
		value := fmt.Sprintf("value_%d", i)
		err := storage.Set([]byte(key), []byte(value))
		if err != nil {
			b.Fatalf("Set failed: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		scan := storage.Scan([]byte("scan_prefix_"))
		count := 0
		for scan.Valid() {
			_, _, err := scan.Item()
			if err != nil {
				b.Fatalf("Item failed: %v", err)
			}
			count++
			scan.Next()
		}
		scan.Close()

		if count != 1000 {
			b.Fatalf("Expected 1000 items, got %d", count)
		}
	}
}

func setupTestBench(b *testing.B) (*KvStorage, func()) {
	tmpDir, err := os.MkdirTemp("", "foxdb_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	storage, err := NewKvStorage(tmpDir)
	if err != nil {
		os.RemoveAll(tmpDir)
		b.Fatalf("Failed to create storage: %v", err)
	}

	cleanup := func() {
		storage.Remove()
		os.RemoveAll(tmpDir)
	}

	return storage, cleanup
}
