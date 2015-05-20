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
	
	balance
		ec                  "entry credit balance of eckey"
		factoid             "factoid balance of factoid"

	buy
		amt #n              "buy n entry credits for 'wallet'"

	eckey
		new                 "generate a new eckey"
		pub                 "print the pubkey from the wallet"
		
	fatoidtx [dest] [amt]   "create and submit a factoid transaction"

	get
		dbinfo "hash"       "get dbinfo by hash"
		dblock "hash"       "get dblock by hash"
		dblocks #from #to   "get dblocks by range"
		eblock "merkelroot" "get eblock by merkel root"
		entry "hash"        "get entry by hash"
		height              "get current height of dblock chain"
		
	help [command]          "print help message for the sub-command"

	mkchain [opt] [name]    "create a new factom chain with 'name'. read"
                            "the data for the first entry from stdin"
		-e externalid       "externalid for the first entry

	put                     "read data from stdin and write to factom"
		-e [externalid]     "specify an exteral id for the factom entry. -e"
                            "can be used multiple times"
		-c [chainid]        "spesify the chain that the entry belongs to"
`

// man returns an usage error string for the specified sub command.
func man(s string) error {
	m := map[string]string{
		"balance":    "factom-cli balance ec|factoid [wallet]",
		"buy":        "factom-cli buy amt",
		"factoidtx":  "factom-cli factoidtx addr amt",
		"get":        "factom-cli get height|dblocks|eblocks|entry",
		"getDirBlockInfo":  "factom-cli get dbinfo [hash]",
		"getDBlock":  "factom-cli get dblock [hash]",
		"getDBlocks": "factom-cli get dblocks #from #to",
		"getEBlock":  "factom-cli get eblock [merkelroot]",
		"getEntry":   "factom-cli get entry [entryhash]",
		"getHeight":  "factom-cli get height",
		"help":       "factom-cli help [subcommand]",
		"mkchain":    "factom-cli mkchain [-e extid ...] name",
		"eckey":      "factom-cli eckey new|pub",
		"put":        "factom-cli put [-e extid ...] <stdin>",
		"default":    usage,
	}

	if m[s] != "" {
		return fmt.Errorf(m[s])
	}
	return fmt.Errorf(m["default"])
}
