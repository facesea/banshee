// Copyright 2015 Eleme Inc. All rights reserved.

package main

import (
	"flag"
	"os"

	"github.com/eleme/banshee/config"
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
	_, err := config.NewConfigWithJsonFile(*fileName)
	if err != nil {
		logger.Fatal("%s", err)
	}
}
