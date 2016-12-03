// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/FactomProject/factom"
)

var signmessage = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli signmessage ADDRESS \"message\""
	cmd.description = "Sign the message with the given Factoid or EntryCredit address"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 2 {
			fmt.Println(cmd.helpMsg)
			return
		}

		addr := args[0]
		msg := args[1]

		pubKeyPrefixed, sig, err := factom.SignMessage(addr, msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Public key: " + pubKeyPrefixed)
		fmt.Println("Signature: " + sig)
		return
	}
	help.Add("signmessage", cmd)
	return cmd
}()

var verifymessage = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli verifymessage <public key> \"signature\" \"message\""
	cmd.description = "Verify the signature for the given message and public key"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) != 3 {
			fmt.Println(cmd.helpMsg)
			return
		}

		addr := args[0]
		sig := args[1]
		msg := args[2]

		result, addrString, err := factom.VerifyMessage(addr, sig, msg)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Address: " + addrString)
		fmt.Println("Result: " + strconv.FormatBool(result))
		return
	}
	help.Add("verifymessage", cmd)
	return cmd
}()
