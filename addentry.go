// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/FactomProject/factom"
)

var addentry = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addentry -c CHAINID [-f -e EXTID1 -e EXTID2 -E" +
		" BEEF1D ...] ECADDRESS <STDIN>"
	cmd.description = "Create a new Factom Entry. Read data for the Entry" +
		" from stdin. Use the Entry Credits from the specified address."
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			cid   = flag.String("c", "", "hex encoded chainid for the entry")
			eAcii extidsASCII
			eHex  extidsHex
		)
		exidCollector = make([][]byte, 0)
		flag.Var(&eAcii, "e", "external id for the entry in ascii")
		flag.Var(&eHex, "E", "external id for the entry in hex")
		fflag := flag.Bool(
			"f",
			false,
			"force the entry to commit and reveal without waiting on any"+
				" acknowledgement checks",
		)
		cdisp := flag.Bool("C", false, "display only the ChainID")
		edisp := flag.Bool("E", false, "display only the Entry Hash")
		tdisp := flag.Bool("T", false, "display only the TxID")

		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		// display normal output iff no display flags are set
		display := true
		if *cdisp || *edisp || *tdisp {
			display = false
		}

		e := new(factom.Entry)

		if *cid == "" {
			fmt.Println(cmd.helpMsg)
			return
		}
		e.ChainID = *cid

		e.ExtIDs = exidCollector

		// Entry.Content is read from stdin
		if p, err := ioutil.ReadAll(os.Stdin); err != nil {
			errorln(err)
			return
		} else if size := len(p); size > 10240 {
			errorln(fmt.Sprintf("Entry of %d bytes is too large", size))
			return
		} else {
			e.Content = p
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}

		if !*fflag {
			if !factom.ChainExists(e.ChainID) {
				errorln("Chain", e.ChainID, "was not found")
				return
			}

			// check ec address balance
			balance, err := factom.GetECBalance(ecpub)
			if err != nil {
				errorln(err)
				return
			}
			if cost, err := factom.EntryCost(e); err != nil {
				errorln(err)
				return
			} else if balance < int64(cost) {
				errorln("Not enough Entry Credits")
				return
			}
		}

		// commit entry
		txid, err := factom.CommitEntry(e, ec)
		if err != nil {
			errorln(err)
			return
		}
		if display {
			fmt.Println("CommitTxID:", txid)
		} else if *tdisp {
			fmt.Println(txid)
		}

		if !*fflag {
			if _, err := waitOnCommitAck(txid); err != nil {
				errorln(err)
				return
			}
		}
		// reveal entry
		hash, err := factom.RevealEntry(e)
		if err != nil {
			errorln(err)
			return
		}
		if !*fflag {
			if _, err := waitOnRevealAck(txid); err != nil {
				errorln(err)
				return
			}
		}
		if display {
			fmt.Println("ChainID:", e.ChainID)
			fmt.Println("Entryhash:", hash)
		} else if *cdisp {
			fmt.Println(e.ChainID)
		} else if *edisp {
			fmt.Println(hash)
		}


	}
	help.Add("addentry", cmd)
	return cmd
}()

var composeentry = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli composeentry -c CHAINID [-f -e EXTID1 -e EXTID2" +
		" -E BEEF1D ...] ECADDRESS <STDIN>"
	cmd.description = "Create API calls to create a new Factom Entry. Read" +
		" data for the Entry from stdin. Use the Entry Credits from the" +
		" specified address."
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			cid   = flag.String("c", "", "hex encoded chainid for the entry")
			eAcii extidsASCII
			eHex  extidsHex
		)
		exidCollector = make([][]byte, 0)
		flag.Var(&eAcii, "e", "external id for the entry in ascii")
		flag.Var(&eHex, "E", "external id for the entry in hex")
		fflag := flag.Bool(
			"f",
			false,
			"force the entry to commit and reveal without waiting on any"+
				" acknowledgement checks",
		)
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		e := new(factom.Entry)

		if *cid == "" {
			fmt.Println(cmd.helpMsg)
			return
		}
		e.ChainID = *cid

		e.ExtIDs = exidCollector

		// Entry.Content is read from stdin
		if p, err := ioutil.ReadAll(os.Stdin); err != nil {
			errorln(err)
			return
		} else if size := len(p); size > 10240 {
			errorln(fmt.Sprintf("Entry of %d bytes is too large", size))
			return
		} else {
			e.Content = p
		}

		commit, reveal, err := factom.WalletComposeEntryCommitReveal(e, ecpub, *fflag)
		if err != nil {
			errorln(err)
			return
		}

		fmt.Println(commit)
		fmt.Println(reveal)
	}
	help.Add("composeentry", cmd)
	return cmd
}()
