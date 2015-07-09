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

func get(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		man("get")
		return
	}

	switch args[0] {
	case "head":
		getHead()
	case "dblock":
		getDBlock(args)
	case "chain":
		getChain(args)
	case "eblock":
		getEBlock(args)
	case "entry":
		getEntry(args)
	default:
		man("get")
	}
}

func getHead() {
	head, err := factom.GetDBlockHead()
	if err != nil {
		errorln(err)
		return
	}
	fmt.Println(head.KeyMR)
}

func getDBlock(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		man("getDBlock")
		return
	}

	keymr := args[0]
	dblock, err := factom.GetDBlock(keymr)
	if err != nil {
		errorln(err)
		return
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
}

func getChain(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		man("getChain")
		return
	}
	
	chainid := args[0]
	chain, err := factom.GetChainHead(chainid)
	if err != nil {
		errorln(err)
		return
	}
	
	fmt.Println(chain.EntryBlockKeyMR)
}

func getEBlock(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		man("getEBlock")
		return
	}

	keymr := args[0]
	eblock, err := factom.GetEBlock(keymr)
	if err != nil {
		errorln(err)
		return
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
}

func getEntry(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		man("getEntry")
		return
	}

	hash := args[0]
	entry, err := factom.GetEntry(hash)
	if err != nil {
		errorln(err)
		return
	}
	
	fmt.Println("ChainID:", entry.ChainID)
	for _, v := range entry.ExtIDs {
		fmt.Println("ExtID:", v)
	}
	fmt.Println("Content:", entry.Content)
}
