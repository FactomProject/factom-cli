package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
)

// newIdentityKey generates a new identity key in the wallet
var newIdentityKey = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli newidentitykey"
	cmd.description = "Generate a new identity key in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		k, err := factom.GenerateIdentityKey()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(k.PubString())
	}
	help.Add("newidentitykey", cmd)
	return cmd
}()

// importIdentityKeys imports identity keys from 1 or more secret keys into the wallet
var importIdentityKeys = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli importidentitykeys SECKEY [SECKEY...]"
	cmd.description = "Import one or more identity keys into the wallet from the specified idsec keys"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		keys, err := factom.ImportIdentityKeys(args...)
		if err != nil {
			errorln(err)
			return
		}
		for _, k := range keys {
			fmt.Println(k)
		}
	}
	help.Add("importidentitykeys", cmd)
	return cmd
}()

// exportIdentityKeys lists the identity key pairs (public and private) stored in the wallet
var exportIdentityKeys = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli exportidentitykeys"
	cmd.description = "List the identity key pairs (public and private) stored in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		keys, err := factom.FetchIdentityKeys()
		if err != nil {
			errorln(err)
			return
		}
		for _, k := range keys {
			fmt.Println(k.SecString(), k.PubString())
		}
	}
	help.Add("exportidentitykeys", cmd)
	return cmd
}()

// listIdentityKeys lists the addresses in the wallet
var listIdentityKeys = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli listidentitykeys"
	cmd.description = "List the public identity keys stored in the wallet"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		keys, err := factom.FetchIdentityKeys()
		if err != nil {
			errorln(err)
			return
		}

		for _, k := range keys {
			fmt.Println(k.PubString())
		}
	}
	help.Add("listaddresses", cmd)
	return cmd
}()

// removeIdentityKey removes an identity key pair from the wallet
var removeIdentityKey = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli rmidentitykey PUBKEY"
	cmd.description = "Removes the identity key pair from the wallet for the specified idpub key."
	cmd.execFunc = func(args []string) {
		if len(args) < 2 {
			fmt.Println(cmd.helpMsg)
			return
		}
		pub := args[1]

		err := factom.RemoveIdentityKey(pub)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}
	help.Add("rmidentitykey", cmd)
	return cmd
}()
