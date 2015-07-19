// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"os"
)

func help(args []string) {
	os.Args = args

	flag.Parse()
	args = flag.Args()
	s := "help"
	if len(args) > 0 {
		s = args[0]
	}

	man(s)
}
