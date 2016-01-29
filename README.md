Factom Command Line Interface
===

Synopsis
------
---
	factom
		-c [ChainID]   # hex encoded chainid for the entry
		-e [extid]     # external id for the entry. -e may be used multiple times
		-h             # display help message
		-s             # path to the factomclient. default: localhost:8083
---
Description
-------
factom-cli takes data from the command arguments and stdin and constructs a json entry for the factomclient to decode and submit to factom.

Examples
-------
---
Getting the Help screen:

	$ factom-cli -h
	
Submiting arbitrary data to factom:

	$ echo "Some data to save in factom" | factom-cli -c <my chainid> -e mydata -e somekey

Submitting a file:

	$ factom-cli -s "factom.org/demo" -c <my chainid> -e filename <file

Getting a list of all Factom and Entry Credit addresses and balances in your wallet:

	$ factom-cli getaddresses
	
Getting the balance of any address (either by name, or by address):

	$ factom-cli balance fct MyAddressName
	$ factom-cli balance fct FA3Y6ZZbiuCjrQ5WKLFq3GaEi9drsyTJQagFVLmUU5pfPH33Dgqg
	$ factom-cli balance ec MyEntryCreditAddress
	$ factom-cli balance ec EC2W1KAv9KevaUKqA25978yMU6EJ7yBfygRmzcdpa3oHzYZRKh17

Creating and submitting a transaction to move 10 factoids from MyAddress to BobsAddress 

	$ factom-cli newtransaction paybob                   # Create a transaction with key 'paybob'
	$ factom-cli addinput paybob MyAddress 10.1          # Add input plus fee to 'paybob'
	$ factom-cli addoutput paybob BobsAddress 10         # Bob gets 10 Factoids to 'paybob'
	$ factom-cli sign paybob                             # sign the 'paybob' transaction
	$ factom-cli submit paybob                           # submit the 'paybob' transaction
	
Purchase 1 factoid worth of Entry Credits and attach them to the Entry Credit Address myEC

	$ factom-cli newtransaction buyec                    # Create a transaction with key 'buyec'
	$ factom-cli addinput buyec MyAddress 1.1            # Add input plus fee to 'buyec'
	$ factom-cli addecoutput buyec myEC 1                # 1 Factoid worth of EC credits to myEC
	$ factom-cli sign buyec                              # sign 'buyec'
	$ factom-cli submit buyec                            # submit 'buyec'
	
