// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"os"
	"strconv"
	
	"github.com/FactomProject/factom"
)

func buy(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("buy")
	}
	
	amt, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	
	pub, err := ecPubKey()
	if err != nil {
		return err
	}
	
	err = factom.BuyTestCredits(hex.EncodeToString(pub[:]), amt)
	if err != nil {
		return err
	}

	return nil
}

//func buy(args []string) error {
//	var (
//		amt  string
//		addr string = wallet
//		api         = "http://" + server + "/v1/buycredit"
//	)
//	if len(args) == 1 {
//		return man("buy")
//	}
//	args = args[1:]
//	amt = args[0]
//
//	data := url.Values{
//		"to":     {addr},
//		"amount": {amt},
//	}
//
//	resp, err := http.PostForm(api, data)
//	if err != nil {
//		return err
//	}
//	p, err := ioutil.ReadAll(resp.Body)
//	resp.Body.Close()
//	fmt.Println(string(p))
//
//	return nil
//}
