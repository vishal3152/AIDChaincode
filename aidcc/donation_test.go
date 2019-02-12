package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

// Verifies scenarios related to donation
func TestDonationRW(t *testing.T) {
	fmt.Println("Executing Test - DonationRW")

	// Test data - Refer to test narratives for test description
	var testTable = []struct {
		txnID          string
		donorName      string
		projectID      string
		itemID         string
		amount         string
		expectedStatus int32
		testNarrative  string
	}{
		{"D101", "vishal", "P101", "Itm001", "100", 200, "Happy scenario"},
		{"D102", "vishal", "P101", "Itm001", "", 500, "Check for postive amount"},
		{"D102", "vishal", "P101", "Itm001", "0", 500, "Check for positive amount"},
		{"D102", "vishal", "P101", "Itm001", "-1", 500, "Check for negative amount"},
		{"", "vishal", "P101", "Itm001", "100", 500, "Check for missing transaction ID"},
		{"D102", "", "P101", "Itm001", "100", 500, "Check for missing donor"},
		{"D102", "vishal", "", "Itm001", "100", 500, "Check for missing projectID"},
		{"D102", "vishal", "P102", "Itm001", "100", 500, "Check for project ID existence"},
		{"D102", "vishal", "P101", "", "100", 500, "Check for missing item ID"},
		{"D102", "vishal", "P101", "Itm002", "100", 500, "Check for item ID existence"},
		{"D101", "vishal", "P101", "Itm001", "100", 500, "Check for duplicate transaction ID"},
	}
	// struct for parsing the shim APIs response
	type dResp struct {
		Code    string   `json:"code"`
		Message string   `json:"message"`
		Payload donation `json:"payload"`
	}
	type pResp struct {
		Code    string  `json:"code"`
		Message string  `json:"message"`
		Payload project `json:"payload"`
	}
	assert := assert.New(t)

	// Instantiate mockStub using AidChaincode as the target chaincode to unit test
	stub := shim.NewMockStub("TestStub", new(AidChaincode))
	//Verify stub is available
	assert.NotNil(stub, "Stub is nil, Test stub creation failed")

	uid := uuid.New().String()

	// Add test project ID i.e. P101 and item ID i.e. Itm101 to state.
	// Donations must be tagged to a project and an item
	// Adding project
	result := stub.MockInvoke(uid,
		[][]byte{[]byte(AddProject),
			[]byte("P101"),
			[]byte("Prj101")})

	assert.EqualValues(shim.OK, result.GetStatus(), "Saving test project to state failed.")

	// Adding item
	result = stub.MockInvoke(uid,
		[][]byte{[]byte(AddItem),
			[]byte("Itm001"),
			[]byte("Item001"),
			[]byte("Medicine")})

	assert.EqualValues(shim.OK, result.GetStatus(), "Saving test item to state failed.")

	// Executing tests
	for _, test := range testTable {
		result := stub.MockInvoke(uid,
			[][]byte{[]byte(AddDonation),
				[]byte(test.txnID),
				[]byte(test.donorName),
				[]byte(test.projectID),
				[]byte(test.itemID),
				[]byte(test.amount)})

		assert.Equal(test.expectedStatus, result.GetStatus(), test.testNarrative+" failed - "+(result.GetMessage()))

		if result.GetStatus() == shim.OK {
			// verify donation read
			result = stub.MockInvoke(uid,
				[][]byte{[]byte(GetDonation),
					[]byte(test.txnID)})

			assert.EqualValues(shim.OK, result.GetStatus(), GetDonation+" failed to read the project data")

			r := &dResp{}
			err := json.Unmarshal(result.GetPayload(), r)
			if err != nil {
				panic(err)
			}
			assert.Equal(test.txnID, r.Payload.TxnID, "Reterived txn ID mismatch")

			// Verify donation amount reflect under project fund
			result = stub.MockInvoke(uid,
				[][]byte{[]byte(GetProject),
					[]byte(test.projectID)})

			assert.EqualValues(shim.OK, result.GetStatus(), GetProject+" failed to read the project data")

			p := &pResp{}
			err = json.Unmarshal(result.GetPayload(), p)
			if err != nil {
				panic(err)
			}
			assert.Equal(test.projectID, p.Payload.ProjectID, "Reterived project ID mismatch")
			exceptedVal, _ := decimal.NewFromString(test.amount)
			assert.EqualValues(exceptedVal, p.Payload.Data.AvlFund, "Donation amount fails to reflect under project fund")
		}
	}
}
