// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// balance prints the current balance of the specified wallet
func balance(args []string) error {
	os.Args = args
	serv := flag.String("s", "localhost:8088",
		"Path to the factom api server")
	flag.Parse()
	args = flag.Args()
	server := "http://" + *serv + "/v1/creditbalance"
	if len(args) == 0 {
		args = []string{"ec", "wallet"}
	}
	if len(args) == 1 {
		args = append(args, "wallet")
	}
	key := args[1]
		
	switch args[0] {
	case "ec":
		return ecBalance(key, server)
	case "factoid":
		return factoidBalance(key, server)
	default:
		return man("balance")
	}
	panic("something went wrong with balance")
}

func ecBalance(pubkey, server string) error {
	data := url.Values{
		"pubkey": {pubkey},
	}
	
	resp, err := http.PostForm(server, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("Entry Credit Balance:", string(p))
	
	return nil
}

func factoidBalance(pubkey, server string) error {
	return fmt.Errorf("Factoid Balance: not implimented")
}