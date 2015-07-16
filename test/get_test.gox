// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"testing"
)

var _ = fmt.Sprint("")

var (
	cfg = ReadConfig().Main
)

func TestGet(t *testing.T) {
	server = cfg.Server
	wallet = cfg.Wallet
	subs := map[string][]string{
		"get":     {"get"},
		"height":  {"get", "height"},
		"dblocks": {"get", "dblocks", "0", "1"},
		"eblock":  {"get", "eblock", "merkle"},
		"entry":   {"get", "entry", "abcd"},
	}
	for i, v := range subs {
		fmt.Printf("%s\n===\n", i)
		err := get(v)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println()
	}
}
