// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
    "bytes"
    "strconv"
    "net/http"
    "encoding/hex"
    "encoding/json"
    "io/ioutil"
    "flag"
    "fmt"
    "os"
    fct "github.com/FactomProject/factoid"
    "github.com/FactomProject/factom"
)

var _ = hex.EncodeToString
var serverFct = "localhost:8089"

func getCmd(cmd string, cmderror string) {
	resp, err := http.Get(cmd)
	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body) 
	
	if err != nil {
        fmt.Println(err)
        os.Exit(1)
	}
	resp.Body.Close()

	type x struct{ 
        Body string
        Success bool 
    }
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil || !b.Success {
		fmt.Println(cmderror)
		fmt.Println("Command Failed: ", string(body))
        os.Exit(1)
	}
    fmt.Println(b.Body)
	return 
}

func postCmd(cmd string) {
	resp, err := http.PostForm(cmd, nil)
	if err != nil {
		fmt.Println(err)
        os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	resp.Body.Close()

	type x struct{ Response string; Success bool }
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil {
        fmt.Printf("Failed to parse the response from factomd: %s\n",body)
        os.Exit(1)
	}
	
	fmt.Println(b.Response)
    
    if !b.Success {
        os.Exit(1)
    }
      
	return 
}

// Generates a new Address
func generateaddress(args []string) {
    
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 2 {
        fmt.Println(man("generatefactoidaddress"))
        os.Exit(1)
    }
    
    var err error
    var Addr string
    switch args[0]{
        case "ec": 
            Addr, err= factom.GenerateEntryCreditAddress(args[1])
        case "fct":
            Addr, err= factom.GenerateFactoidAddress(args[1])
        default:
            panic("Expected ec|fct name")
    }
    
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println(args[0]," = ",Addr)
    return 
    
}

func getaddresses(args []string) {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) > 0 {
        fmt.Println(man("getaddresses"))
        os.Exit(1)
    }
    
    str := fmt.Sprintf("http://%s/v1/factoid-get-addresses/", serverFct)
    getCmd(str,"Error printing addresses")
    
    return 
    
    
}

func fctnewtrans(args []string) {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 1 {
        fmt.Println("Missing Key")
        fmt.Println(man("newtransaction"))
        os.Exit(1)
    } 
    
    str := fmt.Sprintf("http://%s/v1/factoid-new-transaction/%s", serverFct, args[0])
    postCmd(str)
    
    return 
    
}

func fctdeletetrans(args []string) {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 1 {
        fmt.Println("Missing Key")
        fmt.Println(man("deletetransaction"))
        os.Exit(1)
    } 
    
    str := fmt.Sprintf("http://%s/v1/factoid-delete-transaction/%s", serverFct, args[0])
    postCmd(str)
    
    return 
    
}


func fctaddinput(args []string) {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 3 {
        fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
        os.Exit(1)
    } 

    amt,err := fct.ConvertFixedPoint(args[2])
    if err != nil { 
        fmt.Println(err)
        os.Exit(1) 
    }
    
    ramt,err := strconv.ParseInt(amt,0,64)
    if err != nil { 
        fmt.Println(err)
        os.Exit(1) 
    }
 
    _,ok2 := fct.ValidateAmounts(uint64(ramt))
    if !ok2 { 
        fmt.Println("Input is too large or is negative. Command ignored") 
        os.Exit(1)
    }

    str := fmt.Sprintf("http://%s/v1/factoid-add-input/?key=%s&name=%s&amount=%s", 
                       serverFct, args[0],args[1],amt)
    postCmd(str)
    
    return 
}

func fctaddoutput(args []string) {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 3 {
        fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
        os.Exit(1) 
    } 
    // localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>
    
    amt,err := fct.ConvertFixedPoint(args[2])
    if err != nil { 
        fmt.Println("Invalid format for a number: ",args[2])
        fmt.Println(man("addoutput"))
        os.Exit(1) 
    }
    
    ramt,err := strconv.ParseInt(amt,0,64)
    if err != nil { 
        fmt.Println(err)
        os.Exit(1) 
    }

    _,ok2 := fct.ValidateAmounts(uint64(ramt))
    if !ok2 { 
        fmt.Println("Input is too large or is negative. Command ignored") 
        os.Exit(1)
    }
    
    
    str := fmt.Sprintf("http://%s/v1/factoid-add-output/?key=%s&name=%s&amount=%s", 
                       serverFct, args[0],args[1],amt)
    postCmd(str)
    
    return 
}

func fctaddecoutput(args []string) {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 3 {
        fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
        os.Exit(1) 
    } 
    // localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>
    
    amt,err := fct.ConvertFixedPoint(args[2])
    if err != nil { 
        fmt.Println(err)
        os.Exit(1)  
    }
    
    ramt,err := strconv.ParseInt(amt,0,64)
    if err != nil { 
        fmt.Println(err)
        os.Exit(1) 
    }

    _,ok2 := fct.ValidateAmounts(uint64(ramt))
    if !ok2 { 
        fmt.Println("Input is too large or is negative. Command ignored") 
        os.Exit(1)
    }
    
    str := fmt.Sprintf("http://%s/v1/factoid-add-ecoutput/?key=%s&name=%s&amount=%s", 
                       serverFct, args[0],args[1],amt)
    postCmd(str)
    
    return 
}

func fctgetfee(args []string) {
    resp, err := http.Get(fmt.Sprintf("http://%s/v1/factoid-get-fee/",serverFct))
    if err != nil {
        fmt.Println("Command Failed Get")
        os.Exit(1) 
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Command Failed")
        os.Exit(1) 
    }
    resp.Body.Close()
    
    // We pull the fee.  If the fee isn't positive, or if we fail to marshal, then there is a failure
    type x struct { Fee int64 }
    b := new(x)
    b.Fee = -1
    if err := json.Unmarshal(body, b); err != nil || b.Fee == -1 {
        fmt.Println("Command Failed")
        os.Exit(1)
    }
    tv := b.Fee/100000000
    lv := b.Fee-(tv*100000000)
    r := fmt.Sprintf("Fee: %d.%08d",tv,lv)
    var i int; for i=len(r)-1; r[i]=='0'; i-- {}
    if string(r[i])=="." { i +=1 }
    fmt.Println(r[:i+1])
     
}
    
func fctsign(args []string) {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 1 {
        fmt.Println("Missing Key")
        os.Exit(1)
    } 
    
    str := fmt.Sprintf("http://%s/v1/factoid-sign-transaction/%s", serverFct, args[0])
    postCmd(str)
    
}

func fctsubmit(args []string) {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 1 {
        fmt.Println("Missing Key")
        os.Exit(1) 
    } 
            
    s := struct{Transaction string}{args[0]}
    
    jdata, err := json.Marshal(s)
    if err != nil {
        fmt.Println("Submitt failed")
        os.Exit(1) 
    }
    
    str:=fmt.Sprintf("http://%s/v1/factoid-submit/%s", serverFct,bytes.NewBuffer(jdata))
    postCmd(str)
}