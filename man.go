// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

const usage = `
factom-cli [options] [subcommand]

        -c [ChainID]        hex encoded chainid for the entry
        -e [extid]          external id for the entry. -e may be used multiple times
        -h                  display help message
        -s                  path to the factomclient. default: localhost:8083    
        
    get
        head                Get the keymr of the last completed directory block
        height              Get the current directory block height
        dblock keymr        Get dblock contents by merkle root
        chain chainid       Get ebhead by chainid
        eblock keymr        Get eblock by merkle root
        entry hash          Get entry by hash
        firstentry chainid  Get the first entry in a chain
        
    help [command]          Print help message for a sub-command

    mkchain name            Create a new factom chain. read the data for the
                            first entry from stdin.  Use the entry credits
                            from the specified name.
        -e externalid       Externalid for the first entry

    put name                Read data from stdin and write to factom. Use
                            the entry credits from the specified entry credit
                            address
        -e [externalid]     Specify an exteral id for the factom entry. -e
                            can be used multiple times
        -c [chainid]        Specify the chain that the entry belongs to
        
    newaddress or   
    generateaddress         Generate addresses, giving them names.
        ec name             Generate an Entry Credit address, tied to the name
        fct name            Generate a Factoid address, tied to the name
                            Names must be unique, or you will get a
                            Duplicate Name or Invalid Name error.
                            Names are limited to 32 characters
        ec name Es...       Import a secret EC key to the wallet
        fct name Fs...      Import a secret Factoid key to the wallet
        fct name "12 words" Import the Factoid key from the Token Sale
                            
    deletetransaction key   Delete the specified transaction in flight. 
   
    balance key|address     If this is an ec balance, returns number of 
                            entry credits
                            If this is a Factoid balance, returns the 
                            factoids at that address
            
    newtransaction key      Create a new transaction.  The key is used to
                            add inputs, outputs, and ecoutputs (to buy   
                            entry credits).  Once the transaction is built,
                            call validate, and if all is good, submit
                            
    addinput                Add an input to a transaction
        key name amount     Use the name supplied to genfactoidaddr
        key address amount  Use an address

    addoutput               Add an output to a transaction
        key name amount     Use the name supplied to genfactoidaddr
        key address amount  Use an address
        
    addecoutput             Add an ecoutput (purchase of entry credits to 
                            a transaction
        key name amount     Use the name supplied to genfactoidaddr
        key address amount  Use an address
        
    getfee key              Get the current fee required for this 
                            transaction.  If a transaction is specified,
                            then getfee returns the fee due for the 
                            transaction.  If no transaction is provided,
                            then the cost of an Entry Credit is returned.
    
    addfee trans address    Adds the needed fee to the given transaction.
                            The address specified must be an input to
                            the transaction, and it must have a balance
                            able to cover the additional fee. Also, the
                            inputs must exactly balance the outputs, 
                            since the logic to understand what to do
                            otherwise is quite complicated, and prone
                            to odd behavior.
    
    sign key                Sign the transaction specified by the key
    
    submit key              Submit the transaction specified by the key
                            to Factom
    
    balances or
    getaddresses            Returns the list of addresses known to the
                            wallet. Returns the name that can be used
                            tied to each address, as well as the base 58
                            address (which is the actual address).  This
                            command also returns the balances at each 
                            address.
   
    transactions            Prints information about pending transactions.
                            Returns a list of all the transactions being
                            constructed by the user.  It shows the fee
                            required (at this point) as well as the fee 
                            the user will pay.  Some additional error 
                            checking is done as well, with messages 
                            provided to the user.

    list                    List confirmed transactions' details. 
         [transaction id]   Lists the confirmed transactions with the given 
                            transaction id.
         [address]          Dumps all Factoid transactions that use the 
                            given address as an input or an output.
         all                Dumps all Factoid transactions to date.

    properties              Returns information about factomd, fctwallet,
                            the Protocol version, the version of this CLI,
                            and more.
`

// man returns an usage error string for the specified sub command.
func man(s string) {
	m := map[string]string{
		"testcredit":     "factom-cli testcredit [key]",
		"balance":        "factom-cli balance ec|fct [key]",
		"buy":            "factom-cli buy #amt",
		"get":            "factom-cli get head|dblock|height|chain|eblock|entry",
		"getHead":        "factom-cli get head",
		"getDBlock":      "factom-cli get dblock [keymr]",
		"getChain":       "factom-cli get chain [chainid]",
		"getEBlock":      "factom-cli get eblock [keymr]",
		"getEntry":       "factom-cli get entry [hash]",
		"getFirstEntry":  "factom-cli get firstentry [chainid]",
		"getaddresses":   "factom-cli getaddresses|balances",
		"balances":       "factom-cli getaddresses|balances",
		"transactions":   "factom-cli transactions",
		"help":           usage,
		"mkchain":        "factom-cli mkchain [-e extid ...] name <stdin>",
		"genfactoidaddr": "factom-cli genfactoidaddr name",
		"newtransaction": "factom-cli newtransaction key",
		"addinput":       "factom-cli addinput key name|address amount",
		"addoutput":      "factom-cli addoutput key name|address amount",
		"addecoutput":    "factom-cli addecoutput key name|address amount",
		"validate":       "factom-cli validate key",
		"submit":         "factom-cli submit key",
		"properties":     "factom-cli properties",
		"put":            "factom-cli put [-e extid ...] name <stdin>",
		"default":        "More Help can be found by typing:\n\n  factom-cli help",
	}

	if m[s] != "" {
		errorln(m[s])
		return
	}
	errorln(m["default"])
}
