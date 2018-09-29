#! /bin/bash

# cryptogen generate --config crypto-config.yaml

mkdir channel-artifacts
configtxgen -profile OneOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

configtxgen -profile OneOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel

configtxgen -profile OneOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP