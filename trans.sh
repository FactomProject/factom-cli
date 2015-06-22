factom-cli generateaddress fct bill
factom-cli generateaddress fct bob
factom-cli newtransaction PayToBob
factom-cli balance factoid bill
factom-cli balance factoid bob
factom-cli addinput PayToBob bill 10
factom-cli addoutput PayToBob bob 9
factom-cli sign PayToBob
echo "Submit Transaction for Bill to pay Bob"
factom-cli submit PayToBob
factom-cli balance factoid bill
factom-cli balance factoid bob
factom-cli getaddresses

