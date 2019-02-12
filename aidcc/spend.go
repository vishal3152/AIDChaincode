package main

import (
	"encoding/json"
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/shopspring/decimal"
)

// Asset model for spend. All asset models are kept in
// private scope i.e. they are not exported and ramin invisible to other packages
type spend struct {
	ObjectType string       `json:"docType"` // donation Type 'DONOUT'
	TxnID      string       `json:"txnID"`   //asset unique key
	Benficiary string       `json:"donor"`
	Data       donationBase `json:"data"` // composition
}

// Write asset state to ledger
func (s *spend) putState(stub shim.ChaincodeStubInterface) pb.Response {

	// check if affiliated project exists
	c, cErr := checkAsset(stub, s.Data.ProjectID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	} else if !c {
		e := &chainError{"putSpend", s.TxnID, CODENOTFOUND, errors.New("Affiliated project not found")}
		return shim.Error(e.Error())
	}

	// check if affiliated item exists
	c, cErr = checkAsset(stub, s.Data.ItemID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	} else if !c {
		e := &chainError{"putSpend", s.TxnID, CODENOTFOUND, errors.New("Affiliated item not found")}
		return shim.Error(e.Error())
	}

	// check if txnID is unique
	c, cErr = checkAsset(stub, s.TxnID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	} else if c {
		e := &chainError{"putSpend", s.TxnID, CODEAlRDEXIST, errors.New("Asset with key already exists")}
		return shim.Error(e.Error())
	}

	// check if project has funds available before making a spend
	prj, err := stub.GetState(s.Data.ProjectID)
	if err != nil {
		cErr = &chainError{"putSpend", s.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	p := &project{}
	err = json.Unmarshal(prj, p)
	if err != nil {
		cErr = &chainError{"putSpend", s.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	zeroDecimal, _ := decimal.NewFromString("0")
	if p.Data.AvlFund.Sub(s.Data.Amount).LessThanOrEqual(zeroDecimal) {
		cErr = &chainError{"putSpend", s.TxnID, CODENOTALLWD, errors.New("Overwithdrawal for funds not allowed")}
		return shim.Error(cErr.Error())
	}

	// Convert spend struct to []byte
	b, err := json.Marshal(s)
	if err != nil {
		cErr = &chainError{"putSpend", s.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	// Write spend to ledger
	err = stub.PutState(s.TxnID, b)
	if err != nil {
		cErr = &chainError{"putSpend", s.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}

	// Add indexkey for range query :- each spend is stored in form of a delta and aggregated whenever
	// project state is read
	indexName := string(INDXNM)
	indexKey, err := stub.CreateCompositeKey(indexName, []string{"0", s.TxnID, s.Data.Amount.StringFixedBank(int32(FIXEDPT))})
	if err != nil {
		cErr = &chainError{"putSpend", s.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	value := []byte{0x00}

	err = stub.PutState(indexKey, value)
	if err != nil {
		cErr = &chainError{"putSpend", s.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}

	// Emit transaction event for listeners
	stub.SetEvent((s.TxnID + "_AID_SPND_" + s.Data.Amount.StringFixed(FIXEDPT)), nil)
	r := response{CODEALLAOK, s.TxnID, nil}
	return shim.Success((r.formatResponse()))
}

// Read spend state from ledger
func (s *spend) getState(stub shim.ChaincodeStubInterface) pb.Response {

	spend, cErr := queryAsset(stub, s.TxnID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	}
	r := response{CODEALLAOK, "OK", spend}
	return shim.Success((r.formatResponse()))
}
