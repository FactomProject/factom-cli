package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func buy(args []string) error {
	os.Args = args
	
	var (
		serv = flag.String("s", "localhost:8088", "path to the factomclient")
	)
	flag.Parse()
	
	if len(args) < 2 {
		return fmt.Errorf("the ammount of factoids to be transferd must be specified")
	}
	amt := args[1]
	server := "http://" + *serv + "/v1/buycredit"
	data := url.Values{
		"to":      {"wallet"},
		"ammount": {amt},
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
