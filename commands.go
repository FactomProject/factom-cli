package main

import (
	"fmt"
)

func get() {
	fmt.Println("get called")
}

func help(s string) {
	fmt.Println("help called with:", s)
}

func mkchain() {
	fmt.Println("mkchain called")
}
