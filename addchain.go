// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/FactomProject/factom"
)

var exidCollector [][]byte

// extids will be a flag receiver for adding chains and entries
// In ASCII
type extidsAscii []string

func (e *extidsAscii) String() string {
	return fmt.Sprint(*e)
}

func (e *extidsAscii) Set(s string) error {
	*e = append(*e, s)
	exidCollector = append(exidCollector[:], []byte(s))
	return nil
}

// extids will be a flag receiver for adding chains and entries
// In HEX
type extidsHex []string

func (e *extidsHex) String() string {
	return fmt.Sprint(*e)
}

func (e *extidsHex) Set(s string) error {
	*e = append(*e, s)
	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	exidCollector = append(exidCollector[:], b)
	return nil
}

var addchain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addchain [-e EXTID1 -e EXTID2 ...] ECADDRESS <STDIN>"
	cmd.description = "Create a new Factom Chain. Read data for the First Entry from stdin. Use the Entry Credits from the specified address."
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			eAcii extidsAscii
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

		//for _, id := range eAcii {
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
			fmt.Println("Commiting Chain Transaction ID: " + txID)
		}

		// TODO - get commit acknowledgement

		// reveal chain
		if hash, err := factom.RevealChain(c); err != nil {
			errorln(err)
			return
		} else {
			fmt.Println("ChainID  : " + c.ChainID)
			fmt.Println("Entryhash: " + hash)
		}
		// ? get reveal ack
	}
	help.Add("addchain", cmd)
	return cmd
}()
