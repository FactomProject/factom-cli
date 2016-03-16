// Copyright 2016 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/FactomProject/factomd/common/primitives"
	"testing"
)

func Test_Input_Amounts(test *testing.T) {

	str, err := primitives.ConvertFixedPoint("10")
	if err != nil {
		test.Fail()
	}
	if str != "1000000000" {
		test.Fail()
	}

	str, err = primitives.ConvertFixedPoint("10.08")
	if err != nil {
		test.Fail()
	}
	if str != "1008000000" {
		test.Fail()
	}

	str, err = primitives.ConvertFixedPoint("10.08000001")
	if err != nil {
		test.Fail()
	}
	if str != "1008000001" {
		test.Fail()
	}

}
