package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// LandContract provides functions for managing the land registry
type LandContract struct {
	contractapi.Contract
}

type PaginatedQueryResult struct {
	Records             interface{} `json:"records"`
	FetchedRecordsCount int32       `json:"fetchedRecordsCount"`
	Bookmark            string      `json:"bookmark"`
}

type HistoryQueryResult struct {
	Record    interface{} `json:"record"`
	TxId      string      `json:"txId"`
	Timestamp string      `json:"timestamp"`
	IsDelete  bool        `json:"isDelete"`
}

type Asset struct {
	AssetType string `json:"assetType"`
}

type Land struct {
	Asset
	LandID   string `json:"landID"`
	OwnedBy  string `json:"ownedBy"`
	Location string `json:"location"` // Add additional property for location
}

// LandContract.LandExists returns true when the land with the given ID exists in the world state
func (lc *LandContract) LandExists(ctx contractapi.TransactionContextInterface, landID string) (bool, error) {
	data, err := ctx.GetStub().GetState(landID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return data != nil, nil
}

// LandContract.CreateLand creates a new instance of Land
func (lc *LandContract) CreateLand(ctx contractapi.TransactionContextInterface, landID, ownerName, location string) (string, error) {
	land := Land{
		Asset:    Asset{AssetType: "land"},
		LandID:   landID,
		OwnedBy:  ownerName,
		Location: location, // Set the location field
	}

	bytes, _ := json.Marshal(land)

	err := ctx.GetStub().PutState(landID, bytes)
	if err != nil {
		return "", err
	}

	// Emit event for Land creation
	err = ctx.GetStub().SetEvent("LandCreated", bytes)
	if err != nil {
		return "", fmt.Errorf("failed to emit LandCreated event: %v", err)
	}

	return fmt.Sprintf("successfully added land %v", landID), nil
}

// LandContract.UpdateLandOwnership transfers ownership of land to a new owner
func (lc *LandContract) UpdateLandOwnership(ctx contractapi.TransactionContextInterface, landID string, newOwner string) (string, error) {
	land, err := lc.ReadLand(ctx, landID)
	if err != nil {
		return "", err
	}

	land.OwnedBy = newOwner

	bytes, _ := json.Marshal(land)

	err = ctx.GetStub().PutState(landID, bytes)
	if err != nil {
		return "", err
	}

	// Emit event for Ownership update
	err = ctx.GetStub().SetEvent("OwnershipUpdated", bytes)
	if err != nil {
		return "", fmt.Errorf("failed to emit OwnershipUpdated event: %v", err)
	}

	return fmt.Sprintf("successfully updated ownership of land %v to %v", landID, newOwner), nil
}

// LandContract.ReadLand retrieves an instance of Land from the world state
func (lc *LandContract) ReadLand(ctx contractapi.TransactionContextInterface, landID string) (*Land, error) {
	exists, err := lc.LandExists(ctx, landID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state: %v", err)
	} else if !exists {
		return nil, fmt.Errorf("the asset %s does not exist", landID)
	}

	bytes, _ := ctx.GetStub().GetState(landID)

	land := new(Land)

	err = json.Unmarshal(bytes, &land)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal world state data to type Land")
	}

	return land, nil
}

// LandContract.DeleteLand removes the instance of Land from the world state
func (lc *LandContract) DeleteLand(ctx contractapi.TransactionContextInterface, landID string) (string, error) {
	exists, err := lc.LandExists(ctx, landID)
	if err != nil {
		return "", fmt.Errorf("could not read from world state: %v", err)
	} else if !exists {
		return "", fmt.Errorf("the asset %s does not exist", landID)
	}

	err = ctx.GetStub().DelState(landID)
	if err != nil {
		return "", err
	}

	// Emit event for Land deletion
	err = ctx.GetStub().SetEvent("LandDeleted", []byte(landID))
	if err != nil {
		return "", fmt.Errorf("failed to emit LandDeleted event: %v", err)
	}

	return fmt.Sprintf("land with id %v is deleted from the world state.", landID), nil
}

// LandContract.GetLandsByRange gives a range of asset details based on a start key and end key
func (lc *LandContract) GetLandsByRange(ctx contractapi.TransactionContextInterface, startKey, endKey string) ([]*Land, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return landResultIteratorFunction(resultsIterator)
}

// LandContract.GetLandHistory retrieves the history of changes to a land
func (lc *LandContract) GetLandHistory(ctx contractapi.TransactionContextInterface, landID string) ([]*HistoryQueryResult, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(landID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []*HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var land Land
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &land)
			if err != nil {
				return nil, err
			}
		} else {
			land = Land{
				LandID: landID,
			}
		}

		timestamp := response.Timestamp.AsTime()

		formattedTime := timestamp.Format(time.RFC1123)

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: formattedTime,
			Record:    &land,
			IsDelete:  response.IsDelete,
		}
		records = append(records, &record)
	}

	return records, nil
}

// LandContract.GetAllLands retrieves all land instances from the world state
func (lc *LandContract) GetAllLands(ctx contractapi.TransactionContextInterface) ([]*Land, error) {
	queryString := `{"selector":{"assetType":"land"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return landResultIteratorFunction(resultsIterator)
}

// LandContract.GetLandsWithPagination retrieves a paginated list of land instances
func (lc *LandContract) GetLandsWithPagination(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {
	queryString := `{"selector":{"assetType":"land"}}`
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	lands, err := landResultIteratorFunction(resultsIterator)
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             lands,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// Iterator function
func landResultIteratorFunction(resultsIterator shim.StateQueryIteratorInterface) ([]*Land, error) {
	var lands []*Land
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var land Land
		err = json.Unmarshal(queryResult.Value, &land)
		if err != nil {
			return nil, err
		}
		lands = append(lands, &land)
	}

	return lands, nil
}

// Example: RegisterLand registers land to the buyer
func (lc *LandContract) RegisterLand(ctx contractapi.TransactionContextInterface, landID string, ownerName string) (string, error) {
	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	if clientOrgID == "mvd-auto-com" {
		exists, err := lc.LandExists(ctx, landID)
		if err != nil {
			return "", fmt.Errorf("Could not read from world state. %s", err)
		}
		if exists {
			land, _ := lc.ReadLand(ctx, landID)
			land.OwnedBy = ownerName

			bytes, _ := json.Marshal(land)
			err = ctx.GetStub().PutState(landID, bytes)
			if err != nil {
				return "", err
			}

			// Emit event for Land registration
			err = ctx.GetStub().SetEvent("LandRegistered", bytes)
			if err != nil {
				return "", fmt.Errorf("failed to emit LandRegistered event: %v", err)
			}

			return fmt.Sprintf("Land %v successfully registered to %v", landID, ownerName), nil
		} else {
			return "", fmt.Errorf("Land %v does not exist!", landID)
		}
	} else {
		return "", fmt.Errorf("User under following MSPID: %v cannot able to perform this action", clientOrgID)
	}
}
