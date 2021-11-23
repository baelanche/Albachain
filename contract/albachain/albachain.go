package main

import (
	"encoding/json"
	"fmt"
	//"strconv"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)

type Worker struct {
	WorkerId string `json:"WorkerId"`
	WorkerName string `json:"WorkerName"`
	WorkplaceNumber []string `json:"WorkplaceNumberList"`
}

type Employer struct {
	EmployerId string `json:"EmployerId"`
	WorkplaceNumber []string `json:"WorkplaceNumber"`
}

type Workplace struct {
	WorkplaceNumber string `json:"WorkplaceNumber"`
	EmployerId string `json:"EmployerId"`
	Workers []string `json:"WorkerList"`
	Wage int `json:"Wage"`
}

type WorkHistory struct {
	WorkHistoryNumber string `json:"WorkHistoryNumber"`
	WorkerId string `json:"WorkerId"`
	WorkplaceNumber string `json:"WorkplaceNumber"`
	WorkStartTime string `json:"WorkStartTime"`
	WorkFinishTime string `json:"WorkFinishTime"`
	Wage int `json:"Wage"`
	HistoryCreateTime string `json:"HistoryCreateTime"`
	HistoryApprovalTime string `json:"HistoryApprovalTime"`
	Approved bool `json:"Approved"`
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
		t.addWorker(stub, args)
	} else if fn == "getWorker" {
		result, err = t.getWorker(stub, args)
	} else if fn == "addEmployer" {
		result, err = t.addEmployer(stub, args)
	} else if fn == "getEmployer" {
		result, err = t.getEmployer(stub, args)
	} else if fn == "addWorkplace" {
		result, err = t.addWorkplace(stub, args)
	} else if fn == "addWorkHistory" {

	} else {return shim.Error(err.Error())}

	if err != nil {return shim.Error(err.Error())}
	return shim.Success([]byte(result))
}

func (t *Albachain) addWorker(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {return "", fmt.Errorf("Call addWorker failed")}

	/* duplicate check */
	id, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get worker: %s", err)}
	if id != nil {return "", fmt.Errorf("This id already exists")}

	var wNumber []string
	wNumber = append(wNumber, "P001")
	var value = Worker{WorkerId: args[0], WorkerName: args[1], WorkplaceNumber: wNumber}
	valueAsBytes, _ := json.Marshal(value)
	stub.PutState(args[0], valueAsBytes)

	return string(valueAsBytes), nil
}

func (t *Albachain) getWorker(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {return "", fmt.Errorf("Failed to call addWorker")}

	value, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get worker: %s", err)}
	if value == nil {return "", fmt.Errorf("Worker not found: %s", args[0])}
	return string(value), nil
}

func (t *Albachain) addEmployer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {return "", fmt.Errorf("Failed to call addEmployer")}

	/* duplicate check */
	id, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get employer: %s", err)}
	if id != nil {return "", fmt.Errorf("This id already exists")}

	var workplaceNumber []string
	var value = Employer{EmployerId: args[0], WorkplaceNumber: workplaceNumber}
	valueAsBytes, _ := json.Marshal(value)
	err2 := stub.PutState(args[0], valueAsBytes)

	if err2 != nil {return "", fmt.Errorf("Error during addEmployer function")}
	return string(valueAsBytes), nil
}

func (t *Albachain) getEmployer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {return "", fmt.Errorf("Call getEmployer failed")}

	value, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get employer: %s", err)}
	if value == nil {return "", fmt.Errorf("Employer not found: %s", args[0])}
	return string(value), nil
}

func (t *Albachain) addWorkplace(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {return "", fmt.Errorf("Call addWorkplace failed")}

	workerAsBytes, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf(err.Error())}
	if workerAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	worker := Worker{}
	err = json.Unmarshal(workerAsBytes, &worker)
	if err != nil {return "", fmt.Errorf(err.Error())}

	worker.WorkplaceNumber = append(worker.WorkplaceNumber, args[1])
	workerAsBytes, _ = json.Marshal(worker)
	stub.PutState(args[0], workerAsBytes)

	return string(args[1]), nil
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