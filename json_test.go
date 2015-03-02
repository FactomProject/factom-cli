package main

import (
	"encoding/json"
	"fmt"
	"io"
//	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

// anonymous declarations to squash import errors for testing
var _ string = fmt.Sprint()
var _ []string = os.Args

func TestBalance(t *testing.T) {
	fmt.Printf("TestJson\n===\n")
	type balance struct {
		Publickey string
		Credits float64
	}
	server := "http://demo.factom.org:8088/v1/creditbalance"
	
	data := url.Values{
		"pubkey": {"wallet"},
	}
	
	resp, err := http.PostForm(server, data)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer resp.Body.Close()
	
	dec := json.NewDecoder(resp.Body)
	for {
		var bal balance
		if err := dec.Decode(&bal); err == io.EOF {
			break
		} else if err != nil {
			t.Errorf(err.Error())
		}
		fmt.Println("Entry Credit Balance:", bal.Credits)
	}
}
