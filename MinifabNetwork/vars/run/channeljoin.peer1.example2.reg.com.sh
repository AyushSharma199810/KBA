#!/bin/bash
# Script to join a peer to a channel
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=172.31.89.52:7003
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/example2.reg.com/peers/peer1.example2.reg.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=example2-reg-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/example2.reg.com/users/Admin@example2.reg.com/msp
export ORDERER_ADDRESS=172.31.89.52:7006
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/reg.com/orderers/orderer1.reg.com/tls/ca.crt
if [ ! -f "autochannel.genesis.block" ]; then
  peer channel fetch oldest -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA \
  --tls -c autochannel /vars/autochannel.genesis.block
fi

peer channel join -b /vars/autochannel.genesis.block \
  -o $ORDERER_ADDRESS --cafile $ORDERER_TLS_CA --tls
