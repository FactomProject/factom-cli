package main

import (
	"fmt"
)

//func bintx(args []string) error {
//	return fmt.Errorf("bintx called")
//}

func get(args []string) error {
	return fmt.Errorf("get called")
}

func help(arg string) error {
	return fmt.Errorf("help called with: %s", arg)
}

func mkchain(args []string) error {
	return fmt.Errorf("mkchain called")
}
