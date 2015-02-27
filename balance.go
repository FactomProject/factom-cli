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

func balance(args []string) error {
	os.Args = args
	serv := flag.String("s", "localhost:8088",
		"Path to the factom api server")
	server := "http://" + *serv + "/v1/entrycreditbalance"
	flag.Parse()
	args = flag.Args()
	
	if len(args) == 0 {
		args = []string{"ec"}
	}
	
	switch args[0] {
	case "ec":
		return ecBalance(args, server)
	case "factoid":
		return factoidBalance(args, server)
	default:
		return man("balance")
	}
	panic("something went wrong with balance")
}

func ecBalance(args []string, server string) error {
	data := url.Values{
		"pubkey": {"wallet"},
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
	return fmt.Errorf("%s", p)
}

func factoidBalance(args []string, server string) error {
	return fmt.Errorf("Deadend")
}