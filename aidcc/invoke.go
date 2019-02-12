package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// These are function names from Invoke first parameter
const (
	Init        string = "Init"
	AddProject  string = "AddProject"
	AddItem     string = "AddItem"
	AddDonation string = "AddDonation"
	AddSpend    string = "AddSpend"
	GetProject  string = "GetProject"
	GetItem     string = "GetItem"
	GetDonation string = "GetDonation"
	GetSpend    string = "GetSpend"
)

// Invoke - Implements shim.Chaincode interface Invoke() method
func (t *AidChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// Reterive function and input arguments. This is not a recommended approach.
	// Instead, GetArgs() is a more suitable method and works perfectly with protocol buffers
	function, args := stub.GetFunctionAndParameters()
	logger.Info(fmt.Sprintf("Starting Phantom chaincode Invoke for %s and no of argument passed are %d", function, len(args)))

	if function == Init {
		return t.Init(stub)
	} else if function == AddProject {
		return validateProjectW(stub, args)
	} else if function == AddItem {
		return validateItemW(stub, args)
	} else if function == AddDonation {
		return validateDonationW(stub, args)
	} else if function == AddSpend {
		return validateSpendW(stub, args)
	} else if function == GetProject {
		return validateProjectR(stub, args)
	} else if function == GetItem {
		return validateItemR(stub, args)
	} else if function == GetDonation {
		return validateDonationR(stub, args)
	} else if function == GetSpend {
		return validateSpendR(stub, args)
	}

	e := chainError{"Invoke", "", CODEUNKNOWNINVOKE, errors.New("Unknown function invoke")}
	return shim.Error(e.Error())
}
