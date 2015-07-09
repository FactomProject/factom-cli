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
)

// balance prints the current balance of the specified address
func balance(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		fmt.Println(man("balance"))
        return fmt.Errorf("Too Few Arguments")
	}
	
	switch args[0] {
	case "ec":
		return ecbalance(args[1])
	case "fct":
		return fctbalance(args[1])
	default:
        fmt.Println("Must specify an address type, either 'ec' or 'fct'")
		fmt.Println(man("balance"))
        return fmt.Errorf("")
	}

}

func ecbalance(addr string) error {

    if b, err := factom.ECBalance(addr); err != nil {
        fmt.Println(err)
        return err
	} else {
        fmt.Println("Balance of ", addr, " = ", b)
    }
	
	return nil	
}

func fctbalance(addr string) error {


	if b, err := factom.FctBalance(addr); err != nil {
		fmt.Println(err)
        return err
	} else {
        fmt.Println("Balance of ", addr, " = ", strings.TrimSpace(fct.ConvertDecimal(uint64(b))))
	}

	return nil
}
