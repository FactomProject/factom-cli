factom-cli setup trans.sh
factom-cli generateaddress ec app
factom-cli newtransaction newtrans
factom-cli addinput newtrans 01-Fountain 50
factom-cli addecoutput newtrans app 49
factom-cli sign newtrans
factom-cli transactions
factom-cli submit newtrans
factom-cli balances
