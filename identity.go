// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
	"github.com/posener/complete"
)

var identity = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity addchain|addkeyreplacement|addattribute|addattributeendorsement|composechain|" +
		"composekeyreplacement|composeattribute|composeattributeendorsement|getactivekeys|getactivekeysatheight"
	cmd.description = "Create/manage Factom Identity Chains, their currently valid keys, attributes, and" +
		" attribute endorsements"
	cmd.completion = complete.Command{
		Sub: complete.Commands{
			"addchain":                    addIdentityChain.completion,
			"addkeyreplacement":           addIdentityKeyReplacement.completion,
			"addattribute":                addIdentityAttribute.completion,
			"addattributeendorsement":     addIdentityAttributeEndorsement.completion,
			"composechain":                composeIdentityChain.completion,
			"composekeyreplacement":       composeIdentityKeyReplacement.completion,
			"composeattribute":            composeIdentityAttribute.completion,
			"composeattributeendorsement": composeIdentityAttributeEndorsement.completion,
			"getactivekeys":               getActiveIdentityKeys.completion,
			"getactivekeysatheight":       getActiveIdentityKeysAtHeight.completion,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		c := cli.New()
		c.Handle("addchain", addIdentityChain)
		c.Handle("addkeyreplacement", addIdentityKeyReplacement)
		c.Handle("addattribute", addIdentityAttribute)
		c.Handle("addattributeendorsement", addIdentityAttributeEndorsement)
		c.Handle("composechain", composeIdentityChain)
		c.Handle("composekeyreplacement", composeIdentityKeyReplacement)
		c.Handle("composeattribute", composeIdentityAttribute)
		c.Handle("composeattributeendorsement", composeIdentityAttributeEndorsement)
		c.Handle("getactivekeys", getActiveIdentityKeys)
		c.Handle("getactivekeysatheight", getActiveIdentityKeysAtHeight)
		c.HandleDefaultFunc(func(args []string) {
			fmt.Println(cmd.helpMsg)
		})
		c.Execute(args)
	}
	help.Add("identity", cmd)
	return cmd
}()

// 'add' commands: actually submit requests to a factomd instance

var addIdentityChain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity addchain [-fq] [-n NAME1 -n NAME2] [-k PUBKEY1 -k PUBKEY2] [-CET] ECADDRESS"
	cmd.description = "Create a new Identity Chain. Use the Entry Credits from the specified address." +
		"Optional output flags: -C ChainID. -E EntryHash. -T TxID."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,
			"-q": complete.PredictNothing,
			"-n": complete.PredictAnything,
			"-k": complete.PredictAnything,
			"-C": complete.PredictNothing,
			"-E": complete.PredictNothing,
			"-T": complete.PredictNothing,
		},
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		var (
			nAscii namesASCII
			kAscii keysASCII
		)
		os.Args = args
		nameCollector = make([][]byte, 0)
		flag.Var(&nAscii, "n", "Identity Chain name element in ascii. Also is extid of First Entry")
		flag.Var(&kAscii, "k", "A public key string for the identity (decreasing order of priority)")
		fflag := flag.Bool(
			"f",
			false,
			"force the chain to commit and reveal without waiting on any acknowledgement checks",
		)
		cdisp := flag.Bool("C", false, "display only the ChainID")
		edisp := flag.Bool("E", false, "display only the Entry Hash")
		tdisp := flag.Bool("T", false, "display only the TxID")
		qflag := flag.Bool("q", false, "quiet mode; no output")
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		// display normal output iff no display flags are set and quiet is unspecified
		display := true
		if *tdisp || *cdisp || *edisp || *qflag {
			display = false
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}

		c, err := factom.NewIdentityChain(nAscii, kAscii)
		if err != nil {
			errorln("Error composing identity chain struct: ", err.Error())
			return
		}

		if !*fflag {
			if factom.ChainExists(c.ChainID) {
				errorln("Chain", c.ChainID, "already exists")
				return
			}

			// check ec address balance
			balance, err := factom.GetECBalance(ecpub)
			if err != nil {
				errorln(err)
				return
			}
			if cost, err := factom.EntryCost(c.FirstEntry); err != nil {
				errorln(err)
				return
			} else if balance < int64(cost)+10 {
				errorln("Not enough Entry Credits")
				return
			}
		}

		// commit the chain
		var repeated bool
		txid, err := factom.CommitChain(c, ec)
		if err != nil {
			if len(err.Error()) > 15 && err.Error()[:15] != "Repeated Commit" {
				errorln(err)
				return
			}

			fmt.Println("Repeated Commit: A commit with equal or greater payment already exists, skipping commit")
			repeated = true
		}

		if !repeated {
			if display {
				fmt.Println("CommitTxID:", txid)
			} else if *tdisp {
				fmt.Println(txid)
			}

			if !*fflag {
				if _, err := waitOnCommitAck(txid); err != nil {
					errorln(err)
					return
				}
			}
		}

		// reveal chain
		hash, err := factom.RevealChain(c)
		if err != nil {
			errorln(err)
			return
		}
		if display {
			fmt.Println("ChainID:", c.ChainID)
			fmt.Println("Entryhash:", hash)
		} else if *cdisp {
			fmt.Println(c.ChainID)
		} else if *edisp {
			fmt.Println(hash)
		}

		if !*fflag {
			if _, err := waitOnRevealAck(hash); err != nil {
				errorln(err)
				return
			}
		}
	}
	help.Add("identity addchain", cmd)
	return cmd
}()

