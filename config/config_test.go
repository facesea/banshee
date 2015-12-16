package config

import (
	"github.com/eleme/banshee/util/assert"
	"reflect"
	"testing"
)

func TestExampleConfigParsing(t *testing.T) {
	config := New()
	err := config.UpdateWithJSONFile("./exampleConfig.json")
	assert.Ok(t, err == nil)
	defaults := New()
	assert.Ok(t, reflect.DeepEqual(config, defaults))
}
