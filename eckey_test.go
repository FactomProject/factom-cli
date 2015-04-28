package main

import (
	"fmt"
	"os"
	"testing"
)

func TestNewECKey(t *testing.T) {
	fmt.Printf("TestNewECKey\n---\n")
	wallet = os.Getenv("HOME") + "/.factom/ecwallet"

	err := newECKey()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println()
}

func TestPubKey(t *testing.T) {
	fmt.Printf("TestPubKey\n---\n")
	wallet = os.Getenv("HOME") + "/.factom/ecwallet"
	
	err := printPubKey()
	if err != nil {
		t.Error(err)
	}
}