var addIdentityKeyReplacement = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity addkeyreplacement [-fq] [-c CHAINID | -n NAME1 -n NAME2 ... -n NAMEN]" +
		" --oldkey PUBKEY --newkey PUBKEY --signerkey PUBKEY ECADDRESS [-CET]"
	cmd.description = "Create a new Identity Key Replacement Entry using the Entry Credits from the specified address." +
		" The oldkey is replaced by the newkey, and signerkey (same or higher priority as" +
		" oldkey) authorizes the replacement. Optional output flags: -C ChainID. -E EntryHash. -T TxID."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f":          complete.PredictNothing,
			"-q":          complete.PredictNothing,
			"-c":          complete.PredictAnything,
			"-n":          complete.PredictAnything,
			"--oldkey":    predictIdentityKey,
			"--newkey":    complete.PredictAnything,
			"--signerkey": predictIdentityKey,
			"-C":          complete.PredictNothing,
			"-E":          complete.PredictNothing,
			"-T":          complete.PredictNothing,
		},
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			name            namesASCII
			cid             = flag.String("c", "", "hex encoded chainid for the identity of interest")
			oldKeyString    = flag.String("oldkey", "", "identity public key to be replaced")
			newKeyString    = flag.String("newkey", "", "identity public key to take its place")
			signerKeyString = flag.String("signerkey", "", "identity public key to authorize the replacement"+
				"Must be the same or higher priority than the key being replaced. Key must be stored in the wallet)")
		)

		// -n names
		nameCollector = make([][]byte, 0)
		flag.Var(&name, "n", "an element of the identity's name (used if no ChainID is provided with -c)")

		// -f force
		fflag := flag.Bool(
			"f",
			false,
			"force the entry to commit and reveal without waiting on any"+
				" acknowledgement checks",
		)

		// -CET display flags
		cdisp := flag.Bool("C", false, "display only the ChainID")
		edisp := flag.Bool("E", false, "display only the Entry Hash")
		tdisp := flag.Bool("T", false, "display only the TxID")

		// -q quiet flags
		qflag := flag.Bool("q", false, "quiet mode; no output")

		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		// display normal output iff no display flags are set and quiet is unspecified
		display := true
		if *cdisp || *edisp || *tdisp || *qflag {
			display = false
		}

		// set the chainid from -c or from -n
		var identityChainID string
		if *cid != "" {
			identityChainID = *cid
		} else if len(nameCollector) != 0 {
			nameCollector = append([][]byte{[]byte("IdentityChain")}, nameCollector...)
			identityChainID = nametoid(nameCollector)
		} else {
			fmt.Println(cmd.helpMsg)
			return
		}

		signerKey, err := factom.FetchIdentityKey(*signerKeyString)
		if err != nil {
			errorln(fmt.Errorf("Failed to fetch signer key from wallet"))
			return
		}

		e, err := factom.NewIdentityKeyReplacementEntry(identityChainID, *oldKeyString, *newKeyString, signerKey)
		if err != nil {
			errorln("Error composing identity key replacement entry: ", err.Error())
			return
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}

		if !*fflag {
			if !factom.ChainExists(e.ChainID) {
				errorln("Chain", e.ChainID, "was not found")
				return
			}

			// check ec address balance
			balance, err := factom.GetECBalance(ecpub)
			if err != nil {
				errorln(err)
				return
			}
			if cost, err := factom.EntryCost(e); err != nil {
				errorln(err)
				return
			} else if balance < int64(cost) {
				errorln("Not enough Entry Credits")
				return
			}
		}

		// commit entry
		var repeated bool
		txid, err := factom.CommitEntry(e, ec)
		if err != nil {
			if len(err.Error()) > 15 && err.Error()[:15] != "Repeated Commit" {
				errorln(err)
				return
			}

			fmt.Println("Repeated Commit: A commit with equal or greater payment already exists, skipping commit")
			repeated = true
		}

		if !repeated {
			if display {
				fmt.Println("CommitTxID:", txid)
			} else if *tdisp {
				fmt.Println(txid)
			}

			if !*fflag {
				if _, err := waitOnCommitAck(txid); err != nil {
					errorln(err)
					return
				}
			}
		}
		// reveal entry
		hash, err := factom.RevealEntry(e)
		if err != nil {
			errorln(err)
			return
		}
		if !*fflag {
			if _, err := waitOnRevealAck(hash); err != nil {
				errorln(err)
				return
			}
		}
		if display {
			fmt.Println("ChainID:", e.ChainID)
			fmt.Println("Entryhash:", hash)
		} else if *cdisp {
			fmt.Println(e.ChainID)
		} else if *edisp {
			fmt.Println(hash)
		}

	}
	help.Add("identity addkeyreplacement", cmd)
	return cmd
}()

