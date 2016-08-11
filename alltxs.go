// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/factoid"
)

var getAllTransactions = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get alltxs"
	cmd.description = "Get the entire history of transactions in Factom"
	cmd.execFunc = func(args []string) {
		zHash := "0000000000000000000000000000000000000000000000000000000000000000"
		fblockID := "000000000000000000000000000000000000000000000000000000000000000f"

		dbhead, err := factom.GetDBlockHead()
		if err != nil {
			errorln(err)
			return
		}
		dblock, err := factom.GetDBlock(dbhead)
		if err != nil {
			errorln(err)
			return
		}

		var fblockmr string
		for _, eblock := range dblock.EntryBlockList {
			if eblock.ChainID == fblockID {
				fblockmr = eblock.KeyMR
			}
		}
		if fblockmr == "" {
			errorln("no fblock in current dblock")
			return
		}

		// get the most recent block
		p, err := factom.GetRaw(fblockmr)
		if err != nil {
			errorln(err)
			return
		}
		fblock, err := factoid.UnmarshalFBlock(p)
		if err != nil {
			errorln(err)
			return
		}

		for fblock.GetPrevKeyMR().String() != zHash {
			txs := fblock.GetTransactions()
			for _, tx := range txs {
				fmt.Println(tx)
			}
			p, err := factom.GetRaw(fblock.GetPrevKeyMR().String())
			if err != nil {
				errorln(err)
				return
			}
			fblock, err = factoid.UnmarshalFBlock(p)
			if err != nil {
				errorln(err)
				return
			}
		}

		// print the first fblock
		txs := fblock.GetTransactions()
		for _, tx := range txs {
			fmt.Println(tx)
		}
	}
	help.Add("get alltxs", cmd)
	return cmd
}()
