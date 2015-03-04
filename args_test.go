package main

import (
	"flag"
	"fmt"
	"net/http"
	"testing"
	"time"
	"os"
)

// anonymous declarations to squash import errors for testing
var _ = flag.Bool("bool", false, "")
var _ string = fmt.Sprint()
var _ []string = os.Args
var _ = time.Second

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r)
}

func TestMain(t *testing.T) {
	fmt.Printf("TestMain\n===\n")
	os.Args = []string{"test", "help", "get"}
	main()
	fmt.Println()
}

func TestTmp(t *testing.T) {
	fmt.Printf("TestTmp\n===\n")
	var (
		b = flag.Bool("b", false, "boolflag")
		a = flag.Bool("a", false, "a flag")
	)
	os.Args = []string{"test", "-b", "something", "-a", "amt"}
	flag.Parse()
	_ = a
	_ = b
	args := flag.Args()
	fmt.Println(args)
	fmt.Println()
}
