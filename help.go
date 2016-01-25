// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
)

type helper struct {
	topics map[string]*fctCmd
}

func NewHelper() *helper {
	h := new(helper)
	h.topics = make(map[string]*fctCmd)
	return h
}

func (h *helper) Add(s string, c *fctCmd) {
	h.topics[s] = c
}

func (h *helper) All() {
	for _, v := range h.topics {
		fmt.Println(v.description)
		fmt.Println(v.helpMsg)
	}
}

func (h *helper) Execute(args []string) {
	if len(args) < 1 {
		fmt.Println("factom-cli help [subcommand]")
		return
	}
	if len(args) < 2 {
		h.All()
		return
	}
	
	topic := strings.Join(args[1:], " ")
	
	c, ok := h.topics[topic]
	if !ok {
		if c, ok = h.topics[args[1]]; !ok {
			fmt.Println("No help for:", topic)
			return
		}
	}
	fmt.Println(c.description)
	fmt.Println(c.helpMsg)
}

var help = NewHelper()