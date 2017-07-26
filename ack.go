// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/FactomProject/factom"
)

var status = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli status TxID|FullTx"
	cmd.description = "Returns information about a factoid transaction, or an" +
		" entry / entry credit transaction"
	cmd.execFunc = func(args []string) {
		os.Args = args
		tdisp := flag.Bool("T", false, "display the transaction id only")
		sdisp := flag.Bool("S", false, "display the transaction status only")
		udisp := flag.Bool(
			"U",
			false,
			"display the unix time of the transaction",
		)
		ddisp := flag.Bool("D", false, "display the time of the transaction")
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
				t, err := factom.GetTmpTransaction(tx)
				if err != nil {
					errorln(err)
					return
				}
				txID = t.TxID
			} else {
				fullTx = strings.Replace(tx, "\"", "", -1)
			}
		}

		// Putting all 0s indicates an entry. If this fails, we will have to check if  it's a commit
		eack, err := factom.EntryRevealACK(txID, fullTx, "0000000000000000000000000000000000000000000000000000000000000000")
		if err != nil {
			errorln(err)
			return
		}

		// Check if its an entry commit
		ecack, err := factom.EntryCommitACK(txID, fullTx)
		if err != nil {
			errorln(err)
			return
		}

		// If entryack returns unknown, and the commit is found to have the commitTxID
		if ecack != nil && (eack == nil || (eack != nil && eack.EntryData.Status == "Unknown")) {
			if ecack.CommitTxID != "" && ecack.CommitTxID == txID && ecack.CommitData.Status != "Unknown" {
				switch {
				case *tdisp:
					fmt.Println(ecack.CommitTxID)
				case *sdisp:
					fmt.Println(ecack.CommitData.Status)
				case *udisp:
					fmt.Println(ecack.CommitData.TransactionDate)
				case *ddisp:
					fmt.Println(ecack.CommitData.TransactionDateString)
				default:
					fmt.Println(ecack)
				}

				return
			}
		}

		// If it's not a commit, it could be an entry hash that is unknown. We filtered that out earlier, so now include it
		if eack != nil {
			// You searched for an entry by hash
			if eack.EntryHash != "" && eack.EntryHash == txID && (eack.EntryData.Status != "Unknown" || eack.CommitData.Status != "Unknown") {
				switch {
				case *tdisp:
					fmt.Println(eack.EntryHash)
				case *sdisp:
					fmt.Println(eack.EntryData.Status)
				case *udisp:
					fmt.Println(eack.EntryData.TransactionDate)
				case *ddisp:
					fmt.Println(eack.EntryData.TransactionDateString)
				default:
					fmt.Println(eack)
				}
				return
			}
		}

		// Check if its a factoid transaction
		fcack, err := factom.FactoidACK(txID, fullTx)
		if err != nil {
			errorln(err)
			return
		}

		if fcack != nil {
			if fcack.Status != "Unknown" {
				switch {
				case *tdisp:
					fmt.Println(fcack.TxID)
				case *sdisp:
					fmt.Println(fcack.Status)
				case *udisp:
					fmt.Println(fcack.TransactionDate)
				case *ddisp:
					fmt.Println(fcack.TransactionDateString)
				default:
					fmt.Println(fcack)
				}
				return
			}
		}

		fmt.Printf("Entry / transaction not found.\n")
	}
	help.Add("status", cmd)
	return cmd
}()
