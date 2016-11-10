// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/FactomProject/factom"
)

// exidCollector accumulates the external ids from the command line -e and -E
// flags
var exidCollector [][]byte

// extids will be a flag receiver for adding chains and entries
// In ASCII
type extidsASCII []string

func (e *extidsASCII) String() string {
	return fmt.Sprint(*e)
}

func (e *extidsASCII) Set(s string) error {
	*e = append(*e, s)
	exidCollector = append(exidCollector[:], []byte(s))
	return nil
}

// extids will be a flag receiver for adding chains and entries
// In HEX
type extidsHex []string

func (e *extidsHex) String() string {
	return fmt.Sprint(*e)
}

func (e *extidsHex) Set(s string) error {
	*e = append(*e, s)
	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	exidCollector = append(exidCollector[:], b)
	return nil
}

// nameCollector accumulates the components of a chain name from the command
// line -n and -N flags
var nameCollector [][]byte

// namesASCII will be a flag receiver for ASCII chain names.
type namesASCII []string

func (n *namesASCII) String() string {
	return fmt.Sprint(*n)
}

func (n *namesASCII) Set(s string) error {
	*n = append(*n, s)
	nameCollector = append(nameCollector[:], []byte(s))
	return nil
}

// namesHex will be a flag receiver for HEX encoded chain names.
type namesHex []string

func (n *namesHex) String() string {
	return fmt.Sprint(*n)
}

func (n *namesHex) Set(s string) error {
	*n = append(*n, s)
	b, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	nameCollector = append(nameCollector[:], b)
	return nil
}

func factoshiToFactoid(v interface{}) string {
	value, err := strconv.Atoi(fmt.Sprint(v))
	if err != nil {
		return ""
	}
	f := float64(value) / 1e8
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// nametoid computes a chainid from the chain name components
func nametoid(name [][]byte) string {
	hs := sha256.New()
	for _, v := range name {
		h := sha256.Sum256(v)
		hs.Write(h[:])
	}
	return hex.EncodeToString(hs.Sum(nil))
}

// waitOnFctAck blocks while waiting for a factom ack message and returns the
// ack status or times out after 10 seconds.
func waitOnFctAck(txid string) (string, error) {
	stat := make(chan string, 1)
	errchan := make(chan error, 1)

	// poll for the acknowledgement
	go func() {
		for {
			s, err := factom.FactoidACK(txid, "")
			if err != nil {
				errchan <- err
				break
			}
			if s.Status != "Unknown" {
				stat <- s.Status
				break
			}
			time.Sleep(time.Second / 2)
		}
	}()

	// wait for the acknowledgement or timeout after 10 sec
	select {
	case err := <-errchan:
		return "", err
	case s := <-stat:
		return s, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timeout: no acknowledgement found")
	}

	// code should not reach this point
	return "", fmt.Errorf("unknown error")
}

// waitOnCommitAck blocks while waiting for an ack message for an Entry Commit
// and returns the ack status or times out after 10 seconds.
func waitOnCommitAck(txid string) (string, error) {
	stat := make(chan string, 1)
	errchan := make(chan error, 1)

	// poll for the acknowledgement
	go func() {
		for {
			s, err := factom.EntryACK(txid, "")
			if err != nil {
				errchan <- err
				break
			}
			if s.CommitData.Status != "Unknown" {
				stat <- s.CommitData.Status
				break
			}
			time.Sleep(time.Second / 2)
		}
	}()

	// wait for the acknowledgement or timeout after 10 sec
	select {
	case err := <-errchan:
		return "", err
	case s := <-stat:
		return s, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timeout: no acknowledgement found")
	}

	// code should not reach this point
	return "", fmt.Errorf("unknown error")
}

// waitOnRevealAck blocks while waiting for an ack message for an Entry Reveal
// and returns the ack status or times out after 10 seconds.
func waitOnRevealAck(txid string) (string, error) {
	stat := make(chan string, 1)
	errchan := make(chan error, 1)

	// poll for the acknowledgement
	go func() {
		for {
			s, err := factom.EntryACK(txid, "")
			if err != nil {
				errchan <- err
				break
			}
			if s.EntryData.Status != "Unknown" {
				stat <- s.EntryData.Status
				break
			}
			time.Sleep(time.Second / 2)
		}
	}()

	// wait for the acknowledgement or timeout after 10 sec
	select {
	case err := <-errchan:
		return "", err
	case s := <-stat:
		return s, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timeout: no acknowledgement found")
	}

	// code should not reach this point
	return "", fmt.Errorf("unknown error")
}
