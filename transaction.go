// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/FactomProject/factoid"
	"github.com/FactomProject/factom"
	"github.com/FactomProject/fctwallet/Wallet/Utility"
)

var _ = hex.EncodeToString
var serverFct = "localhost:8089"

var badChar, _ = regexp.Compile("[^A-Za-z0-9_-]")

type Response struct {
	Response string
	Success  bool
}

func ValidateKey(key string) (msg string, valid bool) {
	if len(key) > factoid.ADDRESS_LENGTH {
		return "Key is too long.  Keys must be less than or equal to 32 characters", false
	}
	if badChar.FindStringIndex(key) != nil {
		str := fmt.Sprintf("The key or name '%s' contains invalid characters.\n"+
			"Keys and names are restricted to alphanumeric characters,\n"+
			"minuses (dashes), and underscores", key)
		return str, false
	}
	return "", true
}

func getCmd(cmd string, cmderror string) {
	resp, err := http.Get(cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	resp.Body.Close()

	b := new(Response)
	if err := json.Unmarshal(body, b); err != nil || !b.Success {
		fmt.Println(cmderror)
		fmt.Println("Command Failed: ", string(body))
		os.Exit(1)
	}
	fmt.Println(b.Response)
	return
}

func postCmd(cmd string) {
	resp, err := http.PostForm(cmd, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	resp.Body.Close()

	b := new(Response)
	if err := json.Unmarshal(body, b); err != nil {
		fmt.Printf("Failed to parse the response from factomd: %s\n", body)
		os.Exit(1)
	}

	fmt.Println(b.Response)

	if !b.Success {
		os.Exit(1)
	}

	return
}

var transactions = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli transactions"
	cmd.description = "Prints information about pending transactions. Returns a list of all the transactions being constructed by the user.  It shows the fee required (at this point) as well as the fee the user will pay.  Some additional error checking is done as well, with messages provided to the user."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) > 0 {
			fmt.Println(cmd.helpMsg)
			return
		}

		str := fmt.Sprintf("http://%s/v1/factoid-get-transactions/", serverFct)
		getCmd(str, "Error Getting Transactions")
	}
	help.Add("transactions", cmd)
	return cmd
}()

var getlist = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli list TXNAME|ADDRESS|all"
	cmd.description = "List confirmed transactions' details."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Nothing to List.  Consider 'list all', or 'list [address]'.")
			return
		}
		var list string
		if len(args) == 0 {
			list = fmt.Sprintf("http://%s/v1/factoid-get-processed-transactions/", serverFct)
		} else if len(args) == 1 && args[0] == "all" {
			list = fmt.Sprintf("http://%s/v1/factoid-get-processed-transactions/?cmd=all", serverFct)
		} else if len(args) == 1 {
			list = fmt.Sprintf("http://%s/v1/factoid-get-processed-transactions/?address=%s", serverFct, args[0])
		} else {
			fmt.Println("Did not understand the arguments.  Proper syntax is either 'list all' or")
			fmt.Println("'list <address>' where <address> can be a valid Factoid or Entry Credit address")
		}
		postCmd(list)
	}
	help.Add("list", cmd)
	return cmd
}()

var getlistj = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listj TXNAME|ADDRESS|all"
	cmd.description = "List confirmed transactions' details."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Nothing to List.  Consider 'list all', or 'list [address]'.")
			return
		}
		var list string
		if len(args) == 0 {
			list = fmt.Sprintf("http://%s/v1/factoid-get-processed-transactionsj/", serverFct)
		} else if len(args) == 1 && args[0] == "all" {
			list = fmt.Sprintf("http://%s/v1/factoid-get-processed-transactionsj/?cmd=all", serverFct)
		} else if len(args) == 1 {
			list = fmt.Sprintf("http://%s/v1/factoid-get-processed-transactionsj/?address=%s", serverFct, args[0])
		} else {
			fmt.Println("Did not understand the arguments.  Proper syntax is either 'list all' or")
			fmt.Println("'list <address>' where <address> can be a valid Factoid or Entry Credit address")
		}
		postCmd(list)
	}
	help.Add("listj", cmd)
	return cmd
}()

var fctnewtrans = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli newtransaction TXNAME"
	cmd.description = "Create a new transaction. The TXNAME is used to add inputs, outputs, and ecoutputs (to buy entry credits).  Once the transaction is built, call validate, and if all is good, submit"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Missing Name")
			fmt.Println(cmd.helpMsg)
			return
		}
		msg, valid := ValidateKey(args[0])
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}
		str := fmt.Sprintf("http://%s/v1/factoid-new-transaction/%s", serverFct, args[0])
		postCmd(str)
	}
	help.Add("newtransaction", cmd)
	return cmd
}()

