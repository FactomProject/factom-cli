// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/util"
)

// Version of factom-cli
// FactomcliVersion sets the semantic version number of the build
// $ go install -ldflags "-X main.FactomcliVersion=`cat VERSION`" -v
// It also seems to need to have the previous binary deleted if recompiling to have this message show up if no code has changed.

var FactomcliVersion string = "BuiltWithoutVersion"

func main() {
	var (
		walletRpcUser = flag.String(
			"walletuser",
			"",
			"Username for API connections to factom-walletd",
		)
		walletRpcPassword = flag.String(
			"walletpassword",
			"",
			"Password for API connections to factom-walletd",
		)
		factomdRpcUser = flag.String(
			"factomduser",
			"",
			"Username for API connections to factomd",
		)
		factomdRpcPassword = flag.String(
			"factomdpassword",
			"",
			"Password for API connections to factomd",
		)
		factomdLocation = flag.String(
			"s",
			"",
			"IPAddr:port# of factomd API to use to access blockchain (default"+
				" localhost:8088)",
		)
		walletdLocation = flag.String(
			"w",
			"",
			"IPAddr:port# of factom-walletd API to use to create transactions"+
				" (default localhost:8089)",
		)
		walletTLSflag = flag.Bool(
			"wallettls",
			false,
			"Set to true when the wallet API is encrypted",
		)
		walletTLSCert = flag.String(
			"walletcert",
			"",
			"This file is the TLS certificate provided by the factom-walletd"+
				" API. (default ~/.factom/walletAPIpub.cert)",
		)
		factomdTLSflag = flag.Bool(
			"factomdtls",
			false,
			"Set to true when the factomd API is encrypted",
		)
		factomdTLSCert = flag.String(
			"factomdcert",
			"",
			"This file is the TLS certificate provided by the factomd API."+
				" (default ~/.factom/m2/factomdAPIpub.cert)",
		)
	)
	flag.Parse()

	// see if the config file has values which should be used instead of null
	// strings
	filename := util.ConfigFilename()

	//instead of giving warnings, check that the file exists before attempting
	// to read it.
	if _, err := os.Stat(filename); err == nil {
		cfg := util.ReadConfig(filename)

		if *walletRpcUser == "" {
			if cfg.Walletd.WalletRpcUser != "" {
				*walletRpcUser = cfg.Walletd.WalletRpcUser
				*walletRpcPassword = cfg.Walletd.WalletRpcPass
			}
		}

		if *factomdRpcUser == "" {
			if cfg.App.FactomdRpcUser != "" {
				*factomdRpcUser = cfg.App.FactomdRpcUser
				*factomdRpcPassword = cfg.App.FactomdRpcPass
			}
		}

		if *factomdLocation == "" {
			if cfg.Walletd.FactomdLocation != "localhost:8088" {
				*factomdLocation = cfg.Walletd.FactomdLocation
			}
		}

		if *walletdLocation == "" {
			if cfg.Walletd.WalletdLocation != "localhost:8089" {
				*walletdLocation = cfg.Walletd.WalletdLocation
			}
		}

		//if a config file is found, and the wallet will start with TLS,
		// factom-cli should use TLS too
		if cfg.Walletd.WalletTlsEnabled == true {
			*walletTLSflag = true
		}

		// if specified on the command line, don't use the config file
		if *walletTLSCert == "" {
			//otherwise check if the the config file has something new
			if cfg.Walletd.WalletTlsPublicCert != "/full/path/to/walletAPIpub.cert" {
				*walletTLSCert = cfg.Walletd.WalletTlsPublicCert
			}
		}

		//if a config file is found, and the factomd will start with TLS,
		// factom-cli should use TLS too
		if cfg.App.FactomdTlsEnabled == true {
			*factomdTLSflag = true
		}

		//if specified on the command line, don't use the config file
		if *factomdTLSCert == "" {
			//otherwise check if the the config file has something new
			if cfg.App.FactomdTlsPublicCert != "/full/path/to/factomdAPIpub.cert" {
				*factomdTLSCert = cfg.App.FactomdTlsPublicCert
			}
		}
	}

	//if all defaults were specified on both the command line and config file
	if *walletTLSCert == "" {
		*walletTLSCert = fmt.Sprint(
			util.GetHomeDir(),
			"/.factom/walletAPIpub.cert",
		)
	}

	//if all defaults were specified on both the command line and config file
	if *factomdTLSCert == "" {
		*factomdTLSCert = fmt.Sprint(
			util.GetHomeDir(),
			"/.factom/m2/factomdAPIpub.cert",
		)
	}

	//set the default if a config file doesn't exist
	if *factomdLocation == "" {
		*factomdLocation = "localhost:8088"
	}
	if *walletdLocation == "" {
		*walletdLocation = "localhost:8089"
	}

	args := flag.Args()

	factom.SetFactomdServer(*factomdLocation)
	factom.SetWalletServer(*walletdLocation)
	factom.SetFactomdRpcConfig(*factomdRpcUser, *factomdRpcPassword)
	factom.SetWalletRpcConfig(*walletRpcUser, *walletRpcPassword)
	factom.SetWalletEncryption(*walletTLSflag, *walletTLSCert)
	factom.SetFactomdEncryption(*factomdTLSflag, *factomdTLSCert)
	c := cli.New()
	c.Handle("help", help)
	c.Handle("status", status)
	c.Handle("addchain", addchain)
	c.Handle("addentry", addentry)
	c.Handle("backupwallet", backupwallet)
	c.Handle("balance", balance)
	//c.Handle("balancetotals", balancetotals)
	c.Handle("composechain", composechain)
	c.Handle("composeentry", composeentry)
	c.Handle("ecrate", ecrate)
	c.Handle("exportaddresses", exportaddresses)
	c.Handle("get", get)
	c.Handle("importaddress", importaddresses)
	c.Handle("importkoinify", importkoinify)
	c.Handle("listaddresses", listaddresses)
	c.Handle("newecaddress", newecaddress)
	c.Handle("newfctaddress", newfctaddress)
	c.Handle("properties", properties)
	c.Handle("receipt", receipt)
	c.Handle("backupwallet", backupwallet)
	c.Handle("rmaddress", removeAddress)

	// transaction commands
	c.Handle("newtx", newtx)
	c.Handle("rmtx", rmtx)
	c.Handle("listtxs", listtxs)
	c.Handle("addtxinput", addtxinput)
	c.Handle("addtxoutput", addtxoutput)
	c.Handle("addtxecoutput", addtxecoutput)
	c.Handle("addtxfee", addtxfee)
	c.Handle("subtxfee", subtxfee)
	c.Handle("signtx", signtx)
	c.Handle("composetx", composetx)
	c.Handle("sendtx", sendtx)
	c.Handle("sendfct", sendfct)
	c.Handle("buyec", buyec)

	// identity commands
	c.Handle("newidentitykey", newIdentityKey)
	c.Handle("importidentitykeys", importIdentityKeys)
	c.Handle("exportidentitykeys", exportIdentityKeys)
	c.Handle("listidentitykeys", listIdentityKeys)
	c.Handle("rmidentitykey", removeIdentityKey)
	c.Handle("identity", identity)

	c.HandleDefault(help)
	c.Execute(args)
}

func errorln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
