package config

import (
	"github.com/eleme/banshee/util"
	"reflect"
	"testing"
)

func TestExampleConfigParsing(t *testing.T) {
	config := New()
	err := config.UpdateWithJsonFile("./exampleConfig.json")
	util.Assert(t, err == nil)
	defaults := New()
	util.Assert(t, reflect.DeepEqual(config, defaults))
}
