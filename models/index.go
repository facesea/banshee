// Copyright 2015 Eleme Inc. All rights reserved.

package models

// Index is for metric indexing, it records metric name for faster metric
// indexing. And metric latest informations like timstamp, score, average are
// also indexed.
type Index struct {
	// Index may be cached.
	cache
	// Metric name
	Name string
	// Latest stamp for the metric.
	Stamp uint32
	// Latest score for the metric.
	Score float64
	// Latest average for the metric.
	Average float64
}

// CopyIfShared returns a copy if the index is shared.
func (idx *Index) CopyIfShared() *Index {
	if idx.IsShared() {
		return idx.Copy()
	}
	return idx
}

// Copy the index.
func (idx *Index) Copy() *Index {
	i := &Index{}
	idx.CopyTo(i)
	return i
}

// CopyTo copy the index to another.
func (idx *Index) CopyTo(i *Index) {
	idx.RLockIfShared()
	defer idx.RUnlockIfShared()
	i.LockIfShared()
	defer i.UnlockIfShared()
	i.Name = idx.Name
	i.Stamp = idx.Stamp
	i.Score = idx.Score
	i.Average = idx.Average
}

// Equal tests the equality.
func (idx *Index) Equal(i *Index) bool {
	idx.RLockIfShared()
	defer idx.RUnlockIfShared()
	i.RLockIfShared()
	defer i.RUnlockIfShared()
	if i.Name != idx.Name {
		return false
	}
	if i.Stamp != idx.Stamp {
		return false
	}
	if i.Score != idx.Score {
		return false
	}
	if i.Average != idx.Average {
		return false
	}
	return true
}
