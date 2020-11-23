package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	AssetType      string `json:"AssetType"`
	Quantity       int    `json:"Quantity"`
	AppraisedValue int    `json:"AppraisedValue"`
}

type HistoryQuery struct {
	Record    *Asset    `json:"record"`
	TxID      string    `json:"txID"`
	Timestamp time.Time `json:"timestamp"`
	Deleted   bool      `json:"deleted"`
}

// This function adds some predefined values into the blockchain
// transaction so that while viewing we have output
// NOTE: don't use Init function as it is defined in fabric sdk
func (sc *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "Cupset", Owner: "Tomoko", AssetType: "Household", Quantity: 10, AppraisedValue: 300},
		{ID: "Drill Machine", Owner: "Brad", AssetType: "Industry", Quantity: 6, AppraisedValue: 2000},
		{ID: "GPU", Owner: "Jin Soo", AssetType: "Electronics", Quantity: 18, AppraisedValue: 500},
		{ID: "Car", Owner: "Max", AssetType: "Heavy", Quantity: 2, AppraisedValue: 26000},
		{ID: "Chair", Owner: "Adriana", AssetType: "Household", Quantity: 26, AppraisedValue: 75},
		{ID: "Medicines", Owner: "Michel", AssetType: "Fragile", Quantity: 8, AppraisedValue: 10},
	}

	for _, asset := range assets {
		err := sc.CreateAsset(ctx, asset.ID, asset.Owner, asset.AssetType, asset.Quantity, asset.AppraisedValue)
		if err != nil {
			return err
		}
	}
	return nil
}

// Ledger is created using this function
// leverage this function for creation
func (sc *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id, owner, assetType string, quantity, appraisedvalue int) error {
	exists, err := sc.ExistsAsset(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if exists {
		return fmt.Errorf("Asset %s already exits", id)
	}

	asset := &Asset{
		ID:             id,
		Owner:          owner,
		AssetType:      assetType,
		Quantity:       quantity,
		AppraisedValue: appraisedvalue,
	}

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetBytes)
}

func (sc *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read World State: %v", err)
	}

	if assetBytes == nil {
		return nil, fmt.Errorf("Asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (sc *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id, owner string, quantity, appraisedvalue int) error {
	exists, err := sc.ExistsAsset(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("Asset %s already exits", id)
	}

	asset := Asset{
		ID:             id,
		Owner:          owner,
		Quantity:       quantity,
		AppraisedValue: appraisedvalue,
	}

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetBytes)
}

func (sc *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := sc.ExistsAsset(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("Asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

func (sc *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id, newOwner string) error {
	asset, err := sc.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetBytes)
}

func (sc *SmartContract) QueryAssetsById(ctx contractapi.TransactionContextInterface, id string) ([]*Asset, error) {
	queryString := fmt.Sprintf(`{"selector":{"ID":"%s"}}`, id)
	return getQueryResultForQueryString(ctx, queryString)
}

func (sc *SmartContract) QueryAssetsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Asset, error) {
	queryString := fmt.Sprintf(`{"selector":{"Owner": "%s"}}`, owner)
	return getQueryResultForQueryString(ctx, queryString)
}

func (sc *SmartContract) QueryAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset

	for resultsIterator.HasNext() {
		querResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset

		err = json.Unmarshal(querResponse.Value, &asset)
		if err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	return assets, nil
}

func (sc *SmartContract) GetAssetsHistory(ctx contractapi.TransactionContextInterface, id string) ([]HistoryQuery, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQuery
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Asset{
				ID: id,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}
		record := HistoryQuery{
			TxID:      response.TxId,
			Timestamp: timestamp,
			Record:    &asset,
			Deleted:   response.IsDelete,
		}

		records = append(records, record)
	}

	return records, nil
}

func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Asset, error) {
	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset

		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}

		assets = append(assets, &asset)
	}

	return assets, nil
}

func (sc *SmartContract) ExistsAsset(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJson != nil, nil
}

// Code starts executing from here
func main() {
	ReflectChaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		log.Panicf("Error Creating Asset: %v", err)
	}

	err = ReflectChaincode.Start()
	if err != nil {
		log.Panicf("Error Starting chaincode: %v", err)
	}
}
