// Copyright 2015 Eleme Inc. All rights reserved.

package main

import (
	"flag"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/detector"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util"
)

func main() {
	// Argv parsing
	fileName := flag.String("c", "config.json", "config file")
	flag.Parse()
	// Logging
	logger := util.NewLogger("banshee")
	logger.Runtime()
	// Config
	cfg := config.NewWithDefaults()
	if flag.NFlag() != 0 {
		err := cfg.UpdateFromJsonFile(*fileName)
		if err != nil {
			logger.Fatal("%s", err)
		}
	} else {
		logger.Warn("no config file specified, using default..")
	}
	// Storage
	numGrids, gridLen := cfg.Periodicity[0], cfg.Periodicity[1]
	db, err := storage.Open(cfg.Storage.Path, numGrids, gridLen)
	if err != nil {
		logger.Fatal("failed to open %s: %v", cfg.Storage.Path, err)
	}
	// Detector
	detector := detector.New(cfg, db)
	detector.Start()
}
