#!/bin/bash
# Script to instantiate chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=172.31.89.52:7002
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/example.reg.com/peers/peer1.example.reg.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=example-reg-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/example.reg.com/users/Admin@example.reg.com/msp
export ORDERER_ADDRESS=172.31.89.52:7006
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/reg.com/orderers/orderer1.reg.com/tls/ca.crt
SID=$(peer lifecycle chaincode querycommitted -C autochannel -O json \
  | jq -r '.chaincode_definitions|.[]|select(.name=="KBA-Automobile")|.sequence' || true)

if [[ -z $SID ]]; then
  SEQUENCE=1
else
  SEQUENCE=$((1+$SID))
fi

peer lifecycle chaincode commit -o $ORDERER_ADDRESS --channelID autochannel \
  --name KBA-Automobile --version 1.0 --sequence $SEQUENCE \
  --peerAddresses 172.31.89.52:7002 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/example.reg.com/peers/peer1.example.reg.com/tls/ca.crt \
  --peerAddresses 172.31.89.52:7003 \
  --tlsRootCertFiles /vars/keyfiles/peerOrganizations/example2.reg.com/peers/peer1.example2.reg.com/tls/ca.crt \
  --collections-config /vars/KBA-Automobile_collection_config.json \
  --cafile $ORDERER_TLS_CA --tls
