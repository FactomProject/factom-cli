// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
)

var get = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get allentries|chainhead|dblock|eblock|entry|" +
		"firstentry|head|heights|walletheight|pendingentries|" +
		"pendingtransactions|raw|dbheight|abheight|fbheight|ecbheight"
	cmd.description = "Get data about Factom Chains, Entries, and Blocks"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		c := cli.New()
		c.Handle("abheight", abheight)
		c.Handle("allentries", getAllEntries)
		c.Handle("chainhead", getChainHead)
		c.Handle("dbheight", dbheight)
		c.Handle("dblock", getDBlock)
		c.Handle("eblock", getEBlock)
		c.Handle("ecbheight", ecbheight)
		c.Handle("entry", getEntry)
		c.Handle("fbheight", fbheight)
		c.Handle("firstentry", getFirstEntry)
		c.Handle("head", getHead)
		c.Handle("heights", getHeights)
		c.Handle("pendingentries", getPendingEntries)
		c.Handle("pendingtransactions", getPendingTransactions)
		c.Handle("raw", getraw)
		c.Handle("walletheight", getWalletHeight)
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
	cmd.helpMsg = "factom-cli get allentries [-n NAME1 -h HEXNAME2" +
		" ...|CHAINID] [-E]"
	cmd.description = "Get all of the Entries confirmed in a Chain. -n and" +
		" -h to specify the chain name. -E EntryHash."
	cmd.execFunc = func(args []string) {
		var (
			nAcii namesASCII
			nHex  namesHex
		)
		os.Args = args
		edisp := flag.Bool(
			"E",
			false,
			"display only the EntryHashes",
		)
		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name component")
		flag.Var(&nHex, "h", "hex binary name component")
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
				switch {
				case *edisp:
					fmt.Printf("%x\n", e.Hash())
				default:
					fmt.Printf("Entry [%d] {\n%s}\n", i, e)
				}
			}
			errorln(err)
			return
		}

		for i, e := range es {
			switch {
			case *edisp:
				fmt.Printf("%x\n", e.Hash())
			default:
				fmt.Printf("Entry [%d] {\n%s}\n", i, e)
			}
		}
	}
	help.Add("get allentries", cmd)
	return cmd
}()

var getChainHead = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get chainhead [-n NAME1 -h HEXNAME2 ...|CHAINID] [-K]"
	cmd.description = "Get the latest Entry Block of the specified Chain. -n" +
		" and -h to specify the chain name. -K KeyMR."
	cmd.execFunc = func(args []string) {
		var (
			nAcii namesASCII
			nHex  namesHex
		)
		os.Args = args
		kdisp := flag.Bool(
			"K",
			false,
			"display only the KeyMR of the entry block",
		)
		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name component")
		flag.Var(&nHex, "h", "hex binary name component")
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

		switch {
		case *kdisp:
			fmt.Println(head)
		default:
			fmt.Println("EBlock:", head)
			fmt.Println(eblock)
		}
	}
	help.Add("get chainhead", cmd)
	return cmd
}()

var getDBlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get dblock KEYMR"
	cmd.description = "Get dblock contents by Key Merkle Root"
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
	cmd.description = "Get Entry Block by Key Merkle Root"
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
	cmd.description = "Get Entry by Hash"
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
	cmd.helpMsg = "factom-cli get firstentry [-n NAME1 -h HEXNAME2 ...|CHAINID] [-E]"
	cmd.description = "Get the first Entry in a Chain. -E EntryHash"
	cmd.execFunc = func(args []string) {
		var (
			nAcii namesASCII
			nHex  namesHex
		)
		os.Args = args
		edisp := flag.Bool(
			"E",
			false,
			"display only the EntryHash of the first entry",
		)
		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name component")
		flag.Var(&nHex, "h", "hex binary name component")
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

		switch {
		case *edisp:
			fmt.Printf("%x\n", entry.Hash())
		default:
			fmt.Println(entry)
		}

	}
	help.Add("get firstentry", cmd)
	return cmd
}()

var getHead = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get head [-K]"
	cmd.description = "Get the latest completed Directory Block. -K KeyMR."
	cmd.execFunc = func(args []string) {
		os.Args = args
		kdisp := flag.Bool(
			"K",
			false,
			"display only the KeyMR of the directory block",
		)
		flag.Parse()
		args = flag.Args()

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

		switch {
		case *kdisp:
			fmt.Println(head)
		default:
			fmt.Println("DBlock:", head)
			fmt.Println(dblock)
		}
	}
	help.Add("get head", cmd)
	return cmd
}()

