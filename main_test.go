package main

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	err := readConfig()
	if err != nil {
		t.Errorf("Cannot read config file. Error: %v", err)
	}
}
