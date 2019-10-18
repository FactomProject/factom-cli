// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
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

	type JSON2Response struct {
		JSONRPC string      `json:"jsonrpc"` // version string which MUST be "2.0" (version 1.0 didn't have this field)
		ID      interface{} `json:"id"`      // Unique client defined ID associated with incoming request. It may be a number,
		// string, or nil. It is used by the remote to formulate a this response object
		// containing the same ID back to the client. If nil, remote treats request object
		// as a notification, and that the client does not expect a response object back
		Error  *JSONError  `json:"error,omitempty"`  // Must be present if called subroutine had an error (mutually exclusive with Result)
		Result interface{} `json:"result,omitempty"` // Must be present if called subroutine succeeded (mutually exclusive with Error)
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

		resStr, err := json.Marshal(res)
		if err != nil {
			errorln(err)
			return
		}

		response := JSON2Response{}
		json.Unmarshal([]byte(resStr), &response)
		fmt.Println(response.Result.(interface{}).(map[string]interface{})["message"])
	}
	help.Add("replaydbstates", cmd)
	return cmd
}()
