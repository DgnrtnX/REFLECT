package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const index = "AssetType~desc"

type SmartContract struct {
	contractapi.Contract
}

type Asset struct {
	ID             string `json:"ID"`
	Owner          string `json:"Owner"`
	AssetType      string `json:AssetType`
	Quantity       int    `json:"Quantity"`
	AppraisedValue int    `json:"AppraisedValue"`
}

type HistoryQuery struct {
	Record    *Asset    `json:"record"`
	TxID      string    `json:"txID"`
	Timestamp time.Time `json:"timestamp"`
	Deleted   bool      `json:"deleted"`
}

type PaginatedQuery struct {
	Records            []*Asset `json:"records"`
	FetchedRecordCount int32    `json:"fetechedrecordCount"`
	Bookmark           string   `jso:"bookmark"`
}

func (sc *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "asset1", Owner: "Tomoko", AssetType: "household", Quantity: 10, AppraisedValue: 300},
		{ID: "asset2", Owner: "Brad", AssetType: "industry", Quantity: 6, AppraisedValue: 400},
		{ID: "asset3", Owner: "Jin Soo", AssetType: "industry", Quantity: 18, AppraisedValue: 500},
		{ID: "asset4", Owner: "Max", AssetType: "household", Quantity: 2, AppraisedValue: 600},
		{ID: "asset5", Owner: "Adriana", AssetType: "mines", Quantity: 26, AppraisedValue: 700},
		{ID: "asset6", Owner: "Michel", AssetType: "mines", Quantity: 8, AppraisedValue: 800},
	}

	for _, asset := range assets {
		err := sc.CreateAsset(ctx, asset.ID, asset.Owner, asset.AssetType, asset.Quantity, asset.AppraisedValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sc *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id, owner, assettype string, quantity, appraisedvalue int) error {
	exists, err := sc.ExistsAsset(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	if exists {
		return fmt.Errorf("Asset %s already exits", id)
	}

	asset := &Asset{
		ID: id, Owner: owner, AssetType: assettype, Quantity: quantity, AppraisedValue: appraisedvalue,
	}

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	assetTypeIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.AssetType, asset.ID})
	if err != nil {
		return err
	}
	value := []byte{0x00}

	return ctx.GetStub().PutState(assetTypeIndexKey, value)
}

func (sc *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to read World State: %v", err)
	}
	if assetJson == nil {
		return nil, fmt.Errorf("Asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJson, &asset)
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
		ID: id, Owner: owner, Quantity: quantity, AppraisedValue: appraisedvalue,
	}

	assetJson, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJson)
}

func (sc *SmartContract) DelAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := sc.ExistsAsset(ctx, id)
	if err != nil {
		return nil
	}
	if !exists {
		return fmt.Errorf("Asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
	if err != nil {
		return fmt.Errorf("failed to delete asset %s: %v", id, err)
	}

	assetTypeIndexKey, err := ctx.GetStub().CreateCompositeKey(index, []string{asset.AssetType, asset.ID})
	if err != nil {
		return err
	}

	// Delete index entry
	return ctx.GetStub().DelState(colorNameIndexKey)
}

func (sc *SmartContract) TransAsset(ctx contractapi.TransactionContextInterface, id, newOwner string) error {
	asset, err := sc.ReadAsset(id)
	if err != nil {
		return err
	}
	asset.Owner = newOwner
	assetJson, err = json.Marshal(asset)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, assetJson)
}

func (sc *SmartContract) TransAssetByType(ctx contractapi.TransactionContextInterface, assettype, newOwner string) error {
	assetTyperesultsIterator, err := ctx.GetStub().GetStateByPartialCompositeKey(index, []string{assettype})
	if err != nil {
		return err
	}
	defer assetTyperesultsIterator.Close()

	for assetTyperesultsIterator.HasNext() {
		responseRange, err := assetTyperesultsIterator.Next()
		if err != nill {
			return err
		}
		_, compositeKeyParts, err := ctx.GetStub().SliptCompositeKey(responseRange.Key)
		if err != nil {
			return err
		}
		if len(compositeKeyParts) > 1 {
			returnAssetID := compositeKeyParts[1]
			asset, err := sc.ReadAsset(ctx, returnAssetID)
			if err != nil {
				return err
			}
			asset.Owner = newOwner
			assetJson, err := json.Marshal(asset)
			if err != nil {
				return err
			}
			err = ctx.GetStub().PutState(returnAssetID, ass)
			if err != nil {
				return fmt.Errorf("tranfer of ownership failed for asset %s : &v", returnAssetID, err)
			}
		}
	}
	return nil
}

func (sc *SmartContract) QueryResponseFromIterator(resultsIterator shim.SateQueryIteratorInterface) ([]*Asset, error) {
	var asset []*Asset
	for resultsIterator.HasNext() {
		queryResults, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var assets Asset
		err = json.Unmarshal(queryResults.Value, &asset)
		if err != nil {
			return err
		}
		assets = append(assets, &assest)
	}
	return assets, nil
}

func (sc *SimpleChaincode) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

func (sc *SmartContract) GetAssetsByRange(ctx contractapi.TransactionContextInterface, startkey, endkey string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetAssetByRange(startkey, endkey)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return QueryResponseFromIterator(resultsIterator)
}

func (sc *SmartContract) QueryAssetsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]*Asset, error) {
	queryString, err := fmt.Sprintf(`{"selector":{"owner": "%s"}}`, owner)
	return getQueryResultForQueryString(queryString)
}

func (sc *SmartContract) getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub(), GetQueryResult(queryString)
	if err != nil {
		return err
	}
	defer resultsIterator.Close()
	return QueryResponseFromIterator(resultsIterator)
}

func (sc *SmartContract) getAssetsWithPagination(ctx contractapi.TransactionContextInterface, queryString, bookmark string, pageSize int) (*PaginatedQuery, error) {
	return getQueryResultForQueryStringWithPagination(ctx, queryString, bookmark, int32(pageSize))
}

func (sc *SmartContract) getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString, bookmark string, pageSize int32) (*PaginatedQuery, error) {
	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return err
	}
	defer resultsIterator.Close()

	assets, err := QueryResponseFromIterator(resultsIterator)
	if err != nil {
		return err
	}

	return &PaginatedQuery{
		Records:            assets,
		FetchedRecordCount: responseMetadata.FetchedRecordCount,
		Bookmark:           responseMetadata.bookmark,
	}, nil
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
		assets := append(assets, &asset)
	}
	return assets, nil
}

func (sc *SimpleChaincode) GetAssetsHistory(ctx contractapi.TransactionContextInterface, id string) ([]HistoryQuery, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQuery
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			response.Value
		}
		return nil, err
	}

	var asset Asset
	err = json.Unmarshal(response.Value, &asset)
	if err != nil {
		return nil, err
	}

	timestamp, err := ptypes.Timestamp(response.Timestamp)
	if err != nil {
		return nil, err
	}
	record := HistoryQueryResult{
		TxID:      response.TxId,
		Timestamp: timestamp,
		Record:    &asset,
		Deleted:   response.Deleted,
	}
	records = append(records, record)

	return records, nil
}

func (sc *SmartContract) ExistsAsset(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJson, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}
	return assetJson != nil, nil
}
