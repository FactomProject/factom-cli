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

const Version = "0.2.0.0"

func main() {
	var (
		hflag              = flag.Bool("h", false, "help")
		sflag              = flag.String("s", "localhost:8088", "address and port of factomd api")
		wflag              = flag.String("w", "localhost:8089", "address and port of factom-walletd api")
		walletRpcUser      = flag.String("walletuser", "", "Username for API connections to factom-walletd")
		walletRpcPassword  = flag.String("walletpassword", "", "Password for API connections to factom-walletd")
		factomdRpcUser     = flag.String("factomduser", "", "Username for API connections to factomd")
		factomdRpcPassword = flag.String("factomdpassword", "", "Password for API connections to factomd")
	)
	flag.Parse()

	//see if the config file has values which should be used instead of null strings
	filename := util.ConfigFilename() //file name and path to factomd.conf file
	cfg := util.ReadConfig(filename).Rpc
	cfgw := util.ReadConfig(filename).Walletd

	if *walletRpcUser == "" {
		if cfgw.WalletRpcUser != "" {
			fmt.Printf("using factom-walletd API user and password specified in \"%s\" at WalletRpcUser & WalletRpcPass\n", filename)
			*walletRpcUser = cfgw.WalletRpcUser
			*walletRpcPassword = cfgw.WalletRpcPass
		}
	}

	if *factomdRpcUser == "" {
		if cfg.FactomdRpcUser != "" {
			fmt.Printf("using factomd API user and password specified in \"%s\" at FactomdRpcUser & FactomdRpcPass\n", filename)
			*factomdRpcUser = cfg.FactomdRpcUser
			*factomdRpcPassword = cfg.FactomdRpcPass
		}
	}

	args := flag.Args()

	if *hflag {
		args = []string{"help"}
	}
	factom.SetFactomdServer(*sflag)
	factom.SetWalletServer(*wflag)
	factom.SetFactomdRpcConfig(*factomdRpcUser, *factomdRpcPassword)
	factom.SetWalletRpcConfig(*walletRpcUser, *walletRpcPassword)
	c := cli.New()
	c.Handle("help", help)
	c.Handle("ack", ack)
	c.Handle("addchain", addchain)
	c.Handle("addentry", addentry)
	c.Handle("backupwallet", backupwallet)
	c.Handle("balance", balance)
	c.Handle("ecrate", ecrate)
	c.Handle("exportaddresses", exportaddresses)
	c.Handle("get", get)
	c.Handle("importaddress", importaddresses)
	c.Handle("importwords", importwords)
	c.Handle("listaddresses", listaddresses)
	c.Handle("newecaddress", newecaddress)
	c.Handle("newfctaddress", newfctaddress)
	c.Handle("properties", properties)
	c.Handle("receipt", receipt)
	c.Handle("backupwallet", backupwallet)

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
