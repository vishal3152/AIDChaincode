package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/stretchr/testify/assert"
)

// Verify scenarios rlated to a project asset
func TestProjectRW(t *testing.T) {
	fmt.Println("Executing Test - ProjectRW")

	// Test data - Refer to test narratives for test description
	var testTable = []struct {
		projectID      string
		projectName    string
		expectedStatus int32
		testNarrative  string
	}{
		{"P101", "Prj101", 200, "Happy scenario"},
		{"", "Prj101", 500, "Check for missing project ID"},
		{"P102", "", 500, "Check for missing project name"},
		{"P101", "Prj101", 500, "Check for duplicate project ID"},
	}

	// struct for parsing the shim APIs response
	type resp struct {
		Code    string  `json:"code"`
		Message string  `json:"message"`
		Payload project `json:"payload"`
	}

	assert := assert.New(t)

	// Instantiate mockStub using AidChaincode as the target chaincode to unit test
	stub := shim.NewMockStub("TestStub", new(AidChaincode))
	assert.NotNil(stub, "Stub is nil, Test stub creation failed")

	uid := uuid.New().String()

	// Executing tests
	for _, test := range testTable {
		result := stub.MockInvoke(uid,
			[][]byte{[]byte(AddProject),
				[]byte(test.projectID),
				[]byte(test.projectName)})

		assert.Equal(test.expectedStatus, result.GetStatus(), test.testNarrative+" failed.")

		if result.GetStatus() == shim.OK {
			// verify project read
			result = stub.MockInvoke(uid,
				[][]byte{[]byte(GetProject),
					[]byte(test.projectID)})

			assert.EqualValues(shim.OK, result.GetStatus(), GetProject+" failed to read the project data")

			r := &resp{}
			err := json.Unmarshal(result.GetPayload(), r)
			if err != nil {
				panic(err)
			}
			assert.Equal(test.projectID, r.Payload.ProjectID, "Reterived project ID mismatch")
		}
	}
}
