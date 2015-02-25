// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func bintx(args []string) error {
	os.Args = args
	
	var (
		serv = flag.String("s", "localhost:8088", "path to the factomclient")
	)
	flag.Parse()
	args = flag.Args()
	if len(args) < 1 {
		return man("bintx")
	}
	tx := args[0]
	server := "http://" + *serv + "/v1/bintx"
	data := url.Values{
		"tx":      {tx},
	}
	
	resp, err := http.PostForm(server, data)
	if err != nil {
		return err
	}
	p, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Println(string(p))

	return nil
}