var addIdentityAttribute = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity addattribute [-fq] -c CHAINID -creceiver CHAINID -csigner CHAINID" +
		" -signerkey PUBKEY -attribute ATTRIBUTE_JSON_ARRAY ECADDRESS [-CET]"
	cmd.description = "Create a new Identity Attribute Entry using the Entry Credits from the specified address." +
		" Optional output flags: -C ChainID. -E EntryHash. -T TxID."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f":         complete.PredictNothing,
			"-q":         complete.PredictNothing,
			"-c":         complete.PredictAnything,
			"-creceiver": complete.PredictAnything,
			"-csigner":   complete.PredictAnything,
			"-signerkey": predictIdentityKey,
			"-attribute": complete.PredictAnything,
			"-C":         complete.PredictNothing,
			"-E":         complete.PredictNothing,
			"-T":         complete.PredictNothing,
		},
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args

		c := flag.String("c", "", "hex encoded chainid for where the attribute entry is written")
		cReceiver := flag.String("creceiver", "", "hex encoded chainid for the identity receiving the attribute")
		cSigner := flag.String("csigner", "", "hex encoded chainid for the identity signing/giving the attribute")
		signerKeyString := flag.String("signerkey", "", "public identity key that signs the attribute entry"+
			" (must be stored in wallet and should be currently valid for signer's identity)")
		attributesJSON := flag.String("attribute", "", "JSON array describing the attribute to assign"+
			" (must be in the format of '[{\"key\":KEY,\"value\":VALUE},{\"key\":KEY,\"value\":VALUE},...]'")

		// -f force
		fflag := flag.Bool(
			"f",
			false,
			"force the entry to commit and reveal without waiting on any"+
				" acknowledgement checks",
		)

		// -CET display flags
		cdisp := flag.Bool("C", false, "display only the ChainID")
		edisp := flag.Bool("E", false, "display only the Entry Hash")
		tdisp := flag.Bool("T", false, "display only the TxID")

		// -q quiet flags
		qflag := flag.Bool("q", false, "quiet mode; no output")

		flag.Parse()

		// get EC key pair from wallet
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}

		// get signer identity key pair from wallet
		signerKey, err := factom.FetchIdentityKey(*signerKeyString)
		if err != nil {
			errorln(err)
			return
		}

		// check for missing/invalid chain id params
		if len(*cReceiver) != 64 {
			errorln("Missing/invalid receiver ChainID (-creceiver)")
			fmt.Println(cmd.helpMsg)
			return
		} else if len(*c) != 64 {
			errorln("Missing/invalid destination ChainID (-c)")
			fmt.Println(cmd.helpMsg)
			return
		} else if len(*cSigner) != 64 {
			errorln("Missing/invalid signer ChainID (-csigner)")
			fmt.Println(cmd.helpMsg)
			return
		}

		// check that attributes array can be unmarshalled and contains no nil keys or values
		// TODO: move this validation into factom.NewIdentityAttributeEntry() instead
		var attributes []factom.IdentityAttribute
		err = json.Unmarshal([]byte(*attributesJSON), &attributes)
		if err != nil {
			errorln("Invalid attribute array: ", err)
			fmt.Println(cmd.helpMsg)
			return
		}
		for _, attribute := range attributes {
			if attribute.Key == nil {
				errorln("All attribute keys must not be nil")
				fmt.Println(cmd.helpMsg)
				return
			}
			if attribute.Value == nil {
				errorln("All attribute values must not be nil")
				fmt.Println(cmd.helpMsg)
				return
			}
		}

		e := factom.NewIdentityAttributeEntry(*cReceiver, *c, *attributesJSON, signerKey, *cSigner)

		// display normal output iff no display flags are set and quiet is unspecified
		display := true
		if *cdisp || *edisp || *tdisp || *qflag {
			display = false
		}

		if !*fflag {
			if !factom.ChainExists(e.ChainID) {
				errorln("Destination Chain", e.ChainID, "was not found")
				return
			}

			// check ec address balance
			balance, err := factom.GetECBalance(ecpub)
			if err != nil {
				errorln(err)
				return
			}
			if cost, err := factom.EntryCost(e); err != nil {
				errorln(err)
				return
			} else if balance < int64(cost) {
				errorln("Not enough Entry Credits")
				return
			}
		}

		// commit entry
		var repeated bool
		txid, err := factom.CommitEntry(e, ec)
		if err != nil {
			if len(err.Error()) > 15 && err.Error()[:15] != "Repeated Commit" {
				errorln(err)
				return
			}

			fmt.Println("Repeated Commit: A commit with equal or greater payment already exists, skipping commit")
			repeated = true
		}

		if !repeated {
			if display {
				fmt.Println("CommitTxID:", txid)
			} else if *tdisp {
				fmt.Println(txid)
			}

			if !*fflag {
				if _, err := waitOnCommitAck(txid); err != nil {
					errorln(err)
					return
				}
			}
		}
		// reveal entry
		hash, err := factom.RevealEntry(e)
		if err != nil {
			errorln(err)
			return
		}
		if !*fflag {
			if _, err := waitOnRevealAck(hash); err != nil {
				errorln(err)
				return
			}
		}
		if display {
			fmt.Println("ChainID:", e.ChainID)
			fmt.Println("Entryhash:", hash)
		} else if *cdisp {
			fmt.Println(e.ChainID)
		} else if *edisp {
			fmt.Println(hash)
		}

	}
	help.Add("identity addattribute", cmd)
	return cmd
}()

