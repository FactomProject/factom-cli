package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type jsonentry struct {
	ChainID string
	ExtIDs  []string
	Data    string
}

type extids []string

func (e *extids) String() string {
	return fmt.Sprint(*e)
}

func (e *extids) Set(s string) error {
	*e = append(*e, s)
	return nil
}

func put() {
	var (
		cid  = flag.String("c", "", "hex encoded chainid for the entry")
		serv = flag.String("s", "localhost:8088", "path to the factomclient")
		eids extids
	)

	flag.Var(&eids, "e", "external id for the entry")
	flag.Parse()
	
	server := "http://" + *serv + "/v1/submitentry"

	d := make([]byte, 1024)
	n, err := os.Stdin.Read(d)
	if err != nil {
		fmt.Println("put(): ", err)
	}
	d = d[:n]

	e := new(jsonentry)
	e.ChainID = *cid
	for _, v := range eids {
		e.ExtIDs = append(e.ExtIDs, string(v))
	}
	e.Data = hex.EncodeToString(d)

	b, err := json.Marshal(e)
	if err != nil {
		fmt.Println("put(): ", err)
	}
	
	data := url.Values{
		"datatype": {"entry"},
		"format":   {"json"},
		"entry":    {string(b)},
	}
	
	_, err = http.PostForm(server, data)
	if err != nil {
		fmt.Println("Error: ", err)
	}	
}
