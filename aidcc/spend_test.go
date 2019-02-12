package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Verifies scenarios related to spends
func TestSpendRW(t *testing.T) {
	fmt.Println("Executing Test - SpendRW")

	// Test data - Refer to test narratives for test description
	var testTable = []struct {
		txnID          string
		beneficiary    string
		projectID      string
		itemID         string
		amount         string
		expectedStatus int32
		testNarrative  string
	}{
		{"D101", "vishal", "P101", "Itm001", "100", 200, "Happy scenario"},
		{"D111", "vishal", "P101", "Itm001", "150", 500, "Check for over spending"},
		{"D103", "vishal", "P101", "Itm001", "", 500, "Check for postive amount"},
		{"D104", "vishal", "P101", "Itm001", "0", 500, "Check for positive amount"},
		{"D105", "vishal", "P101", "Itm001", "-1", 500, "Check for negative amount"},
		{"", "vishal", "P101", "Itm001", "100", 500, "Check for missing transaction ID"},
		{"D106", "", "P101", "Itm001", "100", 500, "Check for missing beneficiary"},
		{"D107", "vishal", "", "Itm001", "100", 500, "Check for missing projectID"},
		{"D108", "vishal", "P102", "Itm001", "100", 500, "Check for project ID existence"},
		{"D109", "vishal", "P101", "", "100", 500, "Check for missing item ID"},
		{"D100", "vishal", "P101", "Itm002", "100", 500, "Check for item ID existence"},
		{"D101", "vishal", "P101", "Itm001", "100", 500, "Check for duplicate transaction ID"},
	}

	// struct for parsing the shim APIs response
	type sResp struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Payload spend  `json:"payload"`
	}
	type pResp struct {
		Code    string  `json:"code"`
		Message string  `json:"message"`
		Payload project `json:"payload"`
	}

	assert := assert.New(t)

	// Instantiate mockStub using AidChaincode as the target chaincode to unit test
	stub := shim.NewMockStub("TestStub", new(AidChaincode))
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

	//Adding item
	result = stub.MockInvoke(uid,
		[][]byte{[]byte(AddItem),
			[]byte("Itm001"),
			[]byte("Item001"),
			[]byte("Medicine")})

	assert.EqualValues(shim.OK, result.GetStatus(), "Saving test item to state failed.")

	// Add a donation so that spend can happen, overspending check will be
	// performed on this donation
	donationAmount := "200"
	result = stub.MockInvoke(uid,
		[][]byte{[]byte(AddDonation),
			[]byte("D102"),
			[]byte("vishal"),
			[]byte("P101"),
			[]byte("Itm001"),
			[]byte(donationAmount)})
	assert.EqualValues(shim.OK, result.GetStatus(), "Saving test donation entry to state failed.")

	// Refresh Available funds under Project
	result = stub.MockInvoke(uid,
		[][]byte{[]byte(GetProject),
			[]byte("P101")})

	assert.EqualValues(shim.OK, result.GetStatus(), GetProject+" failed to read the project data")

	// Execute tests
	for _, test := range testTable {
		result := stub.MockInvoke(uid,
			[][]byte{[]byte(AddSpend),
				[]byte(test.txnID),
				[]byte(test.beneficiary),
				[]byte(test.projectID),
				[]byte(test.itemID),
				[]byte(test.amount)})

		assert.Equal(test.expectedStatus, result.GetStatus(), test.testNarrative+" failed - "+(result.GetMessage()))

		if result.GetStatus() == shim.OK {
			// Veridy spend read
			result = stub.MockInvoke(uid,
				[][]byte{[]byte(GetSpend),
					[]byte(test.txnID)})

			assert.EqualValues(shim.OK, result.GetStatus(), GetDonation+" failed to read the project data")
			r := &sResp{}
			err := json.Unmarshal(result.GetPayload(), r)
			if err != nil {
				panic(err)
			}
			assert.Equal(test.txnID, r.Payload.TxnID, "Reterived txn ID mismatch")

			// Verify - spend amount should reflect under project funds
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
			initialFundVal, _ := decimal.NewFromString(donationAmount)
			spentAmount, _ := decimal.NewFromString(test.amount)
			expectedFundVal := initialFundVal.Sub(spentAmount)

			assert.EqualValues(expectedFundVal, p.Payload.Data.AvlFund, "Spent amount fails to reflect under project available fund")
			assert.EqualValues(spentAmount, p.Payload.Data.SpentFund, "Spent amount fails to reflect under project spent fund")

		}
	}
}
