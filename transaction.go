// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/FactomProject/factom"
)

// newtx creates a new transaction in the wallet.
var newtx = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli newtx TXNAME"
	cmd.description = "Create a new transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) != 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if err := factom.NewTransaction(args[0]); err != nil {
			errorln(err)
			return
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

func FixPointPrt(value uint64) string {
	whole := value/100000000
	part  := value - (whole*100000000)
	ret := []byte(fmt.Sprintf("%d.%08d",whole,part))
	for string(ret[len(ret)-1])=="0" {
		ret = ret[:len(ret)-1]
	}
	if string(ret[len(ret)-1])=="." {
		ret = ret[:len(ret)-1]
	}
		return string(ret)
}

// listtxs lists the working transactions in the wallet.
var listtxs = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listtxs"
	cmd.description = "List current working transactions in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		txs, err := factom.ListTransactions()
		if err != nil {
			errorln(err)
			return
		}
		for _, tx := range txs {
			fmt.Println("{")
			fmt.Println("	Name:", tx.Name)
			fmt.Println("	TxID:", tx.TxID)
			fmt.Println("	TotalInputs:", FixPointPrt(tx.TotalInputs))
			fmt.Println("	TotalOutputs:", FixPointPrt(tx.TotalOutputs))
			fmt.Println("	TotalECOutputs:", FixPointPrt(tx.TotalECOutputs))
			if tx.TotalInputs <= (tx.TotalOutputs + tx.TotalECOutputs) {
				fmt.Println("	FeesPaid:", 0)
				fmt.Println("	FeesRequired:", FixPointPrt(tx.FeesRequired))
			} else {
				feesPaid := tx.TotalInputs - (tx.TotalOutputs + tx.TotalECOutputs)
				fmt.Println("	FeesPaid:", FixPointPrt(feesPaid))
			}
			fmt.Println("	RawTransaction:", tx.RawTransaction)
			fmt.Println("}")
		}
	}
	help.Add("listtxs", cmd)
	return cmd
}()

// addtxinput adds a factoid input to a transaction in the wallet.
var addtxinput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addtxinput TXNAME ADDRESS AMOUNT"
	cmd.description = "Add a Factoid input to a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
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
		if err := factom.AddTransactionInput(args[0], args[1], amt); err != nil {
			errorln(err)
			return
		}
	}
	help.Add("addtxinput", cmd)
	return cmd
}()

// addtxoutput adds a factoid output to a transaction in the wallet.
var addtxoutput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addtxoutput [-r] TXNAME ADDRESS AMOUNT"
	cmd.description = "Add a Factoid output to a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		var res = flag.Bool("r", false, "resolve dns address")
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
		if err := factom.AddTransactionOutput(args[0], out, amt); err != nil {
			errorln(err)
			return
		}
	}
	help.Add("addtxoutput", cmd)
	return cmd
}()

// addtxecoutput adds an entry credit output to a transaction in the wallet.
var addtxecoutput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addtxecoutput [-r] TXNAME ADDRESS AMOUNT"
	cmd.description = "Add an Entry Credit output to a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		var res = flag.Bool("r", false, "resolve dns address")
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
		if err := factom.AddTransactionECOutput(args[0], out, amt); err != nil {
			errorln(err)
			return
		}
	}
	help.Add("addtxecoutput", cmd)
	return cmd
}()

// addtxfee adds an entry credit output to a transaction in the wallet.
var addtxfee = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addtxfee TXNAME ADDRESS"
	cmd.description = "Add the transaction fee to an input of a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) != 2 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if err := factom.AddTransactionFee(args[0], args[1]); err != nil {
			errorln(err)
			return
		}
	}
	help.Add("addtxfee", cmd)
	return cmd
}()

// subtxfee adds an entry credit output to a transaction in the wallet.
var subtxfee = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli subtxfee TXNAME ADDRESS"
	cmd.description = "Subtract the transaction fee from an output of a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) != 2 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if err := factom.SubTransactionFee(args[0], args[1]); err != nil {
			errorln(err)
			return
		}
	}
	help.Add("subtxfee", cmd)
	return cmd
}()

// signtx signs a transaction in the wallet.
var signtx = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli signtx TXNAME"
	cmd.description = "Sign a transaction in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) != 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if err := factom.SignTransaction(args[0]); err != nil {
			errorln(err)
			return
		}
	}
	help.Add("signtx", cmd)
	return cmd
}()

// composetx composes the signed json rpc object to make a transaction against factomd
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
	cmd.helpMsg = "factom-cli sendfct FROMADDRESS TOADDRESS AMOUNT"
	cmd.description = "Send Factoids between 2 addresses"
	cmd.execFunc = func(args []string) {
		os.Args = args
		var res = flag.Bool("r", false, "resolve dns address")
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
		fmt.Println("TxID:", t)
	}
	help.Add("sendfct", cmd)
	return cmd
}()

// buyec sends factoids between 2 addresses
var buyec = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli buyec FCTADDRESS ECADDRESS ECAMOUNT"
	cmd.description = "Buy entry credits"
	cmd.execFunc = func(args []string) {
		os.Args = args
		var res = flag.Bool("r", false, "resolve dns address")
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
		fmt.Println("TxID:", t)
	}
	help.Add("buyec", cmd)
	return cmd
}()
