// Copyright 2016 Eleme Inc. All rights reserved.

package health

import (
	"github.com/eleme/banshee/storage"
	"sync"
	"sync/atomic"
	"time"
)

// AggregationInterval in seconds, default: 5min
const AggregationInterval int = 5 * 60

// max length for detectionCosts.
const maxDetectionCostsLen = 100 * 1024

// Info is the stats container.
type Info struct {
	lock sync.RWMutex
	// Total
	AggregationInterval int   `json:"aggregationInterval"`
	NumIndexTotal       int   `json:"numIndexTotal"`
	NumClients          int64 `json:"numClients"`
	// Aggregation
	DetectionCost     float64 `json:"detectionCost.5min"` // ms
	NumMetricIncomed  int64   `json:"numMetricIncome.5min"`
	NumMetricDetected int64   `json:"numMetricDetected.5min"`
	NumAlertingEvents int64   `json:"numAlertingEvents.5min"`
}

// Copy info.
func (info *Info) copy() *Info {
	info.lock.RLock()
	defer info.lock.RUnlock()
	return &Info{
		AggregationInterval: AggregationInterval,
		NumIndexTotal:       info.NumIndexTotal,
		NumClients:          info.NumClients,
		DetectionCost:       info.DetectionCost,
		NumMetricIncomed:    info.NumMetricIncomed,
		NumMetricDetected:   info.NumMetricDetected,
		NumAlertingEvents:   info.NumAlertingEvents,
	}
}

// Aggregation hub.
type hub struct {
	// Storage
	db *storage.DB
	// Info
	info Info
	// Total
	numClients int64
	// Aggregation
	detectionCosts     []float64
	detectionCostsLock sync.Mutex
	numMetricIncomed   int64
	numMetricDetected  int64
	numAlertingEvents  int64
}

// Single-ton hub.
var h hub

// Init hub.
func Init(db *storage.DB) {
	h.db = db
}

// Get info.
func Get() *Info {
	return h.info.copy()
}

// IncrNumClients increments NumClients by n.
func IncrNumClients(n int64) {
	atomic.AddInt64(&h.numClients, n)
}

// DecrNumClients decrments NumClients by n.
func DecrNumClients(n int64) {
	IncrNumClients(0 - n)
}

// AddDetectionCost appends cost to DetectionCosts.
func AddDetectionCost(n float64) {
	h.detectionCostsLock.Lock()
	defer h.detectionCostsLock.Unlock()
	if len(h.detectionCosts) < maxDetectionCostsLen {
		h.detectionCosts = append(h.detectionCosts, n)
	}
}

// IncrNumMetricIncomed increments NumMetricIncomed by n.
func IncrNumMetricIncomed(n int64) {
	atomic.AddInt64(&h.numMetricIncomed, n)
}

// IncrNumMetricDetected increments NumMetricDetected by n.
func IncrNumMetricDetected(n int64) {
	atomic.AddInt64(&h.numMetricDetected, n)
}

// IncrNumAlertingEvents increments NumAlertingsEvents by n.
func IncrNumAlertingEvents(n int64) {
	atomic.AddInt64(&h.numAlertingEvents, n)
}

// Refresh NumIndexTotal.
func refreshNumIndexTotal() {
	h.info.lock.Lock()
	defer h.info.lock.Unlock()
	h.info.NumIndexTotal = h.db.Index.Len()
}

// Refresh NumClients.
func refreshNumClients() {
	h.info.lock.Lock()
	defer h.info.lock.Unlock()
	h.info.NumClients = atomic.LoadInt64(&h.numClients)
}

// Aggregate DetectionCost.
func aggregateDetectionCost() {
	h.info.lock.Lock()
	defer h.info.lock.Unlock()
	h.detectionCostsLock.Lock()
	defer h.detectionCostsLock.Unlock()
	var sum float64
	for i := 0; i < len(h.detectionCosts); i++ {
		sum += h.detectionCosts[i]
	}
	h.info.DetectionCost = sum / float64(len(h.detectionCosts))
	h.detectionCosts = h.detectionCosts[:0]
}

// Aggregate NumMetricIncomed.
func aggregateNumMetricIncomed() {
	h.info.lock.Lock()
	defer h.info.lock.Unlock()
	h.info.NumMetricIncomed = atomic.LoadInt64(&h.numMetricIncomed)
	atomic.StoreInt64(&h.numMetricIncomed, 0)
}

// Aggregate NumMetricDetected.
func aggregateNumMetricDetected() {
	h.info.lock.Lock()
	defer h.info.lock.Unlock()
	h.info.NumMetricDetected = atomic.LoadInt64(&h.numMetricDetected)
	atomic.StoreInt64(&h.numMetricDetected, 0)
}

// Aggregate NumAlertingEvents.
func aggregateNumAlertingEvents() {
	h.info.lock.Lock()
	defer h.info.lock.Unlock()
	h.info.NumAlertingEvents = atomic.LoadInt64(&h.numAlertingEvents)
	atomic.StoreInt64(&h.numAlertingEvents, 0)
}

// Start the health aggregator.
func Start() {
	interval := time.Duration(AggregationInterval) * time.Second
	ticker := time.NewTicker(interval)
	for {
		<-ticker.C
		refreshNumIndexTotal()
		refreshNumClients()
		aggregateDetectionCost()
		aggregateNumMetricIncomed()
		aggregateNumMetricDetected()
		aggregateNumAlertingEvents()
	}
}
