package inmemory

import (
	"fmt"
	"github.com/bilalislam/kvdb/pkg/kvdb"
	"log"
	"sync"
)

// Store is a simple key value store which keeps the state in-memory.
// It has the performance characteristics of go:s map structure
type Store struct {
	maxRecordSize int
	logger        *log.Logger

	sync.RWMutex
	table map[string]string
}

// Config contains the configuration properties for the inmemory store
type Config struct {
	MaxRecordSize int
	Logger        *log.Logger
}

// NewStore returns a new memory store
func NewStore(config Config) *Store {
	return &Store{
		maxRecordSize: config.MaxRecordSize,
		logger:        config.Logger,
		table:         map[string]string{},
	}
}

// Get returns the value associated with the key or a kvdb.NotFoundError if the
// key was not found, or any other error encountered
func (s *Store) Get(key string) (string, error) {
	s.RLock()
	v, ok := s.table[key]
	s.RUnlock()

	if !ok {
		return "", kvdb.NewNotFoundError(key)
	}

	return v, nil
}

// Set saves the value to the database and returns any error encountered
func (s *Store) Set(key string, value string) error {
	size := len(key) + len(value)
	if size > s.maxRecordSize {
		msg := fmt.Sprintf("key-value too big,max size: %d", s.maxRecordSize)
		return kvdb.NewBadRequestError(msg)
	}
	s.Lock()
	s.table[key] = value
	s.Unlock()
	return nil
}
