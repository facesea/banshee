// Copyright 2015 Eleme Inc. All rights reserved.

package models

// Index is for metric indexing, it records metric name for faster metric
// indexing. And metric latest informations like timstamp, score, average are
// also indexed.
type Index struct {
	// Index may be cached.
	cache `json:"-"`
	// Metric name
	Name string `json:"name"`
	// Latest stamp for the metric.
	Stamp uint32 `json:"stamp"`
	// Latest score for the metric.
	Score float64 `json:"score"`
	// Latest average for the metric.
	Average float64 `json:"average"`
}

// WriteMetric writes metric to index.
func (idx *Index) WriteMetric(m *Metric) {
	idx.Lock()
	defer idx.Unlock()
	idx.Name = m.Name
	idx.Stamp = m.Stamp
	idx.Score = m.Score
	idx.Average = m.Average
}

// Copy the index.
func (idx *Index) Copy() *Index {
	i := &Index{}
	idx.CopyTo(i)
	return i
}

// CopyTo copy the index to another.
func (idx *Index) CopyTo(i *Index) {
	idx.RLock()
	defer idx.RUnlock()
	i.Lock()
	defer i.Unlock()
	i.Name = idx.Name
	i.Stamp = idx.Stamp
	i.Score = idx.Score
	i.Average = idx.Average
}

// Equal tests the equality.
func (idx *Index) Equal(i *Index) bool {
	idx.RLock()
	defer idx.RUnlock()
	i.RLock()
	defer i.RUnlock()
	return (idx.Name == i.Name &&
		idx.Stamp == i.Stamp &&
		idx.Score == i.Score &&
		idx.Average == i.Average)
}
