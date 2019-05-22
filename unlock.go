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
	"github.com/posener/complete"
)

// unlockwallet creates a new transaction in the wallet.
var unlockwallet = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli unlockwallet [-v] \"passphrase\" <seconds-to-unlock>"
	cmd.description = "Unlock the wallet for some number of seconds; must be an encrypted wallet. -v verbose."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-v": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		vflag := flag.Bool("v", false, "verbose mode; print relock time in UTC epoch seconds when the unlock is successful")
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
		// verbose mode; print when wallet will re-lock
		case *vflag:
			fmt.Println(unix)
		default:
		}
	}
	help.Add("unlockwallet", cmd)
	return cmd
}()
