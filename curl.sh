curl -X POST http://localhost:8089/v1/factoid-new-transaction/t
curl -X POST http://localhost:8089/v1/factoid-add-input/ -d "key=t&name=ps&amount=10000000"
curl -X POST http://localhost:8089/v1/factoid-add-output/  -d "key=t&name=FA2dAYismYSSaT5yopvquNm7e15KG8KVyYVkMDxgs5XrmiY4wERb&amount=10000000"
curl -X POST http://localhost:8089/v1/factoid-add-fee/ -d "key=t&name=ps"
curl -X POST http://localhost:8089/v1/factoid-sign-transaction/t
curl -X POST http://localhost:8089/v1/factoid-submit/\\{\"Transaction\":\"t\"\\}
curl -X GET http://localhost:8089/v1/factoid-get-addresses/ 

