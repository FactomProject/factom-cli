// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FactomProject/factom"
	"github.com/michaelbeam/cli"
	"github.com/posener/complete"
)

var diagnostics = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli diagnostics [server|network|sync|election|authset]"
	cmd.description = "Get diagnostic information about the Factom network"
	cmd.completion = complete.Command{
		Sub: complete.Commands{
			"server":   complete.Command{},
			"network":  complete.Command{},
			"sync":     complete.Command{},
			"election": complete.Command{},
			"authset":  complete.Command{},
		},
	}
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		c := cli.New()
		c.Handle("server", diagnosticsServer)
		c.Handle("network", diagnosticsNetwork)
		c.Handle("sync", diagnosticsSync)
		c.Handle("election", diagnosticsElection)
		c.Handle("authset", diagnosticsAuthSet)
		c.HandleDefaultFunc(func(args []string) {
			if len(args) != 0 {
				fmt.Println(cmd.helpMsg)
				return
			}
			d, err := factom.GetDiagnostics()
			if err != nil {
				errorln(err)
				return
			}
			fmt.Println(d)
		})
		c.Execute(args)
	}
	help.Add("diagnostics", cmd)
	return cmd
}()

var diagnosticsServer = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli diagnostics server [-NIKR]"
	cmd.description = "Get diagnostic information about the Factom API server"
	cmd.execFunc = func(args []string) {
		os.Args = args
		ndisp := flag.Bool("N", false, "display only the name of the API server")
		idisp := flag.Bool("I", false, "display only the ID of the API server")
		kdisp := flag.Bool("K", false, "display only the public key of the API server")
		rdisp := flag.Bool("R", false, "display only the role of the API server")
		flag.Parse()
		args = flag.Args()

		d, err := factom.GetDiagnostics()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *ndisp:
			fmt.Println(d.Name)
		case *idisp:
			fmt.Println(d.ID)
		case *kdisp:
			fmt.Println(d.PublicKey)
		case *rdisp:
			fmt.Println(d.Role)
		default:
			fmt.Println("Name:", d.Name)
			fmt.Println("ID:", d.ID)
			fmt.Println("PublicKey:", d.PublicKey)
			fmt.Println("Role:", d.Role)
		}
	}
	help.Add("diagnostics server", cmd)
	return cmd
}()

var diagnosticsNetwork = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli diagnostics network [-LMDPBTS]"
	cmd.description = "Get diagnostic information about the Factom network"
	cmd.execFunc = func(args []string) {
		os.Args = args
		ldisp := flag.Bool("L", false, "display only the Leader height")
		mdisp := flag.Bool("M", false, "display only the current minute")
		ddisp := flag.Bool("D", false, "display only the current minute duration")
		pdisp := flag.Bool("P", false, "display only the previous minute duration")
		hdisp := flag.Bool("H", false, "display only the balance hash")
		tdisp := flag.Bool("T", false, "display only the temporary balance hash")
		bdisp := flag.Bool("B", false, "display only the last block from the DBState")
		flag.Parse()
		args = flag.Args()

		d, err := factom.GetDiagnostics()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *ldisp:
			fmt.Println(d.LeaderHeight)
		case *mdisp:
			fmt.Println(d.CurrentMinute)
		case *ddisp:
			fmt.Println(d.CurrentMinuteDuration)
		case *pdisp:
			fmt.Println(d.PrevMinuteDuration)
		case *hdisp:
			fmt.Println(d.BalanceHash)
		case *tdisp:
			fmt.Println(d.TempBalanceHash)
		case *bdisp:
			fmt.Println(d.LastBlockFromDBState)
		default:
			fmt.Println("LeaderHeight:", d.LeaderHeight)
			fmt.Println("CurrentMinute:", d.CurrentMinute)
			fmt.Println("CurrentMinuteDuration:", d.CurrentMinuteDuration)
			fmt.Println("PrevMinuteDuration:", d.PrevMinuteDuration)
			fmt.Println("BalanceHash:", d.BalanceHash)
			fmt.Println("TempBalanceHash:", d.TempBalanceHash)
			fmt.Println("LastBlockFromDBState:", d.LastBlockFromDBState)
		}
	}
	help.Add("diagnostics network", cmd)
	return cmd
}()

