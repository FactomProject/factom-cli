// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/FactomProject/factom"
)

func mkchain(args []string) {
	os.Args = args
	var (
		eids extids
	)
	flag.Var(&eids, "e", "external id for the entry")
	flag.Parse()
	args = flag.Args()

	if len(args) < 1 {
		man("mkchain")
		return
	}

	factom.SetServer(server)

	name := args[0]

	e := factom.NewEntry()

	for _, id := range eids {
		e.ExtIDs = append(e.ExtIDs, []byte(id))
	}

	// Entry.Content is read from stdin
	if p, err := ioutil.ReadAll(os.Stdin); err != nil {
		errorln(err)
		return
	} else if size := len(p); size > 10240 {
		errorln(fmt.Errorf("Entry of %d bytes is too large", size))
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

	fmt.Println("Creating Chain:", c.ChainID)
	if err := factom.CommitChain(c, name); err != nil {
		errorln(err)
		return
	}
	time.Sleep(10 * time.Second)
	if err := factom.RevealChain(c); err != nil {
		errorln(err)
		return
	}
}
