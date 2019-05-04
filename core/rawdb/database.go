package rawdb

import (
	"github.com/altair-lab/xoreum/xordb"
	"github.com/altair-lab/xoreum/xordb/leveldb"
	"github.com/altair-lab/xoreum/xordb/memorydb"
)

// NewDatabase creates a high level database on top of a given key-value data
func NewDatabase(db xordb.KeyValueStore) xordb.Database {
	return db
}

// NewMemoryDatabase creates an in-memory key-value database
func NewMemoryDatabase() xordb.Database {
	return NewDatabase(memorydb.New())
}

// NewMemoryDatabaseWithCap creates an in-memory key-value database with
// an initial starting capacity
func NewMemoryDatabaseWithCap(size int) xordb.Database {
	return NewDatabase(memorydb.NewWithCap(size))
}

// NewLevelDBDatabase creates a persistent key-value database
func NewLevelDBDatabase(file string, cache int, handles int, namespace string) (xordb.Database, error) {
	db, err := leveldb.New(file, cache, handles, namespace)
	if err != nil {
		return nil, err
	}
	return NewDatabase(db), nil
}
