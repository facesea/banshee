package tsdb

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Batch is a write batch.
type Batch struct {
	// DB handle
	db *DB
	// LevelDB batch
	b *leveldb.Batch
}

// NewBatch creates a write batch on a DB instance.
func NewBatch(db *DB) *Batch {
	b := new(Batch)
	b.db = db
	b.b = new(leveldb.Batch)
	return b
}

// Reset batch.
func (b *Batch) Reset() {
	b.b.Reset()
}

// Flush write batch to db.
func (b *Batch) Flush() {
	b.db.db.Write(b.b, nil)
}

// Put stamp with value and score to db.
func (b *Batch) PutStamp(name string, stamp uint64, value float64, score float64) {
	key := encodeStampKey(name, stamp)
	val := encodeStampValue(value, score)
	b.b.Put([]byte(key), []byte(val))
}

// Put name with trend and stamp.
func (b *Batch) PutName(name string, trend float64, stamp uint64) {
	key := encodeNameKey(name)
	val := encodeNameValue(trend, stamp)
	b.b.Put([]byte(key), []byte(val))
}
