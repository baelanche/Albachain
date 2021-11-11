package main

import (
	"encoding/json"
	"fmt"
	//"strconv"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)

type Worker struct {
	workerId string `json:"workerId"`
	workplaceNumber []int `json:"workplaceNumberList"`
}

type Employer struct {
	employerId string `json:"employerId"`
	workplaceNumber []int `json:"workplaceNumberList"`
}

type Workplace struct {
	workplaceNumber int `json:"key"`
	employerId string `json:"employerId"`
	worker []string `json:"workerList"`
	wage int `json:"wage"`
}

type WorkHistory struct {
	workplaceNumber int `json:"key"`
	workerId string `json:"workerId"`
	workStartTime string `json:"workStartTime"`
	workFinishTime string `json:"workFinishTime"`
	wage int `json:"wage"`
	historyCreateTime string `json:"historyCreateTime"`
	historyApprovalTime string `json:"historyApprovalTime"`
	approved bool `json:"approved"`
}

type Albachain struct {
}

func (t *Albachain) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (t *Albachain) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "addWorker" {
		result, err = t.addWorker(stub, args)
	} else if fn == "getWorker" {

	} else if fn == "addEmployer" {

	} else if fn == "getEmployer" {

	} else if fn == "addWorkplace" {

	} else if fn == "addWorkhisHistory" {

	}

	if err != nil {return shim.Error(err.Error())}
	return shim.Success([]byte(result))
}

func (t *Albachain) addWorker(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {return "", fmt.Errorf("Call addWorker failed")}

	id := args[0]
	var workplaceNumber []int
	var value = Worker{workerId: args[0], workplaceNumber: workplaceNumber}
	valueAsBytes, _ := json.Marshal(value)
	err := stub.PutState(id, valueAsBytes)

	if err != nil {return "", fmt.Errorf("Error during addWorker function")}
	return string(valueAsBytes), nil
}

func main() {
        if err := shim.Start(new(Albachain)); err != nil {
                fmt.Printf("Error creating new Smart Contract: %s", err)
        }
}

