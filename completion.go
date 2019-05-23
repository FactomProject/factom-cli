// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/FactomProject/factom"
	"github.com/posener/complete"
)

// predictTxName creates command completions from the tx names in
// factom-walletd. It predicts nothing if there is no wallet running.
var predictTxName = complete.PredictFunc(func(args complete.Args) []string {
	argc := 0
	for _, arg := range args.Completed[1:] {
		if string(arg[0]) != "-" {
			argc++
		}
	}
	if argc != 0 {
		return nil
	}

	txs, err := factom.ListTransactionsTmp()
	if err != nil {
		return nil
	}

	s := make([]string, 0)
	for _, tx := range txs {
		s = append(s, tx.Name)
	}
	return s
})

// predictAddress creates command completions from the addresses in
// factom-walletd. It predicts nothing if there is not wallet running.
var predictAddress = complete.PredictFunc(func(args complete.Args) []string {
	argc := 0
	for _, arg := range args.Completed[1:] {
		if string(arg[0]) != "-" {
			argc++
		}
	}
	if argc != 0 {
		return nil
	}

	fcs, ecs, err := factom.FetchAddresses()
	if err != nil {
		return nil
	}

	s := make([]string, 0)
	for _, fc := range fcs {
		s = append(s, fc.String())
	}
    for _, ec := range ecs {
        s = append(s, ec.String())
    }
	return s
})
