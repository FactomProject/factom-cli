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
	cmd.helpMsg = "factom-cli addentry -c CHAINID [-e EXTID1 -e EXTID2 -E BEEF1D ...] ECADDRESS <STDIN>"
	cmd.description = "Create a new Factom Entry. Read data for the Entry from stdin. Use the Entry Credits from the specified address."
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

		//for _, id := range eids {
		//	e.ExtIDs = append(e.ExtIDs, []byte(id))
		//}
		e.ExtIDs = exidCollector

		// Entry.Content is read from stdin
		if p, err := ioutil.ReadAll(os.Stdin); err != nil {
			errorln(err)
			return
		} else if size := len(p); size > 10240 {
			errorln("Entry of %d bytes is too large", size)
			return
		} else {
			e.Content = p
		}

		if !factom.ChainExists(e.ChainID) {
			errorln("Chain", e.ChainID, "was not found")
			return
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}

		// check ec address balance
		balance, err := factom.GetECBalance(ecpub)
		if err != nil {
			errorln(err)
			return
		}
		if balance == 0 {
			errorln("Entry Credit balance is zero")
			return
		}

		// commit the chain
		txID, err := factom.CommitEntry(e, ec)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println("Committing Entry Transaction ID:", txID)

		// TODO - get commit acknowledgement

		// reveal chain
		hash, err := factom.RevealEntry(e)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println("ChainID:", *cid)
		fmt.Println("Entryhash:", hash)

		// ? get reveal ack
	}
	help.Add("addentry", cmd)
	return cmd
}()

type entryRequest struct {
	Entry factom.Entry `json:"entry"`
	ECPub string       `json:"ecpub"`
}

var composeentry = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli composeentry -c CHAINID [-e EXTID1 -e EXTID2 -E BEEF1D ...] ECADDRESS <STDIN>"
	cmd.description = "Create API calls to create a new Factom Entry. Read data for the Entry from stdin. Use the Entry Credits from the specified address."
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

		//for _, id := range eids {
		//	e.ExtIDs = append(e.ExtIDs, []byte(id))
		//}
		e.ExtIDs = exidCollector

		// Entry.Content is read from stdin
		if p, err := ioutil.ReadAll(os.Stdin); err != nil {
			errorln(err)
			return
		} else if size := len(p); size > 10240 {
			errorln("Entry of %d bytes is too large", size)
			return
		} else {
			e.Content = p
		}

		commit, reveal, err := factom.WalletComposeEntryCommitReveal(e, ecpub)
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
