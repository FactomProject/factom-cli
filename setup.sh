
factom-cli generateaddress ec ec01
factom-cli generateaddress ec ec02
factom-cli generateaddress ec ec03
factom-cli generateaddress ec ec04
factom-cli generateaddress ec ec05
factom-cli generateaddress ec ec06
factom-cli generateaddress ec ec07
factom-cli generateaddress ec ec08
factom-cli generateaddress ec ec09

factom-cli generateaddress fct fct01
factom-cli generateaddress fct fct02
factom-cli generateaddress fct fct03
factom-cli generateaddress fct fct04
factom-cli generateaddress fct fct05
factom-cli generateaddress fct fct06
factom-cli generateaddress fct fct07
factom-cli generateaddress fct fct08
factom-cli generateaddress fct fct09
factom-cli generateaddress fct fct10


factom-cli newtransaction newtrans
factom-cli addinput newtrans 01-Fountain 1000
factom-cli addinput newtrans 02-Fountain 1000
factom-cli addinput newtrans 03-Fountain 1000
factom-cli addinput newtrans 04-Fountain 1000
factom-cli addinput newtrans 05-Fountain 1000

factom-cli addecoutput newtrans ec01  1000
factom-cli addoutput   newtrans fct01 1000
factom-cli addoutput   newtrans fct02 1000
factom-cli addoutput   newtrans fct03 1000
factom-cli addoutput   newtrans fct04 1000

factom-cli addfee newtrans 01-Fountain 
factom-cli sign newtrans
factom-cli transactions
factom-cli submit newtrans

factom-cli balances
