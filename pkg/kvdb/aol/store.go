package aol

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"

	"github.com/bilalislam/kvdb/pkg/kvdb"
	"github.com/bilalislam/kvdb/pkg/kvdb/record"
)

const (
	logFile              = "store.db"  // may be passing as a env
	defaultMaxRecordSize = 1024 * 1024 //1Mb
	defaultAsync         = false
)

var voidLogger = log.New(ioutil.Discard, "", log.LstdFlags)

type Config struct {
	BasePath      string
	MaxRecordSize *int
	Async         *bool
	Logger        *log.Logger
}

type Store struct {
	storagePath   string
	maxRecordSize int
	logger        *log.Logger
	async         bool
	index         *index
	writeMutex    sync.Mutex
}

type index struct {
	mutex  sync.RWMutex
	table  map[string]string
	buffer map[string]string
}

func (i *index) get(key string) (string, bool) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	val, ok := i.table[key]
	return val, ok
}

func (i *index) put(key string, value string, bootLoader bool) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.table[key] = value
	if !bootLoader {
		i.buffer[key] = value
	}
}

func buildIndex(filePath string, maxRecordSize int, bootLoader bool) (*index, error) {
	idx := index{
		table:  map[string]string{},
		buffer: map[string]string{},
	}

	f, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	scanner, err := record.NewScanner(f, maxRecordSize)
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		record := scanner.Record()
		idx.put(record.Key(), string(record.Value()), bootLoader)
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("could not scan entry, %w", err)
	}

	return &idx, nil
}

func NewStore(config Config) (*Store, error) {

	var (
		maxRecordSize = defaultMaxRecordSize
		storagePath   = path.Join(config.BasePath, logFile)
		async         = defaultAsync
		logger        = voidLogger
	)

	if _, err := os.OpenFile(storagePath, os.O_CREATE, 0600); err != nil {
		return nil, err
	}

	if config.MaxRecordSize != nil {
		maxRecordSize = *config.MaxRecordSize
	}

	if config.Async != nil {
		async = *config.Async
	}

	if config.Logger != nil {
		logger = config.Logger
	}

	idx, err := buildIndex(storagePath, maxRecordSize, true)
	if err != nil {
		return nil, err
	}

	logger.Printf("Index rebuilt with %d records", len(idx.table))

	return &Store{
		storagePath:   storagePath,
		maxRecordSize: maxRecordSize,
		async:         async,
		logger:        logger,
		index:         idx,
	}, nil
}

func (s *Store) Get(key string) ([]byte, error) {
	value, ok := s.index.get(key)
	if !ok {
		return nil, kvdb.NewNotFoundError(key)
	}

	return []byte(value), nil
}

func (s *Store) Set(key string, value []byte) error {
	s.index.put(key, string(value), false)
	return nil
}

func (s *Store) Sync() error {
	s.index.mutex.RLock()
	defer s.index.mutex.RUnlock()
	for key, value := range s.index.buffer {
		record := record.NewValue(key, []byte(value))
		if record.Size() > s.maxRecordSize {
			msg := fmt.Sprintf("key-value too big,max size : %d", s.maxRecordSize)
			return kvdb.NewBadRequestError(msg)
		}

		s.writeMutex.Lock()
		defer s.writeMutex.Unlock()

		file, err := os.OpenFile(s.storagePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		defer file.Close()
		if err != nil {
			return fmt.Errorf("could not open file : %s for write, %w", s.storagePath, err)
		}

		_, err = record.Write(file)
		if err != nil {
			return fmt.Errorf("could not write record to file %s, %w", s.storagePath, err)
		}

		if !s.async {
			if err := file.Sync(); err != nil {
				return err
			}
		}

		if err := file.Close(); err != nil {
			return err
		}
	}

	s.index.buffer = map[string]string{}
	return nil
}

func (s *Store) Close() error {
	s.logger.Print("Closing database")
	return nil
}

func (s *Store) IsNotFoundError(err error) bool {
	return kvdb.IsNotFoundError(err)
}

func (s *Store) IsBadRequestError(err error) bool {
	return kvdb.IsBadRequestError(err)
}
