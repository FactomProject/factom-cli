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
	cmd.helpMsg = "factom-cli addentry [-e EXTID1 -e EXTID2 ...] ECADDRESS <STDIN>"
	cmd.description = "Create a new Factom Entry. Read data for the Entry from stdin. Use the Entry Credits from the specified address."
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			eids extids
		)
		flag.Var(&eids, "e", "external id for the entry")
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		e := new(factom.Entry)

		for _, id := range eids {
			e.ExtIDs = append(e.ExtIDs, []byte(id))
		}

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

		if _, err := factom.GetChainHead(e.ChainID); err != nil {
			errorln("Chain", e.ChainID, "was not found")
			return
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}
		// commit the chain
		if _, err := factom.CommitEntry(e, ec); err != nil {
			errorln(err)
			return
		}

		// TODO - get commit acknowledgement

		// reveal chain
		if err := factom.RevealEntry(e); err != nil {
			errorln(err)
			return
		}
		// ? get reveal ack
	}
	help.Add("addentry", cmd)
	return cmd
}()
