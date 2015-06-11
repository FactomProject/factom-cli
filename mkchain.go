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

func mkchain(args []string) error {
	os.Args = args
	var (
		eids extids
	)
	flag.Var(&eids, "e", "external id for the entry")
	flag.Parse()
	args = flag.Args()
	
	e := factom.NewEntry()
	
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
	
	priv, err := ecPrivKey()
	if err != nil {
		return err
	}

	c := factom.NewChain(e)
	
	if err := factom.CommitChain(c, priv); err != nil {
		return err
	}
	time.Sleep(10 * time.Second)
	if err := factom.RevealChain(c); err != nil {
		return err
	}
	
	fmt.Println("New Chain:", c.ChainID)

	return nil
}
