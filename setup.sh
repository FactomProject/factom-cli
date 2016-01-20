
echo ">factom-cli newaddress ec zeros Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG"
factom-cli newaddress ec zeros Es2Rf7iM6PdsqfYCo3D1tnAR65SkLENyWJG1deUzpRMQmbh9F3eG

echo ">factom-cli newaddress fct sand Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK"
factom-cli newaddress fct sand Fs3E9gV6DXsYzf7Fqx1fVBQPQXV695eP3k5XbmHEZVRLkMdD9qCK

echo ">factom-cli newtransaction trans1"
factom-cli newtransaction trans1

echo ">factom-cli addinput trans1 sand 10"
factom-cli addinput trans1 sand 10

echo ">factom-cli addecoutput trans1 zeros 10"
factom-cli addecoutput trans1 zeros 10

echo ">factom-cli addfee trans1 sand"
factom-cli addfee trans1 sand

echo ">factom-cli sign trans1"
factom-cli sign trans1

echo ">factom-cli transactions"
factom-cli transactions

echo ">factom-cli submit trans1"
factom-cli submit trans1


echo ">this may take a while, scanning the blockchain for balances"
echo ">factom-cli balances"
factom-cli balances
