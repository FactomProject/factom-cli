// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/FactomProject/factom"
)

var addchain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addchain [-fq] [-n NAME1 -n NAME2 -h HEXNAME3" +
		" ] [-CET] ECADDRESS <STDIN>"
	cmd.description = "Create a new Factom Chain. Read data for the First" +
		" Entry from stdin. Use the Entry Credits from the specified address." +
		" -C ChainID. -E EntryHash. -T TxID."
	cmd.execFunc = func(args []string) {
		var (
			eAcii extidsASCII
			eHex  extidsHex
		)
		os.Args = args
		exidCollector = make([][]byte, 0)
		flag.Var(&eAcii, "n", "Chain name element in ascii. Also is extid of"+
			" First Entry")
		flag.Var(&eHex, "h", "Chain name element in hex. Also is extid of"+
			" First Entry")
		fflag := flag.Bool(
			"f",
			false,
			"force the chain to commit and reveal without waiting on any"+
				" acknowledgement checks",
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

		e := new(factom.Entry)

		e.ExtIDs = exidCollector

		// Entry.Content is read from stdin
		if p, err := ioutil.ReadAll(os.Stdin); err != nil {
			errorln(err)
			return
		} else if size := len(p); size > 10240 {
			errorln(fmt.Sprintf("Entry of %d bytes is too large", size))
			return
		} else {
			e.Content = p
		}

		// get the ec address from the wallet
		ec, err := factom.FetchECAddress(ecpub)
		if err != nil {
			errorln(err)
			return
		}

		c := factom.NewChain(e)

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
		txid, err := factom.CommitChain(c, ec)
		if err != nil {
			errorln(err)
			return
		}
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
			if _, err := waitOnRevealAck(txid); err != nil {
				errorln(err)
				return
			}
		}
	}
	help.Add("addchain", cmd)
	return cmd
}()

var composechain = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli composechain [-f] [-n NAME1 -n NAME2" +
		" -h HEXNAME3 ] ECADDRESS <STDIN>"
	cmd.description = "Create API calls to create a new Factom Chain. Read" +
		" data for the First Entry from stdin. Use the Entry Credits from the" +
		" specified address."
	cmd.execFunc = func(args []string) {
		var (
			eAcii extidsASCII
			eHex  extidsHex
		)
		os.Args = args
		exidCollector = make([][]byte, 0)
		flag.Var(&eAcii, "n", "Chain name element in ascii. Also is extid of"+
			" First Entry")
		flag.Var(&eHex, "h", "Chain name element in hex. Also is extid of"+
			" First Entry")
		fflag := flag.Bool(
			"f",
			false,
			"force the chain to commit and reveal without waiting on any"+
				" acknowledgement checks",
		)
		flag.Parse()
		args = flag.Args()

		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}
		ecpub := args[0]

		e := new(factom.Entry)

		e.ExtIDs = exidCollector

		// Entry.Content is read from stdin
		if p, err := ioutil.ReadAll(os.Stdin); err != nil {
			errorln(err)
			return
		} else if size := len(p); size > 10240 {
			errorln(fmt.Sprintf("Entry of %d bytes is too large", size))
			return
		} else {
			e.Content = p
		}

		c := factom.NewChain(e)

		commit, reveal, err := factom.WalletComposeChainCommitReveal(c, ecpub, *fflag)
		if err != nil {
			errorln(err)
			return
		}

		fmt.Println(
			"curl -X POST --data-binary",
			"'"+commit.String()+"'",
			"-H 'content-type:text/plain;' http://localhost:8088/v2",
		)
		fmt.Println(
			"curl -X POST --data-binary",
			"'"+reveal.String()+"'",
			"-H 'content-type:text/plain;' http://localhost:8088/v2",
		)

	}
	help.Add("composechain", cmd)
	return cmd
}()