var addIdentityAttributeEndorsement = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity addattributeendorsement [-fq] -c CHAINID -csigner CHAINID -signerkey PUBKEY" +
		" -entryhash ENTRYHASH ECADDRESS [-CET]"
	cmd.description = "Create a new Endorsement Entry for the Identity Attribute at the given entry hash. Uses the" +
		" Entry Credits from the specified address. Optional output flags: -C ChainID. -E EntryHash. -T TxID."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f":         complete.PredictNothing,
			"-q":         complete.PredictNothing,
			"-c":         complete.PredictAnything,
			"-csigner":   complete.PredictAnything,
			"-signerkey": predictIdentityKey,
			"-entryhash": complete.PredictAnything,
			"-C":         complete.PredictNothing,
			"-E":         complete.PredictNothing,
			"-T":         complete.PredictNothing,
		},
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args

		c := flag.String("c", "", "hex encoded chainid for where the endorsement entry is written")
		cSigner := flag.String("csigner", "", "hex encoded chainid for the identity signing/giving the endorsement")
		signerKeyString := flag.String("signerkey", "", "public identity key that signs the endorsement entry"+
			" (must be stored in wallet and should be currently valid for signer's identity)")
		entryHash := flag.String("entryhash", "", "hex encoded entry hash for the attribute entry being endorsed")

		// -f force
		fflag := flag.Bool(
			"f",
			false,
			"force the entry to commit and reveal without waiting on any acknowledgement checks",
		)

		// -CET display flags
		cdisp := flag.Bool("C", false, "display only the ChainID")
		edisp := flag.Bool("E", false, "display only the Entry Hash")
		tdisp := flag.Bool("T", false, "display only the TxID")

		// -q quiet flags
		qflag := flag.Bool("q", false, "quiet mode; no output")

		flag.Parse()

		// get EC key pair from wallet
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}

		// get signer identity key pair from wallet
		signerKey, err := factom.FetchIdentityKey(*signerKeyString)
		if err != nil {
			errorln(err)
			return
		}

		// check for missing/invalid chain id and entry hash params
		if len(*c) != 64 {
			errorln("Missing/invalid destination ChainID (-c)")
			fmt.Println(cmd.helpMsg)
			return
		} else if len(*cSigner) != 64 {
			errorln("Missing/invalid signer ChainID (-csigner)")
			fmt.Println(cmd.helpMsg)
			return
		} else if len(*entryHash) != 64 {
			errorln("Missing/invalid entry hash to endorse (-entryhash)")
			fmt.Println(cmd.helpMsg)
			return
		}

		e := factom.NewIdentityAttributeEndorsementEntry(*c, *entryHash, signerKey, *cSigner)

		// display normal output iff no display flags are set and quiet is unspecified
		display := true
		if *cdisp || *edisp || *tdisp || *qflag {
			display = false
		}

		if !*fflag {
			if !factom.ChainExists(e.ChainID) {
				errorln("Destination Chain", e.ChainID, "was not found")
				return
			}

			// check ec address balance
			balance, err := factom.GetECBalance(ecpub)
			if err != nil {
				errorln(err)
				return
			}
			if cost, err := factom.EntryCost(e); err != nil {
				errorln(err)
				return
			} else if balance < int64(cost) {
				errorln("Not enough Entry Credits")
				return
			}
		}

		// commit entry
		var repeated bool
		txid, err := factom.CommitEntry(e, ec)
		if err != nil {
			if len(err.Error()) > 15 && err.Error()[:15] != "Repeated Commit" {
				errorln(err)
				return
			}

			fmt.Println("Repeated Commit: A commit with equal or greater payment already exists, skipping commit")
			repeated = true
		}

		if !repeated {
			if display {
				fmt.Println("CommitTxID:", txid)
			} else if *tdisp {
				fmt.Println(txid)
			}

			if !*fflag {
				if _, err := waitOnCommitAck(txid); err != nil {
					errorln(err)
					return
				}
			}
		}
		// reveal entry
		hash, err := factom.RevealEntry(e)
		if err != nil {
			errorln(err)
			return
		}
		if !*fflag {
			if _, err := waitOnRevealAck(hash); err != nil {
				errorln(err)
				return
			}
		}
		if display {
			fmt.Println("ChainID:", e.ChainID)
			fmt.Println("Entryhash:", hash)
		} else if *cdisp {
			fmt.Println(e.ChainID)
		} else if *edisp {
			fmt.Println(hash)
		}

	}
	help.Add("identity addattributeendorsement", cmd)
	return cmd
}()

