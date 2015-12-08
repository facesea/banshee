package tsdb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Batch is a write batch.
type Batch struct {
	// DB handle
	db *DB
	// LevelDB batch
	lb *leveldb.Batch
}

// NewBatch creates a write batch on a DB instance.
func NewBatch(db *DB) *Batch {
	b := new(Batch)
	b.db = db
	b.lb = new(leveldb.Batch)
	return b
}

// Add Put operation to the write batch.
func (b *Batch) Put(name string, t uint64, v float64) {
	//
}
