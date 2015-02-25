// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}
	
	args := os.Args[1:]

	switch args[0] {
	case "bintx":
		err := bintx(args)
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
