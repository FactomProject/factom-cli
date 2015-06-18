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

    buy #amt                "buy amt entry credits for the pubkey in the wallet"

    eckey
        new                 "generate a new eckey"
        pub                 "print the pubkey from the wallet"
        
    get
        head                "get current dbhead"
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

    genfactoidaddr name     "generate a new address, and give it a name"
        
    genentrycreditaddr name "generate an Entry Credit address, and give it"
                            "a name"
    balance
        ec key|address      "entry credit balance of eckey"
        fct key|address     "factoid balance of factoid"
            
    newtransaction key      "create a new transaction.  The key is used to"
                            "add inputs, outputs, and ecoutputs (to buy   "
                            "entry credits).  Once the transaction is built,"
                            "call validate, and if all is good, submit"
                            
    addinput                "Add an input to a transaction"
        key name amount     "Use the name supplied to genfactoidaddr"
        key address amount  "Use an address"

    addoutput               "Add an output to a transaction"
        key name amount     "Use the name supplied to genfactoidaddr"
        key address amount  "Use an address"
        
    addecoutput             "Add an ecoutput (purchase of entry credits" to "
                            "a transaction"
        key name amount     "Use the name supplied to genfactoidaddr"
        key address amount  "Use an address"
    
    sign key                "Sign the transaction specified by the key"
    
    submit key              "Submit the transaction specified by the key"
                            "to Factom"
`

// man returns an usage error string for the specified sub command.
func man(s string) error {
    m := map[string]string{
        "testcredit":     "factom-cli testcredit [key]",
        "balance":        "factom-cli balance ec|fct [key]",
        "buy":            "factom-cli buy #amt",
        "get":            "factom-cli get head|dblock|chain|eblock|entry",
        "getHead":        "factom-cli get head",
        "getDBlock":      "factom-cli get dblock [keymr]",
        "getChain":       "factom-cli get chain [chainid]",
        "getEBlock":      "factom-cli get eblock [keymr]",
        "getEntry":       "factom-cli get entry [hash]",
        "help":           "factom-cli help [subcommand]",
        "mkchain":        "factom-cli mkchain [-e extid ...] <stdin>",
        "genfactoidaddr": "factom-cli genfactoidaddr name",
        "newtransaction": "factom-cli newtransaction key",
        "addinput":       "factom-cli addinput key name|address amount",
        "addoutput":      "factom-cli addoutput key name|address amount",
        "addecoutput":    "factom-cli addecoutput key name|address amount",
        "validate":       "factom-cli validate key",
        "submit":         "factom-cli submit key",
        "eckey":          "factom-cli eckey new|pub",
        "put":            "factom-cli put [-e extid ...] <stdin>",
        "default":        usage,
    }

    if m[s] != "" {
        return fmt.Errorf(m[s])
    }
    return fmt.Errorf(m["default"])
}
