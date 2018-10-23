package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/FactomProject/factom"
)

// 'add' commands: actually submit requests to a factomd instance

var addIdentityChain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addidentitychain [-fq] [-n NAME1 -n NAME2] [-k PUBKEY1 -k PUBKEY2] [-CET] ECADDRESS"
	cmd.description = "Create a new Identity Chain. Use the Entry Credits from the specified address." +
		"Optional output flags: -C ChainID. -E EntryHash. -T TxID."
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
	help.Add("createidentity", cmd)
	return cmd
}()

var addIdentityKeyReplacement = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addidentitykeyreplacement [-fq] [-c CHAINID | -n NAME1 -n NAME2 ... -n NAMEN]" +
		" --oldkey PUBKEY --newkey PUBKEY --signerkey PUBKEY ECADDRESS [-CET]"
	cmd.description = "Create a new Identity Key Replacement Entry using the Entry Credits from the specified address." +
		" The oldkey is replaced by the newkey, and signerkey (same or higher priority as" +
		" oldkey) authorizes the replacement. Optional output flags: -C ChainID. -E EntryHash. -T TxID."
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
			identityChainID = nametoid(nameCollector)
		} else {
			fmt.Println(cmd.helpMsg)
			return
		}

		signerKey, err := factom.FetchIdentityKey(*signerKeyString)
		if err != nil {
			errorln(fmt.Errorf("Failed to fetch signer key from wallet"))
		}

		e, err := factom.NewIdentityKeyReplacementEntry(identityChainID, *oldKeyString, *newKeyString, signerKey)
		if err != nil {
			errorln("Error composing identity key replacement entry: ", err.Error())
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
	help.Add("addidentitykeyreplacement", cmd)
	return cmd
}()

// 'compose' commands: returns the curl commands needed to be issued

var composeIdentityChain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli composeidentitychain [-f] [-n NAME1 -n NAME2] [-k PUBKEY1 -k PUBKEY2] ECADDRESS"
	cmd.description = "Create API calls to create a new Factom Identity Chain. Use the Entry Credits from the" +
		" specified address."
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
	help.Add("composeidentitychain", cmd)
	return cmd
}()

var composeIdentityKeyReplacement = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli composeidentitykeyreplacement [-f] [-c CHAINID | -n NAME1 -n NAME2 ... -n NAMEN]" +
		" --oldkey PUBKEY --newkey PUBKEY --signerkey PUBKEY ECADDRESS"
	cmd.description = "Create API calls to create a new Identity key replacement entry using the Entry Credits from" +
		" the specified address. The oldkey is replaced by the newkey, and signerkey (same or higher priority as" +
		" oldkey) authorizes the replacement."
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
			identityChainID = nametoid(nameCollector)
		} else {
			fmt.Println(cmd.helpMsg)
			return
		}

		signerKey, err := factom.FetchIdentityKey(*signerKeyString)
		if err != nil {
			errorln(fmt.Errorf("Failed to fetch signer key from wallet"))
		}

		e, err := factom.NewIdentityKeyReplacementEntry(identityChainID, *oldKeyString, *newKeyString, signerKey)
		if err != nil {
			errorln("Error composing identity key replacement entry: ", err.Error())
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
	help.Add("composeidentitykeyreplacement", cmd)
	return cmd
}()

// Other commands

var getIdentityKeysAtHeight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli getidentitykeysatheight [-c CHAINID | -n NAME1 -n NAME2 ... -n NAMEN] HEIGHT"
	cmd.description = "Gets the set of identity public keys that were valid for the given identity chain at the" +
		" specified height."
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
			identityChainID = nametoid(nameCollector)
		} else {
			fmt.Println(cmd.helpMsg)
			return
		}

		i := factom.Identity{}
		i.ChainID = identityChainID
		keys, err := i.GetKeysAtHeight(int64(height))
		if err != nil {
			errorln(err)
		}

		for _, k := range keys {
			fmt.Println(k)
		}
	}
	help.Add("getidentitykeysatheight", cmd)
	return cmd
}()
