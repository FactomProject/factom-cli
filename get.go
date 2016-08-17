// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
)

var get = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get allentries|alltxs|chainhead|dblock|eblock|entry|firstentry|head|height"
	cmd.description = "Get data about Factom Chains, Entries, and Blocks"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		c := cli.New()
		c.Handle("allentries", getAllEntries)
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
	help.Add("get", cmd)
	return cmd
}()

var getAllEntries = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get allentries [-n NAME1 -N HEXNAME2 ...] CHAINID"
	cmd.description = "Get all of the Entries in a Chain"
	cmd.execFunc = func(args []string) {
		var (
			nAcii namesAscii
			nHex  namesHex
		)
		os.Args = args
		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name component")
		flag.Var(&nHex, "N", "hex binary name component")
		flag.Parse()
		args = flag.Args()

		var chainid string

		if len(args) < 1 && len(nameCollector) == 0 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if len(nameCollector) != 0 {
			chainid = nametoid(nameCollector)
		} else {
			chainid = args[0]
		}

		es, err := factom.GetAllChainEntries(chainid)
		if err != nil {
			for i, e := range es {
				fmt.Printf("Entry [%d] {\n%s}\n", i, e)
			}
			errorln(err)
			return
		}

		for i, e := range es {
			fmt.Printf("Entry [%d] {\n%s}\n", i, e)
		}
	}
	help.Add("get allentries", cmd)
	return cmd
}()

var getChainHead = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get chainhead [-n NAME1 -N HEXNAME2 ...] CHAINID"
	cmd.description = "Get ebhead by chainid"
	cmd.execFunc = func(args []string) {
		var (
			nAcii namesAscii
			nHex  namesHex
		)
		os.Args = args
		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name component")
		flag.Var(&nHex, "N", "hex binary name component")
		flag.Parse()
		args = flag.Args()

		var chainid string

		if len(args) < 1 && len(nameCollector) == 0 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if len(nameCollector) != 0 {
			chainid = nametoid(nameCollector)
		} else {
			chainid = args[0]
		}

		head, err := factom.GetChainHead(chainid)
		if err != nil {
			errorln(err)
			return
		}
		eblock, err := factom.GetEBlock(head)
		if err != nil {
			errorln(err)
			return
		}

		fmt.Println("EBlock:", head)
		fmt.Println(eblock)
	}
	help.Add("get chainhead", cmd)
	return cmd
}()

var getDBlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get dblock KEYMR"
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

var getEBlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get eblock KEYMR"
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
	cmd.helpMsg = "factom-cli get entry HASH"
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
	cmd.helpMsg = "factom-cli get firstentry [-n NAME1 -N HEXNAME2 ...] CHAINID"
	cmd.description = "Get the first entry from a chain"
	cmd.execFunc = func(args []string) {
		var (
			nAcii namesAscii
			nHex  namesHex
		)
		os.Args = args
		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name component")
		flag.Var(&nHex, "N", "hex binary name component")
		flag.Parse()
		args = flag.Args()

		var chainid string

		if len(args) < 1 && len(nameCollector) == 0 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if len(nameCollector) != 0 {
			chainid = nametoid(nameCollector)
		} else {
			chainid = args[0]
		}

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

var getHead = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get head"
	cmd.description = "Get the latest completed directory block"
	cmd.execFunc = func(args []string) {
		head, err := factom.GetDBlockHead()
		if err != nil {
			errorln(err)
			return
		}
		dblock, err := factom.GetDBlock(head)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println("DBlock:", head)
		fmt.Println(dblock)
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

var properties = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli properties"
	cmd.description = "Get version information about facotmd and the factom wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		f, a, w, err := factom.GetProperties()
		if err != nil {
			errorln(err)
			return
		}

		fmt.Println("Factomd Version:", f)
		fmt.Println("API Version:", a)
		fmt.Println("Wallet Version:", w)
		fmt.Println("CLI Version:", Version)
	}
	help.Add("properties", cmd)
	return cmd
}()
