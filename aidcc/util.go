package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// respnse struct to have consistent response structure for all chaincode invokes
type response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Payload []byte `json:"payload"`
}

func (r *response) formatResponse() []byte {
	var buffer bytes.Buffer
	buffer.WriteString("{\"code\":")
	buffer.WriteString("\"")
	buffer.WriteString(r.Code)
	buffer.WriteString("\",")
	buffer.WriteString("\"message\":")
	buffer.WriteString("\"")
	buffer.WriteString(r.Message)
	if r.Payload == nil {
		buffer.WriteString("\"}")
	} else {
		buffer.WriteString("\",")
		buffer.WriteString("\"payload\":")
		buffer.Write(r.Payload)
		buffer.WriteString("}")
	}
	return buffer.Bytes()
}

// chainError - custom error
type chainError struct {
	fcn  string // function/method
	key  string // associated key, if any
	code string // Error code
	err  error  // Error.
}

func (e *chainError) Error() string {
	return e.fcn + " " + e.key + " " + e.code + ": " + e.err.Error()
}

// checkAsset - check if an asset ( with a key) is already available on the ledger
func checkAsset(stub shim.ChaincodeStubInterface, assetID string) (bool, *chainError) {

	assetBytes, err := stub.GetState(assetID)

	if err != nil {
		e := &chainError{"checkAsset", assetID, CODEGENEXCEPTION, err}
		return false, e
	} else if assetBytes != nil {
		//e := &chainError{"checkAsset", assetID, CODEAlRDEXIST, errors.New("Asset with key already exists")}
		return true, nil
	}
	return false, nil
}

// queryAsset - return query state from the ledger
func queryAsset(stub shim.ChaincodeStubInterface, assetID string) ([]byte, *chainError) {

	assetBytes, err := stub.GetState(assetID)

	if err != nil {
		e := &chainError{"queryAsset", assetID, CODEGENEXCEPTION, err}
		return nil, e
	} else if assetBytes == nil {
		e := &chainError{"queryAsset", assetID, CODENOTFOUND, errors.New("Asset ID not found")}
		return nil, e
	}
	return assetBytes, nil
}

// getCallerID - reterive caller id from ECert
func getCallerID(stub shim.ChaincodeStubInterface) (string, *chainError) {
	id, err := cid.New(stub)
	callerID, err := id.GetID()
	if err != nil {
		e := &chainError{"getCallerID", "", CODEGENEXCEPTION, err}
		return "", e
	}
	// decode the returned base64 string
	data, err := base64.StdEncoding.DecodeString(callerID)
	if err != nil {
		e := &chainError{"getCallerID", "", CODEGENEXCEPTION, err}
		return "", e
	}
	l := strings.Split(string(data), "::")
	return l[1], nil
}
