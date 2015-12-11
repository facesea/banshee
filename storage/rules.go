// Copyright 2015 Eleme Inc. All rights reserved.

package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Rules file name
const rulesFileName = "rules"

// A rulesDB handles the leveldb for rules storage.
type rulesDB struct {
	// LevelDB
	db *leveldb.DB
}

// Open a rulesDB for rules storage. If the path dosen't exist, it will be
// created.
func openRulesDB(fileName string) (*rulesDB, error) {
	db, err := leveldb.OpenFile(fileName, nil)
	if err != nil {
		return nil, err
	}
	return &rulesDB{db: db}, nil
}

// Close a rulesDB.
func (db *rulesDB) close() error {
	return db.db.Close()
}
