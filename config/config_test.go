package config

import (
"testing"
)

var expected = Config{
	Version:     "UNSET",
	Revision:    "UNSET",
	FlagVerbose: false,
	FlagJSONLog: false,
	FlagAPIPort: 9093,
}

func TestConfigDefaults(t *testing.T) {
	testConfig := Config{}
	testConfig.setDefaults()
	if testConfig != expected {
		t.Error("Expected", expected, "got", testConfig)
	}
}

func TestConfiguration(t *testing.T) {
	c := Configuration()
	if c != expected {
		t.Error("Expected", expected, "got", c)
	}
}