package main

import (
	"fmt"
	"os"
	"testing"
)

var (
	_ = fmt.Sprint("testing")
)

func TestHeights(t *testing.T) {
	// need to have factomd available to really test
	// find way to run and look for specific results
	os.Args[0] = "1"
	os.Args[1] = "-r"

	testEc := Ecbheight

	fmt.Println(testEc)

	testA := Abheight

	fmt.Println(testA)

	testDB := Dbheight

	fmt.Println(testDB)

	testF := Fbheight

	fmt.Println(testF)
}
