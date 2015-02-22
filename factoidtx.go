// Copyright (c) 2015 FactomProject/FactomCode Systems LLC.
// Use of this source code is governed by an ISC
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

func factoidtx(args []string) error {
	os.Args = args
	
	var (
		serv = flag.String("s", "localhost:8088", "path to the factomclient")
		wallet = flag.String("w", "", "Factoid wallet address")
	)
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return fmt.Errorf("the ammount of factoids to be transferd must be specified")
	}
	amt := args[0]
	server := "http://" + *serv + "/v1/factoidtx"
	data := url.Values{
		"to":      {*wallet},
		"ammount": {amt},
	}
	
	resp, err := http.PostForm(server, data)
	if err != nil {
		return err
	}
	p, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(string(p))

	return nil
}
