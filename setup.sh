
echo -e "\t>factom-cli importaddress zeros Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
factom-cli importaddress zeros Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG

echo -e "\t>factom-cli importaddress sand Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"
factom-cli importaddress sand Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK

echo -e "\t>factom-cli newtransaction trans1"
factom-cli newtransaction trans1

echo -e "\t>factom-cli addinput trans1 sand 10"
factom-cli addinput trans1 sand 10

echo -e "\t>factom-cli addecoutput trans1 zeros 10"
factom-cli addecoutput trans1 zeros 10

echo -e "\t>factom-cli addfee trans1 sand"
factom-cli addfee trans1 sand

echo -e "\t>factom-cli sign trans1"
factom-cli sign trans1

echo -e "\t>factom-cli transactions"
factom-cli transactions

echo -e "\t>factom-cli submit trans1"
factom-cli submit trans1


echo -e "\t>this may take a while, scanning the blockchain for balances"
echo -e "\t>factom-cli balances"
factom-cli balances