// 'compose' commands: returns the curl commands needed to be issued

var composeIdentityChain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity composechain [-f] [-n NAME1 -n NAME2] [-k PUBKEY1 -k PUBKEY2] ECADDRESS"
	cmd.description = "Create API calls to create a new Factom Identity Chain. Use the Entry Credits from the" +
		" specified address."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,
			"-n": complete.PredictAnything,
			"-k": complete.PredictAnything,
		},
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		var (
			nAscii namesASCII
			kAscii keysASCII
		)
		os.Args = args
		nameCollector = make([][]byte, 0)
		flag.Var(&nAscii, "n", "Identity name part in ascii. Also is extid of First Entry")
		flag.Var(&kAscii, "k", "A public key string for the identity (decreasing order of priority)")
		fflag := flag.Bool(
			"f",
			false,
			"force the chain to commit and reveal without waiting on any acknowledgement checks",
		)
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		c, err := factom.NewIdentityChain(nAscii, kAscii)
		if err != nil {
			errorln("Error composing identity chain struct: ", err.Error())
			return
		}

		commit, reveal, err := factom.WalletComposeChainCommitReveal(c, ecpub, *fflag)
		if err != nil {
			errorln(err)
			return
		}

		factomdServer := GetFactomdServer()

		fmt.Println(
			"curl -X POST --data-binary",
			"'"+commit.String()+"'",
			"-H 'content-type:text/plain;' http://"+factomdServer+"/v2",
		)
		fmt.Println(
			"curl -X POST --data-binary",
			"'"+reveal.String()+"'",
			"-H 'content-type:text/plain;' http://"+factomdServer+"/v2",
		)

	}
	help.Add("identity composechain", cmd)
	return cmd
}()

