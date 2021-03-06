// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/FactomProject/cli"
	"github.com/FactomProject/factom"
	"github.com/posener/complete"
)

var get = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get allentries|authorities|chainhead|" +
		"currentminute|ablock|dblock|eblock|ecblock|fblock|entry|firstentry|" +
		"head|heights|pendingentries|pendingtransactions|raw|tps|walletheight"
	cmd.description = "Get data about Factom Chains, Entries, and Blocks"
	cmd.completion = complete.Command{
		Sub: complete.Commands{
			"allentries":          getAllEntries.completion,
			"authorities":         complete.Command{},
			"chainhead":           getChainHead.completion,
			"currentminute":       getCurrentMinute.completion,
			"ablock":              getABlock.completion,
			"dblock":              getDBlock.completion,
			"eblock":              getEBlock.completion,
			"ecblock":             getECBlock.completion,
			"fblock":              getFBlock.completion,
			"entry":               getEntry.completion,
			"firstentry":          getFirstEntry.completion,
			"head":                getHead.completion,
			"heights":             getHeights.completion,
			"pendingentries":      getPendingEntries.completion,
			"pendingtransactions": getPendingTransactions.completion,
			"raw":                 getraw.completion,
			"tps":                 getTPS.completion,
			"walletheight":        complete.Command{},
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		c := cli.New()
		c.Handle("allentries", getAllEntries)
		c.Handle("authorities", getAuthorities)
		c.Handle("chainhead", getChainHead)
		c.Handle("currentminute", getCurrentMinute)
		c.Handle("ablock", getABlock)
		c.Handle("dblock", getDBlock)
		c.Handle("eblock", getEBlock)
		c.Handle("ecblock", getECBlock)
		c.Handle("fblock", getFBlock)
		c.Handle("entry", getEntry)
		c.Handle("firstentry", getFirstEntry)
		c.Handle("head", getHead)
		c.Handle("heights", getHeights)
		c.Handle("pendingentries", getPendingEntries)
		c.Handle("pendingtransactions", getPendingTransactions)
		c.Handle("raw", getraw)
		c.Handle("tps", getTPS)
		c.Handle("walletheight", getWalletHeight)

		// Deprecated calls - be backwards compatible
		c.Handle("abheight", Abheight)
		c.Handle("dbheight", Dbheight)
		c.Handle("ecbheight", Ecbheight)
		c.Handle("fbheight", Fbheight)

		c.HandleDefaultFunc(func(args []string) {
			fmt.Println(cmd.helpMsg)
		})
		c.Execute(args)
	}
	help.Add("get", cmd)
	return cmd
}()

var getAllEntries = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get allentries [-n NAME1 -h HEXNAME2" +
		" ...|CHAINID] [-E]"
	cmd.description = "Get all of the Entries confirmed in a Chain. -n and" +
		" -h to specify the chain name. -E EntryHash."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,

			"-E": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		var (
			nAcii namesASCII
			nHex  namesHex
		)
		os.Args = args
		edisp := flag.Bool("E", false, "display only the EntryHashes")
		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name component")
		flag.Var(&nHex, "h", "hex binary name component")
		flag.Parse()
		args = flag.Args()

		var chainid string

		if len(args) < 1 && len(nameCollector) == 0 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if len(nameCollector) != 0 {
			chainid = nametoid(nameCollector)
		} else {
			chainid = args[0]
		}

		es, err := factom.GetAllChainEntries(chainid)
		if err != nil {
			for i, e := range es {
				switch {
				case *edisp:
					fmt.Printf("%x\n", e.Hash())
				default:
					fmt.Printf("Entry [%d] {\n%s}\n", i, e)
				}
			}
			errorln(err)
			return
		}

		for i, e := range es {
			switch {
			case *edisp:
				fmt.Printf("%x\n", e.Hash())
			default:
				fmt.Printf("Entry [%d] {\n%s}\n", i, e)
			}
		}
	}
	help.Add("get allentries", cmd)
	return cmd
}()

