// Copyright 2017 Factom Foundation
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

// unlockwallet creates a new transaction in the wallet.
var unlockwallet = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli walletpassphrase [-q] \"passphrase\" <seconds-to-unlock>"
	cmd.description = "Unlock the wallet for some number of seconds; must be an encrypted wallet. -q quiet."
	cmd.execFunc = func(args []string) {
		os.Args = args
		qflag := flag.Bool("q", false, "quiet mode; no output")
		flag.Parse()
		args = flag.Args()

		if len(args) != 2 {
			fmt.Println(cmd.helpMsg)
			return
		}

		seconds, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			errorln(err)
			return
		}

		unix, err := factom.UnlockWallet(args[0], seconds)
		if err != nil {
			errorln(err)
			return
		}

		// output
		switch {
		// quiet mode; don't print anything
		case *qflag:
		default:
			fmt.Println(unix)
		}
	}
	help.Add("walletpassphrase", cmd)
	return cmd
}()
