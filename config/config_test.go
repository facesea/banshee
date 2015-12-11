package config

import (
	"github.com/eleme/banshee/util"
	"reflect"
	"testing"
)

func TestExampleConfigParsing(t *testing.T) {
	config, err := NewConfigWithJsonFile("./exampleConfig.json")
	util.Assert(t, err == nil)
	defaultC := NewConfigWithDefaults()
	util.Assert(t, reflect.DeepEqual(config, defaultC))
}
