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
var _ = flag.Bool("b", false, "")
var _ string = fmt.Sprint()
var _ []string = os.Args

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r)
}

func server() error {
	http.HandleFunc("/", ServeHTTP)
	err := http.ListenAndServe("localhost:4321", nil)
	if err != nil {
		return err
	}
	return nil
}

func TestMain(t *testing.T) {
	fmt.Printf("TestMain\n===\n")
	os.Args = []string{"test", "help", "get"}
	main()
	fmt.Println()
}

func TestTmp(t *testing.T) {
	fmt.Printf("TestTmp\n===\n")
	os.Args = []string{"test", "-h", "something"}
	var h = flag.Bool("h", false, "help")
	flag.Parse()
	if *h {
		fmt.Println("got help flag")
	}
	fmt.Println(os.Args)
	fmt.Println()
}

func TestPut(t *testing.T) {
	c := make(chan error)
	go func() {
		c <- server()
	}()
	
	time.Sleep(50 * time.Millisecond)

	select {
	case err := <-c:
		t.Errorf(err.Error())
	default:
		args := []string{
			"put",
			"-s", "localhost:4321",
			"-e", "1234",
			"-e", "test",
			"-c", "d5f39e4c4e041c37dfe0d65c7405d215924650891a689425c736e974c88d5ba0",
		}
		
		err := put(args)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}
