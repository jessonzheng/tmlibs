package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cmn "github.com/tendermint/tmlibs/common"
)

func testBackendGetSetDelete(t *testing.T, backend string) {
	// Default
	dir, dirname := cmn.Tempdir(fmt.Sprintf("test_backend_%s_", backend))
	defer dir.Close()
	db := NewDB("testdb", backend, dirname)
	key := []byte("abc")
	require.Nil(t, db.Get(key))

	// Set empty ("")
	db.Set(key, []byte(""))
	require.NotNil(t, db.Get(key))
	require.Empty(t, db.Get(key))

	// Set empty (nil)
	db.Set(key, nil)
	require.NotNil(t, db.Get(key))
	require.Empty(t, db.Get(key))

	// Delete
	db.Delete(key)
	require.Nil(t, db.Get(key))
}

func TestBackendsGetSetDelete(t *testing.T) {
	for dbType, _ := range backends {
		testBackendGetSetDelete(t, dbType)
	}
}

func assertPanics(t *testing.T, dbType, name string, fn func()) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, cmn.Fmt("expecting %s.%s to panic", dbType, name))
	}()

	fn()
}

func TestBackendsNilKeys(t *testing.T) {
	// test all backends
	for dbType, creator := range backends {
		name := cmn.Fmt("test_%x", cmn.RandStr(12))
		db, err := creator(name, "")
		assert.Nil(t, err)

		assertPanics(t, dbType, "get", func() { db.Get(nil) })
		assertPanics(t, dbType, "has", func() { db.Has(nil) })
		assertPanics(t, dbType, "set", func() { db.Set(nil, []byte("abc")) })
		assertPanics(t, dbType, "setsync", func() { db.SetSync(nil, []byte("abc")) })
		assertPanics(t, dbType, "delete", func() { db.Delete(nil) })
		assertPanics(t, dbType, "deletesync", func() { db.DeleteSync(nil) })

		db.Close()
		err = os.RemoveAll(name + ".db")
		assert.Nil(t, err)
	}
}

func TestGoLevelDBBackendStr(t *testing.T) {
	name := cmn.Fmt("test_%x", cmn.RandStr(12))
	db := NewDB(name, LevelDBBackendStr, "")
	defer os.RemoveAll(name + ".db")

	if _, ok := backends[CLevelDBBackendStr]; !ok {
		_, ok := db.(*GoLevelDB)
		assert.True(t, ok)
	}
}
