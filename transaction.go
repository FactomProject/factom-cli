// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
)

// newtx creates a new transaction in the wallet.
var newtx = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli newtx [-q] TXNAME"
	cmd.description = "Create a new transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		qflag := flag.Bool("q", false, "quiet mode; no output")
		flag.Parse()
		args = flag.Args()

		if len(args) != 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		tx, err := factom.NewTransaction(args[0])
		if err != nil {
			errorln(err)
			return
		}

		// output
		switch {
		// quiet mode; don't print anything
		case *qflag:
		default:
			fmt.Println(tx)
		}
	}
	help.Add("newtx", cmd)
	return cmd
}()

// rmtx removes a transaction in the wallet.
var rmtx = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli rmtx TXNAME"
	cmd.description = "Remove a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) != 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if err := factom.DeleteTransaction(args[0]); err != nil {
			errorln(err)
			return
		}
	}
	help.Add("rmtx", cmd)
	return cmd
}()

// addtxinput adds a factoid input to a transaction in the wallet.
var addtxinput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addtxinput [-q] TXNAME ADDRESS AMOUNT"
	cmd.description = "Add a Factoid input to a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		qflag := flag.Bool("q", false, "quiet mode; no output")
		flag.Parse()
		args = flag.Args()
		if len(args) != 3 {
			fmt.Println(cmd.helpMsg)
			return
		}

		var amt uint64
		if i, err := strconv.ParseFloat(args[2], 64); err != nil {
			errorln(err)
		} else if i < 0 {
			errorln("AMOUNT may not be less than 0")
		} else {
			amt = uint64(i * 1e8)
		}

		tx, err := factom.AddTransactionInput(args[0], args[1], amt)
		if err != nil {
			errorln(err)
			return
		}

		// output
		switch {
		// quiet mode; don't print anything
		case *qflag:
		default:
			fmt.Println(tx)
		}
	}
	help.Add("addtxinput", cmd)
	return cmd
}()

// addtxoutput adds a factoid output to a transaction in the wallet.
var addtxoutput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addtxoutput [-rq] TXNAME ADDRESS AMOUNT"
	cmd.description = "Add a Factoid output to a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		res := flag.Bool("r", false, "resolve dns address")
		qflag := flag.Bool("q", false, "quiet mode; no output")
		flag.Parse()
		args = flag.Args()
		if len(args) != 3 {
			fmt.Println(cmd.helpMsg)
			return
		}

		var amt uint64
		if i, err := strconv.ParseFloat(args[2], 64); err != nil {
			errorln(err)
		} else if i < 0 {
			errorln("AMOUNT may not be less than 0")
		} else {
			amt = uint64(i * 1e8)
		}

		out := args[1]
		if *res {
			if f, _, err := factom.ResolveDnsName(args[1]); err != nil {
				errorln(err)
				return
			} else if f == "" {
				errorln("could not resolve factoid address")
			} else {
				out = f
			}
		}

		tx, err := factom.AddTransactionOutput(args[0], out, amt)
		if err != nil {
			errorln(err)
			return
		}

		// output
		switch {
		// quiet mode; don't print anything
		case *qflag:
		default:
			fmt.Println(tx)
		}
	}
	help.Add("addtxoutput", cmd)
	return cmd
}()

// addtxecoutput adds an entry credit output to a transaction in the wallet.
var addtxecoutput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addtxecoutput [-r] TXNAME ADDRESS AMOUNT"
	cmd.description = "Add an Entry Credit output to a transaction in the" +
		" wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		res := flag.Bool("r", false, "resolve dns address")
		qflag := flag.Bool("q", false, "quiet mode; no output")
		flag.Parse()
		args = flag.Args()
		if len(args) != 3 {
			fmt.Println(cmd.helpMsg)
			return
		}

		var amt uint64
		if i, err := strconv.ParseFloat(args[2], 64); err != nil {
			errorln(err)
		} else if i < 0 {
			errorln("AMOUNT may not be less than 0")
		} else {
			amt = uint64(i * 1e8)
		}

		out := args[1]
		if *res {
			if _, e, err := factom.ResolveDnsName(args[1]); err != nil {
				errorln(err)
				return
			} else if e == "" {
				errorln("could not resolve entry credit address")
			} else {
				out = e
			}
		}

		tx, err := factom.AddTransactionECOutput(args[0], out, amt)
		if err != nil {
			errorln(err)
			return
		}

		// output
		switch {
		// quiet mode; don't print anything
		case *qflag:
		default:
			fmt.Println(tx)
		}
	}
	help.Add("addtxecoutput", cmd)
	return cmd
}()

// addtxfee adds an entry credit output to a transaction in the wallet.
var addtxfee = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addtxfee [-q] TXNAME ADDRESS"
	cmd.description = "Add the transaction fee to an input of a transaction" +
		" in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		qflag := flag.Bool("q", false, "quiet mode; no output")
		flag.Parse()
		args = flag.Args()

		if len(args) != 2 {
			fmt.Println(cmd.helpMsg)
			return
		}

		tx, err := factom.AddTransactionFee(args[0], args[1])
		if err != nil {
			errorln(err)
			return
		}

		// output
		switch {
		// quiet mode; don't print anything
		case *qflag:
		default:
			fmt.Println(tx)
		}
	}
	help.Add("addtxfee", cmd)
	return cmd
}()

