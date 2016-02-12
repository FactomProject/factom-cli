// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

type fctCmd struct {
	execFunc    func([]string)
	helpMsg     string
	description string
}

func (c *fctCmd) Execute(args []string) {
	c.execFunc(args)
}
