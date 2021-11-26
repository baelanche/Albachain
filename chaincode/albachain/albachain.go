package main

import (
	"encoding/json"
	"fmt"
	"time"
	//"strconv"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/hyperledger/fabric/protos/peer"
)

type Worker struct {
	WorkerId string `json:"WorkerId"`
	WorkerName string `json:"WorkerName"`
	WorkplaceNumber string `json:"WorkplaceNumber"`
	WorkplaceName string `json:"WorkplaceName"`
	WorkJoinDate string `json:"WorkJoinDate"`
	WorkRetireDate string `json:"WorkRetireDate"`
	Wage string `json:"Wage"`
	Approved bool `json:"Approved"`
}

type Workplace struct {
	WorkplaceNumber string `json:"WorkplaceNumber"`
	WorkplaceName string `json:"WorkplaceName"`
	EmployerId string `json:"EmployerId"`
	EmployerName string `json:"EmployerName"`
	Workers []string `json:"WorkerList"`
	DefaultWage string `json:"DefaultWage"`
	JoinDate string `json:"JoinDate"`
	RetireDate string `json:"RetireDate"`
}

type Employer struct {
	EmployerId string `json:"EmployerId"`
	EmployerName string `json:"EmployerName"`
	WorkplaceList []Workplace `json:"WorkplaceList"`
	JoinDate string `json:"JoinDate"`
}

type WorkHistory struct {
	WorkHistoryNumber string `json:"WorkHistoryNumber"`
	WorkerId string `json:"WorkerId"`
	WorkplaceNumber string `json:"WorkplaceNumber"`
	WorkStartTime string `json:"WorkStartTime"`
	WorkFinishTime string `json:"WorkFinishTime"`
	Wage string `json:"Wage"`
	HistoryCreateTime string `json:"HistoryCreateTime"`
	HistoryApprovalTime string `json:"HistoryApprovalTime"`
	Approved bool `json:"Approved"`
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
	} else if fn == "addWorkplace" {
		result, err = t.addWorkplace(stub, args)
	} else if fn == "deleteWorkplace" {
		result, err = t.deleteWorkplace(stub, args)	
	} else if fn == "addEmployer" {
		result, err = t.addEmployer(stub, args)
	} else if fn == "getEmployer" {
		result, err = t.getEmployer(stub, args)
	} else if fn == "workplaceApproval" {
		result, err = t.workplaceApproval(stub, args)
	} else {return shim.Error(err.Error())}

	if err != nil {return shim.Error(err.Error())}
	return shim.Success([]byte(result))
}

/* 
노동자 가입
param : WorkerId, WorkerName 
*/
func (t *Albachain) addWorker(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {return "", fmt.Errorf("Call addWorker failed")}

	/* duplicate check */
	id, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get worker: %s", err)}
	if id != nil {return "", fmt.Errorf("This id already exists")}

	var value = Worker{WorkerId: args[0], WorkerName: args[1], WorkplaceNumber: "", WorkplaceName: "", WorkJoinDate: "", WorkRetireDate: "", Wage: "0", Approved: false}
	valueAsBytes, _ := json.Marshal(value)
	stub.PutState(args[0], valueAsBytes)

	return string(valueAsBytes), nil
}

/*
노동자 조회
param : WorkerId
*/
func (t *Albachain) getWorker(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {return "", fmt.Errorf("Failed to call addWorker")}

	value, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get worker: %s", err)}
	if value == nil {return "", fmt.Errorf("Worker not found: %s", args[0])}
	return string(value), nil
}

/* 
노동자 근무지 추가
param : WorkerId, WorkplaceNumber, WorkplaceName, Wage 
*/
func (t *Albachain) addWorkplace(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 4 {return "", fmt.Errorf("Call addWorkplace failed")}
	
	workerAsBytes, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf(err.Error())}
	if workerAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	worker := Worker{}
	err = json.Unmarshal(workerAsBytes, &worker)
	if err != nil {return "", fmt.Errorf(err.Error())}

	worker.WorkplaceNumber = args[1]
	worker.WorkplaceName = args[2]
	worker.WorkJoinDate = ""
	worker.Wage = args[3]
	
	workerAsBytes, _ = json.Marshal(worker)
	stub.PutState(args[0], workerAsBytes)

	return string(args[1]), nil
}

/* 
노동자 근무지 삭제
param : WorkerId 
*/ 
func (t *Albachain) deleteWorkplace(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {return "", fmt.Errorf("Call deleteWorkplace failed")}

	workerAsBytes, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf(err.Error())}
	if workerAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	worker := Worker{}
	err = json.Unmarshal(workerAsBytes, &worker)
	if err != nil {return "", fmt.Errorf(err.Error())}

	now := time.Now()
	yms := now.Format("20010101")

	worker.WorkplaceNumber = ""
	worker.WorkplaceName = ""
	worker.WorkRetireDate = yms
	worker.Wage = "0"
	worker.Approved = false

	workerAsBytes, _ = json.Marshal(worker)
	stub.PutState(args[0], workerAsBytes)

	return string(args[0]), nil
}

/*
고용주 가입
param : EmployerId, EmployerName
*/
func (t *Albachain) addEmployer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {return "", fmt.Errorf("Failed to call addEmployer")}

	/* duplicate check */
	id, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get employer: %s", err)}
	if id != nil {return "", fmt.Errorf("This id already exists")}

	now := time.Now()
	yms := now.Format("20010101")

	var workplaceList []Workplace
	var value = Employer{EmployerId: args[0], EmployerName: args[1], WorkplaceList: workplaceList, JoinDate: yms}
	valueAsBytes, _ := json.Marshal(value)
	err2 := stub.PutState(args[0], valueAsBytes)

	if err2 != nil {return "", fmt.Errorf("Error during addEmployer function")}
	return string(valueAsBytes), nil
}

/*
고용주 조회
param : EmployerId
*/
func (t *Albachain) getEmployer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {return "", fmt.Errorf("Call getEmployer failed")}

	value, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get employer: %s", err)}
	if value == nil {return "", fmt.Errorf("Employer not found: %s", args[0])}
	return string(value), nil
}

/*
노동자 근무지 추가 승인
param : WorkerId, Wage
*/
func (t *Albachain) workplaceApproval(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {return "", fmt.Errorf("Call workplaceApproval failed")}
	
	workerAsBytes, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf(err.Error())}
	if workerAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	worker := Worker{}
	err = json.Unmarshal(workerAsBytes, &worker)
	if err != nil {return "", fmt.Errorf(err.Error())}

	now := time.Now()
	yms := now.Format("20010101")

	worker.WorkJoinDate = yms
	worker.Wage = args[1]
	worker.Approved = true
	
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