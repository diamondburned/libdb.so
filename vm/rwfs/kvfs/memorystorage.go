package kvfs

import (
	"io/fs"
	"strings"
	"sync"
)

// MemoryStorage is a simple non-persistent key-value store.
func MemoryStorage() Store {
	return &memoryStorage{
		m: make(map[string]StoredValue),
	}
}

type memoryStorage struct {
	mu sync.RWMutex
	m  map[string]StoredValue
}

var _ Store = (*memoryStorage)(nil)

func (s *memoryStorage) Get(fullpath string) (StoredValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, ok := s.m[fullpath]
	if !ok {
		return nil, fs.ErrNotExist
	}

	return v, nil
}

func (s *memoryStorage) Set(fullpath string, v StoredValue) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[fullpath] = v

	return nil
}

func (s *memoryStorage) Delete(fullpath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.m, fullpath)

	return nil
}

func (s *memoryStorage) List(fullpath string, recursive bool) ([]PathedStoreValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var values []PathedStoreValue
	for filepath, value := range s.m {
		if !strings.HasPrefix(filepath, fullpath) {
			continue
		}

		filename := strings.TrimPrefix(filepath, fullpath)
		if !recursive && strings.Contains(filename, "/") {
			continue
		}

		values = append(values, PathedStoreValue{
			Path:        filepath,
			StoredValue: value,
		})
	}

	return values, nil
}
