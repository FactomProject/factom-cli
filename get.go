// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	
	"github.com/FactomProject/factom"
)

func get(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("get")
	}

	switch args[0] {
	case "head":
		return getHead()
	case "dblock":
		return getDBlock(args)
	case "chain":
		return getChain(args)
	case "eblock":
		return getEBlock(args)
	case "entry":
		return getEntry(args)
	default:
		return man("get")
	}

	panic("Something went really wrong with get!")
}

func getHead() error {
	head, err := factom.GetDBlockHead()
	if err != nil {
		return err
	}
	fmt.Println(head)
	return nil
}

func getDBlock(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("getDBlock")
	}

	keymr := args[0]
	dblock, err := factom.GetDBlock(keymr)
	if err != nil {
		return err
	}

	fmt.Println(dblock)
	return nil
}

func getChain(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("getChain")
	}
	
	chainid := args[0]
	chain, err := factom.GetChainHead(chainid)
	if err != nil {
		return err
	}
	
	fmt.Println(chain)
	return nil
}

func getEBlock(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("getEBlock")
	}

	keymr := args[0]
	eblock, err := factom.GetEBlock(keymr)
	if err != nil {
		return err
	}

	fmt.Println(eblock)
	return nil
}

func getEntry(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("getEntry")
	}

	hash := args[0]
	entry, err := factom.GetEntry(hash)
	if err != nil {
		return err
	}
	
	fmt.Println(entry)
	return nil
}
