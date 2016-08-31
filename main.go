// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/util"
)

const Version = "0.2.0.0"

func main() {
	var (
		hflag = flag.Bool("h", false, "help")
		sflag = flag.String("s", "localhost:8088", "address of api server")
		wflag = flag.String("w", "localhost:8089", "address of wallet api server")
		rpcUserflag     = flag.String("rpcuser", "", "Username for JSON-RPC connections")
		rpcPasswordflag = flag.String("rpcpassword", "", "Password for JSON-RPC connections")		
	)
	flag.Parse()

	filename := util.ConfigFilename()
	cfg := util.ReadConfig(filename).Rpc
	if *rpcUserflag == "" {
		*rpcUserflag = cfg.RpcUser
	}
	if *rpcPasswordflag == "" {
		*rpcPasswordflag = cfg.RpcPass
	}
	if *rpcUserflag == "" || *rpcPasswordflag == "" {
		log.Fatal("Rpc user and password did not set, using -rpcuser and -rpcpassword or config file")
	}
	
	args := flag.Args()

	if *hflag {
		args = []string{"help"}
	}
	factom.SetFactomdServer(*sflag)
	factom.SetWalletServer(*wflag)
	factom.SetRpcConfig(*rpcUserflag, *rpcPasswordflag)
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
