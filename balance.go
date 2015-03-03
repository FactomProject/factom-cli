// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// balance prints the current balance of the specified wallet
func balance(args []string) error {
	api := "http://" + server + "/v1/creditbalance"
	key := wallet

	if len(args) == 1 {
		args = append(args, "ec")
	}
	args = args[1:]

	switch args[0] {
	case "ec":
		return ecBalance(key, api)
	case "factoid":
		return factoidBalance(key, api)
	default:
		return man("balance")
	}
	panic("something went wrong with balance")
}

func ecBalance(pubkey, server string) error {
	type balance struct {
		Publickey string
		Credits   float64
	}

	data := url.Values{
		"pubkey": {pubkey},
	}

	resp, err := http.PostForm(server, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	for {
		var b *balance
		if err := dec.Decode(&b); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		fmt.Println("EC Balance:", b.Credits)
	}

	return nil
}

func factoidBalance(pubkey, server string) error {
	return fmt.Errorf("Factoid Balance: not implimented")
}