var composeIdentityKeyReplacement = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity composekeyreplacement [-f] [-c CHAINID | -n NAME1 -n NAME2 ... -n NAMEN]" +
		" --oldkey PUBKEY --newkey PUBKEY --signerkey PUBKEY ECADDRESS"
	cmd.description = "Create API calls to create a new Identity key replacement entry using the Entry Credits from" +
		" the specified address. The oldkey is replaced by the newkey, and signerkey (same or higher priority as" +
		" oldkey) authorizes the replacement."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f":          complete.PredictNothing,
			"-c":          complete.PredictAnything,
			"-n":          complete.PredictAnything,
			"--oldkey":    predictIdentityKey,
			"--newkey":    complete.PredictAnything,
			"--signerkey": predictIdentityKey,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			name            namesASCII
			cid             = flag.String("c", "", "hex encoded chainid for the identity of interest")
			oldKeyString    = flag.String("oldkey", "", "identity public key to be replaced")
			newKeyString    = flag.String("newkey", "", "identity public key to take its place")
			signerKeyString = flag.String("signerkey", "", "identity public key to authorize the replacement"+
				"Must be the same or higher priority than the key being replaced. Key must be stored in the wallet)")
		)

		// -n names
		nameCollector = make([][]byte, 0)
		flag.Var(&name, "n", "an element of the identity's name (used if no ChainID is provided with -c)")

		// -f force
		fflag := flag.Bool(
			"f",
			false,
			"force the entry to commit and reveal without waiting on any"+
				" acknowledgement checks",
		)

		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		// set the chainid from -c or from -n
		var identityChainID string
		if *cid != "" {
			identityChainID = *cid
		} else if len(nameCollector) != 0 {
			nameCollector = append([][]byte{[]byte("IdentityChain")}, nameCollector...)
			identityChainID = nametoid(nameCollector)
		} else {
			fmt.Println(cmd.helpMsg)
			return
		}

		signerKey, err := factom.FetchIdentityKey(*signerKeyString)
		if err != nil {
			errorln(fmt.Errorf("Failed to fetch signer key from wallet"))
			return
		}

		e, err := factom.NewIdentityKeyReplacementEntry(identityChainID, *oldKeyString, *newKeyString, signerKey)
		if err != nil {
			errorln("Error composing identity key replacement entry: ", err.Error())
			return
		}

		commit, reveal, err := factom.WalletComposeEntryCommitReveal(e, ecpub, *fflag)
		if err != nil {
			errorln(err)
			return
		}

		factomdServer := GetFactomdServer()

		fmt.Println(
			"curl -X POST --data-binary",
			"'"+commit.String()+"'",
			"-H 'content-type:text/plain;' http://"+factomdServer+"/v2",
		)
		fmt.Println(
			"curl -X POST --data-binary",
			"'"+reveal.String()+"'",
			"-H 'content-type:text/plain;' http://"+factomdServer+"/v2",
		)
	}
	help.Add("identity composekeyreplacement", cmd)
	return cmd
}()

