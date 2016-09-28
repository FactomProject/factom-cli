// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

// exidCollector accumulates the external ids from the command line -e and -E
// flags
var exidCollector [][]byte

// extids will be a flag receiver for adding chains and entries
// In ASCII
type extidsAscii []string

func (e *extidsAscii) String() string {
	return fmt.Sprint(*e)
}

func (e *extidsAscii) Set(s string) error {
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

// namesAscii will be a flag receiver for ASCII chain names.
type namesAscii []string

func (n *namesAscii) String() string {
	return fmt.Sprint(*n)
}

func (n *namesAscii) Set(s string) error {
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
	whole := value / 100000000
	part := value - (whole * 100000000)

	ret := []byte(fmt.Sprintf("%d.%08d", whole, part))
	for string(ret[len(ret)-1]) == "0" {
		ret = ret[:len(ret)-1]
	}
	if string(ret[len(ret)-1]) == "." {
		ret = ret[:len(ret)-1]
	}

	return string(ret)
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
