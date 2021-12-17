package aol

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleSetGet(t *testing.T) {
	testPath, err := ioutil.TempDir("./", "test")
	if err != nil {
		t.Fatal(err)
	}

	defer (func() {
		os.RemoveAll(testPath)
	})()

	const (
		testKey   = "testKey"
		testValue = "testValue"
	)

	store, err := NewStore(Config{
		BasePath: testPath,
	})

	require.NoError(t, err)
	err = store.Set(testKey, []byte(testValue))
	require.NoError(t, err)

	data, err := store.Get(testKey)
	require.NoError(t, err)
	require.Equal(t, []byte(testValue), data)
}
