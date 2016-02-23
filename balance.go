// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	fct "github.com/FactomProject/factoid"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
)

// balance prints the current balance of the specified address
var balance = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli balance ADDRESS"
	cmd.description = "If this is an EC Address, returns number of Entry Credits. If this is a Factoid Address, returns the Factoid balance."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		addr := args[0]
		
		if strings.HasPrefix(addr, "FA") {
			if !Utility.IsValidAddress(addr) {
				fmt.Println("Invalid Factoid Address")
			}
			
			if b, err := factom.FctBalance(addr); err != nil {
				fmt.Println("Undefined or invalid address")
			} else {
				fmt.Println(addr, fct.ConvertDecimal(uint64(b)))
			}
		} else if strings.HasPrefix(addr, "EC") {
			if !Utility.IsValidAddress(addr) {
				fmt.Println("Invalid EC Address")
			}
			
			if b, err := factom.ECBalance(addr); err != nil {
				fmt.Println("Undefined or invalid address")
			} else {
				fmt.Println(addr, b)
			}
		} else {
			fmt.Println("Invalid address")
		}
	}
	help.Add("balance", cmd)
	return cmd
}()
