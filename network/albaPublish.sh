#!/bin/bash

set -ev

#chaincode install
docker exec cli peer chaincode install -n albachain -v 1.0 -p github.com/albachain
#chaincode instatiate
docker exec cli peer chaincode instantiate -o orderer.hydrogreen.com:7050 -n albachain -v 1.0 -C mygreen -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member","Org3MSP.member")'

#sleep 2
#docker exec cli peer chaincode invoke -n albachain -C mygreen -c '{"Args":["addWorker","user1", "bob"]}'
#sleep 2
#docker exec cli peer chaincode invoke -n albachain -C mygreen -c '{"Args":["addEmployer","employer1", "sajang"]}'
#sleep 2
#docker exec cli peer chaincode invoke -n albachain -C mygreen -c '{"Args":["addWorkplaceByEmployer","employer1", "sajang", "CU0001", "Kangnam Gu", "7000"]}'
#sleep 2
#docker exec cli peer chaincode invoke -n albachain -C mygreen -c '{"Args":["addWorkplace","user1", "CU0001", "Kangnam Gu", "7000"]}'
#sleep 2
#docker exec cli peer chaincode invoke -n albachain -C mygreen -c '{"Args":["workplaceApproval","user1", "CU0001", "7000"]}'
#sleep 3
#docker exec cli peer chaincode invoke -n albachain -C mygreen -c '{"Args":["addWorkHistory","WH001","user1","bob","CU0001","Kangnam Gu", "20211126 10:00", "20211126 20:00","7000"]}'
#sleep 3
#docker exec cli peer chaincode invoke -n albachain -C mygreen -c '{"Args":["workHistoryApproval","WH001","7000"]}'
#sleep 3
#docker exec cli peer chaincode invoke -n albachain -C mygreen -c '{"Args":["deleteWorkplace","user1", "CU0001"]}'
#sleep 2
#docker exec cli peer chaincode invoke -n albachain -C mygreen -c '{"Args":["getWorker","user1"]}'
#sleep 2

echo '-------------------------------------END-------------------------------------'