var composeIdentityAttribute = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity composeattribute [-f] -c CHAINID -creceiver CHAINID -csigner CHAINID" +
		" -signerkey PUBKEY -attribute ATTRIBUTE_JSON_ARRAY ECADDRESS"
	cmd.description = "Create API calls to create a new Identity Attribute Entry using the Entry Credits from the specified address."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f":         complete.PredictNothing,
			"-c":         complete.PredictAnything,
			"-creceiver": complete.PredictAnything,
			"-csigner":   complete.PredictAnything,
			"-signerkey": predictIdentityKey,
			"-attribute": complete.PredictAnything,
		},
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args

		c := flag.String("c", "", "hex encoded chainid for where the attribute entry is written")
		cReceiver := flag.String("creceiver", "", "hex encoded chainid for the identity receiving the attribute")
		cSigner := flag.String("csigner", "", "hex encoded chainid for the identity signing/giving the attribute")
		signerKeyString := flag.String("signerkey", "", "public identity key that signs the attribute entry"+
			" (must be stored in wallet and should be currently valid for signer's identity)")
		attributesJSON := flag.String("attribute", "", "JSON array describing the attribute to assign"+
			" (must be in the format of '[{\"key\":KEY,\"value\":VALUE},{\"key\":KEY,\"value\":VALUE},...]'")

		// -f force
		fflag := flag.Bool(
			"f",
			false,
			"force the entry to commit and reveal without waiting on any"+
				" acknowledgement checks",
		)

		flag.Parse()

		// get EC key pair from wallet
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		// get signer identity key pair from wallet
		signerKey, err := factom.FetchIdentityKey(*signerKeyString)
		if err != nil {
			errorln(err)
			return
		}

		// check for missing/invalid chain id params
		if len(*cReceiver) != 64 {
			errorln("Missing/invalid receiver ChainID (-creceiver)")
			fmt.Println(cmd.helpMsg)
			return
		} else if len(*c) != 64 {
			errorln("Missing/invalid destination ChainID (-c)")
			fmt.Println(cmd.helpMsg)
			return
		} else if len(*cSigner) != 64 {
			errorln("Missing/invalid signer ChainID (-csigner)")
			fmt.Println(cmd.helpMsg)
			return
		}

		// check that attributes array can be unmarshalled and contains no nil keys or values
		// TODO: move this validation into factom.NewIdentityAttributeEntry() instead
		var attributes []factom.IdentityAttribute
		err = json.Unmarshal([]byte(*attributesJSON), &attributes)
		if err != nil {
			errorln("Invalid attribute array: ", err)
			fmt.Println(cmd.helpMsg)
			return
		}
		for _, attribute := range attributes {
			if attribute.Key == nil {
				errorln("All attribute keys must not be nil")
				fmt.Println(cmd.helpMsg)
				return
			}
			if attribute.Value == nil {
				errorln("All attribute values must not be nil")
				fmt.Println(cmd.helpMsg)
				return
			}
		}

		e := factom.NewIdentityAttributeEntry(*cReceiver, *c, *attributesJSON, signerKey, *cSigner)

		commit, reveal, err := factom.WalletComposeEntryCommitReveal(e, ecpub, *fflag)
		if err != nil {
			errorln(err)
			return
		}

		factomdServer := GetFactomdServer()

		fmt.Println(
			"curl -X POST --data-binary",
			"'"+commit.String()+"'",
			"-H 'content-type:text/plain;' http://"+factomdServer+"/v2",
		)
		fmt.Println(
			"curl -X POST --data-binary",
			"'"+reveal.String()+"'",
			"-H 'content-type:text/plain;' http://"+factomdServer+"/v2",
		)
	}
	help.Add("identity composeattribute", cmd)
	return cmd
}()

