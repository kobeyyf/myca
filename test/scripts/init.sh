#! /bin/bash
ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/dctest.com/orderers/orderer0.dctest.com/msp/tlscacerts/tlsca.dctest.com-cert.pem

peer channel create -o orderer0.dctest.com:7050 -c mychannel -f ./channel-artifacts/channel.tx --tls true --cafile $ORDERER_CA

peer channel join -b mychannel.block

peer channel update -o orderer0.dctest.com:7050 -c mychannel -f ./channel-artifacts/Org1MSPanchors.tx --tls true --cafile $ORDERER_CA

peer chaincode install -n 1_org_1_peer -v 1.0 -p yang.com/fabric/chaincode/

peer chaincode instantiate -o orderer0.dctest.com:7050 --tls true --cafile $ORDERER_CA -C mychannel -n 1_org_1_peer -v 1.0 -c '{"Args":["init","512B"]}' -P "OR ('Org1MSP.member')"

peer chaincode query -C mychannel -n 1_org_1_peer -c '{"Args":["query","0"]}'
