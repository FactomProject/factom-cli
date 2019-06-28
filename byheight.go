// All functions in this file are used to return GetBlockByHeightRaw.
// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/FactomProject/factom"
)

// Abheight - public so it can be accessed by tests
// Deprecated: should use ablock.
var Abheight = func() *fctCmd {
	var supressRawData string
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get abheight HEIGHT -r (to suppress Raw Data)"
	cmd.description = "Get Admin Block by height"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		h := args[0]
		height, err := strconv.ParseInt(h, 0, 64)
		if err != nil {
			fmt.Println(cmd.helpMsg)
			return
		}

		if len(args) > 1 {
			supressRawData = args[1]
		}

		resp, err := factom.GetBlockByHeightRaw("a", height)
		if err != nil {
			errorln(err)
			return
		}

		if supressRawData == "-r" {
			resp.RawData = ""
		}

		data, err := json.Marshal(resp)
		if err != nil {
			errorln(err)
			return
		}

		var out bytes.Buffer
		json.Indent(&out, data, "", "\t")

		fmt.Printf("%s\n", out.Bytes())
	}
	help.Add("get abheight", cmd)
	return cmd
}()

//Deprecated: should use dblock
var Dbheight = func() *fctCmd {
	var supressRawData string
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get dbheight HEIGHT -r (to suppress Raw Data)"
	cmd.description = "Get Directory Block by height"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		h := args[0]
		height, err := strconv.ParseInt(h, 0, 64)
		if err != nil {
			fmt.Println(cmd.helpMsg)
			return
		}

		if len(args) > 1 {
			supressRawData = args[1]
		}

		resp, err := factom.GetBlockByHeightRaw("d", height)
		if err != nil {
			errorln(err)
			return
		}

		if supressRawData == "-r" {
			resp.RawData = ""
		}

		data, err := json.Marshal(resp)
		if err != nil {
			errorln(err)
			return
		}

		var out bytes.Buffer
		json.Indent(&out, data, "", "\t")

		fmt.Printf("%s\n", out.Bytes())
	}
	help.Add("get dbheight", cmd)
	return cmd
}()

//Deprecated: should use ecblock
var Ecbheight = func() *fctCmd {
	var supressRawData string
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get ecbheight HEIGHT -r (to suppress Raw Data)"
	cmd.description = "Get Entry Credit Block by height"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		h := args[0]
		height, err := strconv.ParseInt(h, 0, 64)
		if err != nil {
			fmt.Println(cmd.helpMsg)
			return
		}

		if len(args) > 1 {
			supressRawData = args[1]
		}

		resp, err := factom.GetBlockByHeightRaw("ec", height)
		if err != nil {
			errorln(err)
			return
		}

		if supressRawData == "-r" {
			resp.RawData = ""
		}

		data, err := json.Marshal(resp)
		if err != nil {
			errorln(err)
			return
		}

		var out bytes.Buffer
		json.Indent(&out, data, "", "\t")

		fmt.Printf("%s\n", out.Bytes())
	}
	help.Add("get ecbheight", cmd)
	return cmd
}()
//Deprecated: should use fblock.
var Fbheight = func() *fctCmd {
	var supressRawData string
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get fbheight HEIGHT -r (to suppress Raw Data)"
	cmd.description = "Get Factoid Block by height"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		h := args[0]
		height, err := strconv.ParseInt(h, 0, 64)
		if err != nil {
			fmt.Println(cmd.helpMsg)
			return
		}

		if len(args) > 1 {
			supressRawData = args[1]
		}

		resp, err := factom.GetBlockByHeightRaw("f", height)
		if err != nil {
			errorln(err)
			return
		}

		if supressRawData == "-r" {
			resp.RawData = ""
		}

		data, err := json.Marshal(resp)
		if err != nil {
			errorln(err)
			return
		}

		var out bytes.Buffer
		json.Indent(&out, data, "", "\t")

		fmt.Printf("%s\n", out.Bytes())
	}
	help.Add("get fbheight", cmd)
	return cmd
}()
