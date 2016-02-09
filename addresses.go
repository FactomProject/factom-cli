// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
)

// Generate a new Address
var generateaddress = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli generateaddress fct|ec name"
	cmd.description = "Generate and name a new factoid or ec address"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		c := cli.New()
		c.Handle("ec", ecGenerateAddr)
		c.Handle("fct", fctGenerateAddr)
		c.HandleDefaultFunc(func(args []string) {
			fmt.Println(cmd.helpMsg)
		})
		c.Execute(args)
	}
	help.Add("generateaddress", cmd)
	help.Add("newaddress", cmd)
	return cmd
}()

// Generate a new Entry Credit Address
var ecGenerateAddr = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli generateaddress ec name"
	cmd.description = "Generate and name a new ec address"
	cmd.execFunc = func(args []string) {
		if addr, err := factom.GenerateEntryCreditAddress(args[1]); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(addr)
		}
	}
	help.Add("generateaddress ec", cmd)
	help.Add("newaddress ec", cmd)
	return cmd

}()

// Generate a new Factoid Address
var fctGenerateAddr = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli generateaddress fct name"
	cmd.description = "Generate and name a new factoid address"
	cmd.execFunc = func(args []string) {
		if addr, err := factom.GenerateFactoidAddress(args[1]); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(addr)
		}
	}
	help.Add("generateaddress fct", cmd)
	help.Add("newaddress fct", cmd)
	return cmd
}()

var getaddresses = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli getaddresses|balances"
	cmd.description = "Returns the list of addresses known to the wallet. Returns the name that can be used tied to each address, as well as the base 58 address (which is the actual address). This command also returns the balances at each address."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) > 0 {
			fmt.Println(cmd.helpMsg)
		}

		str := fmt.Sprintf("http://%s/v1/factoid-get-addresses/", serverFct)
		getCmd(str, "Error printing addresses")
	}
	help.Add("getaddress", cmd)
	help.Add("balances", cmd)
	return cmd
}()

// importaddr imports a Factoid or Entry Credit private key and adds the
// address to the wallet database.
var importaddr = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli importaddress name EsKey|FsKey|'12Words'"
	cmd.description = "Import an Entry Credit or Factoid Private Key"
	cmd.execFunc = func(args []string) {
		if len(args) < 3 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if strings.HasPrefix(args[2], "Fs") {
			if addr, err := factom.GenerateFactoidAddressFromHumanReadablePrivateKey(args[1], args[2]); err != nil {
				fmt.Println(err)
				fmt.Println(cmd.helpMsg)
			} else {
				fmt.Println(args[1], addr)
			}
		} else if strings.HasPrefix(args[2], "Es") {
			if addr, err := factom.GenerateEntryCreditAddressFromHumanReadablePrivateKey(args[1], args[2]); err != nil {
				fmt.Println(err)
				fmt.Println(cmd.helpMsg)
			} else {
				fmt.Println(args[1], addr)
			}
		} else {
			if addr, err := factom.GenerateFactoidAddressFromMnemonic(args[1], args[2]); err != nil {
				fmt.Println(err)
				fmt.Println(cmd.helpMsg)
			} else {
				fmt.Println(args[1], addr)
			}
		}
	}
	help.Add("importaddress", cmd)
	return cmd
}()