var diagnosticsSync = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli diagnostics sync [-SRE]"
	cmd.description = "Get diagnostic information about the network syncing"
	cmd.execFunc = func(args []string) {
		os.Args = args
		sdisp := flag.Bool("S", false, "display only the syncing status")
		rdisp := flag.Bool("R", false, "display only the received status")
		edisp := flag.Bool("E", false, "display only the expected status")
		flag.Parse()
		args = flag.Args()

		d, err := factom.GetDiagnostics()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *sdisp:
			fmt.Println(d.SyncInfo.Status)
		case *rdisp:
			fmt.Println(d.SyncInfo.Received)
		case *edisp:
			fmt.Println(d.SyncInfo.Expected)
		default:
			fmt.Println("Status:", d.SyncInfo.Status)
			fmt.Println("Received:", d.SyncInfo.Received)
			fmt.Println("Expected:", d.SyncInfo.Expected)
			for _, m := range d.SyncInfo.Missing {
				fmt.Println("Missing:", m)
			}
		}
	}
	help.Add("diagnostics sync", cmd)
	return cmd
}()

var diagnosticsElection = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli diagnostics election [-PVFIR]"
	cmd.description = "Get diagnostic information about the Factom network election process"
	cmd.execFunc = func(args []string) {
		os.Args = args
		pdisp := flag.Bool("P", false, "display only the progress status")
		vdisp := flag.Bool("V", false, "display only the VM Index")
		fdisp := flag.Bool("F", false, "display only the Federated Index")
		idisp := flag.Bool("I", false, "display only the Federated ID")
		rdisp := flag.Bool("R", false, "display only the current round")
		flag.Parse()
		args = flag.Args()

		d, err := factom.GetDiagnostics()
		if err != nil {
			errorln(err)
			return
		}

		switch {
		case *pdisp:
			fmt.Println(d.ElectionInfo.InProgress)
		case *vdisp:
			fmt.Println(d.ElectionInfo.VMIndex)
		case *fdisp:
			fmt.Println(d.ElectionInfo.FedIndex)
		case *idisp:
			fmt.Println(d.ElectionInfo.FedID)
		case *rdisp:
			fmt.Println(d.ElectionInfo.Round)
		default:
			fmt.Println("InProgress:", d.ElectionInfo.InProgress)
			fmt.Println("VMIndex:", d.ElectionInfo.VMIndex)
			fmt.Println("FedIndex:", d.ElectionInfo.FedIndex)
			fmt.Println("FedID:", d.ElectionInfo.FedID)
			fmt.Println("Round:", d.ElectionInfo.Round)
		}
	}
	help.Add("diagnostics election", cmd)
	return cmd
}()

var diagnosticsAuthSet = func() *fctCmd {
	cmd := new(fctCmd)
	cmd.helpMsg = "factom-cli diagnostics authset"
	cmd.description = "Get diagnostic information about the Factom authorized servers"
	cmd.execFunc = func(args []string) {
		os.Args = args
		flag.Parse()
		args = flag.Args()

		d, err := factom.GetDiagnostics()
		if err != nil {
			errorln(err)
			return
		}

		fmt.Println("Leaders {")
		for _, v := range d.AuthSet.Leaders {
			fmt.Println(" ID:", v.ID)
			fmt.Println(" VM:", v.VM)
			fmt.Println(" ProcessListHeight:", v.ProcessListHeight)
			fmt.Println(" ListLength:", v.ListLength)
			fmt.Println(" NextNil:", v.NextNil)
		}
		fmt.Println("}") // Leaders
		fmt.Println("Audits {")
		for _, v := range d.AuthSet.Audits {
			fmt.Println(" ID:", v.ID)
			fmt.Println(" Online:", v.Online)
		}
		fmt.Println("}") // Audits
	}
	help.Add("diagnostics authset", cmd)
	return cmd
}()