var fctdeletetrans = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli deletetransaction TXNAME"
	cmd.description = "Delete the specified transaction in flight."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Missing Name")
			fmt.Println(cmd.helpMsg)
			return
		}
		msg, valid := ValidateKey(args[0])
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}
		str := fmt.Sprintf("http://%s/v1/factoid-delete-transaction/%s", serverFct, args[0])
		postCmd(str)
	}
	help.Add("deletetransaction", cmd)
	return cmd
}()

var fctaddfee = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addfee TXNAME FCADDRESS"
	cmd.description = "Adds the needed fee to the given transaction. The Factoid Address specified must be an input to the transaction, and it must have a balance able to cover the additional fee. Also, the inputs must exactly balance the outputs,  since the logic to understand what to do otherwise is quite complicated, and prone to odd behavior."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 2 {
			fmt.Println("Was expecting a transaction name, and an address used as an input to that transaction.")
			fmt.Println(cmd.helpMsg)
			return
		}

		msg, valid := ValidateKey(args[0])
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}

		str := fmt.Sprintf("http://%s/v1/factoid-add-fee/?key=%s&name=%s",
			serverFct, args[0], args[1])
		postCmd(str)
	}
	help.Add("addfee", cmd)
	return cmd
}()

var fctsubfee = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli subfee TXNAME FCTADDRESS"
	cmd.description = "Subtracts the needed fee to the given transaction. The Factoid Address specified must be an output to the transaction. Also, the inputs must exactly balance the outputs,  since the logic to understand what to do otherwise is quite complicated, and prone to odd behavior."
	cmd.execFunc = func(args []string) {
		var res = flag.Bool("r", false, "resolve dns address")
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 2 {
			fmt.Println("Was expecting a transaction name, and an address used as an input to that transaction.")
			fmt.Println(cmd.helpMsg)
			return
		}

		msg, valid := ValidateKey(args[0])
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}

		addr := args[1]

		if *res {
			f, e, err := factom.ResolveDnsName(addr)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else if f == "" {
				if e == "" {
					fmt.Println("Could not resolve address")
					os.Exit(1)
				}
				addr = e
				return
			}
			addr = f
		}

		str := fmt.Sprintf("http://%s/v1/factoid-sub-fee/?key=%s&name=%s",
			serverFct, args[0], addr)
		postCmd(str)
	}
	help.Add("subfee", cmd)
	return cmd
}()

var fctaddinput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addinput TXNAME NAME|FCADDRESS AMOUNT"
	cmd.description = "Add an input to a transaction."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 3 {
			fmt.Println("Expecting a 1) transaction name, 2) an Address or Address name, and 3) an amount.")
			fmt.Println(cmd.helpMsg)
			return
		}

		msg, valid := ValidateKey(args[0])
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}

		amt, err := factoid.ConvertFixedPoint(args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ramt, err := strconv.ParseInt(amt, 10, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = factoid.ValidateAmounts(uint64(ramt))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		str := fmt.Sprintf("http://%s/v1/factoid-add-input/?key=%s&name=%s&amount=%s",
			serverFct, args[0], args[1], amt)
		postCmd(str)
	}
	help.Add("addinput", cmd)
	return cmd
}()

var fctaddoutput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addoutput [-r] TXNAME NAME|FCADDRESS|DNSADDRESS AMOUNT"
	cmd.description = "Add an output to a transaction."
	cmd.execFunc = func(args []string) {
		var res = flag.Bool("r", false, "resolve dns address")
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 3 {
			fmt.Println("Expecting a 1) transaction name, 2) an Address or Address name, and 3) an amount.")
			fmt.Println(cmd.helpMsg)
			return
		}
		// localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>

		addr := args[1]

		if *res {
			f, _, err := factom.ResolveDnsName(addr)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else if f == "" {
				fmt.Println("Could not resolve address")
				os.Exit(1)
			}

			addr = f
		}

		msg, valid := ValidateKey(args[0])
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}

		amt, err := factoid.ConvertFixedPoint(args[2])
		if err != nil {
			fmt.Println("Invalid format for a number: ", args[2])
			fmt.Println(cmd.helpMsg)
			return
		}

		ramt, err := strconv.ParseInt(amt, 10, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = factoid.ValidateAmounts(uint64(ramt))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		str := fmt.Sprintf("http://%s/v1/factoid-add-output/?key=%s&name=%s&amount=%s",
			serverFct, args[0], addr, amt)
		postCmd(str)
	}
	help.Add("addoutput", cmd)
	return cmd
}()

