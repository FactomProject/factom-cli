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

	"github.com/FactomProject/factom"
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
		return man("eckey")
	}
	
	panic("Something went really wrong with eckey!")
}

func newECKey() error {
	key := factom.NewECKey()
	fmt.Printf("%x", *key)
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
