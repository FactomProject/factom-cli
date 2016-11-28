// Copyright 2016 Factom Foundation
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
const Version = "0.2.0.0"

func main() {
	var (
		hflag              = flag.Bool("h", false, "help")
		walletRpcUser      = flag.String("walletuser", "", "Username for API connections to factom-walletd")
		walletRpcPassword  = flag.String("walletpassword", "", "Password for API connections to factom-walletd")
		factomdRpcUser     = flag.String("factomduser", "", "Username for API connections to factomd")
		factomdRpcPassword = flag.String("factomdpassword", "", "Password for API connections to factomd")

		factomdLocation = flag.String("s", "", "IPAddr:port# of factomd API to use to access blockchain (default localhost:8088)")
		walletdLocation = flag.String("w", "", "IPAddr:port# of factom-walletd API to use to create transactions (default localhost:8089)")

		walletTLSflag = flag.Bool("wallettls", false, "Set to true when the wallet API is encrypted")
		walletTLSCert = flag.String("walletcert", "", "This file is the TLS certificate provided by the factom-walletd API. (default ~/.factom/walletAPIpub.cert)")

		factomdTLSflag = flag.Bool("factomdtls", false, "Set to true when the factomd API is encrypted")
		factomdTLSCert = flag.String("factomdcert", "", "This file is the TLS certificate provided by the factomd API. (default ~/.factom/m2/factomdAPIpub.cert)")
	)
	flag.Parse()

	//see if the config file has values which should be used instead of null strings
	filename := util.ConfigFilename() //file name and path to factomd.conf file
	//if the config file doesn't exist, it gives lots of warnings when util.ReadConfig is called.
	//instead of giving warnings, check that the file exists before attempting to read it.
	//if it doesn't exist, silently ignore the file
	if _, err := os.Stat(filename); err == nil {
		cfg := util.ReadConfig(filename)

		if *walletRpcUser == "" {
			if cfg.Walletd.WalletRpcUser != "" {
				//fmt.Printf("using factom-walletd API user and password specified in \"%s\" at WalletRpcUser & WalletRpcPass\n", filename)
				*walletRpcUser = cfg.Walletd.WalletRpcUser
				*walletRpcPassword = cfg.Walletd.WalletRpcPass
			}
		}

		if *factomdRpcUser == "" {
			if cfg.App.FactomdRpcUser != "" {
				//fmt.Printf("using factomd API user and password specified in \"%s\" at FactomdRpcUser & FactomdRpcPass\n", filename)
				*factomdRpcUser = cfg.App.FactomdRpcUser
				*factomdRpcPassword = cfg.App.FactomdRpcPass
			}
		}

		if *factomdLocation == "" {
			if cfg.Walletd.FactomdLocation != "localhost:8088" {
				//fmt.Printf("using factomd location specified in \"%s\" as FactomdLocation = \"%s\"\n", filename, cfg.Walletd.FactomdLocation)
				*factomdLocation = cfg.Walletd.FactomdLocation
			}
		}

		if *walletdLocation == "" {
			if cfg.Walletd.WalletdLocation != "localhost:8089" {
				//fmt.Printf("using factom-walletd location specified in \"%s\" as WalletdLocation = \"%s\"\n", filename, cfg.Walletd.WalletdLocation)
				*walletdLocation = cfg.Walletd.WalletdLocation
			}
		}

		if cfg.Walletd.WalletTlsEnabled == true { //if a config file is found, and the wallet will start with TLS, factom-cli should use TLS too
			*walletTLSflag = true
		}

		if *walletTLSCert == "" { //if specified on the command line, don't use the config file
			if cfg.Walletd.WalletTlsPublicCert != "/full/path/to/walletAPIpub.cert" { //otherwise check if the the config file has something new
				//fmt.Printf("using wallet TLS certificate file specified in \"%s\" at WalletTlsPublicCert = \"%s\"\n", filename, cfg.Walletd.WalletTlsPublicCert)
				*walletTLSCert = cfg.Walletd.WalletTlsPublicCert
			}
		}

		if cfg.App.FactomdTlsEnabled == true { //if a config file is found, and the factomd will start with TLS, factom-cli should use TLS too
			*factomdTLSflag = true
		}

		if *factomdTLSCert == "" { //if specified on the command line, don't use the config file
			if cfg.App.FactomdTlsPublicCert != "/full/path/to/factomdAPIpub.cert" { //otherwise check if the the config file has something new
				//fmt.Printf("using wallet TLS certificate file specified in \"%s\" at FactomdTlsPublicCert = \"%s\"\n", filename, cfg.App.FactomdTlsPublicCert)
				*factomdTLSCert = cfg.App.FactomdTlsPublicCert
			}
		}
	}

	if *walletTLSCert == "" { //if all defaults were specified on both the command line and config file
		*walletTLSCert = fmt.Sprint(util.GetHomeDir(), "/.factom/walletAPIpub.cert")
		//fmt.Printf("using default wallet TLS certificate file \"%s\"\n", *walletTLSCert)
	}
	if *factomdTLSCert == "" { //if all defaults were specified on both the command line and config file
		*factomdTLSCert = fmt.Sprint(util.GetHomeDir(), "/.factom/m2/factomdAPIpub.cert")
		//fmt.Printf("using default factomd TLS certificate file \"%s\"\n", *factomdTLSCert)
	}

	if *factomdLocation == "" { //set the default if a config file doesn't exist
		*factomdLocation = "localhost:8088"
	}
	if *walletdLocation == "" { //set the default if a config file doesn't exist
		*walletdLocation = "localhost:8089"
	}

	args := flag.Args()

	if *hflag {
		args = []string{"help"}
	}
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

	c.HandleDefault(help)
	c.Execute(args)
}

func errorln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
