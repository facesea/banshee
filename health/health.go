// Copyright 2016 Eleme Inc. All rights reserved.

package health

import (
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util/mathutil"
	"sync"
	"sync/atomic"
	"time"
)

// AggregationInterval in seconds, default: 5min
const AggregationInterval int = 5 * 60

// max length for detectionCosts.
const maxDetectionCostsLen = 100 * 1024

// max length for filterCosets.
const maxFilterCostsLen = 100 * 1024

// Info is the stats container.
type Info struct {
	lock sync.RWMutex
	// Total
	AggregationInterval int   `json:"aggregationInterval"`
	NumIndexTotal       int   `json:"numIndexTotal"`
	NumClients          int64 `json:"numClients"`
	NumRules            int   `json:"numRules"`
	// Aggregation
	DetectionCost     float64 `json:"detectionCost"` // ms
	FilterCost        float64 `json:"filterCost"`    // ms
	NumMetricIncomed  int64   `json:"numMetricIncomed"`
	NumMetricDetected int64   `json:"numMetricDetected"`
	NumAlertingEvents int64   `json:"numAlertingEvents"`
}

// Copy info.
func (info *Info) copy() *Info {
	info.lock.RLock()
	defer info.lock.RUnlock()
	return &Info{
		AggregationInterval: AggregationInterval,
		NumIndexTotal:       info.NumIndexTotal,
		NumClients:          info.NumClients,
		NumRules:            info.NumRules,
		DetectionCost:       info.DetectionCost,
		FilterCost:          info.FilterCost,
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
	filterCosts        []float64
	filterCostsLock    sync.Mutex
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

// AddFilterCost appends cost to FilterCosts.
func AddFilterCost(n float64) {
	h.filterCostsLock.Lock()
	defer h.filterCostsLock.Unlock()
	if len(h.filterCosts) < maxFilterCostsLen {
		h.filterCosts = append(h.filterCosts, n)
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

// Refresh NumRules.
func refreshNumRules() {
	h.info.lock.Lock()
	defer h.info.lock.Unlock()
	h.info.NumRules = h.db.Admin.RulesCache.Len()
}

// Aggregate DetectionCost.
func aggregateDetectionCost() {
	h.info.lock.Lock()
	defer h.info.lock.Unlock()
	h.detectionCostsLock.Lock()
	defer h.detectionCostsLock.Unlock()
	h.info.DetectionCost = mathutil.Average(h.detectionCosts)
	h.detectionCosts = h.detectionCosts[:0]
}

// Aggregate FilterCost.
func aggregationFilterCost() {
	h.info.lock.Lock()
	defer h.info.lock.Unlock()
	h.filterCostsLock.Lock()
	defer h.filterCostsLock.Unlock()
	h.info.FilterCost = mathutil.Average(h.filterCosts)
	h.filterCosts = h.filterCosts[:0]
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
		refreshNumRules()
		aggregateDetectionCost()
		aggregateNumMetricIncomed()
		aggregateNumMetricDetected()
		aggregateNumAlertingEvents()
		aggregationFilterCost()
	}
}
