// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

const usage = `factom-cli [options] [subcommand]
	-s [server]             "address for the api server"
	-w [wallet]             "wallet file"
	
	testcredit [key]        "add 100 test credits to the key. use the wallet"
	                        "file if no key is specified"
	balance
		ec                  "entry credit balance of eckey"
		factoid             "factoid balance of factoid"

	buy #amt                "buy amt entry credits for the pubkey in the wallet"

	eckey
		new                 "generate a new eckey"
		pub                 "print the pubkey from the wallet"
		
	get
		head				"get current dbhead"
		dblock "keymr"      "get dblock by merkel root"
		chain "chainid"     "get ebhead by chainid"
		eblock "keymr"      "get eblock by merkel root"
		entry "hash"        "get entry by hash"
		
	help [command]          "print help message for a sub-command"

	mkchain                 "create a new factom chain. read the data for the"
	                        "first entry from stdin"
		-e externalid       "externalid for the first entry

	put                     "read data from stdin and write to factom"
		-e [externalid]     "specify an exteral id for the factom entry. -e"
                            "can be used multiple times"
		-c [chainid]        "specify the chain that the entry belongs to"
`

// man returns an usage error string for the specified sub command.
func man(s string) error {
	m := map[string]string{
		"testcredit": "factom-cli testcredit [key]",
		"balance":    "factom-cli balance ec|factoid [wallet]",
		"buy":        "factom-cli buy #amt",
		"get":        "factom-cli get head|dblock|chain|eblock|entry",
		"getHead":    "factom-cli get head",
		"getDBlock":  "factom-cli get dblock [keymr]",
		"getChain":   "factom-cli get chain [chainid]",
		"getEBlock":  "factom-cli get eblock [keymr]",
		"getEntry":   "factom-cli get entry [hash]",
		"help":       "factom-cli help [subcommand]",
		"mkchain":    "factom-cli mkchain [-e extid ...] <stdin>",
		"eckey":      "factom-cli eckey new|pub",
		"put":        "factom-cli put [-e extid ...] <stdin>",
		"default":    usage,
	}

	if m[s] != "" {
		return fmt.Errorf(m[s])
	}
	return fmt.Errorf(m["default"])
}
