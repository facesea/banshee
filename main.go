// Copyright 2015 Eleme Inc. All rights reserved.

package main

import (
	"flag"
	"log"
	"github.com/eleme/banshee/config
)

func main() {
	fileName := flag.String("c", "config.json", "config file")
	flag.Parse()
	_, err := config.NewConfigWithJsonFile(*fileName)
	if err != nil {
		log.Fatal(err)
	}
}
