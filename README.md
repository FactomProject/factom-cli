Factom Command Line Interface
===

Synopsis
---
	factom-cli [options] [subcommand]
		-h						"print usage information"
		-s [server]				"address for the api server"
		-w [wallet]				"address for the wallet"
		
		balance
			ec					"entry credit balance of 'wallet'"
			factoid				"factoid balance of 'wallet'"

		buy
			amt #n				"buy n entry credits for 'wallet'"

		fatoidtx [addr] [amt]	"create and submit a factoid transaction"

		get						"not yet defined"

		help [command]			"print help message for the sub-command"

		mkchain [opt] [name]	"create a new factom chain with 'name'. read"
								"the data for the first entry from stdin"
			-e externalid		"externalid for the first entry

		put						"read data from stdin and write to factom"
			-e [externalid]		"specify an exteral id for the factom entry. -e" 								"can be used multiple times"
			-c [chainid]		"spesify the chain that the entry belongs to"
			
Description
---
factom-cli is the command line interface to the factom api

Examples
---
	# Submit arbitrary data to factom
	echo "Some data to save in factom" | factom-cli put -c [chainid] -e mydata -e somekey
	
	# Submit a file
	factom-cli -s "facom.org/demo" put -c [chainid] -e filename <file
	
	factom-cli balance ec
	factom-cli factoidtx [address] 100

Files
---
	factom-cli.conf	"factom-cli will try and read the conf file from
					$HOME/.factom/ If the conf file does not exist it will use
					the default configuration. Command line flags will
					overwrite the configurations."