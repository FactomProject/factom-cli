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

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
)

var get = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get head|dblock|height|chain|eblock|entry|firstentry"
	cmd.description = "get Block or Entry data from factomd"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		c := cli.New()
		c.Handle("head", getHead)
		c.Handle("height", getHeight)
		c.Handle("dblock", getDBlock)
		c.Handle("chain", getChainHead)
		c.Handle("eblock", getEBlock)
		c.Handle("entry", getEntry)
		c.Handle("firstentry", getFirstEntry)
		c.HandleDefaultFunc(func(args []string) {
			fmt.Println(cmd.helpMsg)
		})
		c.Execute(args)
	}
	help.Add("get", cmd)
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
	help.Add("get head", cmd)
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
	help.Add("get height", cmd)
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
		fmt.Println(dblock)
	}
	help.Add("get dblock", cmd)
	return cmd
}()

var getChainHead = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get chain [chainid]"
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
	help.Add("get chain", cmd)
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
		fmt.Println(eblock)
	}
	help.Add("get eblock", cmd)
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
		fmt.Println(entry)
	}
	help.Add("get entry", cmd)
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
		fmt.Println(entry)
	}
	help.Add("get firstentry", cmd)
	return cmd
}()