var composeIdentityAttributeEndorsement = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity composeattributeendorsement [-f] -c CHAINID -csigner CHAINID -signerkey PUBKEY" +
		" -entryhash ENTRYHASH ECADDRESS"
	cmd.description = "Compose API calls to create a new Endorsement Entry for the Identity Attribute at the given" +
		" entry hash. Uses the Entry Credits from the specified address."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f":         complete.PredictNothing,
			"-c":         complete.PredictAnything,
			"-csigner":   complete.PredictAnything,
			"-signerkey": predictIdentityKey,
			"-entryhash": complete.PredictAnything,
		},
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args

		c := flag.String("c", "", "hex encoded chainid for where the endorsement entry is written")
		cSigner := flag.String("csigner", "", "hex encoded chainid for the identity signing/giving the endorsement")
		signerKeyString := flag.String("signerkey", "", "public identity key that signs the endorsement entry"+
			" (must be stored in wallet and should be currently valid for signer's identity)")
		entryHash := flag.String("entryhash", "", "hex encoded entry hash for the attribute entry being endorsed")

		// -f force
		fflag := flag.Bool(
			"f",
			false,
			"force the entry to commit and reveal without waiting on any acknowledgement checks",
		)

		flag.Parse()

		// get EC key pair from wallet
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		// get signer identity key pair from wallet
		signerKey, err := factom.FetchIdentityKey(*signerKeyString)
		if err != nil {
			errorln(err)
			return
		}

		// check for missing/invalid chain id and entry hash params
		if len(*c) != 64 {
			errorln("Missing/invalid destination ChainID (-c)")
			fmt.Println(cmd.helpMsg)
			return
		} else if len(*cSigner) != 64 {
			errorln("Missing/invalid signer ChainID (-csigner)")
			fmt.Println(cmd.helpMsg)
			return
		} else if len(*entryHash) != 64 {
			errorln("Missing/invalid entry hash to endorse (-entryhash)")
			fmt.Println(cmd.helpMsg)
			return
		}

		e := factom.NewIdentityAttributeEndorsementEntry(*c, *entryHash, signerKey, *cSigner)

		commit, reveal, err := factom.WalletComposeEntryCommitReveal(e, ecpub, *fflag)
		if err != nil {
			errorln(err)
			return
		}

		factomdServer := GetFactomdServer()

		fmt.Println(
			"curl -X POST --data-binary",
			"'"+commit.String()+"'",
			"-H 'content-type:text/plain;' http://"+factomdServer+"/v2",
		)
		fmt.Println(
			"curl -X POST --data-binary",
			"'"+reveal.String()+"'",
			"-H 'content-type:text/plain;' http://"+factomdServer+"/v2",
		)
	}
	help.Add("identity composeattributeendorsement", cmd)
	return cmd
}()

// Other commands

var getActiveIdentityKeys = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity getactivekeys [-c CHAINID | -n NAME1 -n NAME2 ... -n NAMEN]"
	cmd.description = "Gets the set of identity public keys that are active for the given identity chain at the" +
		" highest known block height."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-c": complete.PredictAnything,
			"-n": complete.PredictAnything,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			name namesASCII
			cid  = flag.String("c", "", "hex encoded chainid for the identity of interest")
		)

		// -n names
		nameCollector = make([][]byte, 0)
		flag.Var(&name, "n", "an element of the identity's name (used if no ChainID is provided with -c)")

		flag.Parse()

		// set the chainid from -c or from -n
		var identityChainID string
		if *cid != "" {
			identityChainID = *cid
		} else if len(nameCollector) != 0 {
			nameCollector = append([][]byte{[]byte("IdentityChain")}, nameCollector...)
			identityChainID = nametoid(nameCollector)
		} else {
			fmt.Println(cmd.helpMsg)
			return
		}

		keys, _, err := factom.GetActiveIdentityKeys(identityChainID)
		if err != nil {
			errorln(err)
			return
		}

		for _, k := range keys {
			fmt.Println(k)
		}
	}
	help.Add("identity getactivekeys", cmd)
	return cmd
}()

var getActiveIdentityKeysAtHeight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli identity getactivekeysatheight [-c CHAINID | -n NAME1 -n NAME2 ... -n NAMEN] HEIGHT"
	cmd.description = "Gets the set of identity public keys that were valid for the given identity chain at the" +
		" specified block height."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-c": complete.PredictAnything,
			"-n": complete.PredictAnything,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			name namesASCII
			cid  = flag.String("c", "", "hex encoded chainid for the identity of interest")
		)

		// -n names
		nameCollector = make([][]byte, 0)
		flag.Var(&name, "n", "an element of the identity's name (used if no ChainID is provided with -c)")

		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		height, err := strconv.Atoi(args[0])
		if err != nil || height < 0 {
			errorln("Height must be a valid non-negative integer")
			return
		}

		// set the chainid from -c or from -n
		var identityChainID string
		if *cid != "" {
			identityChainID = *cid
		} else if len(nameCollector) != 0 {
			nameCollector = append([][]byte{[]byte("IdentityChain")}, nameCollector...)
			identityChainID = nametoid(nameCollector)
		} else {
			fmt.Println(cmd.helpMsg)
			return
		}

		keys, err := factom.GetActiveIdentityKeysAtHeight(identityChainID, int64(height))
		if err != nil {
			errorln(err)
			return
		}

		for _, k := range keys {
			fmt.Println(k)
		}
	}
	help.Add("identity getactivekeysatheight", cmd)
	return cmd
}()
