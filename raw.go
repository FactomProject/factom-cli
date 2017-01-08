// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/FactomProject/factom"
)

var getraw = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get raw HASH"
	cmd.description = "Returns a raw hex representation of a block, transaction, entry, or commit"
	cmd.execFunc = func(args []string) {
		if len(args) < 2 {
			fmt.Println(cmd.helpMsg)
			return
		}

		h := args[1]
		hx, err := hex.DecodeString(strings.Replace(h, "\"", "", -1))
		if err != nil {
			errorln("Error reading hash")
			return
		}
		if len(hx) != 32 {
			errorln("Invalid argument length - should be 64 characters (32 bytes) long")
			return
		}

		raw, err := factom.GetRaw(fmt.Sprintf("%x", hx))
		if err != nil {
			errorln(err)
			return
		}
		if len(raw) > 0 {
			fmt.Printf("%x\n", raw)
			return
		}

		fmt.Printf("Block, transaction or entry not found.\n")
	}
	help.Add("get raw", cmd)
	return cmd
}()
