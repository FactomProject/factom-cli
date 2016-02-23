// Copyright 2016 Factom Foundation
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

type extids []string

func (e *extids) String() string {
	return fmt.Sprint(*e)
}

func (e *extids) Set(s string) error {
	*e = append(*e, s)
	return nil
}

// put commits then reveals an entry to factomd
var put = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli put -c CHAINID [-e EXTID1 -e EXTID2 ...] NAME <STDIN>"
	cmd.description = "Read data from stdin and write to factom using the named entry credit address."
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			cid  = flag.String("c", "", "hex encoded chainid for the entry")
			eids extids
		)
		flag.Var(&eids, "e", "external id for the entry")
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
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
			e.ExtIDs = append(e.ExtIDs, []byte(econf.Extid))
		}

		for _, v := range eids {
			e.ExtIDs = append(e.ExtIDs, []byte(v))
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

		// Make sure the Chain exists before writing the Entry
		if _, err := factom.GetChainHead(e.ChainID); err != nil {
			errorln("Chain:", e.ChainID, "does not exist")
			return
		}

		fmt.Printf("Creating Entry: %x\n", e.Hash())
		if err := factom.CommitEntry(e, name); err != nil {
			errorln(err)
			return
		}
		time.Sleep(10 * time.Second)
		if err := factom.RevealEntry(e); err != nil {
			errorln(err)
			return
		}

	}
	help.Add("put", cmd)
	return cmd
}()
