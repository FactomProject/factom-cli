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
	cmd.helpMsg = "factom-cli addchain [-e EXTID1 -e EXTID2 -E BINEXTID3 ...] ECADDRESS <STDIN>"
	cmd.description = "Create a new Factom Chain. Read data for the First Entry from stdin. Use the Entry Credits from the specified address."
	cmd.execFunc = func(args []string) {
		var (
			eAcii extidsASCII
			eHex  extidsHex
		)
		os.Args = args
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

		c := factom.NewChain(e)

		if factom.ChainExists(c.ChainID) {
			errorln("Chain", c.ChainID, "already exists")
			return
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}
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
		txID, err := factom.CommitChain(c, ec)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println("Commiting Chain Transaction ID:", txID)

		// TODO - get commit acknowledgement

		// reveal chain
		hash, err := factom.RevealChain(c)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println("ChainID:", c.ChainID)
		fmt.Println("Entryhash:", hash)

		// ? get reveal ack
	}
	help.Add("addchain", cmd)
	return cmd
}()

var composechain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli composechain [-e EXTID1 -e EXTID2 -E BINEXTID3 ...] ECADDRESS <STDIN>"
	cmd.description = "Create API calls to create a new Factom Chain. Read data for the First Entry from stdin. Use the Entry Credits from the specified address."
	cmd.execFunc = func(args []string) {
		var (
			eAcii extidsASCII
			eHex  extidsHex
		)
		os.Args = args
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

		c := factom.NewChain(e)

		if factom.ChainExists(c.ChainID) {
			errorln("Chain", c.ChainID, "already exists")
			//return
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}
		balance, err := factom.GetECBalance(ecpub)
		if err != nil {
			errorln(err)
			//return
		}
		if balance == 0 {
			errorln("Entry Credit balance is zero")
			//return
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
