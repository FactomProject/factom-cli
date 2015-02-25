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

func buy(args []string) error {
	os.Args = args

	var (
		amt  string = "0"
		serv        = flag.String("s", "localhost:8088", "path to the factomclient")
	)
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("buy")
	}
	server := "http://" + *serv + "/v1/buycredit"
	data := url.Values{
		"to":     {"wallet"},
		"amount": {amt},
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
