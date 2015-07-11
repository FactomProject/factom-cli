#!/bin/bash
# example script using factom-cli to write data into factom and get it back

echo "Using the cli to write new data into Factom"

echo "creating factoid and entry credit addresses"
factom-cli generateaddress fct michael
factom-cli generateaddress ec app

echo "buying entry credits with factoid address"
factom-cli newtransaction ecpurchase
factom-cli addinput ecpurchase michael 2
factom-cli addecoutput ecpurchase app 1.8
factom-cli sign ecpurchase
factom-cli submit ecpurchase
echo "sleep for 1 minute to wait for block creation (production system will have 10 minute blocks)"
sleep 1m
factom-cli getaddresses

echo "creating a new chain for my entries"
factom-cli mkchain -e app1 -e mychain app <<FIRSTENTRY
Hello Factom!
FIRSTENTRY

echo "sleep for 1 minute to wait for the chain to be included in the block (procution system will have 10 minute blocks)"
sleep 1m

echo "adding another entry to my chain"
factom-cli put -c 9e54c63c6ccf2f1e7bb6e86a4e026b63c5665dca2b649c1cb407d2e39d7e83f3 -e entry -e 2 app <<ENTRY
Hello again!
ENTRY

echo "sleep for 1 minute to wait for the entry to be included in the block (procution system will have 10 minute blocks)"
sleep 1m

echo "get the newest entry block for my chain"
eblock=$(factom-cli get chain 9e54c63c6ccf2f1e7bb6e86a4e026b63c5665dca2b649c1cb407d2e39d7e83f3)
factom-cli get eblock $eblock
ehash=$(factom-cli get eblock $eblock | grep EntryHash | awk '{print $2}')

echo "get the entry out of the entry block"
factom-cli get entry $ehash
