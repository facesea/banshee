// Copyright 2015 Eleme Inc. All rights reserved.

package main

import (
	"flag"
	"os"

	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/detector"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util"
)

func main() {
	// Argv parsing
	fileName := flag.String("c", "config.json", "config file")
	flag.Parse()
	if flag.NFlag() != 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Config
	logger := util.NewLogger("banshee")
	logger.Runtime(nil)
	cfg, err := config.NewConfigWithJsonFile(*fileName)
	if err != nil {
		logger.Fatal("%s", err)
	}
	// Storage
	db, err := storage.Open(cfg)
	if err != nil {
		logger.Fatal("failed to open %s: %v", cfg.Storage.Path, err)
	}
	// Detector
	detector := detector.New(cfg, db)
	detector.Start()
}
