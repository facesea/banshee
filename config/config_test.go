package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExampleConfigParsing(t *testing.T) {
	config, err := NewConfigWithJsonFile("./exampleConfig.json")
	assert.Nil(t, err)
	defaultC := NewConfigWithDefaults()
	// example config should be default
	assert.Equal(t, config.Debug, defaultC.Debug)
	assert.Equal(t, config.Interval, defaultC.Interval)
	assert.Equal(t, config.Periodicity, defaultC.Periodicity)
	assert.Equal(t, config.Storage.Path, defaultC.Storage.Path)
	assert.Equal(t, config.Detector.Port, defaultC.Detector.Port)
	assert.Equal(t, config.Detector.TrendFactor, defaultC.Detector.TrendFactor)
	assert.Equal(t, config.Detector.Strict, defaultC.Detector.Strict)
	assert.Equal(t, config.Detector.BlackList, defaultC.Detector.BlackList)
	assert.Equal(t, config.Detector.StartSize, defaultC.Detector.StartSize)
	assert.Equal(t, config.Webapp.Port, defaultC.Webapp.Port)
	assert.Equal(t, config.Webapp.Auth, defaultC.Webapp.Auth)
	assert.Equal(t, config.Alerter.Command, defaultC.Alerter.Command)
}
