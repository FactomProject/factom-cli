// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
    "strings"
    "strconv"
    "bytes"
    "net/http"
    "encoding/hex"
    "encoding/json"
    "io/ioutil"
    "flag"
    "fmt"
    "os"
    
    "github.com/FactomProject/factom"
)

var _ = hex.EncodeToString
var serverFct = "localhost:8089"

func getCmd(cmd string, cmderror string) error {
	resp, err := http.Get(cmd)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body) 
	
	if err != nil {
		return err
	}
	resp.Body.Close()

	type x struct{ 
        Body string
        Success bool 
    }
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil || !b.Success {
		fmt.Println(cmderror)
		return fmt.Errorf("Command Failed: ", string(body))
	}
    fmt.Println(b.Body)
	return nil
}

func postCmd(cmd string, cmderror string) error {
	resp, err := http.PostForm(cmd, nil)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()

	type x struct{ Success bool }
	b := new(x)
	if err := json.Unmarshal(body, b); err != nil || !b.Success {
		fmt.Println(cmderror)
		return fmt.Errorf("Command Failed: ", string(body))
	}

	return nil
}

// Generates a new Address
func generateaddress(args []string) (err error) {
    
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 2 {
        return man("generatefactoidaddress")
    }
    
    var Addr string
    switch args[0]{
        case "ec": 
            Addr,err = factom.GenerateEntryCreditAddress(args[1])
        case "fct":
            Addr,err = factom.GenerateFactoidAddress(args[1])
        default:
            panic("Expected ec|fct name")
    }
    
    if err != nil {
        panic(err)
    }
    fmt.Println(args[0]," = ",Addr)

    return nil
    
}

func getaddresses(args []string) (err error){
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) > 0 {
        return man("getaddresses")
    }
    
    str := fmt.Sprintf("http://%s/v1/factoid-get-addresses/", serverFct)
    err  = getCmd(str,"Error printing addresses")
    
    return err
    
    
}

func fctnewtrans(args []string) error {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 1 {
        fmt.Println("Missing Key")
        return fmt.Errorf("Missing Key")
    } 
    
    str := fmt.Sprintf("http://%s/v1/factoid-new-transaction/%s", serverFct, args[0])
    err := postCmd(str, "Duplicate or bad key")
    
    return err
    
}

func convertFixedPoint(amt string) (string, error) {
    var v int64
    var err error
    index := strings.Index(amt,".");
    if  index < 0 {
        v, err =strconv.ParseInt(amt,10,64)
        if err != nil {
            fmt.Println("1:",err)
            return "", err
        }
        v *= 100000000      // Convert to Factoshis
    }else{
        tp := amt[:index]
        v, err =strconv.ParseInt(tp,10,64)
        if err != nil {
            fmt.Println("2:",err)
            return "", err
        }
        v = v*100000000      // Convert to Factoshis
        
        bp := amt[index+1:]
        if len(bp)>8 {
            fmt.Println("3: Bad length")
            return "", fmt.Errorf("Format Error in amount")
        }
        bpv, err :=strconv.ParseInt(bp,10,64)
        if err != nil {
            fmt.Println("4:",err)
            return "", err
        }
        for i:=0; i<8-len(bp); i++ {
            bpv *= 10
        }
        v += bpv
    }
    return strconv.FormatInt(v,10),nil
}

func fctaddinput(args []string) error {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 3 {
        fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
        return fmt.Errorf("Missing Arguments")
    } 
    // localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>
    
    
    amt,err := convertFixedPoint(args[2])
    if err != nil { return err }
    str := fmt.Sprintf("http://%s/v1/factoid-add-input/?key=%s&name=%s&amount=%s", 
                       serverFct, args[0],args[1],amt)
    err = postCmd(str,"Failed to add input")
    
    return err
}

func fctaddoutput(args []string) error {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 3 {
        fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
        return fmt.Errorf("Missing Arguments")
    } 
    // localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>
    
    amt,err := convertFixedPoint(args[2])
    if err != nil { return err }
    str := fmt.Sprintf("http://%s/v1/factoid-add-output/?key=%s&name=%s&amount=%s", 
                       serverFct, args[0],args[1],amt)
    err = postCmd(str,"Failed to add output")
    
    return err
}

func fctaddecoutput(args []string) error {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 3 {
        fmt.Println("Expecting a 1) transaction key, 2) an Address or Address name, and 3) an amount.")
        return fmt.Errorf("Missing Arguments")
    } 
    // localhost:8089/v1/factoid-add-input/?key=<key>&name=<name or address>&amount=<amount>
    
    amt,err := convertFixedPoint(args[2])
    if err != nil { return err }
    str := fmt.Sprintf("http://%s/v1/factoid-add-ecoutput/?key=%s&name=%s&amount=%s", 
                       serverFct, args[0],args[1],amt)
    err = postCmd(str,"Failed to add Entry Credit output")
    
    return err
}

func fctgetfee(args []string) error {
    resp, err := http.Get(fmt.Sprintf("http://%s/v1/factoid-get-fee/",serverFct))
    if err != nil {
        fmt.Println("Command Failed Get")
        return err
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Command Failed")
        return err
    }
    resp.Body.Close()
    
    // We pull the fee.  If the fee isn't positive, or if we fail to marshal, then there is a failure
    type x struct { Fee int64 }
    b := new(x)
    b.Fee = -1
    if err := json.Unmarshal(body, b); err != nil || b.Fee == -1 {
        fmt.Println("Command Failed")
        return fmt.Errorf("Command Failed")
    }
    tv := b.Fee/100000000
    lv := b.Fee-(tv*100000000)
    r := fmt.Sprintf("Fee: %d.%08d",tv,lv)
    var i int; for i=len(r)-1; r[i]=='0'; i-- {}
    if string(r[i])=="." { i +=1 }
    fmt.Println(r[:i+1])
    return nil
}
    
func fctsign(args []string) error {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 1 {
        fmt.Println("Missing Key")
        return fmt.Errorf("Missing Key")
    } 
    
    str := fmt.Sprintf("http://%s/v1/factoid-sign-transaction/%s", serverFct, args[0])
    err := postCmd(str,"Cannot sign transaction.  Check balances, inputs, transaction fees")
    
    return err
}

func fctsubmit(args []string) error {
    os.Args = args
    flag.Parse()
    args = flag.Args()
    if len(args) < 1 {
        fmt.Println("Missing Key")
        return fmt.Errorf("Missing Key")
    } 
            
    s := struct{Transaction string}{args[0]}
    
    jdata, err := json.Marshal(s)
    if err != nil {
        fmt.Println("Submitt failed")
        return fmt.Errorf("Submit failed")
    }
    
    resp, err := http.Post(
        fmt.Sprintf("http://%s/v1/factoid-submit/", serverFct),
                           "application/json",
                           bytes.NewBuffer(jdata))
    if err != nil {
        fmt.Println("Submitt failed")
        return fmt.Errorf("Error returned by fctwallet")
    }
    resp.Body.Close()
    return nil
}