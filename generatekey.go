// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
    "flag"
    "fmt"
    "os"
    "github.com/FactomProject/factom"
)

// Generates a new Address
func generateAddress(args []string) error {
    
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 1 {
        return man("generatefactoidaddress")
    }
    
    Addr,err := factom.GenerateFactoidAddress(args[0])
    if err != nil {
        fmt.Println(err)
    }else{
        fmt.Println(args[0]," = ",Addr)
    }
    
    return nil
    
}