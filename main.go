// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	server string
	wallet string
)

func main() {
	cfg := ReadConfig().Main
	server = cfg.Server
	wallet = cfg.Wallet
	
	// command line flags overwirte conf file
	var (
		hflag = flag.Bool("h", false, "help")
		sflag = flag.String("s", "", "address of api server")
		wflag = flag.String("w", "", "wallet address")
	)
	flag.Parse()
	args := flag.Args()
	if *sflag != "" {
		server = *sflag
	}
	if *wflag != "" {
		wallet = *wflag
	}
	if *hflag {
		args = []string{"help"}
	}
	if len(args) == 0 {
		args = append(args, "help")
	}

	switch args[0] {
	//TODO remove testcredit before production
	case "testcredit":
		err := testcredit(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "balance":
		err := balance(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "eckey":
		err := eckey(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "get":
		err := get(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "help":
		err := help(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "mkchain":
		err := mkchain(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
    case "put":
        err := put(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
        
    case "generateaddress":
        err := generateaddress(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    case "getaddresses":
        err := getaddresses(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    case "newtransaction":
        err := fctnewtrans(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        } 
    case "addinput":
        err := fctaddinput(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    case "addoutput":
        err := fctaddoutput(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    case "addecoutput":
        err := fctaddecoutput(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    case "sign":
        err := fctsign(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    case "submit":
        err := fctsubmit(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    case "getfee":
        err := fctgetfee(args)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }
    default:
        fmt.Println("Command not found")
		man("default")
	}
}
