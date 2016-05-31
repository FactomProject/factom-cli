// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
)

var ack = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli ack fct|e TxID|FullTx"
	cmd.description = "Returns information about a factoid transaction, or an entry / entry credit transaction"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 2 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ackType := args[0]
		tx := args[1]

		txID := ""
		fullTx := ""

		if len(tx) == 64 {
			txID = tx
		} else {
			fullTx = tx
		}

		if ackType == "fct" {
			resp, err := factom.FactoidACK(txID, fullTx)
			if err != nil {
				errorln(err)
				return
			}
			str, err := json.MarshalIndent(resp, "", "\t")
			if err != nil {
				errorln(err)
				return
			}
			fmt.Printf("%s\n", str)
		}
		if ackType == "e" {
			resp, err := factom.EntryACK(txID, fullTx)
			if err != nil {
				errorln(err)
				return
			}
			str, err := json.MarshalIndent(resp, "", "\t")
			if err != nil {
				errorln(err)
				return
			}
			fmt.Printf("%s\n", str)
		}
	}
	help.Add("ack", cmd)
	return cmd
}()
