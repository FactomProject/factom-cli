// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/FactomProject/factom"
)

type extids []string

func (e *extids) String() string {
	return fmt.Sprint(*e)
}

func (e *extids) Set(s string) error {
	*e = append(*e, s)
	return nil
}

// put commits then reveals an entry to factomd
func put(args []string) error {
	os.Args = args
	var (
		cid  = flag.String("c", "", "hex encoded chainid for the entry")
		eids extids
	)
	flag.Var(&eids, "e", "external id for the entry")
	flag.Parse()
	args = flag.Args()
	
	if len(args) < 1 {
		return man("put")
	}
	name := args[0]
	
	e := factom.NewEntry()
	
	// use the default chainid and extids from the config file
	econf := ReadConfig().Entry
	if econf.Chainid != "" {
		e.ChainID = econf.Chainid
	}
	if *cid != "" {
		e.ChainID = *cid
	}
	if econf.Extid != "" {
		e.ExtIDs = append(e.ExtIDs, econf.Extid)
	}
	
	for _, v := range eids {
		e.ExtIDs = append(e.ExtIDs, hex.EncodeToString([]byte(v)))
	}

	// Entry.Content is read from stdin
	if p, err := ioutil.ReadAll(os.Stdin); err != nil {
		return err
	} else if size := len(p); size > 10240 {
		return fmt.Errorf("Entry of %d bytes is too large", size)
	} else {
		e.Content = hex.EncodeToString(p)
	}
	
	if err := factom.CommitEntry(e, name); err != nil {
		return err
	}
	time.Sleep(10 * time.Second)
	if err := factom.RevealEntry(e); err != nil {
		return err
	}

	return nil
}
