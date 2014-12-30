package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
)

const usage = "factom -s [server] -c [chainid] -e [extid] -e [extid2] <data"

type entry struct {
	ChainID	string
	ExtIDs	[]string
	Data	string
}

type extids []string
func (e *extids) String() string {
	return fmt.Sprint(*e)
}
func (e *extids) Set(s string) error {
	*e = append(*e, s)
	return nil
}

func main() {
	var (
		help = flag.Bool("h", false, usage)
		cid = flag.String("c", "", "hex encoded chainid for the entry")
		server = flag.String("s", "localhost:8083", "path to the factomclient")
		eids extids
	)
	
	flag.Var(&eids, "e", "external id for the entry")
	flag.Parse()
	
	if *help {
		fmt.Println(usage)
		return
	}
	
	data := make([]byte, 1024)
	n, _ := os.Stdin.Read(data)
	data = data[:n]
	
	e := new(entry)
	e.ChainID = *cid
	for _, v := range eids {
		e.ExtIDs = append(e.ExtIDs, string(v))
	}
	e.Data = string(data)
	
	body, err := json.Marshal(e)
	if err != nil {
		return
	}
	
	_, err = http.Post(*server, "application/json", bytes.NewReader(body))
	if err != nil {
		return
	}
	
	return
}
