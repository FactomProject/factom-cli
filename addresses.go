// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
)

// balance prints the current balance of the specified address
var balance = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli balance [-r] ADDRESS"
	cmd.description = "If this is an EC Address, returns number of Entry Credits. If this is a Factoid Address, returns the Factoid balance."
	cmd.execFunc = func(args []string) {
		os.Args = args
		var res = flag.Bool("r", false, "resolve dns address")
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		addr := args[0]

		switch factom.AddressStringType(addr) {
		case factom.FactoidPub:
			b, err := factom.GetFactoidBalance(addr)
			if err != nil {
				errorln(err)
			}
			fmt.Println(float64(b) / 1e8)
			return
		case factom.ECPub:
			c, err := factom.GetECBalance(addr)
			if err != nil {
				errorln(err)
			}
			fmt.Println(c)
			return
		}

		// if -r flag is present, resolve dns address then get the fct and ec
		// blance
		if *res {
			f, e, err := factom.GetDnsBalance(addr)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(addr, "fct", float64(f)/1e8)
			fmt.Println(addr, "ec", e)
		} else {
			fmt.Println("Undefined or invalid address")
		}
	}
	help.Add("balance", cmd)
	return cmd
}()

// ecrate shows the entry credit conversion rate in factoids
var ecrate = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli ecrate"
	cmd.description = "Show the current Entry Credit conversion rate in the Factom Network"
	cmd.execFunc = func(args []string) {
		rate, err := factom.GetRate()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(float64(rate) / 1e8)

	}
	help.Add("ecrate", cmd)
	return cmd
}()

// exportaddresses lists the private addresses from the wallet
var exportaddresses = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli exportaddresses"
	cmd.description = "List the private addresses stored in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		fs, es, err := factom.FetchAddresses()
		if err != nil {
			errorln(err)
			return
		}
		for _, a := range fs {
			fmt.Println(a.SecString(), a.String())
		}
		for _, a := range es {
			fmt.Println(a.SecString(), a.String())
		}
	}
	help.Add("exportaddresses", cmd)
	return cmd
}()

// importaddresses imports addresses from 1 or more secret keys into the wallet
var importaddresses = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli importaddresses ADDRESS [ADDRESS...]"
	cmd.description = "Import one or more secret keys into the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		fs, es, err := factom.ImportAddresses(args...)
		if err != nil {
			errorln(err)
			return
		}
		for _, a := range fs {
			fmt.Println(a)
		}
		for _, a := range es {
			fmt.Println(a)
		}
	}
	help.Add("importaddresses", cmd)
	help.Add("importaddress", cmd)
	return cmd
}()

// newecaddress generates a new ec address in the wallet
var newecaddress = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli newecaddress"
	cmd.description = "Generate a new Entry Credit Address in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		a, err := factom.GenerateECAddress()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(a)
	}
	help.Add("newecaddress", cmd)
	return cmd
}()

// newfctaddress generates a new ec address in the wallet
var newfctaddress = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli newfctaddress"
	cmd.description = "Generate a new Factoid Address in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		a, err := factom.GenerateFactoidAddress()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(a)
	}
	help.Add("newfctaddress", cmd)
	return cmd
}()

// listaddresses lists the addresses in the wallet
var listaddresses = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listaddresses"
	cmd.description = "List the addresses stored in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		fs, es, err := factom.FetchAddresses()
		if err != nil {
			errorln(err)
			return
		}
		for _, a := range fs {
			b, err := factom.GetFactoidBalance(a.String())
			if err != nil {
				errorln(err)
			}
			fmt.Println(a, float64(b)/1e8)
		}
		for _, a := range es {
			c, err := factom.GetECBalance(a.String())
			if err != nil {
				errorln(err)
			}
			fmt.Println(a, c)
		}
	}
	help.Add("listaddresses", cmd)
	return cmd
}()
