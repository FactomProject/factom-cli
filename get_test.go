// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

var _ = fmt.Sprint("")

func TestGet(t *testing.T) {
	fmt.Printf("Get\n===\n")
	args := []string{"get"}
	err := get(args)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetDBlocks(t *testing.T) {
	fmt.Printf("GetDBlocks\n===\n")
	args := []string{"get", "dblocks", "0", "1"}
	err := get(args)
	if err != nil {
		fmt.Println(err)
	}
}