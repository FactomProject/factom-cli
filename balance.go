// Copyright 2015 Factom Foundation
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
func balance(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 2 {
		fmt.Println("Too few arguments")
		man("balance")
		return fmt.Errorf("Too Few Arguments")
	}

	switch args[0] {
	case "ec":
		return ecbalance(args[1])
	case "fct":
		return fctbalance(args[1])
	default:
		fmt.Println("Must specify an address type, either 'ec' or 'fct'")
		man("balance")
		return fmt.Errorf("")
	}

}

func ecbalance(addr string) error {
	factom.SetServer(server)

	if Utility.IsValidAddress(addr) && strings.HasPrefix(addr,"FA") {
		fmt.Println("Not a valid Entry Credit Address")
		return fmt.Errorf("Not a valid Entry Credit Address")
	}
	if b, err := factom.ECBalance(addr); err != nil {
		fmt.Println("Address undefined or invalid")
		return err
	} else {
		fmt.Println("Balance of ", addr, " = ", b)
	}

	return nil
}

func fctbalance(addr string) error {
	factom.SetServer(server)

	if Utility.IsValidAddress(addr) && strings.HasPrefix(addr,"EC") {
		fmt.Println("Not a valid Entry Credit Address")
		return fmt.Errorf("Not a valid Entry Credit Address")
	}
	
	if b, err := factom.FctBalance(addr); err != nil {
		fmt.Println("Address undefined or invalid: "+err.Error())
		return err
	} else {
		fmt.Println("Balance of ", addr, " = ", fct.ConvertDecimal(uint64(b)))
	}

	return nil
}
