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
	workplaceNumber int `json:"workplaceNumber"`
	employerId string `json:"employerId"`
	worker []string `json:"workerList"`
	wage int `json:"wage"`
}

type WorkHistory struct {
	workerId string `json:"workerId"`
	workplaceNumber int `json:"workplaceNumber"`
	workStartTime string `json:"workStartTime"`
	workFinishTime string `json:"workFinishTime"`
	wage int `json:"wage"`
	historyCreateTime string `json:"historyCreateTime"`
	historyApprovalTime string `json:"historyApprovalTime"`
	approved bool `json:"approved"`
}

type Albachain struct {}

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
		result, err = t.getWorker(stub, args)
	} else if fn == "addEmployer" {

	} else if fn == "getEmployer" {

	} else if fn == "addWorkplace" {

	} else if fn == "addWorkHistory" {

	}

	if err != nil {return shim.Error(err.Error())}
	return shim.Success([]byte(result))
}

func (t *Albachain) addWorker(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {return "", fmt.Errorf("Call addWorker failed")}

	/* duplicate check */
	id, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get worker: %s", err)}
	if id != nil {return "", fmt.Errorf("This id already exists")}

	var workplaceNumber []int
	var value = Worker{workerId: args[0], workplaceNumber: workplaceNumber}
	valueAsBytes, _ := json.Marshal(value)
	err2 := stub.PutState(args[0], valueAsBytes)

	if err2 != nil {return "", fmt.Errorf("Error during addWorker function")}
	return string(valueAsBytes), nil
}

func (t *Albachain) getWorker(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {return "", fmt.Errorf("Call getWorker failed")}

	value, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get worker: %s", err)}
	if value == nil {return "", fmt.Errorf("Worker not found: %s", args[0])}
	return string(value), nil
}

/*
func (t *Albachain) getAllWorker(stub shim.ChaincodeStubInterface) []string {
	start := ""
	end := ""

	resultIterater, err := stub.GetStateByRange(start, end)
	if err != nil {return "", fmt.Errorf(err.Error())}
	defer resultIterater.Close()

	var i int = 0
	var idList []string
	for resultIterater.HasNext() {
		id, err := resultIterater.Next()
		if err != nil {return "", fmt.Errorf(err.Error())}
		idList = append(idList, id.Key)
		i++
		fmt.Printf(idList[i])
	}

	return result
}
*/

func main() {
        if err := shim.Start(new(Albachain)); err != nil {fmt.Printf("Error creating new Albachain: %s", err)}
}