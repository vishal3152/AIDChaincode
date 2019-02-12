package main

import (
	"errors"
	"os"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/shopspring/decimal"
)

// Validations to check for accuracy of input data before any action can be performed on the chain code.

func validateProjectW(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		cErr := &chainError{"validateProjectW", "", CODEUNPROCESSABLEENTITY, errors.New("Incorrect no of input args, excepting 2")}
		return shim.Error(cErr.Error())
	}
	if len(args[0]) == 0 {
		cErr := &chainError{"validateProjectW", "", CODEUNPROCESSABLEENTITY, errors.New("Project ID can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[1]) == 0 {
		cErr := &chainError{"validateProjectW", "", CODEUNPROCESSABLEENTITY, errors.New("Project name can not be empty")}
		//logger.Error(cErr.Error())
		return shim.Error(cErr.Error())
	}

	epochTime, _ := stub.GetTxTimestamp()
	startDt := time.Unix(epochTime.GetSeconds(), 0)
	// Bypass whilst running unit test
	var callerID string
	if os.Getenv("MODE") != "TEST" {
		var err *chainError
		callerID, err = getCallerID(stub)
		if err != nil {
			return shim.Error(err.Error())
		}
	} else {
		callerID = "Test Caller"
	}

	prjBase := projectBase{ProjectName: args[1], StartDt: startDt, RunBy: callerID}
	p := &project{ObjectType: GPRJCT, ProjectID: args[0], Data: prjBase}
	return saveAsset(stub, p)
}

func validateItemW(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		cErr := &chainError{"validateItemW", "", CODEUNPROCESSABLEENTITY, errors.New("Incorrect no of input args, excepting 3")}
		return shim.Error(cErr.Error())
	}
	if len(args[0]) == 0 {
		cErr := &chainError{"validateItemW", "", CODEUNPROCESSABLEENTITY, errors.New("Item ID can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[1]) == 0 {
		cErr := &chainError{"validateItemW", "", CODEUNPROCESSABLEENTITY, errors.New("Item type can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[2]) == 0 {
		cErr := &chainError{"validateItemW", "", CODEUNPROCESSABLEENTITY, errors.New("Item narrative can not be empty")}
		return shim.Error(cErr.Error())
	}

	itmBase := itemBase{ItemType: args[1], Narrative: args[2]}
	it := &item{ObjectType: GITEM, ItemID: args[0], Data: itmBase}

	return saveAsset(stub, it)
}

func validateDonationW(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 5 {
		cErr := &chainError{"validateDonationW", "", CODEUNPROCESSABLEENTITY, errors.New("Incorrect no of input args, excepting 5")}
		return shim.Error(cErr.Error())
	}
	if len(args[0]) == 0 {
		cErr := &chainError{"validateDonationW", "", CODEUNPROCESSABLEENTITY, errors.New("Txn ID can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[1]) == 0 {
		cErr := &chainError{"validateDonationW", "", CODEUNPROCESSABLEENTITY, errors.New("Donor name can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[2]) == 0 {
		cErr := &chainError{"validateDonationW", "", CODEUNPROCESSABLEENTITY, errors.New("Affiliated project ID can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[3]) == 0 {
		cErr := &chainError{"validateDonationW", "", CODEUNPROCESSABLEENTITY, errors.New("Item ID can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[4]) == 0 {
		cErr := &chainError{"validateDonationW", "", CODEUNPROCESSABLEENTITY, errors.New("Donation amount can not be empty")}
		return shim.Error(cErr.Error())
	}
	zeroDecimal, _ := decimal.NewFromString("0")
	amount, _ := decimal.NewFromString(args[4])
	if amount.LessThanOrEqual(zeroDecimal) {
		cErr := &chainError{"validateDonationW", "", CODEUNPROCESSABLEENTITY, errors.New("Donation amount can not less than or equal to zero")}
		return shim.Error(cErr.Error())
	}

	epochTime, _ := stub.GetTxTimestamp()
	timeStamp := time.Unix(epochTime.GetSeconds(), 0)

	donBase := donationBase{ProjectID: args[2], ItemID: args[3], Amount: amount.RoundBank(FIXEDPT), TimeStamp: timeStamp}
	d := &donation{ObjectType: DONIN, TxnID: args[0], Donor: args[1], Data: donBase}

	return saveAsset(stub, d)
}

func validateSpendW(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 5 {
		cErr := &chainError{"validateSpendW", "", CODEUNPROCESSABLEENTITY, errors.New("Incorrect no of input args, excepting 5")}
		return shim.Error(cErr.Error())
	}
	if len(args[0]) == 0 {
		cErr := &chainError{"validateSpendW", "", CODEUNPROCESSABLEENTITY, errors.New("Txn ID can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[1]) == 0 {
		cErr := &chainError{"validateSpendW", "", CODEUNPROCESSABLEENTITY, errors.New("Beneficiary can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[2]) == 0 {
		cErr := &chainError{"validateSpendW", "", CODEUNPROCESSABLEENTITY, errors.New("Affiliated project ID can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[3]) == 0 {
		cErr := &chainError{"validateSpendW", "", CODEUNPROCESSABLEENTITY, errors.New("Item ID can not be empty")}
		return shim.Error(cErr.Error())
	}
	if len(args[4]) == 0 {
		cErr := &chainError{"validateSpendW", "", CODEUNPROCESSABLEENTITY, errors.New("Donation amount can not be empty")}
		return shim.Error(cErr.Error())
	}
	zeroDecimal, _ := decimal.NewFromString("0")
	amount, _ := decimal.NewFromString(args[4])
	if amount.LessThanOrEqual(zeroDecimal) {
		cErr := &chainError{"validateSpendW", "", CODEUNPROCESSABLEENTITY, errors.New("Donation amount can not less than or equal to zero")}
		return shim.Error(cErr.Error())
	}

	epochTime, _ := stub.GetTxTimestamp()
	timeStamp := time.Unix(epochTime.GetSeconds(), 0)

	spendData := donationBase{ProjectID: args[2], ItemID: args[3], Amount: amount.RoundBank(FIXEDPT), TimeStamp: timeStamp}
	s := &spend{ObjectType: DONOUT, TxnID: args[0], Benficiary: args[1], Data: spendData}

	return saveAsset(stub, s)
}

func validateProjectR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		cErr := &chainError{"validateProjectR", "", CODEUNPROCESSABLEENTITY, errors.New("Incorrect no of input args, excepting project ID only")}
		return shim.Error(cErr.Error())
	}

	if len(args[0]) == 0 {
		cErr := &chainError{"validateProjectR", "", CODEUNPROCESSABLEENTITY, errors.New("Project ID can not be empty")}
		return shim.Error(cErr.Error())
	}

	p := &project{ProjectID: args[0]}
	return readAsset(stub, p)
}

func validateItemR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		cErr := &chainError{"validateItemR", "", CODEUNPROCESSABLEENTITY, errors.New("Incorrect no of input args, excepting Item ID only")}
		return shim.Error(cErr.Error())
	}
	if len(args[0]) == 0 {
		cErr := &chainError{"validateItemR", "", CODEUNPROCESSABLEENTITY, errors.New("Item ID can not be empty")}
		return shim.Error(cErr.Error())
	}

	it := &item{ItemID: args[0]}
	return readAsset(stub, it)
}

func validateDonationR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		cErr := &chainError{"validateDonationR", "", CODEUNPROCESSABLEENTITY, errors.New("Incorrect no of input args, excepting Transaction ID only")}
		return shim.Error(cErr.Error())
	}
	if len(args[0]) == 0 {
		cErr := &chainError{"validateDonationR", "", CODEUNPROCESSABLEENTITY, errors.New("Txn ID can not be empty")}
		return shim.Error(cErr.Error())
	}

	d := &donation{TxnID: args[0]}
	return readAsset(stub, d)
}

func validateSpendR(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		cErr := &chainError{"validateSpendR", "", CODEUNPROCESSABLEENTITY, errors.New("Incorrect no of input args, excepting Transaction ID only")}
		return shim.Error(cErr.Error())
	}
	if len(args[0]) == 0 {
		cErr := &chainError{"validateSpendR", "", CODEUNPROCESSABLEENTITY, errors.New("Txn ID can not be empty")}
		return shim.Error(cErr.Error())
	}

	s := &spend{TxnID: args[0]}
	return readAsset(stub, s)
}
