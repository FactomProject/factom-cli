package main

import (
	"os"
)

func main() {
	command := "help"
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}
	switch(command) {
	case "help":
		a := "help"
		if len(os.Args) >= 3 {
			a = os.Args[2]
		}
		help(a)			
	case "put":
		put()
	case "get":
		get()
	case "mkchain":
		mkchain()
	}
}
