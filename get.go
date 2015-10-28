// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
//	"crypto/sha256"
//	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/cli"
)

var get = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get head|dblock|height|chainhead|eblock|entry"
	cmd.description = "get Block or Entry data from factomd"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		
		c := cli.New()
		c.Handle("head", getHead)
		c.Handle("height", getHeight)
		c.Handle("dblock", getDBlock)
		c.Handle("chainhead", getChainHead)
		c.Handle("eblock", getEBlock)
		c.Handle("entry", getEntry)
		c.Handle("firstentry", getFirstEntry)
		c.HandleDefaultFunc(func(args []string) {
			fmt.Println(cmd.helpMsg)
		})
		c.Execute(args)
	}
	return cmd
}()

var getHead = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get head"
	cmd.description = "Get the keymr of the last completed directory block"
	cmd.execFunc = func(args []string) {
		head, err := factom.GetDBlockHead()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(head.KeyMR)
	}
	return cmd
}()

var getHeight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get height"
	cmd.description = "Get the current directory block height"
	cmd.execFunc = func(args []string) {	
		height, err := factom.GetDBlockHeight()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(height)
	}
	return cmd
}()

var getDBlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get dblock [keymr]"
	cmd.description = "Get dblock contents by merkle root"
	cmd.execFunc = func(args []string) {	
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
	
		keymr := args[0]
		dblock, err := factom.GetDBlock(keymr)
		if err != nil {
			errorln(err)
			return
		}
	
		fmt.Println("PrevBlockKeyMR:", dblock.Header.PrevBlockKeyMR)
		fmt.Println("Timestamp:", dblock.Header.Timestamp)
		fmt.Println("SequenceNumber:", dblock.Header.SequenceNumber)
	
		for _, v := range dblock.EntryBlockList {
			fmt.Println("EntryBlock {")
			fmt.Println("	ChainID", v.ChainID)
			fmt.Println("	KeyMR", v.KeyMR)
			fmt.Println("}")
		}
	}
	return cmd
}()

var getChainHead = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get chainhead [chainid]"
	cmd.description = "Get ebhead by chainid"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
	
		chainid := args[0]
		chain, err := factom.GetChainHead(chainid)
		if err != nil {
			errorln(err)
			return
		}
	
		fmt.Println(chain.ChainHead)
	}
	return cmd
}()

var getEBlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get eblock [keymr]"
	cmd.description = "Get eblock by merkle root"
	cmd.execFunc = func(args []string) {	
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
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
		fmt.Println("Timestamp:", eblock.Header.Timestamp)
	
		for _, v := range eblock.EntryList {
			fmt.Println("EBEntry {")
			fmt.Println("	Timestamp", v.Timestamp)
			fmt.Println("	EntryHash", v.EntryHash)
			fmt.Println("}")
		}
	}
	return cmd
}()

var getEntry = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get entry [hash]"
	cmd.description = "Get entry by hash"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
	
		hash := args[0]
		entry, err := factom.GetEntry(hash)
		if err != nil {
			errorln(err)
			return
		}
	
		printEntry(entry)
	}
	return cmd
}()

var getFirstEntry = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get firstentry [chainid]"
	cmd.description = "Get the first entry from a chain"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
	
		chainid := args[0]
		entry, err := factom.GetFirstEntry(chainid)
		if err != nil {
			errorln(err)
			return
		}
		printEntry(entry)
	}
	return cmd
}()

func printEntry(e *factom.Entry) {
	fmt.Println("ChainID:", e.ChainID)
	for _, id := range e.ExtIDs {
		fmt.Println("ExtID:", string(id))
	}
	
	fmt.Println("Content:")
	fmt.Println(string(e.Content))
}

// TODO - replace getChainId with something
// We expect each element to be its own part in a chain ID
//func getChainId(args []string) {
//	if len(args) < 2 {
//		fmt.Printf("No Chain Specification provided.  See help")
//	}
//	sum := sha256.New()
//	fmt.Println("The chain components:")
//	for i, str := range args {
//		if i > 0 {
//			fmt.Println("    ", str)
//			x := sha256.Sum256([]byte(str))
//			sum.Write(x[:])
//		}
//	}
//	chainId := sum.Sum(nil)
//	fmt.Println("produce the ChainID:")
//
//	fmt.Println("    ", hex.EncodeToString(chainId))
//}
