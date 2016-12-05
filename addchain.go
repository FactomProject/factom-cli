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

var addchain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addchain [-f -e EXTID1 -e EXTID2 -x BINEXTID3" +
		" ...] ECADDRESS <STDIN>"
	cmd.description = "Create a new Factom Chain. Read data for the First" +
		" Entry from stdin. Use the Entry Credits from the specified address."
	cmd.execFunc = func(args []string) {
		var (
			eAcii extidsASCII
			eHex  extidsHex
		)
		os.Args = args
		exidCollector = make([][]byte, 0)
		flag.Var(&eAcii, "e", "external id for the entry in ascii")
		flag.Var(&eHex, "x", "external id for the entry in hex")
		fflag := flag.Bool(
			"f",
			false,
			"force the chain to commit and reveal without waiting on any"+
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

		c := factom.NewChain(e)

		if !*fflag {
			if factom.ChainExists(c.ChainID) {
				errorln("Chain", c.ChainID, "already exists")
				return
			}

			// check ec address balance
			balance, err := factom.GetECBalance(ecpub)
			if err != nil {
				errorln(err)
				return
			}
			if cost, err := factom.EntryCost(c.FirstEntry); err != nil {
				errorln(err)
				return
			} else if balance < int64(cost)+10 {
				errorln("Not enough Entry Credits")
				return
			}
		}

		// commit the chain
		txid, err := factom.CommitChain(c, ec)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println("CommitTxID:", txid)

		if !*fflag {
			if _, err := waitOnCommitAck(txid); err != nil {
				errorln(err)
				return
			}
		}

		// reveal chain
		hash, err := factom.RevealChain(c)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println("ChainID:", c.ChainID)
		fmt.Println("Entryhash:", hash)

		if !*fflag {
			if _, err := waitOnRevealAck(txid); err != nil {
				errorln(err)
				return
			}
		}
	}
	help.Add("addchain", cmd)
	return cmd
}()

var composechain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli composechain [-f -e EXTID1 -e EXTID2 -x" +
		" BINEXTID3 ...] ECADDRESS <STDIN>"
	cmd.description = "Create API calls to create a new Factom Chain. Read" +
		" data for the First Entry from stdin. Use the Entry Credits from the" +
		" specified address."
	cmd.execFunc = func(args []string) {
		var (
			eAcii extidsASCII
			eHex  extidsHex
		)
		os.Args = args
		exidCollector = make([][]byte, 0)
		flag.Var(&eAcii, "e", "external id for the entry in ascii")
		flag.Var(&eHex, "x", "external id for the entry in hex")
		fflag := flag.Bool(
			"f",
			false,
			"force the chain to commit and reveal without waiting on any"+
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

		c := factom.NewChain(e)

		if !*fflag {
			if factom.ChainExists(c.ChainID) {
				errorln("Chain", c.ChainID, "already exists")
				return
			}

			// check ec address balance
			balance, err := factom.GetECBalance(ecpub)
			if err != nil {
				errorln(err)
				return
			}
			if cost, err := factom.EntryCost(c.FirstEntry); err != nil {
				errorln(err)
				return
			} else if balance < int64(cost)+10 {
				errorln("Not enough Entry Credits")
				return
			}
		}

		j, err := factom.ComposeChainCommit(c, ec)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(j)

		j, err = factom.ComposeChainReveal(c)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(j)
	}
	help.Add("composechain", cmd)
	return cmd
}()
