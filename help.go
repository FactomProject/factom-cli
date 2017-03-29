// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type helper struct {
	topics map[string]*fctCmd
}

// NewHelper creates a new helper object containing the help messages for a set
// of commands.
func newHelper() *helper {
	h := new(helper)
	h.topics = make(map[string]*fctCmd)
	return h
}

func (h *helper) Add(s string, c *fctCmd) {
	h.topics[s] = c
}

func (h *helper) All() {
	flag.VisitAll(func(f *flag.Flag) {
		m, err := regexp.MatchString("^test", f.Name)
		if err != nil {
			errorln(err)
		}
		if !m {
			fmt.Printf("%s\n\t%s\n\n", "-"+f.Name, f.Usage)
		}
	})

	keys := make([]string, 0)
	for k := range h.topics {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, v := range keys {
		fmt.Printf("%s\n\t%s\n\n", h.topics[v].helpMsg, h.topics[v].description)
	}
}

func (h *helper) Execute(args []string) {
	if len(args) < 1 {
		fmt.Println("factom-cli help [subcommand]")
		return
	}

	if args[0] == "help" {
		if len(args) == 1 {
			help.All()
			return
		}
		args = args[1:]
	}

	topic := strings.Join(args[:], " ")

	c, ok := h.topics[topic]
	if !ok {
		if c, ok = h.topics[args[0]]; !ok {
			fmt.Println("No help for:", topic)
			return
		}
	}
	fmt.Printf("%s\n\t%s\n", c.helpMsg, c.description)
}

var help = newHelper()
