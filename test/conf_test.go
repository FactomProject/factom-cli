package main

import (
	"fmt"
	"testing"
)

func TestReadConfig(t *testing.T) {
	fmt.Printf("TestReadConfig\n===\n")
	cfg := ReadConfig()
	fmt.Printf("%#v\n", cfg)
}