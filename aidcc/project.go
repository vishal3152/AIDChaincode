package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/shopspring/decimal"
)

// Asset model for project. All asset models are kept in
// private scope i.e. they are not exported and ramin invisible to other packages
type project struct {
	ObjectType string      `json:"docType"`   // Project Type 'GPRJCT'
	ProjectID  string      `json:"projectID"` // asset unique key
	Data       projectBase `json:"data"`      // composition
}
type projectBase struct {
	ProjectName string          `json:"projectName"`
	RunBy       string          `json:"runBy"`
	StartDt     time.Time       `json:"startDt"`
	AvlFund     decimal.Decimal `json:"avlFund"`
	SpentFund   decimal.Decimal `json:"spentFund"`
}

// Write asset state to ledger
func (p *project) putState(stub shim.ChaincodeStubInterface) pb.Response {

	// check if projectID already exists
	c, cErr := checkAsset(stub, p.ProjectID)
	if cErr != nil {
		return shim.Error(cErr.Error())
	}
	if c {
		e := &chainError{"putProject", p.ProjectID, CODEAlRDEXIST, errors.New("Asset with key already exists")}
		return shim.Error(e.Error())
	}

	// Marshal the project struct ti []byte
	b, err := json.Marshal(p)
	if err != nil {
		cErr = &chainError{"putProject", p.ProjectID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	// Write key-value to ledger
	err = stub.PutState(p.ProjectID, b)
	if err != nil {
		cErr = &chainError{"putProject", p.ProjectID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}

	// Emit transaction event for listeners
	txID := stub.GetTxID()
	stub.SetEvent((p.ProjectID + "_AID_PRJADD_" + txID), nil)
	r := response{CODEALLAOK, p.ProjectID, nil}
	return shim.Success((r.formatResponse()))

}

func (p *project) getState(stub shim.ChaincodeStubInterface) pb.Response {

	// Aggregate donation & spend deltas:
	// Donations amount do not reflect immediately reflected under project's available fund.
	// Each donation or Spend is stored in the form of a delts (range index). Deltas are aggregated and
	// added to the current value of available funds whenever this operation is called.
	// This is a work around to avoid transactiona colloisons (MVCC R/W conflicts) during a high throughput scenario
	// This operation should be invoked by client application at regular interval to prevent overspending of
	// available funds under teh project.

	donationAggregate, _ := decimal.NewFromString("0")
	spendAggregate, _ := decimal.NewFromString("0")
	//Do a range query on donations (incoming, bitmask "1")
	indexName := string(INDXNM)
	donItr, err := stub.GetStateByPartialCompositeKey(indexName, []string{"1"})
	if err != nil {
		cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	//Close itrerator when done reading
	defer donItr.Close()
	for donItr.HasNext() {
		rangeItem, err := donItr.Next()
		if err != nil {
			cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
			return shim.Error(cErr.Error())
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(rangeItem.Key)
		if err != nil {
			cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
			return shim.Error(cErr.Error())
		}
		// compositeKeyParts[2] represents transaction amount
		txAmount, _ := decimal.NewFromString(compositeKeyParts[2])
		donationAggregate = donationAggregate.Add(txAmount)

		// Delete the key from index after its delta has been aggregated
		err = stub.DelState(rangeItem.Key)
		if err != nil {
			cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
			return shim.Error(cErr.Error())
		}
	}

	//Aggreate spends
	spendItr, err := stub.GetStateByPartialCompositeKey(indexName, []string{"0"})
	if err != nil {
		cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	//Close itrerator when done reading
	defer spendItr.Close()
	for spendItr.HasNext() {
		rangeItem, err := spendItr.Next()
		if err != nil {
			cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
			return shim.Error(cErr.Error())
		}
		_, compositeKeyParts, err := stub.SplitCompositeKey(rangeItem.Key)
		if err != nil {
			cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
			return shim.Error(cErr.Error())
		}
		// Delete the key from index after its delta has been aggregated
		txAmount, _ := decimal.NewFromString(compositeKeyParts[2])
		spendAggregate = spendAggregate.Add(txAmount)

		// Delete the key from index after its delta has been aggregated
		err = stub.DelState(rangeItem.Key)
		if err != nil {
			cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
			return shim.Error(cErr.Error())
		}
	}

	// Get current value of available fund under the project
	prjBytes, err := stub.GetState(p.ProjectID)
	if err != nil {
		cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	prj := &project{}
	err = json.Unmarshal(prjBytes, prj)
	if err != nil {
		cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	// calculate the new value
	prj.Data.AvlFund = prj.Data.AvlFund.Add(donationAggregate.Sub(spendAggregate))
	prj.Data.SpentFund = prj.Data.SpentFund.Add(spendAggregate)

	prjBytes, err = json.Marshal(prj)
	if err != nil {
		cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}

	// Write the new value of available funds on ledger
	err = stub.PutState(p.ProjectID, prjBytes)
	if err != nil {
		cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}

	// Read and return updated project data
	prjBytes, err = stub.GetState(p.ProjectID)
	if err != nil {
		cErr := &chainError{"readProject", p.ProjectID, CODEGENEXCEPTION, err}
		return shim.Error(cErr.Error())
	}
	r := response{CODEALLAOK, p.ProjectID, prjBytes}
	return shim.Success((r.formatResponse()))
}
