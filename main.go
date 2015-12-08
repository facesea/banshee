// Copyright 2015 Eleme Inc. All rights reserved.

package main

import (
	"flag"
	"log"
)

func main() {
	fileName := flag.String("c", "config.json", "config file")
	flag.Parse()
	_, err := NewConfigWithJsonFile(*fileName)
	if err != nil {
		log.Fatal(err)
	}
}
