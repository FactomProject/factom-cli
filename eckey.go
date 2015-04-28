// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
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
	case "pub":
		return printPubKey()
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

	// private key is [32]byte private section + [32]byte public key
	_, priv, err := ed.GenerateKey(rand)
	if err != nil {
		return err
	}
	fmt.Print(hex.EncodeToString(priv[:]))

	return nil
}

func printPubKey() error {
	pub, err := ecPubKey()
	if err != nil {
		return err
	}
	fmt.Println(hex.EncodeToString(pub[:]))

	return nil
}

func ecPrivKey() (*[64]byte, error) {
	priv := new([64]byte)
	
	in, err := os.Open(wallet)
	if err != nil {
		return nil, err
	}
	p, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	key, err := hex.DecodeString(string(p))
	if err != nil {
		return nil, err
	}
	copy(priv[:], key)
	
	return priv, nil
}

func ecPubKey() (*[32]byte, error) {
	pub := new([32]byte)
	
	in, err := os.Open(wallet)
	if err != nil {
		return nil, err
	}
	p, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	key, err := hex.DecodeString(string(p))
	if err != nil {
		return nil, err
	}
	copy(pub[:], key[32:64])
	
	return pub, nil
}
