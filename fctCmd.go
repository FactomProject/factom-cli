// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/posener/complete"
)

type fctCmd struct {
	execFunc    func([]string)
	helpMsg     string
	description string
	completion  complete.Command
}

func (c *fctCmd) Execute(args []string) {
	c.execFunc(args)
}
