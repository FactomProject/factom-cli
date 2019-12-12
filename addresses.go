// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
	"github.com/posener/complete"
)

// balance prints the current balance of the specified address
var balance = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli balance [-r] ADDRESS"
	cmd.description = "If this is an EC Address, returns number of Entry Credits. If this is a Factoid Address, returns the Factoid balance."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-r": complete.PredictNothing,
		},
		Args: predictAddress,
	}
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
		case factom.EthFA:
			addr = addr[1:]
			fallthrough
		case factom.FactoidPub:
			b, err := factom.GetFactoidBalance(addr)
			if err != nil {
				errorln(err)
			} else {
				fmt.Println(factoshiToFactoid(b))
			}
			return
		case factom.ECPub:
			c, err := factom.GetECBalance(addr)
			if err != nil {
				errorln(err)
			} else {
				fmt.Println(c)
			}
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
			fmt.Println(addr, "fct", factoshiToFactoid(f))
			fmt.Println(addr, "ec", e)
		} else {
			fmt.Println("Undefined or invalid address")
		}
	}
	help.Add("balance", cmd)
	return cmd
}()

// balancetotals shows the total balance of all of the Factoid and Entry Credts
// in the wallet
var balancetotals = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli balancetotals [-FS -FA -ES -EA]"
	cmd.description = "This is the total number of Factoids and Entry Credits in the wallet"
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-FS": complete.PredictNothing,
			"-FA": complete.PredictNothing,
			"-ES": complete.PredictNothing,
			"-EA": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		var fsdisp = flag.Bool("FS", false, "Display only the savedFCT value")
		var fadisp = flag.Bool("FA", false, "Display only the ackFCT value")
		var esdisp = flag.Bool("ES", false, "Display only the savedEC value")
		var eadisp = flag.Bool("EA", false, "Display only the ackEC value")
		flag.Parse()
		args = flag.Args()

		fs, fa, es, ea, err := factom.GetBalanceTotals()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *fsdisp:
			fmt.Println(factoshiToFactoid(fs))
		case *fadisp:
			fmt.Println(factoshiToFactoid(fa))
		case *esdisp:
			fmt.Println(es)
		case *eadisp:
			fmt.Println(ea)
		default:
			fmt.Println("savedFCT:", factoshiToFactoid(fs))
			fmt.Println("ackFCT:", factoshiToFactoid(fa))
			fmt.Println("savedEC:", es)
			fmt.Println("ackEC:", ea)
		}
	}
	help.Add("balancetotals", cmd)
	return cmd
}()

// ecrate shows the entry credit conversion rate in factoids
var ecrate = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli ecrate"
	cmd.description = "It takes this many Factoids to buy an Entry Credit.  Displays the larger between current and future rates. Also used to set Factoid fees."
	cmd.execFunc = func(args []string) {
		rate, err := factom.GetECRate()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(factoshiToFactoid(rate))

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

		fs, es, eths, err := factom.FetchAllAddressTypes()
		if err != nil {
			errorln(err)
			return
		}
		for _, a := range fs {
			fmt.Println(a.SecString(), a.String())
		}
		for _, e := range eths {
			fmt.Println(e.SecString(), e.String())
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
	cmd.helpMsg = "factom-cli importaddress ADDRESS [ADDRESS...]"
	cmd.description = "Import one or more secret keys into the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		fs, es, eths, err := factom.ImportAddresses(args...)
		if err != nil {
			errorln(err)
			return
		}
		for _, a := range fs {
			fmt.Println(a)
		}
		for _, a := range eths {
			fmt.Println(a)
		}
		for _, a := range es {
			fmt.Println(a)
		}
	}
	help.Add("importaddress", cmd)
	return cmd
}()

var importkoinify = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli importkoinify '12WORDS'"
	cmd.description = "Import 12 words from Koinify sale into the Wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 || len(args) > 1 {
			fmt.Println(cmd.helpMsg, "  Note, 12 words must be in quotes")
			return
		}
		f, err := factom.ImportKoinify(args[0])
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(f)
	}
	help.Add("importkoinify", cmd)
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

