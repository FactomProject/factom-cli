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
	"strings"

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

		_, err := hex.DecodeString(strings.Replace(tx, "\"", "", -1))
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
				fullTx = strings.Replace(tx, "\"", "", -1)
			}
		}

		resp1, err1 := factom.FactoidACK(txID, fullTx)
		resp2, err2 := factom.EntryACK(txID, fullTx)
		if err1 != nil && err2 != nil {
			errorln(err1)
			return
		}

		if resp1 != nil {
			if resp1.Status != "Unknown" {
				str, err := json.MarshalIndent(resp1, "", "\t")
				if err != nil {
					errorln(err)
					return
				}
				fmt.Printf("%s\n", str)
				return
			}
		}
		if resp2 != nil {
			if resp2.CommitTxID == "" && resp2.EntryHash == "" {
			} else {
				str, err := json.MarshalIndent(resp2, "", "\t")
				if err != nil {
					errorln(err)
					return
				}
				fmt.Printf("%s\n", str)
				return
			}
		}
		fmt.Printf("Entry / transaction not found.\n")
	}
	help.Add("ack", cmd)
	return cmd
}()
