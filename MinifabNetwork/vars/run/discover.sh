#!/bin/bash
# Script to discover endorsers and channel config
cd /vars

export PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/example.reg.com/users/Admin@example.reg.com/tls/ca.crt
export ADMINPRIVATEKEY=/vars/keyfiles/peerOrganizations/example.reg.com/users/Admin@example.reg.com/msp/keystore/priv_sk
export ADMINCERT=/vars/keyfiles/peerOrganizations/example.reg.com/users/Admin@example.reg.com/msp/signcerts/Admin@example.reg.com-cert.pem

discover endorsers --peerTLSCA $PEER_TLS_ROOTCERT_FILE \
  --userKey $ADMINPRIVATEKEY \
  --userCert $ADMINCERT \
  --MSP example-reg-com --channel autochannel \
  --server 172.31.89.52:7002 \
  --chaincode KBA-Automobile | jq '.[0]' | \
  jq 'del(.. | .Identity?)' | jq 'del(.. | .LedgerHeight?)' \
  > /vars/discover/autochannel_KBA-Automobile_endorsers.json

discover config --peerTLSCA $PEER_TLS_ROOTCERT_FILE \
  --userKey $ADMINPRIVATEKEY \
  --userCert $ADMINCERT \
  --MSP example-reg-com --channel autochannel \
  --server 172.31.89.52:7002 > /vars/discover/autochannel_config.json
