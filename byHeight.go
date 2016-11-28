// Copyright 2016 Factom Foundation
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

var dblockByHeight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli dblock-by-height height"
	cmd.description = "Returns dblock by height"
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

		resp, err := factom.GetBlockByHeightRaw("d", height)
		if err != nil {
			errorln(err)
			return
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
	help.Add("dblock-by-height", cmd)
	return cmd
}()

var ablockByHeight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli ablock-by-height height"
	cmd.description = "Returns ablock by height"
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

		resp, err := factom.GetBlockByHeightRaw("a", height)
		if err != nil {
			errorln(err)
			return
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
	help.Add("ablock-by-height", cmd)
	return cmd
}()

var ecblockByHeight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli ecblock-by-height height"
	cmd.description = "Returns ecblock by height"
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

		resp, err := factom.GetBlockByHeightRaw("ec", height)
		if err != nil {
			errorln(err)
			return
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
	help.Add("ecblock-by-height", cmd)
	return cmd
}()

var fblockByHeight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli fblock-by-height height"
	cmd.description = "Returns fblock by height"
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

		resp, err := factom.GetBlockByHeightRaw("f", height)
		if err != nil {
			errorln(err)
			return
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
	help.Add("fblock-by-height", cmd)
	return cmd
}()
