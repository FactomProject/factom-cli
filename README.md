Factom Command Line Interface
===

Synopsis
---
	factom
		-c [ChainID]	"hex encoded chainid for the entry"
		-e [extid]		"external id for the entry"			# -e may be used multiple times
		-h				"display help message"
		-s				"path to the factomclient"			# default: localhost:8083

Description
---
factom-cli takes data from the command arguments and stdin and constructs a json entry for the factomclient to decode and submit to factom.

Examples
---
Submiting arbitrary data to factom:
	$ echo "Some data to save in factom" | factom-cli -c <my chainid> -e mydata -e somekey
Submitting a file:
	$ factom-cli -s "facom.org/demo" -c <my chainid> -e filename <file