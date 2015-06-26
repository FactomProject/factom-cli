echo "#   If the addresses bill, bob, and jill are already created, this" 
echo "#   script will give a Duplicate or Invalid address error."
echo "#   This error does not prevent the construction and submission"
echo "#   of the transaction"

echo "#   Building Addresses for bill and bob, and entry credit address for jill"
factom-cli generateaddress fct bill
factom-cli generateaddress fct bob
factom-cli generateaddress ec  jill

echo "#   Looking at the balance for bill's address" 
factom-cli balance fct bill

echo "#   Looking at the balance for bob's address"
factom-cli balance fct bob

echo "#   Creating a new Transaction"
factom-cli newtransaction newtrans

echo "#   Add an input to the PayToBob transaction from bill"
factom-cli addinput newtrans bill 20

echo "#   Add an output to the PayToBob transaction to bob"
factom-cli addoutput newtrans bob 9

echo "#   By some Entry Credits for Jill"
factom-cli addecoutput newtrans jill 9

echo "#   Sign PayToBob"
factom-cli sign newtrans

echo "#   Submit PayToBob"
factom-cli submit newtrans

echo "#   Get the balance for bill"
factom-cli balance fct bill

echo "#   Get the balance for bob"
factom-cli balance fct bob

echo "#   Get the balance for jill"
factom-cli balance ec jill

echo "#   Get all balances for addresses in this wallet"
factom-cli getaddresses

