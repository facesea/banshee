// Copyright 2015 Eleme Inc. All rights reserved.

package detector

import (
	"strconv"
	"strings"

	"github.com/eleme/banshee/models"
)

// parseMetric parses protocol line string into a Metric, the line protocol is:
//   NAME  TIMESTAMP  VALUE \n
// and the example:
//   foo   1449481993 3.145 \n
//
func parseMetric(line string) (*models.Metric, error) {
	line = strings.TrimSpace(line)
	words := strings.Fields(line)
	if len(words) != 3 {
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
