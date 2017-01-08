package main

import (
	"fmt"
	"testing"
)

var (
	_ = fmt.Sprint("testing")
)

func TestPrint(t *testing.T) {
	vs1 := []uint64{100000000, 10000000, 1000000, 100000, 10000, 1000, 100, 10, 1}
	vs2 := []int64{100000000, 10000000, 1000000, 100000, 10000, 1000, 100, 10, 1,
		-100000000, -10000000, -1000000, -100000, -10000, -1000, -100, -10, -1}
	ans1 := []string{"1", "0.1", "0.01", "0.001", "0.0001", "0.00001",
		"0.000001", "0.0000001", "0.00000001"}
	ans2 := []string{"1", "0.1", "0.01", "0.001", "0.0001", "0.00001",
		"0.000001", "0.0000001", "0.00000001", "-1", "-0.1", "-0.01", "-0.001",
		"-0.0001", "-0.00001", "-0.000001", "-0.0000001", "-0.00000001"}
	for i, v := range vs1 {
		a := factoshiToFactoid(v)
		if a != ans1[i] {
			t.Errorf("The value %d is %s not %s as expected", v, a, ans1[i])
		}
	}
	for i, v := range vs2 {
		a := factoshiToFactoid(v)
		if a != ans2[i] {
			t.Errorf("The value %d is %s not %s as expected", v, a, ans2[i])
		}
	}
}
