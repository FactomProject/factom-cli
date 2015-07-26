echo "
#=====================================================================#
#   Setting up the Wallet                                             #
#=====================================================================#"
factom-cli setup trans.sh

echo "
#=====================================================================#
#   If the addresses bill, bob, dan, sally and jane are already       #
#   created, this script will give 'Generation Failed' errors.        #
#                                                                     #
#   These error(s) do not prevent the construction and submission     # 
#   of the transaction                                                #
#=====================================================================#
factom-cli generateaddress fct bill
factom-cli generateaddress fct bob
factom-cli generateaddress fct sally
factom-cli generateaddress fct george
factom-cli generateaddress ec  dan
factom-cli generateaddress ec  jane
factom-cli getaddresses"
factom-cli generateaddress fct bill
factom-cli generateaddress fct bob
factom-cli generateaddress fct sally
factom-cli generateaddress fct george
factom-cli generateaddress ec  dan
factom-cli generateaddress ec  jane
factom-cli getaddresses

echo "
#=====================================================================#
# Fund a few addresses                                                #
#=====================================================================#"
factom-cli deletetransaction fund
factom-cli newtransaction fund
factom-cli addoutput fund bill 50
factom-cli addoutput fund george 50

factom-cli addinput fund 05-Fountain 100
factom-cli addfee fund 05-Fountain
factom-cli sign fund
factom-cli submit fund

echo "
#=====================================================================#
#   bill sends some money to bob, and also buys some Entry Credits    #
#   for jane and dan                                                  #
#                                                                     #
#   A transaction under construction has a key.  In this              #
#   case, the key is 'newtrans'.  Multiple transactions               #
#   can be under construction at the same time, so you have           #
#   to use the key when you add inputs and outputs.                   #
#   You can add multiple inputs and outputs to the same transaction.  #
#=====================================================================#
factom-cli deletetransaction newtrans
factom-cli newtransaction newtrans
factom-cli addinput newtrans bill 21.54321
factom-cli addoutput newtrans bob 1.5
factom-cli addecoutput newtrans jane 9.321
factom-cli addecoutput newtrans dan 9.321
factom-cli sign newtrans
factom-cli submit newtrans
factom-cli getaddresses"
factom-cli deletetransaction newtrans
factom-cli newtransaction newtrans
factom-cli addinput newtrans bill 21.54321
factom-cli addoutput newtrans bob 1.5
factom-cli addecoutput newtrans jane 9.321
factom-cli addecoutput newtrans dan 9.321
factom-cli sign newtrans
factom-cli submit newtrans
factom-cli getaddresses

echo "
#=====================================================================#
#  george wants to spread the wealth.  A pure Factoid transaction,    #
#  he sends bill, bob, and sally some factoids.                       #
#                                                                     #
#  This transaction is named '2ndTrans'                               #
#=====================================================================#
factom-cli deletetransaction 2ndTrans
factom-cli newtransaction 2ndTrans
factom-cli addinput 2ndTrans george 5
factom-cli addoutput 2ndTrans bill 2
factom-cli addoutput 2ndTrans bob 0.00005
factom-cli addoutput 2ndTrans sally 2
factom-cli sign 2ndTrans
factom-cli submit 2ndTrans
factom-cli getaddresses"
factom-cli deletetransaction 2ndTrans
factom-cli newtransaction 2ndTrans
factom-cli addinput 2ndTrans george 5
factom-cli addoutput 2ndTrans bill 2
factom-cli addoutput 2ndTrans bob 0.00005
factom-cli addoutput 2ndTrans sally 2
factom-cli sign 2ndTrans
factom-cli submit 2ndTrans
factom-cli getaddresses

