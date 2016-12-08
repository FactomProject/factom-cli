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

var abheight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get abheight HEIGHT"
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
	help.Add("get abheight", cmd)
	return cmd
}()

var dbheight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get dbheight HEIGHT"
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
	help.Add("get dbheight", cmd)
	return cmd
}()

var ecbheight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get ecbheight HEIGHT"
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
	help.Add("get ecbheight", cmd)
	return cmd
}()

var fbheight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get fbheight HEIGHT"
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
	help.Add("get fbheight", cmd)
	return cmd
}()