var getAuthorities = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get authorities"
	cmd.description = "Get information about the authority servers on the " +
		"Factom network"
	cmd.execFunc = func(args []string) {
		as, err := factom.GetAuthorities()
		if err != nil {
			errorln(err)
			return
		}
		for _, a := range as {
			fmt.Println(a)
		}
	}
	help.Add("get authorities", cmd)
	return cmd
}()

var getChainHead = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get chainhead [-n NAME1 -h HEXNAME2 ...|CHAINID] [-K]"
	cmd.description = "Get the latest Entry Block of the specified Chain. -n " +
		"and -h to specify the chain name. -K KeyMR."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,

			"-K": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		var (
			nAcii namesASCII
			nHex  namesHex
		)
		os.Args = args
		kdisp := flag.Bool("K", false, "display only the Entry Block Key Merkel Root")

		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name component")
		flag.Var(&nHex, "h", "hex binary name component")
		flag.Parse()
		args = flag.Args()

		var chainid string

		if len(args) < 1 && len(nameCollector) == 0 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if len(nameCollector) != 0 {
			chainid = nametoid(nameCollector)
		} else {
			chainid = args[0]
		}

		head, inPL, err := factom.GetChainHead(chainid)
		if err != nil {
			errorln(err)
			return
		}

		if head == "" && inPL {
			errorln(factom.ErrChainPending)
			return
		}

		eblock, err := factom.GetEBlock(head)
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *kdisp:
			fmt.Println(head)
		default:
			fmt.Println("EBlock:", head)
			fmt.Println(eblock)
		}
	}
	help.Add("get chainhead", cmd)
	return cmd
}()

var getABlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get ablock [-RDBPL] HEIGHT|KEYMR"
	cmd.description = "Get an Admin Block from factom by its Key Merkel Root " +
		"or by its Height"
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-R": complete.PredictNothing,
			"-D": complete.PredictNothing,
			"-B": complete.PredictNothing,
			"-P": complete.PredictNothing,
			"-L": complete.PredictNothing,
		},
		Args: complete.PredictNothing,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		rdisp := flag.Bool("R", false, "display the hex encoding of the raw Directory Block")
		ddisp := flag.Bool("D", false, "display only the Directory Block height")
		bdisp := flag.Bool("B", false, "display only the Backreference Hash")
		pdisp := flag.Bool("P", false, "display only the Previous Backreference Hash")
		ldisp := flag.Bool("L", false, "display only the Lookup Hash")
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		ablock, raw, err := func() (ablock *factom.ABlock, raw []byte, err error) {
			// By KMR
			if len(args[0]) == 64 {
				ablock, err = factom.GetABlock(args[0])
			} else {
				// By height
				i, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return nil, nil, err
				}
				ablock, err = factom.GetABlockByHeight(i)
			}
			if *rdisp && err == nil {
				raw, err = factom.GetRaw(ablock.LookupHash)
			}

			return ablock, raw, err
		}()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *rdisp:
			fmt.Printf("%x\n", raw)
		case *ddisp:
			fmt.Println(ablock.DBHeight)
		case *bdisp:
			fmt.Println(ablock.BackReferenceHash)
		case *pdisp:
			fmt.Println(ablock.PrevBackreferenceHash)
		case *ldisp:
			fmt.Println(ablock.LookupHash)
		default:
			fmt.Println(ablock)
		}
	}
	help.Add("get ablock", cmd)
	return cmd
}()

