package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// const (
// 	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	letterIdxBits = 6                    // 6 bits to represent a letter index
// 	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
// 	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
// )

func main() {
	// inputBuf = RandStringBytesMaskImpr(1024)

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 1 {
		return shim.Error("error")
	}
	argSize, err := parseSizeStr(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	inputBuf := RandStringBytesMaskImpr(argSize, "0", args[0])

	err = stub.PutState("0", inputBuf)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()

	if function == "put" {
		return t.put(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}
	return shim.Error("Invalid invoke function name. Expecting \"put\" \"query\"")
}

func (t *SimpleChaincode) put(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	count := len(args)
	if count != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to put")
	}
	argSize, err := parseSizeStr(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	inputBuf := RandStringBytesMaskImpr(argSize, args[0], args[1])
	err = stub.PutState(args[0], inputBuf)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	// jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	// fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func parseSizeStr(sizeStr string) (bsize int, err error) {
	var breakIndex int
	for i, v := range []byte(sizeStr) {
		if v < 48 || v > 57 {
			breakIndex = i
			break
		}
	}
	size, err := strconv.Atoi(string(sizeStr[:breakIndex]))
	if err != nil {
		return 0, err
	}
	unit := sizeStr[breakIndex:]
	switch unit {
	case "B", "b":
		return size, nil
	case "KB", "K", "kb", "k":
		return size * 1024, nil
	case "MB", "mb", "M", "m":
		return size * 1024 * 1024, nil
	}
	return size, fmt.Errorf("parse %s failed:unsupport unit", sizeStr)
}

func RandStringBytesMaskImpr(n int, index string, size string) []byte {
	count := len(size)
	tmp := n / count
	tmp += 1
	str := index
	for i := 0; i < tmp; i++ {
		str += size
	}
	return []byte(str)[:n]
}
