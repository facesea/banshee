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

var (
	// Arguments
	debug       = flag.Bool("d", false, "debug mode")
	fileName    = flag.String("c", "config.json", "config file path")
	showVersion = flag.Bool("v", false, "show version")
	// Variables
	cfg             = config.New()
	db  *storage.DB = nil
	flt             = filter.New()
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: banshee [-c config] [-d] [-v]\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "copyright eleme https://github.com/eleme/banshee.\n")
	os.Exit(2)
}

func initLog() {
	log.SetName("banshee")
	if *debug {
		log.SetLevel(log.DEBUG)
	}
	goVs := runtime.Version()
	nCpu := runtime.GOMAXPROCS(-1)
	vers := version.Version
	log.Debug("banshee%s %s %d cpu", vers, goVs, nCpu)
}

func initConfig() {
	// Config parsing.
	if flag.NFlag() == 0 || (flag.NFlag() == 1 && *debug) {
		// Case ./program [-d]
		log.Warn("no config specified, using default..")
	} else {
		// Update config.
		err := cfg.UpdateWithJSONFile(*fileName)
		if err != nil {
			log.Fatal("failed to load %s, %s", *fileName, err)
		}
	}
	// Config validation.
	err := cfg.Validate()
	if err == config.ErrAlerterCommandEmpty {
		// Ignore alerter command empty.
		log.Warn("config: %s", err)
	} else {
		log.Fatal("config: %s", err)
	}
}

func initDB() {
	options := &storage.Options{
		NumGrid: cfg.Period[0],
		GridLen: cfg.Period[1],
	}
	path := cfg.Storage.Path
	var err error
	db, err = storage.Open(path, options)
	if err != nil {
		log.Fatal("failed to open %s: %v", path, err)
	}
}

func initFilter() {
	flt.Init(db, cfg)
}

func main() {
	// Arguments
	flag.Usage = usage
	flag.Parse()
	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Version)
		os.Exit(1)
	}
	// Init
	initLog()
	initConfig()
	initDB()
	initFilter()

	// Service
	cleaner := cleaner.New(db, cfg.Period[0]*cfg.Period[1])
	go cleaner.Start()

	alerter := alerter.New(cfg, db, flt)
	alerter.Start()

	go webapp.Start(cfg, db)

	detector := detector.New(cfg, db, flt)
	detector.Out(alerter.In)
	detector.Start()
}
