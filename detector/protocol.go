// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"github.com/eleme/banshee/models"
	"strconv"
	"strings"
)

// Parse input line text into a metric.
//
// The detector's net protocol is, with an example:
//	Name	Stamp		Value	\n
//	foo		1449481993	3.145	\n
//
// Input line will be trimed at first before being processed.
func parseMetric(line string) (*models.Metric, error) {
	// Clean spaces.
	line = strings.TrimSpace(line)
	// Split fields.
	words := strings.Fields(line)
	if len(words) != 3 {
		// Wrong number of fields.
		return nil, ErrProtocol
	}
	var err error
	m := &models.Metric{}
	// Name is a string
	m.Name = words[0]
	num, err := strconv.ParseUint(words[1], 10, 32)
	if err != nil {
		return nil, err
	}
	// Stamp is a uint32.
	m.Stamp = uint32(num)
	// Value is a float64.
	m.Value, err = strconv.ParseFloat(words[2], 64)
	if err != nil {
		return nil, err
	}
	return m, nil
}
