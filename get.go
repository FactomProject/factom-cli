// Copyright 2017 Factom Foundation
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
		// c.Handle("abheight", Abheight)
		c.Handle("allentries", getAllEntries)
		c.Handle("chainhead", getChainHead)
		// c.Handle("dbheight", Dbheight)
		c.Handle("dblock", getDBlock)
		c.Handle("eblock", getEBlock)
		// c.Handle("ecbheight", Ecbheight)
		c.Handle("entry", getEntry)
		// c.Handle("fbheight", Fbheight)
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

		head, inPL, err := factom.GetChainHead(chainid)
		if err != nil {
			errorln(err)
			return
		}

		if head == "" && inPL {
			errorln(factom.ErrChainPending)
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
	cmd.helpMsg = "factom-cli get dblock [-r] KEYMR"
	cmd.description = "Get dblock contents by Key Merkle Root."
	cmd.execFunc = func(args []string) {
		os.Args = args
		rawout := flag.Bool(
			"r",
			false,
			"display the hex encoding of the raw Directory Block",
		)
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		keymr := args[0]
		dblock, raw, err := factom.GetDBlock(keymr)
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *rawout:
			fmt.Printf("%x\n", raw)
		default:
			fmt.Println(dblock)
		}
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
	cmd.helpMsg = "factom-cli get head [-RHKAVNBPFTDC]"
	cmd.description = "Get the latest completed Directory Block"
	cmd.execFunc = func(args []string) {
		os.Args = args
		rdisp := flag.Bool("R", false, "display the hex encoding of the raw Directory Block")
		hdisp := flag.Bool("H", false, "display only the Directory Block Hash")
		kdisp := flag.Bool("K", false, "display only the Directory Block Key Merkel Root")
		adisp := flag.Bool("A", false, "display only the Directory Block Header Hash")
		vdisp := flag.Bool("V", false, "display only the Directory Block Header Version")
		ndisp := flag.Bool("N", false, "display only the Network ID")
		bdisp := flag.Bool("B", false, "display only the Directory Block Body Merkel Root")
		pdisp := flag.Bool("P", false, "display only the Previous Directory Block Key Merkel Root")
		fdisp := flag.Bool("F", false, "display only the Previous Directory Block Full Hash")
		tdisp := flag.Bool("T", false, "display only the Directory Block Timestamp")
		ddisp := flag.Bool("D", false, "display only the Directory Block Height")
		cdisp := flag.Bool("C", false, "display only the Directory Block Count")
		flag.Parse()
		args = flag.Args()

		head, err := factom.GetDBlockHead()
		if err != nil {
			errorln(err)
			return
		}
		dblock, raw, err := factom.GetDBlock(head)
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *rdisp:
			fmt.Printf("%x\n", raw)
		case *hdisp:
			fmt.Println(dblock.DBHash)
		case *kdisp:
			fmt.Println(dblock.KeyMR)
		case *adisp:
			fmt.Println(dblock.HeaderHash)
		case *vdisp:
			fmt.Println(dblock.Header.Version)
		case *ndisp:
			fmt.Println(dblock.Header.NetworkID)
		case *bdisp:
			fmt.Println(dblock.Header.BodyMR)
		case *pdisp:
			fmt.Println(dblock.Header.PrevKeyMR)
		case *fdisp:
			fmt.Println(dblock.Header.PrevFullHash)
		case *tdisp:
			fmt.Println(dblock.Header.Timestamp)
		case *ddisp:
			fmt.Println(dblock.Header.DBHeight)
		case *cdisp:
			fmt.Println(dblock.Header.BlockCount)
		default:
			fmt.Println(dblock)
		}
	}
	help.Add("get head", cmd)
	return cmd
}()

var getHeights = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get heights [-DLBE]"
	cmd.description = "Get the current heights of various items in factomd."
	cmd.execFunc = func(args []string) {
		os.Args = args
		ddisp := flag.Bool("D", false, "display only the Directory Block height")
		ldisp := flag.Bool("L", false, "display only the Leader height")
		bdisp := flag.Bool("B", false, "display only the EntryBlock height")
		edisp := flag.Bool("E", false, "display only the Entry height")
		flag.Parse()
		args = flag.Args()

		heights, err := factom.GetHeights()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *ddisp:
			fmt.Println(heights.DirectoryBlockHeight)
		case *ldisp:
			fmt.Println(heights.LeaderHeight)
		case *bdisp:
			fmt.Println(heights.EntryBlockHeight)
		case *edisp:
			fmt.Println(heights.EntryHeight)
		default:
			fmt.Print(heights.String())
		}
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
	cmd.helpMsg = "factom-cli properties [-CFAWL]"
	cmd.description = "Get version information about factomd and the factom wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		cdisp := flag.Bool("C", false, "display only the CLI version")
		fdisp := flag.Bool("F", false, "display only the factomd version")
		adisp := flag.Bool("A", false, "display only the factomd API version")
		wdisp := flag.Bool("W", false, "display only the factom-wallet version")
		ldisp := flag.Bool("L", false, "display only the wallet API version")
		flag.Parse()
		args = flag.Args()

		props, err := factom.GetProperties()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *cdisp:
			fmt.Println(FactomcliVersion)
		case *fdisp:
			fmt.Println(props.FactomdVersion)
		case *adisp:
			fmt.Println(props.FactomdAPIVersion)
		case *wdisp:
			fmt.Println(props.WalletVersion)
		case *ldisp:
			fmt.Println(props.WalletAPIVersion)
		default:
			fmt.Println("CLI Version:", FactomcliVersion)
			fmt.Println(props)
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
		Fees          int64
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
						fmt.Println("Input:", in.UserAddress, in.Amount/1e9)
					}

					if len(tran.Outputs) != 0 {

						for _, out := range tran.Outputs {
							fmt.Println("Output:", out.UserAddress, out.Amount/1e9)
						}
					}

					if len(tran.ECOutputs) != 0 {
						for _, ecout := range tran.ECOutputs {
							fmt.Println("ECOutput:", ecout.UserAddress, ecout.Amount/1e9)
						}
					}
					if tran.Fees != 0 {
						fmt.Printf("Fees: %8.8f", float64(tran.Fees)/1e9)
					}
					fmt.Println("")
				}
			}
		}

	}
	help.Add("get pendingtransactions", cmd)
	return cmd
}()
