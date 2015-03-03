// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func buy(args []string) error {
	var (
		amt  string
		addr string = wallet
		api         = "http://" + server + "/v1/buycredit"
	)
	if len(args) == 1 {
		return man("buy")
	}
	args = args[1:]
	amt = args[0]

	data := url.Values{
		"to":     {addr},
		"amount": {amt},
	}

	resp, err := http.PostForm(api, data)
	if err != nil {
		return err
	}
	p, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(string(p))

	return nil
}
