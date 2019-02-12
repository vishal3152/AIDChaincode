package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/shopspring/decimal"
)

// Asset model for donation. All asset models are kept in
// private scope i.e. they are not exported and ramin invisible to other packages
type donationBase struct {
	ProjectID string          `json:"projectID"`
	ItemID    string          `json:"itemID"` //ItemID
	Amount    decimal.Decimal `json:"amount"`
	TimeStamp time.Time       `json:"timeStamp"`
}

type donation struct {
	ObjectType string       `json:"docType"` // donation Type 'DONIN'
	TxnID      string       `json:"txnID"`   // asset unique key
	Donor      string       `json:"donor"`
	Data       donationBase `json:"data"` // composition
}

// Write donation to ledger
func (d *donation) putState(stub shim.ChaincodeStubInterface) pb.Response {

	// check if affiliated project exists
	c, cErr := checkAsset(stub, d.Data.ProjectID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	}
	if !c {
		e := &chainError{"putDonation", d.TxnID, CODENOTFOUND, errors.New("Affiliated project not found")}
		return shim.Error(e.Error())
	}

	// check if affiliated item exists
	c, cErr = checkAsset(stub, d.Data.ItemID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	} else if !c {
		e := &chainError{"putDonation", d.TxnID, CODENOTFOUND, errors.New("Affiliated item not found")}
		return shim.Error(e.Error())
	}

	// check if txnID is unique
	c, cErr = checkAsset(stub, d.TxnID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	} else if c {
		e := &chainError{"putDonation", d.TxnID, CODEAlRDEXIST, errors.New("Asset with key already exists")}
		return shim.Error(e.Error())
	}

	// Marshal the donation struct to []byte
	b, err := json.Marshal(d)
	if err != nil {
		cErr = &chainError{"putDonation", d.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}

	// Write key value to ledger
	err = stub.PutState(d.TxnID, b)

	if err != nil {
		cErr = &chainError{"putDonation", d.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}

	// Add indexkey for range query :- each donation  is stored in form of a delta and aggregated whenever
	// project state is read
	indexName := INDXNM
	indexKey, err := stub.CreateCompositeKey(indexName, []string{"1", d.TxnID, d.Data.Amount.StringFixedBank(int32(FIXEDPT))})
	if err != nil {
		cErr = &chainError{"putDonation", d.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	value := []byte{0x00}

	err = stub.PutState(indexKey, value)
	if err != nil {
		cErr = &chainError{"putDonation", d.TxnID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}

	// Emit transaction event for listeners
	stub.SetEvent((d.TxnID + "_AID_DON_" + d.Data.Amount.StringFixed(int32(FIXEDPT))), nil)

	r := response{CODEALLAOK, d.TxnID, nil}
	return shim.Success((r.formatResponse()))
}

// Read donation state from the ledger
func (d *donation) getState(stub shim.ChaincodeStubInterface) pb.Response {

	donation, cErr := queryAsset(stub, d.TxnID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	}
	r := response{CODEALLAOK, "OK", donation}
	return shim.Success((r.formatResponse()))
}
