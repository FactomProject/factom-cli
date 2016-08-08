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
			eAcii extidsAscii
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

		if _, err := factom.GetChainHead(c.ChainID); err == nil {
			// no error means the client found the chain
			errorln("Chain", c.ChainID, "already exists")
			return
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}
		// commit the chain
		if txID, err := factom.CommitChain(c, ec); err != nil {
			errorln(err)
			return
		} else {
			fmt.Println("Commiting Chain Transaction ID:", txID)
		}

		// TODO - get commit acknowledgement

		// reveal chain
		if hash, err := factom.RevealChain(c); err != nil {
			errorln(err)
			return
		} else {
			fmt.Println("ChainID:", c.ChainID)
			fmt.Println("Entryhash:", hash)
		}
		// ? get reveal ack
	}
	help.Add("addchain", cmd)
	return cmd
}()
