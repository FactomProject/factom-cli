// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
)

var ack = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli ack TxID|FullTx"
	cmd.description = "Returns information about a factoid transaction, or an entry / entry credit transaction"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		tx := args[0]

		txID := ""
		fullTx := ""

		_, err := hex.DecodeString(tx)
		if len(tx) == 64 && err == nil {
			txID = tx
		} else {
			if len(tx) < 64 || err != nil {
				h, err := factom.TransactionHash(tx)
				if err != nil {
					errorln(err)
					return
				}
				txID = h
			} else {
				fullTx = tx
			}
		}

		resp1, err1 := factom.FactoidACK(txID, fullTx)
		resp2, err2 := factom.EntryACK(txID, fullTx)

		if err1 != nil && err2 != nil {
			errorln(err1)
			return
		}
		if err1 != nil {
			str, err := json.MarshalIndent(resp2, "", "\t")
			if err != nil {
				errorln(err)
				return
			}
			fmt.Printf("%s\n", str)
			return
		}
		if err2 != nil {
			str, err := json.MarshalIndent(resp1, "", "\t")
			if err != nil {
				errorln(err)
				return
			}
			fmt.Printf("%s\n", str)
			return
		}
	}
	help.Add("ack", cmd)
	return cmd
}()
