echo "If the addresses bill and bob are already created, this" 
echo "script will give a Duplicate or Invalid address error."
echo "This error does not prevent the construction and submission"
echo "of the transaction"
echo ""
echo "Building Addresses for bill and bob"
factom-cli generateaddress fct bill
factom-cli generateaddress fct bob
echo "Looking at the balance for bill's address" 
factom-cli balance bill
echo "Looking at the balance for bob's address"
factom-cli balance bob
echo "Creating a new Transaction"
factom-cli newtransaction PayToBob
echo "Add an input to the PayToBob transaction from bill"
factom-cli addinput PayToBob bill 10
echo "Add an output to the PayToBob transaction to bob"
factom-cli addoutput PayToBob bob 9
echo "Sign PayToBob"
factom-cli sign PayToBob
echo "Submit PayToBob"
factom-cli submit PayToBob
echo "Get the balance for bill"
factom-cli balance bill
echo "Get the balance for bob"
factom-cli balance bob
echo "Get all balances for addresses in this wallet"
factom-cli getaddresses

