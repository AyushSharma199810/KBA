package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kba-chf/kba-automobile/chaincode/contracts"
)

func main() {
	LandContract := new(contracts.LandContract)
	orderContarct := new(contracts.OrderContract)

	chaincode, err := contractapi.NewChaincode(LandContract, orderContarct)

	if err != nil {
		log.Panicf("Could not create chaincode." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		log.Panicf("Failed to start chaincode. " + err.Error())
	}
}
