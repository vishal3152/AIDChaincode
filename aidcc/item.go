package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Asset model for item. All asset models are kept in
// private scope i.e. they are not exported and ramin invisible to other packages
type itemBase struct {
	ItemType  string `json:"itemType"`
	Narrative string `json:"narrative"`
}
type item struct {
	ObjectType string   `json:"docType"` // item Type 'GITEM'
	ItemID     string   `json:"itemID"`  // asset unique key
	Data       itemBase `json:"data"`    // composition
}

// Write asset state to ledger
func (it *item) putState(stub shim.ChaincodeStubInterface) pb.Response {

	// check if itemID already exists
	c, cErr := checkAsset(stub, it.ItemID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	}
	if c {
		e := &chainError{"putItem", it.ItemID, CODEAlRDEXIST, errors.New("Asset with key already exists")}
		return shim.Error(e.Error())
	}

	// Marshal the Item struct to []byte
	b, err := json.Marshal(it)
	if err != nil {
		cErr = &chainError{"putItem", it.ItemID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	// Write key-value to ledger
	err = stub.PutState(it.ItemID, b)
	if err != nil {
		cErr = &chainError{"putItem", it.ItemID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}

	// Emit transaction event for listeners
	txID := stub.GetTxID()
	stub.SetEvent((it.ItemID + "_AID_ITMADD_" + txID), nil)
	r := response{CODEALLAOK, it.ItemID, nil}
	return shim.Success((r.formatResponse()))
}

// Read item state from the ledger
func (it *item) getState(stub shim.ChaincodeStubInterface) pb.Response {

	item, cErr := queryAsset(stub, it.ItemID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	}
	r := response{CODEALLAOK, "OK", item}
	return shim.Success((r.formatResponse()))
}
