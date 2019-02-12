package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Do setup before tests are run. Each test initialize a new MockStub.
// This is equivalent to resetting the fabric storage before each run
func TestMain(m *testing.M) {
	logger.SetLevel(shim.LogError)
	os.Setenv("MODE", "TEST")
	fmt.Printf("\n\n*--------------* RUN Mode set to %s *-----------------* \n \n", os.Getenv("MODE"))
	os.Exit(m.Run())
}