// listtxs lists transactions from the wallet or the Factoid Chain.
var listtxs = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listtxs [address|all|id|name|tmp|range]"
	cmd.description = "List transactions from the wallet or the Factoid Chain"
	cmd.execFunc = func(args []string) {
		if len(args) > 1 {
			args = args[1:]
		}
				
		c := cli.New()
		c.Handle("all", listtxsall)
		c.Handle("address", listtxsaddress)
		c.Handle("id", listtxsid)
		c.Handle("range", listtxsrange)
		c.Handle("tmp", listtxstmp)
		c.Handle("name", listtxsname)
		c.HandleDefaultFunc(func(args []string) {
			tmp := []string{"all"}
			args = append(tmp, args...)
			listtxsall.execFunc(args)
		})
		c.Execute(args)
	}
	help.Add("listtxs", cmd)
	return cmd
}()

// listtxsall lists all transactions from the Factoid Chain
var listtxsall = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listtxs [all] [-T]"
	cmd.description = "List all transactions from the Factoid Chain"
	cmd.execFunc = func(args []string) {
		os.Args = args
		tdisp := flag.Bool("T", false, "display only the TxID")
		flag.Parse()
		args = flag.Args()
		if len(args) > 0 {
			fmt.Println(cmd.helpMsg)
			return
		}

		txs, err := factom.ListTransactionsAll()
		if err != nil {
			errorln(err)
			return
		}
		for _, tx := range txs {
			switch {
			case *tdisp:
				fmt.Println(tx.TxID)
			default:
				fmt.Println(tx)
			}
		}
	}
	help.Add("listtxs all", cmd)
	return cmd
}()

// listtxsaddress lists transactions from the Factoid Chain with matching
// address
var listtxsaddress = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listtxs address [-T] ECADDRESS|FCTADDRESS"
	cmd.description = "List transaction from the Factoid Chain with a" +
		" specific address"
	cmd.execFunc = func(args []string) {
		os.Args = args
		tdisp := flag.Bool("T", false, "display only the TxID")
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		txs, err := factom.ListTransactionsAddress(args[0])
		if err != nil {
			errorln(err)
			return
		}
		for _, tx := range txs {
			switch {
			case *tdisp:
				fmt.Println(tx.TxID)
			default:
				fmt.Println(tx)
			}
		}
	}
	help.Add("listtxs address", cmd)
	return cmd
}()

// listtxsid lists transactions from the Factoid Chain with matching id
var listtxsid = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listtxs id TXID"
	cmd.description = "List transaction from the Factoid Chain"
	cmd.execFunc = func(args []string) {
		os.Args = args
		tdisp := flag.Bool("T", false, "display only the TxID")
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		txs, err := factom.ListTransactionsID(args[0])
		if err != nil {
			errorln(err)
			return
		}
		for _, tx := range txs {
			switch {
			case *tdisp:
				fmt.Println(tx.TxID)
			default:
				fmt.Println(tx)
			}
		}
	}
	help.Add("listtxs id", cmd)
	return cmd
}()

// listtxsname get a working transaction from the wallet.
var listtxsname = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listtxs name TXNAME"
	cmd.description = "Show a current working transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		tdisp := flag.Bool("T", false, "display transaction txid only")
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		name := args[0]

		txs, err := factom.ListTransactionsTmp()
		if err != nil {
			errorln(err)
			return
		}
		for _, tx := range txs {
			if tx.Name == name {
				if *tdisp {
					fmt.Println(tx.TxID)
				} else {
					fmt.Println(tx)
				}
			}
		}
	}
	help.Add("listtxs tmp", cmd)
	return cmd
}()

// listtxsrange lists the transactions from the Factoid Chain within the
// specified block height range
var listtxsrange = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listtxs range [-T] START END"
	cmd.description = "List the transactions from the Factoid Chain within" +
		" the specified range"
	cmd.execFunc = func(args []string) {
		os.Args = args
		tdisp := flag.Bool("T", false, "display only the TxID")
		flag.Parse()
		args = flag.Args()
		if len(args) < 2 {
			fmt.Println(cmd.helpMsg)
			return
		}

		start, err := strconv.Atoi(args[0])
		if err != nil {
			errorln(err)
			return
		}
		end, err := strconv.Atoi(args[1])
		if err != nil {
			errorln(err)
			return
		}

		txs, err := factom.ListTransactionsRange(start, end)
		if err != nil {
			errorln(err)
			return
		}
		for _, tx := range txs {
			switch {
			case *tdisp:
				fmt.Println(tx.TxID)
			default:
				fmt.Println(tx)
			}
		}
	}
	help.Add("listtxs range", cmd)
	return cmd
}()

