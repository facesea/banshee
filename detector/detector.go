// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"bufio"
	"fmt"
	"net"
	"path/filepath"
	"time"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/filter"
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/storage/indexdb"
	"github.com/eleme/banshee/util/log"
)

// Detection timed out in ms.
const detectionTimedOut = 100

// Detector is a tcp server to detect anomalies.
type Detector struct {
	// Config
	cfg *config.Config
	// Storage
	db *storage.DB
	// Filter
	flt *filter.Filter
	// Output
	outs []chan *models.Metric
}

// New creates a detector.
func New(cfg *config.Config, db *storage.DB, flt *filter.Filter) *Detector {
	d := new(Detector)
	d.cfg = cfg
	d.db = db
	d.flt = flt
	d.outs = make([]chan *models.Metric, 0)
	return d
}

// Out outputs detected metrics to given channel.
func (d *Detector) Out(ch chan *models.Metric) {
	d.outs = append(d.outs, ch)
}

// output detected metrics to outs.
func (d *Detector) output(m *models.Metric) {
	for _, ch := range d.outs {
		select {
		case ch <- m:
		default:
			log.Error("output channel is full, skipping..")
		}
	}
}

// Start the detector server.
func (d *Detector) Start() {
	// Listen
	addr := fmt.Sprintf("0.0.0.0:%d", d.cfg.Detector.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("cannot bind %s: %v", addr, err)
	}
	log.Info("detector is listening on %s", addr)
	// Accept
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error("cannot accept conn: %v, skipping..", err)
			continue
		}
		go d.handle(conn)
	}
}

// handle a connection, it will filter the metrics by rules and detect whether
// the metrics are anomalies.
func (d *Detector) handle(conn net.Conn) {
	// Conn
	addr := conn.RemoteAddr()
	defer func() {
		conn.Close()
		log.Info("conn %s disconnected", addr)
	}()
	log.Info("conn %s established", addr)
	// Scan lines.
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			//  Read err.
			log.Error("cannot read conn: %v, closing..", err)
			break
		}
		startAt := time.Now()
		// Parse
		line := scanner.Text()
		m, err := parseMetric(line)
		if err != nil {
			if len(line) > 10 {
				line = line[:10]
			}
			log.Error("cannot parse: %s: %v, skipping..", line, err)
			continue
		}
		// Validate
		if err := validateMetric(m); err != nil {
			log.Error("invalid metric: %s, %v", m.Name, err)
			continue
		}
		ok, rules := d.match(m)
		if ok {
			// TODO
		}
	}
}

// Test whether a metric matches the rules.
func (d *Detector) match(m *models.Metric) (bool, []*models.Rule) {
	// Check rules.
	rules := d.flt.MatchedRules(m)
	if len(rules) == 0 {
		// Hit no rules.
		return false, rules
	}
	// Check blacklist.
	for _, pattern := range d.cfg.Detector.BlackList {
		ok, err := filepath.Match(pattern, m.Name)
		if err != nil {
			log.Error("invalid black pattern: %s, %v", pattern, err)
			continue
		}
		if err == nil && ok {
			// Hit blacklist.
			log.Debug("%s hit black pattern %s", m.Name, pattern)
			return false, rules
		}
	}
	return true, rules
}

// detect incoming metric with history values.
func (d *Detector) detect(m *models.Metric) error {
	// Get index.
	idx, err := d.db.Index.Get(m.Name)
	if err != nil {
		if err == indexdb.ErrNotFound {
			idx = nil
		} else {
			return err
		}
	}
	// TODO
}

// save detected metric and its index.
func (d *Detector) save(m *models.Metric, idx *models.Index) error {

}

// getIndex returns index by metric.
func (d *Detector) getIndex(m *models.Metric) (*models.Index, error) {
	// Index
	idx, err := d.db.Index.Get(m.Name)
	if err != nil {
		if err == indexdb.ErrNotFound {
			// Nil on not found.
			return nil, nil
		}
		return nil, err
	}
	return idx, nil
}

// getValues returns history values by metric.
func (d *Detector) getValues(m *models.Metric, idx *models.Index) (values []float64, err error) {
	span := uint32(d.cfg.Detector.FilterOffset * float64(d.cfg.Period))
	fz := d.needFillBlankZeros(m, idx)
	for stamp := m.Stamp; stamp+d.cfg.Expiration > m.Stamp; stamp -= d.cfg.Period {
		start := stamp - span
		stop := stamp + span
		l, err := d.db.Metric.Get(m.Name, start, stop)
		if err != nil {
			return
		}
		if !fz {
			// No need to fill
			for i := 0; i < len(l); i++ {
				values = append(values, l[i].Value)
			}
		} else {
			values = append(values, d.fillBlankZeros(l, start, stop)...)
		}
	}
	// Push current value.
	values = append(values, m.Value)
	return
}

// test whether need to fill blanks with zeros.
func (d *Detector) needFillBlankZeros(m *models.Metric, idx *models.Index) bool {
	for _, pattern := range d.cfg.Detector.FillBlankZeros {
		ok, err := filepath.Match(pattern, m.Name)
		if err != nil {
			log.Error("invalid fillBlankZeros pattern: %s, %v", pattern, err)
			continue
		}
		if ok {
			if idx != nil {
				// Only fill blanks with zeros if the metric isn't a new one.
				return true
			}
			return false
		}
	}
	return false
}

// fill blanks with zeros.
func (d *Detector) fillBlankZeros(l []*models.Metric, start, stop uint32) (values []float64) {
	i := 0
	step := d.cfg.Interval
	for start < stop {
		if i < len(l) {
			m := l[i]
			// Fill zeros if the start is too small.
			for ; start < m.Stamp; start += step {
				values = append(values, 0)
			}
			values = append(values, m.Value)
			i += 1
		} else {
			// End of history real-metrics.
			values = append(values, 0)
		}
		start += step
	}
	return
}

// div3sigma returns the score by values.
func (d *Detector) div3sigma(values []float64) (score float64, mean float64) {
	if len(values) == 0 {
		// Empty
		return 0, 0
	}
	// Mean
	mean = getMean(values)
	// Std
	std := getStd(values, mean)
	if len(values) < int(d.cfg.Detector.LeastCount) {
		return
	}
	// Latest value
	tail := values[len(values)-1]
	if std == 0 {
		if tail != mean {
			score = 1
		}
		return
	}
	score = (tail - mean) / (3 * std)
	return
}
