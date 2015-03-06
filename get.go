// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func get(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("get")
	}

	switch args[0] {
	case "dblocks":
		return getDBlocks(args)
	case "eblock":
		return getEBlock(args)
	case "entry":
		return getEntry(args)
	case "height":
		return getHeight()
	default:
		return man("get")
	}

	panic("something went really wrong with get!")
}

func getDBlocks(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 2 {
		return man("getDBlocks")
	}

	from, to := args[0], args[1]
	api := fmt.Sprintf("http://%s/v1/dblocksbyrange/%s/%s", server, from, to)

	resp, err := http.Get(api)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(p))

	return nil
}

func getEBlock(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("getEBlocks")
	}

	mr := args[0]
	api := fmt.Sprintf("http://%s/v1/eblockbymr/%s", server, mr)

	resp, err := http.Get(api)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(p))

	return nil
}

func getEntry(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("getEntry")
	}

	hash := args[0]
	api := fmt.Sprintf("http://%s/v1/entry/%s", server, hash)

	resp, err := http.Get(api)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(p))

	return nil
}

func getHeight() error {
	resp, err := http.Get("http://" + server + "/v1/dblockheight")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(p))

	return nil
}
