package config

import (
	"github.com/eleme/banshee/util"
	"reflect"
	"testing"
)

func TestExampleConfigParsing(t *testing.T) {
	config, err := NewWithJsonFile("./exampleConfig.json")
	util.Assert(t, err == nil)
	defaults := NewWithDefaults()
	util.Assert(t, reflect.DeepEqual(config, defaults))
}
