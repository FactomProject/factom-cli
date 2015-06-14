// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
)

// balance prints the current balance of the specified wallet
func balance(args []string) error {
    
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("balance")
	}
	switch args[0] {
	case "ec":
		return ecbalance(args)
    case "factoid":
        fallthrough
    case "fct":
        return fctbalance(args)
	default:
		return man("balance")
	}

	panic("Something went really wrong with balance!")
}

func ecbalance(args []string) error {
	var eckey string
	if p, err := ecPubKey(); err != nil {
		return err
	} else {
		eckey = hex.EncodeToString(p[:])
	}
	
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) > 0 {
		eckey = args[0]
	}
	
	if b, err := factom.ECBalance(eckey); err != nil {
		return err
	} else {
		fmt.Println(b)
	}
	
	return nil	
}

func fctbalance(args []string) error {
    type balance struct {
        Balance int64
    }
    
    var fctAdr string
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) > 0 {
       fctAdr = args[0]
    }else{
        fmt.Println("Missing Address")
    } 
    
    if b, err := factom.FctBalance(fctAdr); err != nil {
        return err
    } else {
        fmt.Println("Balance of ",fctAdr," = ", b.Balance)
    }
    
    return nil
    
}