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