var getCurrentMinute = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get currentminute [-LDMBNTSXFR]"
	cmd.description = "Get information about the current minute and other properties of the factom network."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-L": complete.PredictNothing,
			"-D": complete.PredictNothing,
			"-M": complete.PredictNothing,
			"-B": complete.PredictNothing,
			"-N": complete.PredictNothing,
			"-T": complete.PredictNothing,
			"-S": complete.PredictNothing,
			"-X": complete.PredictNothing,
			"-F": complete.PredictNothing,
			"-R": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		ldisp := flag.Bool("L", false, "display only the Leader height")
		ddisp := flag.Bool("D", false, "display only the Directory Block height")
		mdisp := flag.Bool("M", false, "display only the current minute")
		bdisp := flag.Bool("B", false, "display only the Block start time")
		ndisp := flag.Bool("N", false, "display only the minute start time")
		tdisp := flag.Bool("T", false, "display only the current time")
		sdisp := flag.Bool("S", false, "display only the Directorty Block in seconds")
		xdisp := flag.Bool("X", false, "display only the stall detected value")
		fdisp := flag.Bool("F", false, "display only the Fault timeout")
		rdisp := flag.Bool("R", false, "display only the round timeout")
		flag.Parse()
		args = flag.Args()

		info, err := factom.GetCurrentMinute()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *ldisp:
			fmt.Println(info.LeaderHeight)
		case *ddisp:
			fmt.Println(info.DirectoryBlockHeight)
		case *mdisp:
			fmt.Println(info.Minute)
		case *bdisp:
			fmt.Println(info.CurrentBlockStartTime)
		case *ndisp:
			fmt.Println(info.CurrentMinuteStartTime)
		case *tdisp:
			fmt.Println(info.CurrentTime)
		case *sdisp:
			fmt.Println(info.DirectoryBlockInSeconds)
		case *xdisp:
			fmt.Println(info.StallDetected)
		case *fdisp:
			fmt.Println(info.FaultTimeout)
		case *rdisp:
			fmt.Println(info.RoundTimeout)
		default:
			fmt.Println(info)
		}
	}
	help.Add("get currentminute", cmd)
	return cmd
}()

var getDBlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get dblock [-RHKAVNBPFTDC] HEIGHT|KEYMR"
	cmd.description = "Get a Directory Block from factom by its Key Merkel " +
		"Root or by its Height"
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-R": complete.PredictNothing,
			"-H": complete.PredictNothing,
			"-K": complete.PredictNothing,
			"-A": complete.PredictNothing,
			"-V": complete.PredictNothing,
			"-N": complete.PredictNothing,
			"-B": complete.PredictNothing,
			"-P": complete.PredictNothing,
			"-F": complete.PredictNothing,
			"-T": complete.PredictNothing,
			"-D": complete.PredictNothing,
			"-C": complete.PredictNothing,
		},
		Args: complete.PredictNothing,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		rdisp := flag.Bool("R", false, "display the hex encoding of the raw Directory Block")
		hdisp := flag.Bool("H", false, "display only the Directory Block Hash")
		kdisp := flag.Bool("K", false, "display only the Directory Block Key Merkel Root")
		adisp := flag.Bool("A", false, "display only the Directory Block Header Hash")
		vdisp := flag.Bool("V", false, "display only the Directory Block Header Version")
		ndisp := flag.Bool("N", false, "display only the Network ID")
		bdisp := flag.Bool("B", false, "display only the Directory Block Body Merkel Root")
		pdisp := flag.Bool("P", false, "display only the Previous Directory Block Key Merkel Root")
		fdisp := flag.Bool("F", false, "display only the Previous Directory Block Full Hash")
		tdisp := flag.Bool("T", false, "display only the Directory Block Timestamp")
		ddisp := flag.Bool("D", false, "display only the Directory Block Height")
		cdisp := flag.Bool("C", false, "display only the Directory Block Count")
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		dblock, raw, err := func() (dblock *factom.DBlock, raw []byte, err error) {
			if len(args[0]) == 64 {
				dblock, err = factom.GetDBlock(args[0])
			} else {
				i, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return nil, nil, err
				}
				dblock, err = factom.GetDBlockByHeight(i)
			}
			if *rdisp && err == nil {
				raw, err = factom.GetRaw(dblock.KeyMR)
			}

			return dblock, raw, err
		}()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *rdisp:
			fmt.Printf("%x\n", raw)
		case *hdisp:
			fmt.Println(dblock.DBHash)
		case *kdisp:
			fmt.Println(dblock.KeyMR)
		case *adisp:
			fmt.Println(dblock.HeaderHash)
		case *vdisp:
			fmt.Println(dblock.Header.Version)
		case *ndisp:
			fmt.Println(dblock.Header.NetworkID)
		case *bdisp:
			fmt.Println(dblock.Header.BodyMR)
		case *pdisp:
			fmt.Println(dblock.Header.PrevKeyMR)
		case *fdisp:
			fmt.Println(dblock.Header.PrevFullHash)
		case *tdisp:
			fmt.Println(dblock.Header.Timestamp)
		case *ddisp:
			fmt.Println(dblock.Header.DBHeight)
		case *cdisp:
			fmt.Println(dblock.Header.BlockCount)
		default:
			fmt.Println(dblock)
		}
	}
	help.Add("get dblock", cmd)
	return cmd
}()

var getEBlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get eblock KEYMR"
	cmd.description = "Get Entry Block by Key Merkle Root"
	cmd.completion = complete.Command{
		Args: complete.PredictNothing,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		keymr := args[0]
		eblock, err := factom.GetEBlock(keymr)
		if err != nil {
			errorln(err)
			return
		}
		fmt.Println(eblock)
	}
	help.Add("get eblock", cmd)
	return cmd
}()

var getECBlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get ecblock [-RBPLDAHF] HEIGHT|KEYMR"
	cmd.description = "Get an Entry Credit Block by Key Merkle Root or by " +
		"height"
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-R": complete.PredictNothing,
			"-B": complete.PredictNothing,
			"-P": complete.PredictNothing,
			"-L": complete.PredictNothing,
			"-D": complete.PredictNothing,
			"-A": complete.PredictNothing,
			"-H": complete.PredictNothing,
			"-F": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		rdisp := flag.Bool("R", false, "display the hex encoding of the raw Entry Credit Block")
		bdisp := flag.Bool("B", false, "display only the Body Hash")
		pdisp := flag.Bool("P", false, "display only the Previous Header Hash")
		ldisp := flag.Bool("L", false, "display only the Previous Full Hash")
		ddisp := flag.Bool("D", false, "display only the Directory Block Height")
		adisp := flag.Bool("A", false, "display only the Head Expansion Area")
		hdisp := flag.Bool("H", false, "display only the Header Hash")
		fdisp := flag.Bool("F", false, "display only the Full Hash")
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		ecblock, raw, err := func() (ecblock *factom.ECBlock, raw []byte, err error) {
			if len(args[0]) == 64 {
				ecblock, err = factom.GetECBlock(args[0])
			} else {
				i, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return nil, nil, err
				}
				ecblock, err = factom.GetECBlockByHeight(i)
			}
			if *rdisp && err == nil {
				raw, err = factom.GetRaw(ecblock.FullHash)
			}

			return ecblock, raw, err
		}()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *rdisp:
			fmt.Printf("%x\n", raw)
		case *bdisp:
			fmt.Println(ecblock.Header.BodyHash)
		case *pdisp:
			fmt.Println(ecblock.Header.PrevHeaderHash)
		case *ldisp:
			fmt.Println(ecblock.Header.PrevFullHash)
		case *ddisp:
			fmt.Println(ecblock.Header.DBHeight)
		case *adisp:
			fmt.Println(ecblock.Header.HeaderExpansionArea)
		case *hdisp:
			fmt.Println(ecblock.HeaderHash)
		case *fdisp:
			fmt.Println(ecblock.FullHash)
		default:
			fmt.Println(ecblock)
		}
	}
	help.Add("get ecblock", cmd)
	return cmd
}()

