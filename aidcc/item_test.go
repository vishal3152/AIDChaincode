package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

// Verify all scenarios related to item.
func TestItemRW(t *testing.T) {
	fmt.Println("Executing Test - ItemRW")

	// Test data - Refer to test narratives for test description
	var testTable = []struct {
		itemID         string
		itemType       string
		itemNarrative  string
		expectedStatus int32
		testNarrative  string
	}{
		{"I101", "Itm101", "Medicine", 200, "Happy scenario"},
		{"", "Itm101", "Medicine", 500, "Check for missing item ID"},
		{"I102", "", "Medicine", 500, "Check for missing item type"},
		{"I102", "Itm101", "", 500, "Check for missing item narrtaive"},
		{"I101", "Itm101", "Medicine", 500, "Check for duplicate item ID"},
	}

	// struct for parsing the shim APIs response
	type resp struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Payload item   `json:"payload"`
	}
	assert := assert.New(t)

	// Instantiate mockStub using AidChaincode as the target chaincode to unit test
	stub := shim.NewMockStub("TestStub", new(AidChaincode))

	assert.NotNil(stub, "Stub is nil, Test stub creation failed")

	uid := uuid.New().String()

	// Executing tests
	for _, test := range testTable {
		result := stub.MockInvoke(uid,
			[][]byte{[]byte(AddItem),
				[]byte(test.itemID),
				[]byte(test.itemType),
				[]byte(test.itemNarrative)})

		assert.Equal(test.expectedStatus, result.GetStatus(), test.testNarrative+" failed.")

		if result.GetStatus() == shim.OK {
			// Verify Item read
			result = stub.MockInvoke(uid,
				[][]byte{[]byte(GetItem),
					[]byte(test.itemID)})

			assert.EqualValues(shim.OK, result.GetStatus(), GetItem+" failed to read the item data")

			r := &resp{}
			err := json.Unmarshal(result.GetPayload(), r)
			if err != nil {
				panic(err)
			}
			assert.Equal(test.itemID, r.Payload.ItemID, "Reterived item ID mismatch")
		}
	}
}
