// Copyright 2015 Factom Foundation
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

	"github.com/FactomProject/cli"
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

// VerifyAddressType returns a bool success and a string describing the address
// type
func VerifyAddressType(address string)(string, bool) {
  	var resp string = "Not a Valid Factoid Address"
	var pass bool = false
  
  	if (strings.HasPrefix(address,"FA")) {
		if (factoid.ValidateFUserStr(address)) {
			resp = "Factoid - Public"
			pass = true
		}
	} else if (strings.HasPrefix(address,"EC")) {
		if (factoid.ValidateECUserStr(address)) {
			resp = "Entry Credit - Public"
			pass = true
		}
	} else if (strings.HasPrefix(address,"Fs")) {
		if (factoid.ValidateFPrivateUserStr(address)) {
			resp = "Factoid - Private"
			pass = true
		}
	} else if (strings.HasPrefix(address,"Es")) {
		if (factoid.ValidateECPrivateUserStr(address)) {
			resp = "Entry Credit - Private"
			pass = true
		}
	} 



	//  Add Netki resolution here
	//else if (checkNetki) {
	//	if (factoid.ValidateECPrivateUserStr(address)) {
	//		resp = "{\"AddressType\":\"Factoid - Public\", \"TypeCode\":4 ,\"Success\":true}"
	//	}
	//} 


	return resp, pass
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

var importaddr = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli import name key"
	cmd.description = "Import an Entry Credit or Factoid Private Key"
	cmd.execFunc = func(args []string) {
		if len(args) < 3 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if strings.HasPrefix(args[2], "Fs") {
			if addr, err := factom.GenerateFactoidAddressFromHumanReadablePrivateKey(args[1], args[2]); err != nil {
				fmt.Println(err)
				return
			} else {
				fmt.Println(args[1], addr)
			}
		} else if strings.HasPrefix(args[2], "Es") {
			if addr, err := factom.GenerateEntryCreditAddressFromHumanReadablePrivateKey(args[1], args[2]); err != nil {
				fmt.Println(err)
				return
			} else {
				fmt.Println(args[1], addr)
			}
		} else {
			fmt.Println("Invalid Key")
			fmt.Println(cmd.helpMsg)
		}
	}
	help.Add("importaddress", cmd)
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
	cmd.helpMsg = "factom-cli list [transaction id|address|all]"
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
	cmd.helpMsg = "factom-cli "
	cmd.description = ""
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
	cmd.helpMsg = "factom-cli newtransaction key"
	cmd.description = "Create a new transaction. The key is used to add inputs, outputs, and ecoutputs (to buy entry credits).  Once the transaction is built, call validate, and if all is good, submit"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Missing Key")
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
	cmd.helpMsg = "factom-cli deletetransaction transaction"
	cmd.description = "Delete the specified transaction in flight."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Missing Key")
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
	cmd.helpMsg = "factom-cli addfee transaction fctaddress"
	cmd.description = "Adds the needed fee to the given transaction. The address specified must be an input to the transaction, and it must have a balance able to cover the additional fee. Also, the inputs must exactly balance the outputs,  since the logic to understand what to do otherwise is quite complicated, and prone to odd behavior."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 2 {
			fmt.Println("Was expecting a transaction key, and an address used as an input to that transaction.")
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

var fctaddinput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addinput key name|address amount"
	cmd.description = "Add an input to a transaction."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 3 {
			fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
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
	cmd.helpMsg = "factom-cli addoutput key name|address amount"
	cmd.description = "Add an output to a transaction."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 3 {
			fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
			fmt.Println(cmd.helpMsg)
			return
		}
		// localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>

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
			serverFct, args[0], args[1], amt)
		postCmd(str)
	}
	help.Add("addoutput", cmd)
	return cmd
}()

var fctaddecoutput = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addecoutput key name|address amount"
	cmd.description = "Add an ecoutput (purchase of entry credits to a transaction."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 3 {
			fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
			fmt.Println(cmd.helpMsg)
			return
		}
		// localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>

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
			serverFct, args[0], args[1], amt)
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
		getproperties := fmt.Sprintf("http://%s/v1/properties/", serverFct)

		resp, err := http.Get(getproperties)

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

		top := Version / 1000000
		mid := (Version % 1000000) / 1000
		low := Version % 1000

		ret := b.Response + fmt.Sprintf("factom-cli Version: %d.%d.%d\n", top, mid, low)

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
	cmd.helpMsg = "factom-cli sign transaction"
	cmd.description = "Sign the transaction specified by the key."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Missing Key")
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
	cmd.helpMsg = "factom-cli submit transaction"
	cmd.description = "Submit the transaction specified by the key to Factom."
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println("Missing Key")
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
			fmt.Println("Submitt failed")
			os.Exit(1)
		}

		str := fmt.Sprintf("http://%s/v1/factoid-submit/%s", serverFct, bytes.NewBuffer(jdata))
		postCmd(str)
	}
	help.Add("submit", cmd)
	return cmd
}()