var getFBlock = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get fblock [-RBPLED] KEYMR|HEIGHT"
	cmd.description = "Get a Factoid Block by its Key Merkle Root or Height"
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-R": complete.PredictNothing,
			"-B": complete.PredictNothing,
			"-P": complete.PredictNothing,
			"-L": complete.PredictNothing,
			"-E": complete.PredictNothing,
			"-D": complete.PredictNothing,
		},
		Args: complete.PredictNothing,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		rdisp := flag.Bool("R", false, "display the hex encoding of the raw Factoid Block")
		bdisp := flag.Bool("B", false, "display only the Body Merkel Root")
		pdisp := flag.Bool("P", false, "display only the Previous Key Merkel Root")
		ldisp := flag.Bool("L", false, "display only the Previous Ledger Key Merkel Root")
		edisp := flag.Bool("E", false, "display only the Exchange Rate")
		ddisp := flag.Bool("D", false, "display only the Directory Block Height")
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		fblock, raw, err := func() (fblock *factom.FBlock, raw []byte, err error) {
			if len(args[0]) == 64 {
				fblock, err = factom.GetFBlock(args[0])
			} else {
				i, err := strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return nil, nil, err
				}
				fblock, err = factom.GetFBlockByHeight(i)
			}
			if *rdisp && err == nil {
				raw, err = factom.GetRaw(fblock.KeyMR)
			}

			return fblock, raw, err
		}()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *rdisp:
			fmt.Printf("%x\n", raw)
		case *bdisp:
			fmt.Println(fblock.BodyMR)
		case *pdisp:
			fmt.Println(fblock.PrevKeyMR)
		case *ldisp:
			fmt.Println(fblock.PrevLedgerKeyMR)
		case *edisp:
			fmt.Println(fblock.ExchRate)
		case *ddisp:
			fmt.Println(fblock.DBHeight)
		default:
			data, err := json.Marshal(fblock)
			if err != nil {
				errorln(err)
				return
			}
			var out bytes.Buffer
			json.Indent(&out, data, "", "\t")
			fmt.Printf("%s\n", out.Bytes())
		}
	}
	help.Add("get fblock", cmd)
	return cmd
}()

var getEntry = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get entry [-RC] HASH"
	cmd.description = "Get Entry by Hash. -R raw entry. -C ChainID"
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-R": complete.PredictNothing,
			"-C": complete.PredictNothing,
		},
		Args: complete.PredictNothing,
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		rdisp := flag.Bool("R", false, "display the hex encoding of the raw Entry")
		cdisp := flag.Bool("C", false, "display only the Chain ID")
		flag.Parse()
		args = flag.Args()
		if len(args) < 1 {
			fmt.Println(cmd.helpMsg)
			return
		}

		hash := args[0]
		entry, err := factom.GetEntry(hash)
		if err != nil {
			errorln(err)
			return
		}
		switch {
		case *rdisp:
			b, err := entry.MarshalBinary()
			if err != nil {
				errorln(err)
				return
			}
			fmt.Printf("%x\n", b)
		case *cdisp:
			fmt.Println(entry.ChainID)
		default:
			fmt.Println(entry)
		}
	}
	help.Add("get entry", cmd)
	return cmd
}()

var getFirstEntry = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get firstentry [-n NAME1 -h HEXNAME2 ...|CHAINID] [-REC]"
	cmd.description = "Get the first Entry in a Chain. -R RawEntry. -E EntryHash. -C ChainID."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,

			"-R": complete.PredictNothing,
			"-E": complete.PredictNothing,
			"-C": complete.PredictNothing,
		},
		Args: complete.PredictNothing,
	}
	cmd.execFunc = func(args []string) {
		var (
			nAcii namesASCII
			nHex  namesHex
		)
		os.Args = args
		rdisp := flag.Bool("R", false, "display the hex encoding of the raw Entry")
		edisp := flag.Bool("E", false, "display only the EntryHash")
		cdisp := flag.Bool("C", false, "display only the ChainID")

		nameCollector = make([][]byte, 0)
		flag.Var(&nAcii, "n", "ascii name component")
		flag.Var(&nHex, "h", "hex binary name component")
		flag.Parse()
		args = flag.Args()

		var chainid string

		if len(args) < 1 && len(nameCollector) == 0 {
			fmt.Println(cmd.helpMsg)
			return
		}
		if len(nameCollector) != 0 {
			chainid = nametoid(nameCollector)
		} else {
			chainid = args[0]
		}

		entry, err := factom.GetFirstEntry(chainid)
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *rdisp:
			b, err := entry.MarshalBinary()
			if err != nil {
				errorln(err)
				return
			}
			fmt.Printf("%x\n", b)
		case *edisp:
			fmt.Printf("%x\n", entry.Hash())
		case *cdisp:
			fmt.Println(entry.ChainID)
		default:
			fmt.Println(entry)
		}
	}
	help.Add("get firstentry", cmd)
	return cmd
}()

