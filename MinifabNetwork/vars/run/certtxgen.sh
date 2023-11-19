#!/bin/bash
cd $FABRIC_CFG_PATH
# cryptogen generate --config crypto-config.yaml --output keyfiles
configtxgen -profile OrdererGenesis -outputBlock genesis.block -channelID systemchannel

configtxgen -printOrg example-reg-com > JoinRequest_example-reg-com.json
configtxgen -printOrg example2-reg-com > JoinRequest_example2-reg-com.json
