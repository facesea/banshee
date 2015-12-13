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
	debug := flag.Bool("d", false, "debug mode")
	flag.Parse()
	// Logging
	logger := util.NewLogger("banshee")
	if *debug {
		logger.SetLevel(util.LOG_DEBUG)
	}
	logger.Runtime()
	// Config
	cfg := config.NewWithDefaults()
	if flag.NFlag() == 1 && *debug == false {
		err := cfg.UpdateWithJsonFile(*fileName)
		if err != nil {
			logger.Fatal("failed to open %s: %s", *fileName, err)
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
	detector := detector.New(*debug, cfg, db)
	detector.Start()
}
