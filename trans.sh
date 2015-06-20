factom-cli genfactoidaddr bill
factom-cli genfactoidaddr bob
factom-cli newtransaction PayToBob
factom-cli balance factoid bill
factom-cli balance factoid bob
factom-cli addinput PayToBob bill 1000000
factom-cli addoutput PayToBob bob 500000
factom-cli sign PayToBob
echo "Submit Transaction for Bill to pay Bob"
factom-cli submit PayToBob
factom-cli balance factoid bill
factom-cli balance factoid bob

