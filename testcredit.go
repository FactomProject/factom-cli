package main

import (
	"encoding/hex"
	"flag"
	"os"

	"github.com/FactomProject/factom"
)

func testcredit(args []string) error {
	os.Args = args
	flag.Parse()
	args = flag.Args()

	var key string
	if p, err := ecPubKey(); err != nil {
		// TODO: should fail only if no wallet file AND no key specified
		return err
	} else {
		key = hex.EncodeToString(p[:])
	}
	if args != nil {
		key = args[0]
	}
	
	if err := factom.TestCredit(key); err != nil {
		return err
	}
	return nil
}
