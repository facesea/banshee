// Copyright 2015 Eleme Inc. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/eleme/banshee/alerter"
	"github.com/eleme/banshee/cleaner"
	"github.com/eleme/banshee/config"
	"github.com/eleme/banshee/detector"
	"github.com/eleme/banshee/filter"
	"github.com/eleme/banshee/storage"
	"github.com/eleme/banshee/util/log"
	"github.com/eleme/banshee/version"
	"github.com/eleme/banshee/webapp"
)

func main() {
	// Arguments
	fileName := flag.String("c", "config.json", "config file")
	debug := flag.Bool("d", false, "debug mode")
	vers := flag.Bool("v", false, "version")
	flag.Parse()
	// Version
	if *vers {
		fmt.Fprintln(os.Stderr, version.Version)
		os.Exit(1)
	}
	// Logging
	log.SetName("banshee")
	if *debug {
		log.SetLevel(log.DEBUG)
	}
	log.Debug("using %s, max %d cpu", runtime.Version(), runtime.GOMAXPROCS(-1))
	// Config
	cfg := config.New()
	if flag.NFlag() == 0 || (flag.NFlag() == 1 && *debug == true) {
		log.Warn("no config file specified, using default..")
	} else {
		err := cfg.UpdateWithJSONFile(*fileName)
		if err != nil {
			log.Fatal("failed to load %s, %s", *fileName, err)
		}
	}
	// Storage
	options := &storage.Options{
		NumGrid: cfg.Period[0],
		GridLen: cfg.Period[1],
	}
	db, err := storage.Open(cfg.Storage.Path, options)
	if err != nil {
		log.Fatal("failed to open %s: %v", cfg.Storage.Path, err)
	}
	// Cleaner
	cleaner := cleaner.New(db, cfg.Period[0]*cfg.Period[1])
	go cleaner.Start()
	// Filter
	filter := filter.New()
	filter.Init(db)
	// Alerter
	alerter := alerter.New(cfg, db, filter)
	alerter.Start()
	// Webapp
	go webapp.Start(cfg, db)
	// Detector
	detector := detector.New(cfg, db, filter)
	detector.Out(alerter.In)
	detector.Start()
}
