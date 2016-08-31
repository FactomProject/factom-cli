// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
)

var receipt = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli receipt ENTRYHASH"
	cmd.description = "Returns a Receipt for a given Entry"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) != 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		txID := args[0]

		resp, err := factom.GetReceipt(txID)
		if err != nil {
			errorln(err)
			return
		}
		str, err := json.MarshalIndent(resp, "", "\t")
		if err != nil {
			errorln(err)
			return
		}
		fmt.Printf("%s\n", str)
	}
	help.Add("receipt", cmd)
	return cmd
}()
