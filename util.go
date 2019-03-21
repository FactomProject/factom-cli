// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/FactomProject/btcutil/base58"
	"github.com/FactomProject/factom"
)

// GetFactomdServer returns the current factomd server set in the factom package.
// This is used for compose functions to put the appropriate url in the sample curl
//		localhost:8088 is returned if the factom package does not have one set
func GetFactomdServer() string {
	if factom.RpcConfig != nil && factom.RpcConfig.FactomdServer != "" {
		return factom.RpcConfig.FactomdServer
	}
	return "localhost:8088"
}

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

// keysASCII will be a flag receiver for ASCII identity keys.
type keysASCII []string

func (n *keysASCII) String() string {
	return fmt.Sprint(*n)
}

func (k *keysASCII) Set(s string) error {
	*k = append(*k, s)
	if factom.IdentityKeyStringType(s) != factom.IDPub {
		return fmt.Errorf("Provided key string not a valid public identity key: %s", s)
	}
	b := base58.Decode(s)
	key := factom.NewIdentityKey()
	copy(key.Pub[:], b[factom.IDKeyPrefixLength:factom.IDKeyBodyLength])
	return nil
}

func factoshiToFactoid(v interface{}) string {
	value, err := strconv.Atoi(fmt.Sprint(v))
	if err != nil {
		return ""
	}
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	d := value / 1e8
	r := value % 1e8
	ds := fmt.Sprintf("%s%d", sign, d)
	rs := fmt.Sprintf("%08d", r)
	rs = strings.TrimRight(rs, "0")
	if len(rs) > 0 {
		ds = ds + "."
	}
	return fmt.Sprintf("%s%s", ds, rs)
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
			if (s.Status != "Unknown") && (s.Status != "NotConfirmed") {
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
	case <-time.After(60 * time.Second):
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
			s, err := factom.EntryCommitACK(txid, "")
			if err != nil {
				errchan <- err
				break
			}

			if (s.CommitData.Status != "Unknown") && (s.CommitData.Status != "NotConfirmed") {
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
	case <-time.After(60 * time.Second):
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
			// All 0s signals an entry
			s, err := factom.EntryRevealACK(txid, "", "0000000000000000000000000000000000000000000000000000000000000000")
			if err != nil {
				errchan <- err
				break
			}

			if (s.EntryData.Status != "Unknown") && (s.EntryData.Status != "NotConfirmed") {
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
	case <-time.After(60 * time.Second):
		return "", fmt.Errorf("timeout: no acknowledgement found")
	}

	// code should not reach this point
	return "", fmt.Errorf("unknown error")
}
