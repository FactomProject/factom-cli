// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
)

// backupwallet returns the wallet seed and all of the addresses from the
// wallet.
var backupwallet = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli backupwallet"
	cmd.description = "Backup the running wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		s, err := factom.BackupWallet()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Print(s)
	}
	help.Add("backupwallet", cmd)
	return cmd
}()
