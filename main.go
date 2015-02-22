// Copyright (c) 2015 FactomProject/FactomCode Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.
package main

import (
	"log"
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
			log.Fatalln(err)
		}
	case "buy":
		err := buy(args)
		if err != nil {
			log.Fatalln(err)
		}
	case "factoidtx":
		err := factoidtx(args)
		if err != nil {
			log.Fatalln(err)
		}
	case "get":
		err := get(args)
		if err != nil {
			log.Fatalln(err)
		}
	case "help":
		if len(args) < 2 {
			err := help("help")			
			if err != nil {
				log.Fatalln(err)
			}
		}
		err := help(args[1])			
		if err != nil {
			log.Fatalln(err)
		}
	case "mkchain":
		err := mkchain(args)
		if err != nil {
			log.Fatalln(err)
		}
	case "put":
		err := put(args)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		help(args[0])
	}
}
