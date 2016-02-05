// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/FactomProject/cli"
	fct "github.com/FactomProject/factoid"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
)

// balance prints the current balance of the specified address
var balance = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli balance ec|fct [key]"
	cmd.description = "If this is an ec balance, returns number of entry credits. If this is a Factoid balance, returns the factoids at that address."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		c := cli.New()
		c.Handle("ec", ecBalance)
		c.Handle("fct", fctBalance)
		c.HandleDefaultFunc(func(args []string) {
			fmt.Println(cmd.helpMsg)
		})
		c.Execute(args)
	}
	help.Add("balance", cmd)
	return cmd
}()

var ecBalance = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli balance ec [addr]"
	cmd.description = "Return number of entry credits at the address"
	cmd.execFunc = func(args []string) {
		var addr string
		if len(args) >= 2 {
			addr = args[1]
		}

		if Utility.IsValidAddress(addr) && strings.HasPrefix(addr, "FA") {
			fmt.Println("Not a valid Entry Credit Address")
		}
		if b, err := factom.ECBalance(addr); err != nil {
			fmt.Println("Address undefined or invalid")
		} else {
			fmt.Println("Balance of ", addr, " = ", b)
		}
	}
	help.Add("balance ec", cmd)
	return cmd
}()

var fctBalance = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli balance fct [addr]"
	cmd.description = "Return number of factoids at the address"
	cmd.execFunc = func(args []string) {
		var addr string
		if len(args) >= 2 {
			addr = args[1]
		}

		if Utility.IsValidAddress(addr) && strings.HasPrefix(addr, "EC") {
			fmt.Println("Not a valid Entry Credit Address")
		}

		if b, err := factom.FctBalance(addr); err != nil {
			fmt.Println("Address undefined or invalid")
		} else {
			fmt.Println("Balance of ", addr, " = ", fct.ConvertDecimal(uint64(b)))
		}
	}
	help.Add("balance fct", cmd)
	return cmd
}()