// listtxstmp lists the working transactions in the wallet.
var listtxstmp = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listtxs tmp"
	cmd.description = "List current working transactions in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		ndisp := flag.Bool("N", false, "display transaction names only")
		flag.Parse()
		args = flag.Args()

		txs, err := factom.ListTransactionsTmp()
		if err != nil {
			errorln(err)
			return
		}
		for _, tx := range txs {
			if *ndisp {
				fmt.Println(tx.Name)
			} else {
				fmt.Println(tx)
			}
		}
	}
	help.Add("listtxs tmp", cmd)
	return cmd
}()

// subtxfee adds an entry credit output to a transaction in the wallet.
var subtxfee = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli subtxfee [-q] TXNAME ADDRESS"
	cmd.description = "Subtract the transaction fee from an output of a" +
		" transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		qflag := flag.Bool("q", false, "quiet mode; no output")
		flag.Parse()
		args = flag.Args()
		if len(args) != 2 {
			fmt.Println(cmd.helpMsg)
			return
		}

		tx, err := factom.SubTransactionFee(args[0], args[1])
		if err != nil {
			errorln(err)
			return
		}

		// output
		switch {
		// quiet mode; don't print anything
		case *qflag:
		default:
			fmt.Println(tx)
		}
	}
	help.Add("subtxfee", cmd)
	return cmd
}()

// signtx signs a transaction in the wallet
var signtx = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli signtx [-qT] TXNAME"
	cmd.description = "Sign a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		qflag := flag.Bool("q", false, "quiet mode; no output")
		tdisp := flag.Bool("T", false, "display only the TxID")
		flag.Parse()
		args = flag.Args()
		if len(args) != 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		tx, err := factom.SignTransaction(args[0])
		if err != nil {
			errorln(err)
			return
		}

		// output
		switch {
		// quiet mode; don't print anything
		case *qflag:
		case *tdisp:
			fmt.Println(tx.TxID)
		default:
			fmt.Println(tx)
		}
	}
	help.Add("signtx", cmd)
	return cmd
}()

// composetx composes the signed json rpc object to make a transaction against
// factomd
var composetx = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli composetx TXNAME"
	cmd.description = "Compose a wallet transaction into a json rpc object"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) != 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		p, err := factom.ComposeTransaction(args[0])
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(string(p))
	}
	help.Add("composetx", cmd)
	return cmd
}()

// sendtx composes and sends the signed transaction to factomd
var sendtx = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli sendtx TXNAME"
	cmd.description = "Send a Transaction to Factom"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) != 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		t, err := factom.SendTransaction(args[0])
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println("TxID:", t)
	}
	help.Add("sendtx", cmd)
	return cmd
}()

// sendfct sends factoids between 2 addresses
var sendfct = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli sendfct [-T] FROMADDRESS TOADDRESS AMOUNT"
	cmd.description = "Send Factoids between 2 addresses"
	cmd.execFunc = func(args []string) {
		os.Args = args
		res := flag.Bool("r", false, "resolve dns address")
		tdisp := flag.Bool("T", false, "display only the TxID")
		flag.Parse()
		args = flag.Args()
		if len(args) != 3 {
			fmt.Println(cmd.helpMsg)
			return
		}

		tofc := args[1]

		// if -r flag is present resolve the ec address from the dns name.
		if *res {
			f, _, err := factom.ResolveDnsName(tofc)
			if err != nil {
				errorln(err)
				return
			}
			tofc = f
		}

		var amt uint64
		if i, err := strconv.ParseFloat(args[2], 64); err != nil {
			errorln(err)
		} else if i < 0 {
			errorln("AMOUNT may not be less than 0")
		} else {
			amt = uint64(i * 1e8)
		}

		t, err := factom.SendFactoid(args[0], tofc, amt)
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *tdisp:
			fmt.Println(t)
		default:
			fmt.Println("TxID:", t)
		}
	}
	help.Add("sendfct", cmd)
	return cmd
}()

// buyec sends factoids between 2 addresses
var buyec = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli buyec [-T] FCTADDRESS ECADDRESS ECAMOUNT"
	cmd.description = "Buy entry credits"
	cmd.execFunc = func(args []string) {
		os.Args = args
		res := flag.Bool("r", false, "resolve dns address")
		tdisp := flag.Bool("T", false, "display only the TxID")
		flag.Parse()
		args = flag.Args()
		if len(args) != 3 {
			fmt.Println(cmd.helpMsg)
			return
		}

		toec := args[1]

		// if -r flag is present resolve the ec address from the dns name.
		if *res {
			_, e, err := factom.ResolveDnsName(toec)
			if err != nil {
				errorln(err)
				return
			}
			toec = e
		}

		var amt uint64
		if i, err := strconv.Atoi(args[2]); err != nil {
			errorln(err)
			return
		} else if i < 0 {
			errorln("AMOUNT may not be less than 0")
			return
		} else {
			rate, err := factom.GetRate()
			if err != nil {
				errorln(err)
			}
			amt = uint64(i) * rate
		}

		t, err := factom.BuyEC(args[0], toec, amt)
		if err != nil {
			errorln(err)
			return
		}
		switch {
		case *tdisp:
			fmt.Println(t)
		default:
			fmt.Println("TxID:", t)
		}
	}
	help.Add("buyec", cmd)
	return cmd
}()
