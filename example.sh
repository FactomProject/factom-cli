#!/bin/bash
# example script using factom-cli to write data into factom and get it back

if [[ -z $1 ]]; then
echo "
************************************************
*  Defaulting to 1 minute waits.  But if you
*  are on a 1 minute block server, you can run:
*
*  ./example.sh 10s
*
*  Which will only wait 10 seconds.  Or if you
*  are local, you can specify 0s, which will run
*  way faster!
*
*  Time can be specified like 10x where x can be:
*  s =seconds  m =minute  h = hours   d = days 
************************************************"

st=60s
else
st=$1
fi

echo "Using the cli to write new data into Factom"

echo "Setup the Wallet"
factom-cli setup Setup_For_example.sh

echo "creating factoid and entry credit addresses"
factom-cli generateaddress ec app

echo "buying entry credits with factoid address"
factom-cli newtransaction ecpurchase
factom-cli addinput ecpurchase 01-Fountain 2
factom-cli addecoutput ecpurchase app 2
factom-cli addfee ecpurchase 01-Fountain
factom-cli sign ecpurchase
factom-cli submit ecpurchase
echo "
sleep for " $st " to wait for block creation 
(production system will have 10 minute blocks)
"
sleep $st
factom-cli getaddresses

echo "creating a new chain for my entries"
factom-cli mkchain -e app1 -e mychain app <<FIRSTENTRY
Hello Factom!
FIRSTENTRY

echo "
sleep for " $st " to wait for the chain to be included 
in the block (production system will have 10 minute blocks)
"
sleep $st

echo "adding another entry to my chain"
factom-cli put -c 9e54c63c6ccf2f1e7bb6e86a4e026b63c5665dca2b649c1cb407d2e39d7e83f3 -e entry -e 2 app <<ENTRY
Hello again!
ENTRY

echo "
sleep for " $st " to wait for the entry to be included 
in the block (production system will have 10 minute blocks)
"
sleep $st

echo "get the newest entry block for my chain"
eblock=$(factom-cli get chain 9e54c63c6ccf2f1e7bb6e86a4e026b63c5665dca2b649c1cb407d2e39d7e83f3)
factom-cli get eblock $eblock
ehash=$(factom-cli get eblock $eblock | grep EntryHash | awk '{print $2}')

echo "get the entry out of the entry block"
factom-cli get entry $ehash

