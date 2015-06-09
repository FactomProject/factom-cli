// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"code.google.com/p/gcfg"
)

type CliConf struct {
	Main struct {
		Server	string
		Wallet	string
	}
	Entry struct {
		Chainid	string
		Extid	string
	}
}

const defaultConf = `
[main]
Server	= localhost:8088
Wallet	= ""
[entry]
Chainid	= ""
Extid	= ""
`

func ReadConfig() *CliConf {
	cfg := new(CliConf)
	filename := os.Getenv("HOME")+"/.factom/factom-cli.conf"
	err := gcfg.ReadFileInto(cfg, filename)
	if err != nil {
		gcfg.ReadStringInto(cfg, defaultConf)
	}
	return cfg
}
