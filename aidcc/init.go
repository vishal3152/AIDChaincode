package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// AidChaincode - define a type for chaincode.
// AidChaincode type must implements shim.Chaincode interface
type AidChaincode struct {
}

// Constants
const (
	// Error codes
	CODEALLAOK              string = "P2001" // Success
	CODENOTFOUND            string = "P4001" // resource not found
	CODEUNKNOWNINVOKE       string = "P4002" // Unknown invoke
	CODEUNPROCESSABLEENTITY string = "P4003" // Invalid input
	CODEGENEXCEPTION        string = "P5001" // Unknown exception
	CODEAlRDEXIST           string = "P5002" // Not unique
	CODENOTALLWD            string = "P4004" // Operation not allowed

	// Couch DB Doc types for asset
	GPRJCT string = "GPRJCT"
	GITEM  string = "GITEM"
	DONIN  string = "DONIN"
	DONOUT string = "DONOUT"

	// Range index name - to perform range queries
	INDXNM string = "bitmask~txnID~amount" //bitmask is "0" for donation (spending) & "1" donation(incoming)

	FIXEDPT int32 = 4 // All currency values rounded off to 4 decimals i.e. 0.0000
)

// Init - Implements shim.Chaincode interface Init() method
func (t *AidChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	r := response{(CODEALLAOK), "AIDcc started", nil}
	return shim.Success((r.formatResponse()))
}
