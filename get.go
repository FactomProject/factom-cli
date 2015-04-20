// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	
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
	case "dblock":
		return getDBlock(args)
	case "dblocks":
		return getDBlocks(args)
	case "eblock":
		return getEBlock(args)
	case "entry":
		return getEntry(args)
	case "height":
		return getHeight()
	default:
		return man("get")
	}

	panic("something went really wrong with get!")
}

func getDBlock(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("getDBlock")
	}

	hash := args[0]
	dblock, err := factom.GetDBlock(hash)
	if err != nil {
		return err
	}

	fmt.Println(dblock)
	return nil
}

func getDBlocks(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 2 {
		return man("getDBlocks")
	}

	from, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	to, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}
	
	blocks, err := factom.GetDBlocks(from, to)
	if err != nil {
		return err
	}
	
	fmt.Println(blocks)
	return nil
}

func getEBlock(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("getEBlocks")
	}

	mr := args[0]
	eblock, err := factom.GetEBlock(mr)
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

func getHeight() error {
	n, err := factom.GetBlockHeight()
	if err != nil {
		return err
	}

	fmt.Println(n)
	return nil
}
