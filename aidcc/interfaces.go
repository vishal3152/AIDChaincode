package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// AidAssetInterface interface must be implemented by all asset types defined
// under AidChaincode. The interface enforces that all asset should implement atleast two methods i.e.:
// putState() -> Writes asset state to ledger
// getState() -> Read asset state from ledger
type AidAssetInterface interface {
	putState(stub shim.ChaincodeStubInterface) pb.Response
	getState(stub shim.ChaincodeStubInterface) pb.Response
}

func saveAsset(stub shim.ChaincodeStubInterface, i AidAssetInterface) pb.Response {
	return i.putState(stub)
}
func readAsset(stub shim.ChaincodeStubInterface, i AidAssetInterface) pb.Response {
	return i.getState(stub)
}
