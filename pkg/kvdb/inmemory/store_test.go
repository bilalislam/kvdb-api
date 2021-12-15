package inmemory

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestStore(t *testing.T) {
	var (
		key = "key"
		val = "value"
	)

	store := NewStore(Config{
		MaxRecordSize: 100,
		Logger:        &log.Logger{},
	})

	err := store.Set(key, val)
	require.NoError(t, err)

	getVal, err := store.Get(key)
	require.NoError(t, err)
	require.Equal(t, val, getVal)
}
