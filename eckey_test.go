package main

import (
	"testing"
)

func TestNewECKey(t *testing.T) {
	err := newECKey()
	if err != nil {
		t.Error(err)
	}
}