var fctaddecoutput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addecoutput [-r] TXNAME NAME|ECADDRESS|DNSADDRESS AMOUNT"
	cmd.description = "Add an ecoutput (purchase of entry credits to a transaction. Amount is denominated in factoids"
	cmd.execFunc = func(args []string) {
		var res = flag.Bool("r", false, "resolve dns address")
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 3 {
			fmt.Println("Expecting a 1) transaction name, 2) an Address or Address name, and 3) an amount.")
			fmt.Println(cmd.helpMsg)
			return
		}
		// localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>

		addr := args[1]

		if *res {
			_, e, err := factom.ResolveDnsName(addr)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else if e == "" {
				fmt.Println("Could not resolve address")
				os.Exit(1)
			}

			addr = e
		}

		msg, valid := ValidateKey(args[0])
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}

		amt, err := factoid.ConvertFixedPoint(args[2])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ramt, err := strconv.ParseInt(amt, 10, 64)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = factoid.ValidateAmounts(uint64(ramt))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		str := fmt.Sprintf("http://%s/v1/factoid-add-ecoutput/?key=%s&name=%s&amount=%s",
			serverFct, args[0], addr, amt)
		postCmd(str)
	}
	help.Add("addecoutput", cmd)
	return cmd
}()

var fctgetfee = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli getfee"
	cmd.description = "Get the current fee required for this  transaction. If a transaction is specified, then getfee returns the fee due for the  transaction. If no transaction is provided, then the cost of an Entry Credit is returned."
	cmd.execFunc = func(args []string) {
		var getfeereq string
		if len(args) > 1 {
			getfeereq = fmt.Sprintf("http://%s/v1/factoid-get-fee/?key=%s", serverFct, args[1])
		} else {
			getfeereq = fmt.Sprintf("http://%s/v1/factoid-get-fee/", serverFct)
		}

		resp, err := http.Get(getfeereq)

		if err != nil {
			fmt.Println("Command Failed Get")
			fmt.Println(cmd.helpMsg)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to understand response from fctwallet")
			os.Exit(1)
		}
		resp.Body.Close()

		// We pull the fee.  If the fee isn't positive, or if we fail to marshal, then there is a failure
		type x struct {
			Response string
			Success  bool
		}

		b := new(x)
		if err := json.Unmarshal(body, b); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else if !b.Success {
			fmt.Println(b.Response)
			os.Exit(1)
		}
		if len(args) < 2 {
			fmt.Printf("Currently, Entry Credits are %s Factoids each\n", b.Response)
		} else {
			fmt.Printf("The fee due for this transaction is %s Factoids\n", b.Response)
		}
	}
	help.Add("getfee", cmd)
	return cmd
}()

var fctproperties = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli properties"
	cmd.description = "Returns information about factomd, fctwallet, the Protocol version, the version of this CLI, and more." // TODO
	cmd.execFunc = func(args []string) {
		resp, err := http.Get(fmt.Sprintf("http://%s/v1/properties/", serverFct))

		if err != nil {
			fmt.Println("Get Properties failed")
			os.Exit(1)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Failed to understand response from fctwallet")
			os.Exit(1)
		}
		resp.Body.Close()

		// We pull the fee.  If the fee isn't positive, or if we fail to marshal, then there is a failure
		type x struct {
			Response string
			Success  bool
		}

		b := new(x)
		if err := json.Unmarshal(body, b); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else if !b.Success {
			fmt.Println(b.Response)
			os.Exit(1)
		}

		ret := b.Response + fmt.Sprintln("factom-cli Version:", Version)

		total, err := Utility.TotalFactoids()
		if err == nil {
			ret = ret + fmt.Sprintf("    Total Factoids: %s", strings.TrimSpace(factoid.ConvertDecimal(total)))
		}

		fmt.Println(ret)
	}
	help.Add("properties", cmd)
	return cmd
}()

var fctsign = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli sign TXNAME"
	cmd.description = "Sign the transaction specified by the TXNAME."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Missing Name")
			fmt.Println(cmd.helpMsg)
			return
		}

		msg, valid := ValidateKey(args[0])
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}

		str := fmt.Sprintf("http://%s/v1/factoid-sign-transaction/%s", serverFct, args[0])
		postCmd(str)

	}
	help.Add("sign", cmd)
	return cmd
}()

var fctsubmit = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli submit TXNAME"
	cmd.description = "Submit the transaction specified by the TXNAME to Factom."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Missing Name")
			fmt.Println(cmd.helpMsg)
			return
		}

		msg, valid := ValidateKey(args[0])
		if !valid {
			fmt.Println(msg)
			os.Exit(1)
		}

		s := struct{ Transaction string }{args[0]}

		jdata, err := json.Marshal(s)
		if err != nil {
			fmt.Println("Submit failed")
			os.Exit(1)
		}

		str := fmt.Sprintf("http://%s/v1/factoid-submit/%s", serverFct, bytes.NewBuffer(jdata))
		postCmd(str)
	}
	help.Add("submit", cmd)
	return cmd
}()
