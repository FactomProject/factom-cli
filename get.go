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
	fmt.Println(head.KeyMR)
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

	fmt.Println("PrevBlockKeyMR:", dblock.Header.PrevBlockKeyMR)
	fmt.Println("TimeStamp:", dblock.Header.TimeStamp)
	fmt.Println("SequenceNumber:", dblock.Header.SequenceNumber)
	
	for _, v := range dblock.EntryBlockList {
		fmt.Println("EntryBlock {")
		fmt.Println("	ChainID", v.ChainID)
		fmt.Println("	KeyMR", v.KeyMR)
		fmt.Println("}")
	}
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
	
	fmt.Println(chain.EntryBlockKeyMR)
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

	fmt.Println("BlockSequenceNumber:", eblock.Header.BlockSequenceNumber)
	fmt.Println("ChainID:", eblock.Header.ChainID)
	fmt.Println("PrevKeyMR:", eblock.Header.PrevKeyMR)
	fmt.Println("TimeStamp:", eblock.Header.TimeStamp)
	
	for _, v := range eblock.EntryList {
		fmt.Println("EBEntry {")
		fmt.Println("	TimeStamp", v.TimeStamp)
		fmt.Println("	EntryHash", v.EntryHash)
		fmt.Println("}")
	}
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
	
	fmt.Println("ChainID:", entry.ChainID)
	for _, v := range entry.ExtIDs {
		fmt.Println("ExtID:", v)
	}
	fmt.Println("Content:", entry.Content)
	return nil
}
