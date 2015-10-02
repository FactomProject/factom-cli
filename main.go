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

const Version = 1002


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

	case "get":
		get(args)
	case "help":
		help(args)
	case "mkchain":
		mkchain(args)
	case "put":
		put(args)
	// two commands for the same thing
	case "newaddress":
		generateaddress(args)
	case "generateaddress":
		generateaddress(args)
	// two commands for the same thing
	case "balances":
		getaddresses(args)
	case "getaddresses":
		getaddresses(args)
	case "transactions":
		gettransactions(args)
	case "balance":
		balance(args)
	case "newtransaction":
		fctnewtrans(args)
	case "deletetransaction":
		fctdeletetrans(args)
	case "addinput":
		fctaddinput(args)
	case "addoutput":
		fctaddoutput(args)
	case "addecoutput":
		fctaddecoutput(args)
	case "sign":
		fctsign(args)
	case "submit":
		fctsubmit(args)
	case "getfee":
		fctgetfee(args)
	case "addfee":
		fctaddfee(args)
	case "properties":
		fctproperties(args)
	case "list":
		getlist(args)
	case "listj":
		getlistj(args)
	default:
		fmt.Println("Command not found")
		man("default")
	}
}

func errorln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
