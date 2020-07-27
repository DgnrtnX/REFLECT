package main

import(
	"log"

	"github.com/hyperledger/fabric-contract-api/contractapi"
	"github.com/hyperledger/Reflect/backend/chaincode"
)

func main(){
	ReflectChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil{
		log.Panicf("Error Creating Asset: %v", err)
	}
	err = ReflectChaincode.Start(); if err != nil{
		log.Panicf("Error Starting chaincode: %v", err)
	}
}
