#!/bin/bash

set -ev

#chaincode install
docker exec cli peer chaincode install -n albachain -v 1.0 -p github.com/albachain
#chaincode instatiate
docker exec cli peer chaincode instantiate -o orderer.example.com:7050 -n albachain -v 1.0 -C mychannel -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member","Org3MSP.member")'
sleep 5
#chaincode invoke user1
docker exec cli peer chaincode invoke -n albachain -C mychannel -c '{"Args":["addWorker","user1", "bob"]}'
sleep 5
#chaincode query user1
docker exec cli peer chaincode query -n albachain -C mychannel -c '{"Args":["getWorker","user1"]}'

#chaincode invoke add rating
#docker exec cli peer chaincode invoke -n albachain -C mychannel -c '{"Args":["addWorkplace","user1","p001"]}'
sleep 5

echo '-------------------------------------END-------------------------------------'