var getHead = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get head [-RHKAVNBPFTDC]"
	cmd.description = "Get the latest completed Directory Block"
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-R": complete.PredictNothing,
			"-H": complete.PredictNothing,
			"-K": complete.PredictNothing,
			"-A": complete.PredictNothing,
			"-V": complete.PredictNothing,
			"-N": complete.PredictNothing,
			"-B": complete.PredictNothing,
			"-P": complete.PredictNothing,
			"-F": complete.PredictNothing,
			"-T": complete.PredictNothing,
			"-D": complete.PredictNothing,
			"-C": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		rdisp := flag.Bool("R", false, "display the hex encoding of the raw Directory Block")
		hdisp := flag.Bool("H", false, "display only the Directory Block Hash")
		kdisp := flag.Bool("K", false, "display only the Directory Block Key Merkel Root")
		adisp := flag.Bool("A", false, "display only the Directory Block Header Hash")
		vdisp := flag.Bool("V", false, "display only the Directory Block Header Version")
		ndisp := flag.Bool("N", false, "display only the Network ID")
		bdisp := flag.Bool("B", false, "display only the Directory Block Body Merkel Root")
		pdisp := flag.Bool("P", false, "display only the Previous Directory Block Key Merkel Root")
		fdisp := flag.Bool("F", false, "display only the Previous Directory Block Full Hash")
		tdisp := flag.Bool("T", false, "display only the Directory Block Timestamp")
		ddisp := flag.Bool("D", false, "display only the Directory Block Height")
		cdisp := flag.Bool("C", false, "display only the Directory Block Count")
		flag.Parse()
		args = flag.Args()

		head, err := factom.GetDBlockHead()
		if err != nil {
			errorln(err)
			return
		}
		dblock, err := factom.GetDBlock(head)
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *rdisp:
			raw, err := factom.GetRaw(head)
			if err != nil {
				errorln(err)
				return
			}
			fmt.Printf("%x\n", raw)
		case *hdisp:
			fmt.Println(dblock.DBHash)
		case *kdisp:
			fmt.Println(dblock.KeyMR)
		case *adisp:
			fmt.Println(dblock.HeaderHash)
		case *vdisp:
			fmt.Println(dblock.Header.Version)
		case *ndisp:
			fmt.Println(dblock.Header.NetworkID)
		case *bdisp:
			fmt.Println(dblock.Header.BodyMR)
		case *pdisp:
			fmt.Println(dblock.Header.PrevKeyMR)
		case *fdisp:
			fmt.Println(dblock.Header.PrevFullHash)
		case *tdisp:
			fmt.Println(dblock.Header.Timestamp)
		case *ddisp:
			fmt.Println(dblock.Header.DBHeight)
		case *cdisp:
			fmt.Println(dblock.Header.BlockCount)
		default:
			fmt.Println(dblock)
		}
	}
	help.Add("get head", cmd)
	return cmd
}()

var getHeights = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get heights [-DLBE]"
	cmd.description = "Get the current heights of various items in factomd."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-D": complete.PredictNothing,
			"-L": complete.PredictNothing,
			"-B": complete.PredictNothing,
			"-E": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		ddisp := flag.Bool("D", false, "display only the Directory Block height")
		ldisp := flag.Bool("L", false, "display only the Leader height")
		bdisp := flag.Bool("B", false, "display only the EntryBlock height")
		edisp := flag.Bool("E", false, "display only the Entry height")
		flag.Parse()
		args = flag.Args()

		heights, err := factom.GetHeights()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *ddisp:
			fmt.Println(heights.DirectoryBlockHeight)
		case *ldisp:
			fmt.Println(heights.LeaderHeight)
		case *bdisp:
			fmt.Println(heights.EntryBlockHeight)
		case *edisp:
			fmt.Println(heights.EntryHeight)
		default:
			fmt.Print(heights.String())
		}
	}
	help.Add("get heights", cmd)
	return cmd
}()

var getWalletHeight = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get walletheight"
	cmd.description = "Get the number of factoid blocks factom-walletd has cached"
	cmd.execFunc = func(args []string) {
		height, err := factom.GetWalletHeight()
		if err != nil {
			errorln(err)
			return
		}
		fmt.Printf("WalletHeight: %v\n", height)
	}
	help.Add("get walletheight", cmd)
	return cmd
}()

