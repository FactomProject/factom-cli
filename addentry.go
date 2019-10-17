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
	"github.com/posener/complete"
)

var addentry = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli addentry [-fq] [-n NAME1 -h HEXNAME2" +
		" ...|-c CHAINID] [-e EXTID1 -e EXTID2 -x HEXEXTID ...] [-CET]" +
		" ECADDRESS <STDIN>"
	cmd.description = "Create a new Factom Entry. Read data for the Entry" +
		" from stdin. Use the Entry Credits from the specified address." +
		" -C ChainID. -E EntryHash. -T TxID."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,
			"-q": complete.PredictNothing,

			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,

			"-C": complete.PredictNothing,
			"-E": complete.PredictNothing,
			"-T": complete.PredictNothing,
		},
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			cid   = flag.String("c", "", "hex encoded chainid for the entry")
			eAcii extidsASCII
			eHex  extidsHex
			nAcii namesASCII
			nHex  namesHex
		)

		// -e -x extids
		exidCollector = make([][]byte, 0)
		flag.Var(&eAcii, "e", "external id for the entry in ascii")
		flag.Var(&eHex, "x", "external id for the entry in hex")

		// -n -h names
		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name element")
		flag.Var(&nHex, "h", "hex binary name element")

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

		e := new(factom.Entry)

		// set the chainid from -c or from -n -h
		if *cid != "" {
			e.ChainID = *cid
		} else if len(nameCollector) != 0 {
			e.ChainID = nametoid(nameCollector)
		} else {
			fmt.Println(cmd.helpMsg)
			return
		}

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
	help.Add("addentry", cmd)
	return cmd
}()

var composeentry = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli composeentry [-f] [-n NAME1 -h HEXNAME2" +
		" ...|-c CHAINID]  [-e EXTID1 -e EXTID2 -x HEXEXTID ...] ECADDRESS" +
		" <STDIN>"
	cmd.description = "Create API calls to create a new Factom Entry. Read" +
		" data for the Entry from stdin. Use the Entry Credits from the" +
		" specified address."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,

			"-c": complete.PredictAnything,
			"-e": complete.PredictAnything,
			"-x": complete.PredictAnything,
			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,
		},
		Args: predictAddress,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		var (
			cid   = flag.String("c", "", "hex encoded chainid for the entry")
			eAcii extidsASCII
			eHex  extidsHex
			nAcii namesASCII
			nHex  namesHex
		)

		// -e -x extids
		exidCollector = make([][]byte, 0)
		flag.Var(&eAcii, "e", "external id for the entry in ascii")
		flag.Var(&eHex, "x", "external id for the entry in hex")

		// -n -h names
		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name element")
		flag.Var(&nHex, "h", "hex binary name element")

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

		e := new(factom.Entry)

		// set the chainid from -c or from -n -h
		if *cid != "" {
			e.ChainID = *cid
		} else if len(nameCollector) != 0 {
			e.ChainID = nametoid(nameCollector)
		} else {
			fmt.Println(cmd.helpMsg)
			return
		}

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
	help.Add("composeentry", cmd)
	return cmd
}()
