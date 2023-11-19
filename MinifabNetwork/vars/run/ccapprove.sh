#!/bin/bash
# Script to approve chaincode
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ID=cli
export CORE_PEER_ADDRESS=172.31.89.52:7002
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/example.reg.com/peers/peer1.example.reg.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=example-reg-com
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/example.reg.com/users/Admin@example.reg.com/msp
export ORDERER_ADDRESS=172.31.89.52:7006
export ORDERER_TLS_CA=/vars/keyfiles/ordererOrganizations/reg.com/orderers/orderer1.reg.com/tls/ca.crt

peer lifecycle chaincode queryinstalled -O json | jq -r '.installed_chaincodes | .[] | select(.package_id|startswith("KBA-Automobile_1.0:"))' > ccstatus.json

PKID=$(jq '.package_id' ccstatus.json | xargs)
REF=$(jq '.references.autochannel' ccstatus.json)

SID=$(peer lifecycle chaincode querycommitted -C autochannel -O json \
  | jq -r '.chaincode_definitions|.[]|select(.name=="KBA-Automobile")|.sequence' || true)
if [[ -z $SID ]]; then
  SEQUENCE=1
elif [[ -z $REF ]]; then
  SEQUENCE=$SID
else
  SEQUENCE=$((1+$SID))
fi


export CORE_PEER_LOCALMSPID=example-reg-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/example.reg.com/peers/peer1.example.reg.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/example.reg.com/users/Admin@example.reg.com/msp
export CORE_PEER_ADDRESS=172.31.89.52:7002

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID autochannel \
#   --name KBA-Automobile --version 1.0 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.example-reg-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID autochannel --name KBA-Automobile \
    --version 1.0 --package-id $PKID \
    --collections-config /vars/KBA-Automobile_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi

export CORE_PEER_LOCALMSPID=example2-reg-com
export CORE_PEER_TLS_ROOTCERT_FILE=/vars/keyfiles/peerOrganizations/example2.reg.com/peers/peer1.example2.reg.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/vars/keyfiles/peerOrganizations/example2.reg.com/users/Admin@example2.reg.com/msp
export CORE_PEER_ADDRESS=172.31.89.52:7003

# approved=$(peer lifecycle chaincode checkcommitreadiness --channelID autochannel \
#   --name KBA-Automobile --version 1.0 --init-required --sequence $SEQUENCE --tls \
#   --cafile $ORDERER_TLS_CA --output json | jq -r '.approvals.example2-reg-com')

# if [[ "$approved" == "false" ]]; then
  peer lifecycle chaincode approveformyorg --channelID autochannel --name KBA-Automobile \
    --version 1.0 --package-id $PKID \
    --collections-config /vars/KBA-Automobile_collection_config.json \
    --sequence $SEQUENCE -o $ORDERER_ADDRESS --tls --cafile $ORDERER_TLS_CA
# fi
