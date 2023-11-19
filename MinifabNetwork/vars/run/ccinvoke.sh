#!/bin/bash
# Script to invoke chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=172.31.89.52:7002
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/example.reg.com/peers/peer1.example.reg.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=example-reg-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/example.reg.com/users/Admin@example.reg.com/msp
export ORDERER_ADDRESS=172.31.89.52:7006
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/reg.com/orderers/orderer1.reg.com/tls/ca.crt
peer chaincode invoke -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA \
  --tls -C autochannel -n KBA-Automobile  \
  --peerAddresses 172.31.89.52:7003 \
  --tlsRootCertFiles /vars/discover/autochannel/example2-reg-com/tlscert \
  --peerAddresses 172.31.89.52:7002 \
  --tlsRootCertFiles /vars/discover/autochannel/example-reg-com/tlscert \
  -c '{"Args":["CreateLand","land02","Johny Depp"]}'
