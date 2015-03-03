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
	server = "localhost:8088"
	wallet = "wallet"
)

func main() {
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
	case "bintx":
		err := bintx(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "balance":
		err := balance(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "buy":
		err := buy(args)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	case "factoidtx":
		err := factoidtx(args)
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
	default:
		man("default")
	}
}
