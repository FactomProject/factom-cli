// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
)

var (
	server string
	wallet string
)

const Version = "0.1.6.0"

func main() {
	cfg := ReadConfig().Main
	server = cfg.Server
	wallet = cfg.Wallet

	// command line flags overwirte conf file
	var (
		hflag = flag.Bool("h", false, "help")
		sflag = flag.String("s", "", "address of api server")
		wflag = flag.String("w", "", "address of wallet api server")
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

	factom.SetServer(server)
	factom.SetWallet(wallet)
	serverFct = wallet

	c := cli.New()
	c.Handle("help", help)
	c.Handle("get", get)
	c.Handle("mkchain", mkchain)
	c.Handle("put", put)
	c.Handle("importaddress", importaddr)
	c.Handle("newaddress", generateaddress)
	c.Handle("generateaddress", generateaddress)
	c.Handle("balances", getaddresses)
	c.Handle("balance", balance)
	c.Handle("getaddresses", getaddresses)
	c.Handle("transactions", transactions)
	c.Handle("newtransaction", fctnewtrans)
	c.Handle("deletetransaction", fctdeletetrans)
	c.Handle("addinput", fctaddinput)
	c.Handle("addoutput", fctaddoutput)
	c.Handle("addecoutput", fctaddecoutput)
	c.Handle("sign", fctsign)
	c.Handle("submit", fctsubmit)
	c.Handle("getfee", fctgetfee)
	c.Handle("addfee", fctaddfee)
	c.Handle("subfee", fctsubfee)
	c.Handle("properties", fctproperties)
	c.Handle("list", getlist)
	c.Handle("listj", getlistj)
	c.HandleDefault(help)
	c.Execute(args)
}

func errorln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
