#!/bin/bash
# Script to install chaincode onto a peer node
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=172.31.89.52:7002
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/example.reg.com/peers/peer1.example.reg.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=example-reg-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/example.reg.com/users/Admin@example.reg.com/msp
cd /go/src/github.com/chaincode/KBA-Automobile


if [ ! -f "KBA-Automobile_go_1.0.tar.gz" ]; then
  cd go
  GO111MODULE=on
  go mod vendor
  cd -
  peer lifecycle chaincode package KBA-Automobile_go_1.0.tar.gz \
    -p /go/src/github.com/chaincode/KBA-Automobile/go/ \
    --lang golang --label KBA-Automobile_1.0
fi

peer lifecycle chaincode install KBA-Automobile_go_1.0.tar.gz
