package test

import (
	"fmt"
	"testing"
	"os"
	"encoding/base64"
	"encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/kba-chf/kba-automobile/chaincode/contracts"
	"github.com/kba-chf/kba-automobile/chaincode/test/mocks"
	"github.com/stretchr/testify/require"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext ./car-contract_test.go transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub ./car-contract_test.go chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator ./car-contract_test.go stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

//go:generate counterfeiter -o mocks/clientIdentity.go -fake-name ClientIdentity . clientIdentity
type clientIdentity interface {
	cid.ClientIdentity
}

const myOrg1Msp string = "Org1MSP"
const myOrg1Clientid string = "Org1MSP"

func prepMocksAsOrg1() (*mocks.TransactionContext, *mocks.ChaincodeStub) {
	return prepMocks(myOrg1Msp, myOrg1Clientid)
}

func prepMocks(orgMSP, clientId string) (*mocks.TransactionContext, *mocks.ChaincodeStub) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	clientIdentity := &mocks.ClientIdentity{}
	clientIdentity.GetMSPIDReturns(orgMSP, nil)
	clientIdentity.GetIDReturns(base64.StdEncoding.EncodeToString([]byte(clientId)), nil)
	// set matching msp ID using peer shim env variable
	os.Setenv("CORE_PEER_LOCALMSPID", orgMSP)
	transactionContext.GetClientIdentityReturns(clientIdentity)
	return transactionContext, chaincodeStub
}

func TestCreateAsset(t *testing.T) {
	transactionContext, chaincodeStub := prepMocksAsOrg1()

	// chaincodeStub := &mocks.ChaincodeStub{}
	// transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	carAsset := contracts.CarContract{}
	_, err := carAsset.CreateCar(transactionContext, "car1", "", "", "", "", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	_, err = carAsset.CreateCar(transactionContext, "car1", "", "", "", "", "")
	require.EqualError(t, err, "the asset car1 already exists")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = carAsset.CreateCar(transactionContext, "car1", "", "", "", "", "")
	require.EqualError(t, err, "could not read from world state. failed to read from world state: unable to retrieve asset")
}

func TestReadAsset(t *testing.T) {
	transactionContext, chaincodeStub := prepMocksAsOrg1()
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &contracts.Car{CarId: "car1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	assetTransfer := contracts.CarContract{}
	asset, err := assetTransfer.ReadCar(transactionContext, "car1")
	require.NoError(t, err)
	require.Equal(t, expectedAsset, asset)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = assetTransfer.ReadCar(transactionContext, "")
	require.EqualError(t, err, "could not read from world state. failed to read from world state: unable to retrieve asset")

	chaincodeStub.GetStateReturns(nil, nil)
	asset, err = assetTransfer.ReadCar(transactionContext, "asset1")
	require.EqualError(t, err, "the asset asset1 does not exist")
	require.Nil(t, asset)
}

func TestDeleteAsset(t *testing.T) {
	transactionContext, chaincodeStub := prepMocksAsOrg1()
	transactionContext.GetStubReturns(chaincodeStub)

	asset := &contracts.Car{CarId: "car1"}
	bytes, err := json.Marshal(asset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	// chaincodeStub.DelStateReturns(nil)
	assetTransfer := contracts.CarContract{}
	_, err = assetTransfer.DeleteCar(transactionContext, "car1")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	_, err = assetTransfer.DeleteCar(transactionContext, "car1")
	require.EqualError(t, err, "the asset car1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = assetTransfer.DeleteCar(transactionContext, "")
	require.EqualError(t, err, "could not read from world state. failed to read from world state: unable to retrieve asset")
}