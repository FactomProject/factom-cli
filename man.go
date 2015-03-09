// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

const usage = `factom-cli [options] [subcommand]
	-s [server]				"address for the api server"
	-w [wallet]				"address for the wallet"
	
	balance
		ec					"entry credit balance of 'wallet'"
		factoid				"factoid balance of 'wallet'"

	buy
		amt #n				"buy n entry credits for 'wallet'"

	fatoidtx [dest] [amt]	"create and submit a factoid transaction"

	get
		height				"get current height of dblock chain"
		dblocks #from #to	"get dblocks by range"
		eblock "merkelroot"	"get eblock by merkel root"
		entry "hash"		"get entry by hash"
		
	help [command]			"print help message for the sub-command"

	mkchain [opt] [name]	"create a new factom chain with 'name'. read"
							"the data for the first entry from stdin"
		-e externalid		"externalid for the first entry

	put						"read data from stdin and write to factom"
		-e [externalid]		"specify an exteral id for the factom entry. -e" 								"can be used multiple times"
		-c [chainid]		"spesify the chain that the entry belongs to"
`

// man returns an usage error string for the specified sub command.
func man(s string) error {
	m := map[string]string{
		"balance":    "factom-cli balance ec|factoid [wallet]",
		"buy":        "factom-cli buy amt",
		"factoidtx":  "factom-cli factoidtx addr amt",
		"get":        "factom-cli get height|dblocks|eblocks|entry",
		"getDBlocks": "factom-cli get dblocks #from #to",
		"getEBlock":  "factom-cli get eblocks [merkelroot]",
		"getEntry":   "factom-cli get entry [entryhash]",
		"getHeight":  "factom-cli get height",
		"help":       "factom-cli help [subcommand]",
		"mkchain":    "factom-cli mkchain [-e extid ...] name",
		"put":        "factom-cli put [-e extid ...] <stdin>",
		"default":    usage,
	}

	if m[s] != "" {
		return fmt.Errorf(m[s])
	}
	return fmt.Errorf(m["default"])
}