var properties = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli properties [-CFAWL]"
	cmd.description = "Get version information about factomd and the factom wallet"
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-C": complete.PredictNothing,
			"-F": complete.PredictNothing,
			"-A": complete.PredictNothing,
			"-W": complete.PredictNothing,
			"-L": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		cdisp := flag.Bool("C", false, "display only the CLI version")
		fdisp := flag.Bool("F", false, "display only the factomd version")
		adisp := flag.Bool("A", false, "display only the factomd API version")
		wdisp := flag.Bool("W", false, "display only the factom-wallet version")
		ldisp := flag.Bool("L", false, "display only the wallet API version")
		flag.Parse()
		args = flag.Args()

		props, err := factom.GetProperties()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *cdisp:
			fmt.Println(FactomcliVersion)
		case *fdisp:
			fmt.Println(props.FactomdVersion)
		case *adisp:
			fmt.Println(props.FactomdAPIVersion)
		case *wdisp:
			fmt.Println(props.WalletVersion)
		case *ldisp:
			fmt.Println(props.WalletAPIVersion)
		default:
			fmt.Println("CLI Version:", FactomcliVersion)
			fmt.Println(props)
		}
	}
	help.Add("properties", cmd)
	return cmd
}()

var getPendingEntries = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get pendingentries [-E]"
	cmd.description = "Get all pending entries, which may not yet be written" +
		" to blockchain. -E EntryHash."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-E": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		edisp := flag.Bool(
			"E",
			false,
			"display only the Entry Hashes",
		)
		flag.Parse()
		args = flag.Args()

		entries, err := factom.GetPendingEntries()
		if err != nil {
			errorln(err)
			return
		}

		for _, ents := range entries {
			switch {
			case *edisp:
				fmt.Println(ents.EntryHash)
			default:
				fmt.Println("ChainID:", ents.ChainID)
				fmt.Println("Entryhash:", ents.EntryHash)
				fmt.Println("")
			}
		}
	}
	help.Add("get pendingentries", cmd)
	return cmd
}()

var getPendingTransactions = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get pendingtransactions [-T]"
	cmd.description = "Get all pending factoid transacitons, which may not yet be written to blockchain. -T TxID."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-T": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		tdisp := flag.Bool(
			"T",
			false,
			"display only the Transaction IDs",
		)
		flag.Parse()
		args = flag.Args()

		trans, err := factom.GetPendingTransactions()
		if err != nil {
			errorln(err)
			return
		}

		for _, tran := range trans {
			if len(tran.Inputs) != 0 {
				switch {
				case *tdisp:
					fmt.Println(tran.TxID)
				default:
					fmt.Println("TxID:", tran.TxID)
					for _, in := range tran.Inputs {
						fmt.Println("Input:", in.Address, in.Amount/1e9)
					}

					if len(tran.Outputs) != 0 {

						for _, out := range tran.Outputs {
							fmt.Println("Output:", out.Address, out.Amount/1e9)
						}
					}

					if len(tran.ECOutputs) != 0 {
						for _, ecout := range tran.ECOutputs {
							fmt.Println("ECOutput:", ecout.Address, ecout.Amount/1e9)
						}
					}
					if tran.Fees != 0 {
						fmt.Printf("Fees: %8.8f", float64(tran.Fees)/1e9)
					}
					fmt.Println("")
				}
			}
		}

	}
	help.Add("get pendingtransactions", cmd)
	return cmd
}()

var getTPS = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli get tps [-IT]"
	cmd.description = "Get the current instant and total average rate of " +
		"Transactions Per Second."
	cmd.completion = complete.Command{
		Flags: complete.Flags{
			"-I": complete.PredictNothing,
			"-T": complete.PredictNothing,
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		idisp := flag.Bool("I", false, "display only the instant TPS rate")
		tdisp := flag.Bool("T", false, "display only the total averaged TPS rate")
		flag.Parse()

		i, t, err := factom.GetTPS()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *idisp:
			fmt.Println(i)
		case *tdisp:
			fmt.Println(t)
		default:
			fmt.Printf("Instant: %0.2f\n", i)
			fmt.Printf("Total: %0.2f\n", t)
		}
	}
	help.Add("get tps", cmd)
	return cmd
}()
