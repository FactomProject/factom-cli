package main

import (
	"fmt"
	"testing"
)

var (
	_ = fmt.Sprint("testing")
)

func TestPrint(t *testing.T) {
	values := []uint64{100000000, 10000000, 1000000, 100000, 10000, 1000, 100, 10, 1}
	ans := []string{"1", "0.1", "0.01", "0.001", "0.0001", "0.00001", "0.000001", "0.0000001", "0.00000001"}
	for i, v := range values {
		a := factoshiToFactoid(v)
		if a != ans[i] {
			t.Errorf("The value %d is %s not %s as expected", v, a, ans[i])
		}
	}
}
