// Copyright 2015 Eleme Inc. All rights reserved.

package indexdb

import (
	"fmt"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util"
)

// encode encodes db value from index.
func encode(idx *models.Index) []byte {
	// Value format is Stamp:Score:Average
	score := util.ToFixed(idx.Score, 5)
	average := util.ToFixed(idx.Average, 5)
	s := fmt.Sprintf("%d:%s:%s", idx.Stamp, score, average)
	return []byte(s)
}

// decode decodes db value into index.
func decode(value []byte, idx *models.Index) error {
	n, err := fmt.Sscanf(string(value), "%d:%f:%f", &idx.Stamp, &idx.Score, &idx.Average)
	if err != nil {
		return err
	}
	if n != 3 {
		return ErrCorrupted
	}
	return nil
}
