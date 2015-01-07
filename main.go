package main

import (
	"os"
)

func main() {
	command := "help"
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}
	switch(command) {
	case "help":
		a := "help"
		if len(os.Args) >= 3 {
			a = os.Args[2]
		}
		help(a)			
	case "put":
		put()
	case "get":
		get()
	case "mkchain":
		mkchain()
	}
}

//import (
//	"encoding/hex"
//	"encoding/json"
//	"flag"
//	"fmt"
//	"net/http"
//	"net/url"
//	"os"
//)
//const usage = "factom -s [server] -c [chainid] -e [extid] -e [extid2] <data"
//
//type entry struct {
//	ChainID string
//	ExtIDs  []string
//	Data    string
//}
//
//type extids []string
//
//func (e *extids) String() string {
//	return fmt.Sprint(*e)
//}
//
//func (e *extids) Set(s string) error {
//	*e = append(*e, s)
//	return nil
//}
//
//func main() {
//	var (
//		help = flag.Bool("h", false, usage)
//		cid  = flag.String("c", "", "hex encoded chainid for the entry")
//		serv = flag.String("s", "localhost:8088", "path to the factomclient")
//		eids extids
//	)
//
//	flag.Var(&eids, "e", "external id for the entry")
//	flag.Parse()
//
//	if *help {
//		fmt.Println(usage)
//		return
//	}
//
//	server := "http://" + *serv + "/v1/submitentry"
//
//	fmt.Println(server)
//
//	d := make([]byte, 1024)
//	n, _ := os.Stdin.Read(d)
//	d = d[:n]
//
//	e := new(entry)
//	e.ChainID = *cid
//	for _, v := range eids {
//		e.ExtIDs = append(e.ExtIDs, string(v))
//	}
//	e.Data = hex.EncodeToString(d)
//
//	b, err := json.Marshal(e)
//	if err != nil {
//		return
//	}
//	
//	data := url.Values{
//		"datatype": {"entry"},
//		"format":   {"json"},
//		"entry":    {string(b)},
//	}
//	
//	_, err = http.PostForm(server, data)
//	if err != nil {
//		fmt.Println("Error: ", err)
//	}
//
//	return
//}
