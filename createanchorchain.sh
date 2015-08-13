#!/bin/bash
# example script using factom-cli to write data into factom and get it back

echo "creating the anchor chain"
factom-cli mkchain -e FactomAnchorChain dan <<FIRSTENTRY
This is the Factom anchor chain, which records the anchors Factom puts on Bitcoin and other networks.
FIRSTENTRY

