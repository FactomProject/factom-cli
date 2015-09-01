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
	fct "github.com/FactomProject/factoid"
	"github.com/FactomProject/factom"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

var _ = hex.EncodeToString
var serverFct = "localhost:8089"

var badChar, _ = regexp.Compile("[^A-Za-z0-9_-]")

type Response struct {
	Response string
	Success  bool
}

func ValidateKey(key string) (msg string, valid bool) {
	if len(key) > fct.ADDRESS_LENGTH {
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

// Generates a new Address
func generateaddress(args []string) {

	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 2 {
		man("generatefactoidaddress")
		os.Exit(1)
	}

	msg, valid := ValidateKey(args[1])
	if !valid {
		fmt.Println(msg)
		os.Exit(1)
	}
	
	var err error
	var Addr string
	if len(args) == 2 {
		switch args[0] {
		case "ec":
			Addr, err = factom.GenerateEntryCreditAddress(args[1])
		case "fct":
			Addr, err = factom.GenerateFactoidAddress(args[1])
		default:
			panic("Expected ec|fct name")
		}
	} else {
		switch args[0] {
		case "ec":
			Addr, err = factom.GenerateEntryCreditAddressFromHumanReadablePrivateKey(args[1], args[2])
			if err == nil {
				break
			}
			Addr, err = factom.GenerateEntryCreditAddressFromPrivateKey(args[1], args[2])
			if err == nil {
				break
			}
		case "fct":
			Addr, err = factom.GenerateFactoidAddressFromHumanReadablePrivateKey(args[1], args[2])
			if err == nil {
				break
			}
			Addr, err = factom.GenerateFactoidAddressFromMnemonic(args[1], args[2])
			if err == nil {
				break
			}
			Addr, err = factom.GenerateFactoidAddressFromPrivateKey(args[1], args[2])
			if err == nil {
				break
			}
		default:
			panic("Expected ec|fct name")
		}
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(args[0], " = ", Addr)
	return

}

func getaddresses(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) > 0 {
		man("getaddresses")
		os.Exit(1)
	}

	str := fmt.Sprintf("http://%s/v1/factoid-get-addresses/", serverFct)
	getCmd(str, "Error printing addresses")

	return
}

func gettransactions(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) > 0 {
		man("transactions")
		os.Exit(1)
	}

	str := fmt.Sprintf("http://%s/v1/factoid-get-transactions/", serverFct)
	getCmd(str, "Error Getting Transactions")

	return
}

func fctnewtrans(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		fmt.Println("Missing Key")
		man("newtransaction")
		os.Exit(1)
	}
	msg, valid := ValidateKey(args[0])
	if !valid {
		fmt.Println(msg)
		os.Exit(1)
	}
	str := fmt.Sprintf("http://%s/v1/factoid-new-transaction/%s", serverFct, args[0])
	postCmd(str)

	return

}

func fctdeletetrans(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		fmt.Println("Missing Key")
		man("deletetransaction")
		os.Exit(1)
	}
	msg, valid := ValidateKey(args[0])
	if !valid {
		fmt.Println(msg)
		os.Exit(1)
	}
	str := fmt.Sprintf("http://%s/v1/factoid-delete-transaction/%s", serverFct, args[0])
	postCmd(str)

	return

}

func fctaddfee(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 2 {
		fmt.Println("Was expecting a transaction key, and an address used as an input to that transaction.")
		os.Exit(1)
	}

	msg, valid := ValidateKey(args[0])
	if !valid {
		fmt.Println(msg)
		os.Exit(1)
	}

	str := fmt.Sprintf("http://%s/v1/factoid-add-fee/?key=%s&name=%s",
		serverFct, args[0], args[1])
	postCmd(str)

	return
}

func fctaddinput(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 3 {
		fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
		os.Exit(1)
	}

	msg, valid := ValidateKey(args[0])
	if !valid {
		fmt.Println(msg)
		os.Exit(1)
	}

	amt, err := fct.ConvertFixedPoint(args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ramt, err := strconv.ParseInt(amt, 10, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = fct.ValidateAmounts(uint64(ramt))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	str := fmt.Sprintf("http://%s/v1/factoid-add-input/?key=%s&name=%s&amount=%s",
		serverFct, args[0], args[1], amt)
	postCmd(str)

	return
}

func fctaddoutput(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 3 {
		fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
		os.Exit(1)
	}
	// localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>

	msg, valid := ValidateKey(args[0])
	if !valid {
		fmt.Println(msg)
		os.Exit(1)
	}

	amt, err := fct.ConvertFixedPoint(args[2])
	if err != nil {
		fmt.Println("Invalid format for a number: ", args[2])
		man("addoutput")
		os.Exit(1)
	}

	ramt, err := strconv.ParseInt(amt, 10, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = fct.ValidateAmounts(uint64(ramt))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	str := fmt.Sprintf("http://%s/v1/factoid-add-output/?key=%s&name=%s&amount=%s",
		serverFct, args[0], args[1], amt)
	postCmd(str)

	return
}

func fctaddecoutput(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 3 {
		fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
		os.Exit(1)
	}
	// localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>

	msg, valid := ValidateKey(args[0])
	if !valid {
		fmt.Println(msg)
		os.Exit(1)
	}

	amt, err := fct.ConvertFixedPoint(args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ramt, err := strconv.ParseInt(amt, 10, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, err = fct.ValidateAmounts(uint64(ramt))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	str := fmt.Sprintf("http://%s/v1/factoid-add-ecoutput/?key=%s&name=%s&amount=%s",
		serverFct, args[0], args[1], amt)
	postCmd(str)

	return
}

func fctgetfee(args []string) {
	var getfeereq string
	if len(args) > 1 {
		getfeereq = fmt.Sprintf("http://%s/v1/factoid-get-fee/?key=%s", serverFct, args[1])
	} else {
		getfeereq = fmt.Sprintf("http://%s/v1/factoid-get-fee/", serverFct)
	}

	resp, err := http.Get(getfeereq)

	if err != nil {
		fmt.Println("Command Failed Get")
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
	if len(args) < 2 {
		fmt.Printf("Currently, Entry Credits are %s Factoids each\n", b.Response)
	} else {
		fmt.Printf("The fee due for this transaction is %s Factoids\n", b.Response)
	}
}

func fctsign(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		fmt.Println("Missing Key")
		os.Exit(1)
	}

	msg, valid := ValidateKey(args[0])
	if !valid {
		fmt.Println(msg)
		os.Exit(1)
	}

	str := fmt.Sprintf("http://%s/v1/factoid-sign-transaction/%s", serverFct, args[0])
	postCmd(str)

}

func fctsubmit(args []string) {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		fmt.Println("Missing Key")
		os.Exit(1)
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

