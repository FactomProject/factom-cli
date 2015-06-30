// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import(
    "testing"
)
    
func Test_Input_Amounts(test *testing.T) {
    
    str,err := convertFixedPoint("10") 
    if err != nil { test.Fail() }
    if str != "1000000000" { test.Fail() }
    
    str,err = convertFixedPoint("10.08") 
    if err != nil { test.Fail() }
    if str != "1008000000" { test.Fail() }
    
    str,err = convertFixedPoint("10.08000001") 
    if err != nil { test.Fail() }
    if str != "1008000001" {  test.Fail() }
    
}