// newethaddress generates a new eth address in the wallet
var newethaddress = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli newethaddress"
	cmd.description = "Generate a new Ether linked Address in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		a, err := factom.GenerateEtherAddress()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(a)
	}
	help.Add("newethaddress", cmd)
	return cmd
}()

// listaddresses lists the addresses in the wallet
var listaddresses = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listaddresses"
	cmd.description = "List the addresses stored in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		adisp := flag.Bool(
			"A",
			false,
			"display only the address without looking up the balance",
		)
		flag.Parse()
		args = flag.Args()

		fs, es, eths, err := factom.FetchAllAddressTypes()
		if err != nil {
			errorln(err)
			return
		}

		if *adisp {
			for _, a := range fs {
				fmt.Println(a)
			}
			for _, e := range eths {
				fmt.Println(e)
			}
			for _, a := range es {
				fmt.Println(a)
			}
		} else {
			for _, a := range fs {
				b, err := factom.GetFactoidBalance(a.String())
				if err != nil {
					errorln(err)
					fmt.Println(a)
				} else {
					fmt.Println(a, factoshiToFactoid(b))
				}
			}
			for _, e := range eths {
				b, err := factom.GetFactoidBalance(e.FAString())
				if err != nil {
					errorln(err)
					fmt.Println(e)
				} else {
					fmt.Println(e, factoshiToFactoid(b))
				}
			}
			for _, a := range es {
				c, err := factom.GetECBalance(a.String())
				if err != nil {
					errorln(err)
					fmt.Println(a)
				} else {
					fmt.Println(a, c)
				}
			}
		}
	}
	help.Add("listaddresses", cmd)
	return cmd
}()

// alternativeEncoding is a developer function. It returns all encoding types
// for ethereum addresses
var alternativeEncoding = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli alternativeEncoding"
	cmd.description = "This function is intended for developer use. There is" +
		"no use case for this function outside development purposes and regular" +
		"users should not need to use this. It returns all the encoding options" +
		"for addresses."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		fas, ecs, eths, err := factom.FetchAllAddressTypes()
		if err != nil {
			errorln(err)
			return
		}

		fmt.Println("Eth Linked Addresses")
		format := "%-53s %-53s %-53s\n"
		fmt.Printf(format, "FA", "Fe", "FE")
		for _, e := range eths {
			fmt.Printf(format, e.FAString(), e.String(), e.GateWayAddress())
		}

		fmt.Println()
		fmt.Println("Pure Factom Addresses")
		format = "%-53s\n"
		fmt.Printf(format, "FA/EC")
		for _, a := range fas {
			fmt.Printf(format, a.String())
		}

		for _, a := range ecs {
			fmt.Printf(format, a.String())
		}
	}
	help.Add("altencoding", cmd)
	return cmd
}()

//  linkedaddresses lists the linked addresses on other blockchains
var linkedaddresses = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli linkedaddresses"
	cmd.description = "List the addresses that have links to other addresses " +
		"on other blockchains. These addresses share the private key with another" +
		"blockchains."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		fas, ecs, eths, err := factom.FetchAllAddressTypes()
		if err != nil {
			errorln(err)
			return
		}

		format := "%-53s %-42s\n"
		fmt.Printf(format, "Factom", "Etheruem")
		for _, e := range eths {
			fmt.Printf(format, e.String(), e.EthAddress())
		}

		fmt.Printf(format, "Factom", "")
		for _, a := range fas {
			fmt.Printf(format, a.String(), "")
		}

		for _, a := range ecs {
			fmt.Printf(format, a.String(), "")
		}
	}
	help.Add("linkedaddresses", cmd)
	return cmd
}()

// Removes an address
var removeAddress = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli rmaddress ADDRESS"
	cmd.description = "Removes the public and private key from the wallet for the address specified."
	cmd.completion = complete.Command{
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		if len(args) < 2 {
			fmt.Println(cmd.helpMsg)
			return
		}
		addr := args[1]

		err := factom.RemoveAddress(addr)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
	help.Add("rmaddress", cmd)
	return cmd
}()
