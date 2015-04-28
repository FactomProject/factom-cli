package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewECKey(t *testing.T) {
	wallet = os.Getenv("HOME") + "/.factom/ecwallet"

	err := newECKey()
	if err != nil {
		t.Fatal(err)
	}
	
	wfile, err := os.Open(wallet)
	if err != nil {
		t.Fatal(err)
	}
	p, err := ioutil.ReadAll(wfile)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(p))
}

func TestPubKey(t *testing.T) {
	wallet = os.Getenv("HOME") + "/.factom/ecwallet"
	
	err := pubKey()
	if err != nil {
		t.Error(err)
	}
}
