// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

//import (
//	"flag"
//	"fmt"
//	"os"
//
//    fct "github.com/FactomProject/factoid"
//    "github.com/FactomProject/factom"
//)
//
//// balance prints the current balance of the specified wallet
//func balance(args []string) error {
//	os.Args = args
//	flag.Parse()
//	args = flag.Args()
//	if len(args) < 1 {
//		return man("balance")
//	}
//	
//	switch args[0] {
//	case "ec":
//		return ecbalance(args[1])
//	case "fct":
//		return fctbalance(args[1])
//	default:
//		return man("balance")
//	}
//
//}
//
//func ecbalance(addr string) error {
//
//    if b, err := factom.ECBalance(addr); err != nil {
//		return err
//	} else {
//        fmt.Println("Balance of ", addr, " = ", b)
//    }
//	
//	return nil	
//}
//
//func fctbalance(addr string) error {
//
//
//	if b, err := factom.FctBalance(addr); err != nil {
//		return err
//	} else {
//        fmt.Println("Balance of ", addr, " = ", fct.ConvertDecimal(uint64(b)))
//	}
//
//	return nil
//}
