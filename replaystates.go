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
	type JSONError struct {
		Code    int         `json:"code"`           // The error code associated with the error type
		Message string      `json:"message"`        // The error message as a concise single sentence
		Data    interface{} `json:"data,omitempty"` // Optional data object containing additional information about the error
	}

	type replayResponse struct {
		Message string `json:"message"`
		Start   int64  `json:"startheight"`
		End     int64  `json:"endheight"`
	}

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
			errorln(err)
			return
		}
		startheight := startHeightStr

		var endheight int64
		endheight = 0
		if len(args) == 2 {
			endHeightStr, err := strconv.ParseInt(args[1], 10, 32)
			if err != nil {
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

		fmt.Println(res.Message)
	}
	help.Add("replaydbstates", cmd)
	return cmd
}()
