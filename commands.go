package main

import (
	"fmt"
)

func buy(args []string) error {
	return fmt.Errorf("buy called")
}

func get(args []string) error {
	return fmt.Errorf("get called")
}

func help(args []string) error {
	return fmt.Errorf("help called with: %v", args)
}

func mkchain(args []string) error {
	return fmt.Errorf("mkchain called")
}
