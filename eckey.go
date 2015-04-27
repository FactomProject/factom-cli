// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	ed "github.com/agl/ed25519"
)

var (
	_ = hex.EncodeToString
	_ = fmt.Sprint("testing")
)

func eckey(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()

	switch args[0] {
	case "new":
		return newECKey()
	default:
		return nil
	}
	
	panic("Something went really wrong with eckey!")
}

func newECKey() error {
	rand, err := os.Open("/dev/random")
	if err != nil {
		return err
	}

	_, priv, err := ed.GenerateKey(rand)
	if err != nil {
		return err
	}
	fmt.Println(hex.EncodeToString(priv[:]))

	return nil
}
