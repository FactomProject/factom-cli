// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
)

func man(s string) error {
	m := map[string]string{
		"buy": "factom-cli buy [-s server ] amt",
		"factoidtx": "factom-cli factoidtx [-s server ] addr amt",
		"get": "no help for get",
		"help":	"factom-cli help [subcommand]",
		"mkchain": "factom-cli mkchain [-s server]",
		"put": "factom-cli put [-s server] [-e extid ...] [file]",
		"default": "factom-cli [subcommand] [options]",
	}
	
	if m[s] != "" {
		return fmt.Errorf(m[s])
	}
	return fmt.Errorf(m["default"])
}