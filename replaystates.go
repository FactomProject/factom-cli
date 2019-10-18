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

var replaydbstates = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli replaydbstates STARTHEIGHT ENDHEIGHT"
	cmd.description = "Emit DBStateMsgs over the LiveFeed API between two specifed block heights"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 || len(args) > 2 {
			fmt.Println(cmd.helpMsg)
			return
		}

		startHeightStr, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			// fmt.Println(cmd.helpMsg)
			errorln(err)
			return
		}
		startheight := startHeightStr

		var endheight int64
		endheight = 0
		if len(args) == 2 {
			endHeightStr, err := strconv.ParseInt(args[1], 10, 32)
			if err != nil {
				// fmt.Println(cmd.helpMsg)
				errorln(err)
				return
			}
			endheight = endHeightStr
		}
		res, err := factom.ReplayDBlockFromHeight(startheight, endheight)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(res)
	}
	help.Add("replaydbstates", cmd)
	return cmd
}()