var getHeights = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get heights"
	cmd.description = "Get the current heights of various items in factomd"
	cmd.execFunc = func(args []string) {
		height, err := factom.GetHeights()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(height.String())
	}
	help.Add("get heights", cmd)
	return cmd
}()

var getWalletHeight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get walletheight"
	cmd.description = "Get the number of factoid blocks factom-walletd has cached"
	cmd.execFunc = func(args []string) {
		height, err := factom.GetWalletHeight()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Printf("WalletHeight: %v\n", height)
	}
	help.Add("get walletheight", cmd)
	return cmd
}()

var properties = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli properties"
	cmd.description = "Get version information about factomd and the factom wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		fdv, fdverr, fdapiv, fdapiverr, fwv, fwverr, fwapiv, fwapiverr := factom.GetProperties()

		fmt.Println("CLI Version:", Version)
		if fdverr == "" {
			fmt.Println("Factomd Version:", fdv)
		} else {
			fmt.Println("Factomd Version Unavailable:", fdverr)
		}

		if fdapiverr == "" {
			fmt.Println("Factomd API Version:", fdapiv)
		} else {
			fmt.Println("Factomd API Version Unavailable:", fdapiverr)
		}

		if fwverr == "" {
			fmt.Println("Wallet Version:", fwv)
		} else {
			fmt.Println("Wallet Version Unavailable:", fwverr)
		}
		if fwapiverr == "" {
			fmt.Println("Wallet API Version:", fwapiv)
		} else {
			fmt.Println("Wallet API Version Unavailable:", fwapiverr)
		}

	}
	help.Add("properties", cmd)
	return cmd
}()

var getPendingEntries = func() *fctCmd {
	type pendingEntry struct {
		EntryHash string
		ChainID   string
	}

	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get pendingentries [-E]"
	cmd.description = "Get all pending entries, which may not yet be written" +
		" to blockchain. -E EntryHash."
	cmd.execFunc = func(args []string) {
		os.Args = args
		edisp := flag.Bool(
			"E",
			false,
			"display only the Entry Hashes",
		)
		flag.Parse()
		args = flag.Args()

		entries, err := factom.GetPendingEntries()
		if err != nil {
			errorln(err)
			return
		}

		var entList []pendingEntry
		err = json.Unmarshal([]byte(entries), &entList)
		if err != nil {
			errorln(err)
			return
		}
		for _, ents := range entList {
			switch {
			case *edisp:
				fmt.Println(ents.EntryHash)
			default:
				fmt.Println("ChainID:", ents.ChainID)
				fmt.Println("Entryhash:", ents.EntryHash)
				fmt.Println("")
			}
		}
	}
	help.Add("get pendingentries", cmd)
	return cmd
}()

var getPendingTransactions = func() *fctCmd {
	type LineItem struct {
		Amount      float64
		Address     string
		UserAddress string
	}

	type pendingTransaction struct {
		TransactionID string
		Inputs        []LineItem
		Outputs       []LineItem
		ECOutputs     []LineItem
	}

	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get pendingtransactions [-T]"
	cmd.description = "Get all pending factoid transacitons, which may not yet be written to blockchain. -T TxID."
	cmd.execFunc = func(args []string) {
		os.Args = args
		tdisp := flag.Bool(
			"T",
			false,
			"display only the Transaction IDs",
		)
		flag.Parse()
		args = flag.Args()

		trans, err := factom.GetPendingTransactions()
		if err != nil {
			errorln(err)
			return
		}

		var transList []pendingTransaction
		err = json.Unmarshal([]byte(trans), &transList)
		if err != nil {
			errorln(err)
			return
		}
		for _, tran := range transList {
			if len(tran.Inputs) != 0 {
				switch {
				case *tdisp:
					fmt.Println(tran.TransactionID)
				default:
					fmt.Println("TxID:", tran.TransactionID)
					for _, in := range tran.Inputs {
						fmt.Println("Input:", in.UserAddress, in.Amount)
					}

					if len(tran.Outputs) != 0 {

						for _, out := range tran.Outputs {
							fmt.Println("Output:", out.UserAddress, out.Amount)
						}
					}

					if len(tran.ECOutputs) != 0 {
						for _, ecout := range tran.ECOutputs {
							fmt.Println("ECOutput:", ecout.UserAddress, ecout.Amount)
						}
					}
					fmt.Println("")
				}
			}
		}

	}
	help.Add("get pendingtransactions", cmd)
	return cmd
}()
