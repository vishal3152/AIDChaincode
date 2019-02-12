package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("AidChaincode")

// Chaincode entry point
func main() {
	logger.SetLevel(shim.LogInfo)
	err := shim.Start(new(AidChaincode))
	if err != nil {
		logger.Error("Error starting AidChaincode - ", err)
	}

}
