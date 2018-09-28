#! /bin/bash

# ../../fabric-test/utils/bin/cryptogen generate --config crypto-config.yaml

mkdir channel-artifacts
../../fabric-test/utils/bin/configtxgen -profile OneOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block

../../fabric-test/utils/bin/configtxgen -profile OneOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID mychannel

../../fabric-test/utils/bin/configtxgen -profile OneOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID mychannel -asOrg Org1MSP