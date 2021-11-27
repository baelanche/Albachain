package main

import (
	"encoding/json"
	"fmt"
	"time"
	"bytes"

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
	WorkerList []string `json:"WorkerList"`
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
	WorkerName string `json:"WorkerName"`
	WorkplaceNumber string `json:"WorkplaceNumber"`
	WorkplaceName string `json:"WorkplaceName"`
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
		return t.getWorker(stub, args)
	} else if fn == "addEmployer" {
		result, err = t.addEmployer(stub, args)
	} else if fn == "getEmployer" {
		return t.getEmployer(stub, args)
	} else if fn == "getWorkplace" {
		return t.getWorkplace(stub, args)
	} else if fn == "addWorkHistory" {
		result, err = t.addWorkHistory(stub, args)
	} else if fn == "getAllWorkHistory" {
		return t.getAllWorkHistory(stub, args)
	} else if fn == "workHistoryApproval" {
		result, err = t.workHistoryApproval(stub, args)
	} else {return shim.Error(err.Error())}

	if err != nil {return shim.Error(err.Error())}
	return shim.Success([]byte(result))
}

/* 
노동자 가입
param : WorkerId, WorkerName, WorkplaceNumber, WorkplaceName, Wage
*/
func (t *Albachain) addWorker(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 5 {return "", fmt.Errorf("Call addWorker failed")}

	/* duplicate check */
	id, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get worker: %s", err)}
	if id != nil {return "", fmt.Errorf("This id already exists")}

	var value = Worker{WorkerId: args[0], WorkerName: args[1], WorkplaceNumber: args[2], WorkplaceName: args[3], WorkJoinDate: "", WorkRetireDate: "", Wage: args[4], Approved: true}
	valueAsBytes, _ := json.Marshal(value)
	stub.PutState(args[0], valueAsBytes)

	wpAsBytes, err := stub.GetState(args[2])
	if err != nil {return "", fmt.Errorf("Failed to get worker: %s", err)}

	wp := Workplace{}
	err = json.Unmarshal(wpAsBytes, &wp)
	if err != nil {return "", fmt.Errorf(err.Error())}

	wp.WorkerList = append(wp.WorkerList, args[1]);
	wpAsBytes, _ = json.Marshal(wp)
	stub.PutState(args[2], wpAsBytes)
	return string(valueAsBytes), nil
}

/*
노동자 조회
param : WorkerId
*/
func (t *Albachain) getWorker(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {shim.Error("Failed to call getWorker")}

	value, err := stub.GetState(args[0])
	if err != nil {shim.Error("Failed to get worker")}
	if value == nil {shim.Error("Worker not found")}
	return shim.Success([]byte(value))
}

/*
고용주 가입
param : EmployerId, EmployerName, WorkplaceNumber, WorkplaceName, JoinDate, Wage
*/
func (t *Albachain) addEmployer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 6 {return "", fmt.Errorf("Call addEmployer failed")}

	/* duplicate check */
	id, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get worker: %s", err)}
	if id != nil {return "", fmt.Errorf("This id already exists")}

	var workerList []string
	var workplaceL []Workplace
	var workplace = Workplace{WorkplaceNumber: args[2], WorkplaceName: args[3], EmployerId: args[0], EmployerName: args[1], WorkerList: workerList, DefaultWage: args[5], JoinDate: args[4], RetireDate: ""}
	workplaceL = append(workplaceL, workplace)
	var value = Employer{EmployerId: args[0], EmployerName: args[1], WorkplaceList: workplaceL, JoinDate: args[4]}
	valueAsBytes, _ := json.Marshal(value)
	stub.PutState(args[0], valueAsBytes)
	wpAsBytes, _ := json.Marshal(workplace)
	stub.PutState(args[2], wpAsBytes)

	return string(valueAsBytes), nil
}

/*
고용주 조회
param : EmployerId
*/
func (t *Albachain) getEmployer(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {shim.Error("Failed to call getEmployer")}

	value, err := stub.GetState(args[0])
	if err != nil {shim.Error("Failed to get employer")}
	if value == nil {shim.Error("Employer not found")}
	return shim.Success([]byte(value))
}

/*
근무지 조회
param : WorkplaceNumber
*/
func (t *Albachain) getWorkplace(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {shim.Error("Failed to call getWorkplace")}

	value, err := stub.GetState(args[0])
	if err != nil {shim.Error("Failed to get workplace")}
	if value == nil {shim.Error("Employer not found")}
	return shim.Success([]byte(value))
}

/* 
노동자 근무지 삭제
param : WorkerId , WorkplaceNumber
*/ 
func (t *Albachain) deleteWorkplace(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {return "", fmt.Errorf("Call deleteWorkplace failed")}

	workerAsBytes, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf(err.Error())}
	if workerAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	worker := Worker{}
	err = json.Unmarshal(workerAsBytes, &worker)
	if err != nil {return "", fmt.Errorf(err.Error())}

	now := time.Now()
	ymd := now.Format("20010101")

	worker.WorkplaceNumber = ""
	worker.WorkplaceName = ""
	worker.WorkRetireDate = ymd
	worker.Wage = "0"
	worker.Approved = false

	workerAsBytes, _ = json.Marshal(worker)
	stub.PutState(args[0], workerAsBytes)

	workplaceAsBytes, err2 := stub.GetState(args[1])
	if err2 != nil {return "", fmt.Errorf(err.Error())}
	if workplaceAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	workplace := Workplace{}
	err = json.Unmarshal(workplaceAsBytes, &workplace)
	if err != nil {return "", fmt.Errorf(err.Error())}

	remove := []string{args[0]}
	
loop:
	for i:= 0; i<len(workplace.WorkerList); i++ {
		w := workplace.WorkerList[i]
		for _, rem := range remove {
			if w == rem {
				workplace.WorkerList = append(workplace.WorkerList[:i], workplace.WorkerList[i+1:]...)
				i--
				continue loop
			}
		}
	}

	workplaceAsBytes, _ = json.Marshal(workplace)
	stub.PutState(args[1], workplaceAsBytes)

	return string(args[0]), nil
}

/*
근무 기록 추가
param : WorkHistoryNumber, WorkerId, WorkerName, WorkplaceNumber, WorkplaceName, WorkStartTime, WorkFinishTime, Wage, HistoryCreateTime
*/
func (t *Albachain) addWorkHistory(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 9 {return "", fmt.Errorf("Failed to call addWorkHistory")}

	/* duplicate check */
	historyNumber, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get workhistory: %s", err)}
	if historyNumber != nil {{return "", fmt.Errorf("This history id already exists")}}

	var value = WorkHistory{WorkHistoryNumber: args[0], WorkerId: args[1], WorkerName: args[2], WorkplaceNumber: args[3], WorkplaceName: args[4], WorkStartTime: args[5], WorkFinishTime: args[6], Wage: args[7], HistoryCreateTime: args[8], HistoryApprovalTime: "", Approved: false}
	valueAsBytes, _ := json.Marshal(value)
	err2 := stub.PutState(args[0], valueAsBytes)

	if err2 != nil {return "", fmt.Errorf("Error during addWorkHistory function")}
	return string(valueAsBytes), nil
}

func (t *Albachain) getAllWorkHistory(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	end := args[0]
	end = end[:11] + "9999"
	startKey := args[0]
	endKey := end

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

/*
노동자 근무 기록 추가 승인
param : WorkHistoryNumber, Wage
*/
func (t *Albachain) workHistoryApproval(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {return "", fmt.Errorf("Call workHistoryApproval failed")}
	
	workHistoryAsBytes, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf(err.Error())}
	if workHistoryAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	workHistory := WorkHistory{}
	err = json.Unmarshal(workHistoryAsBytes, &workHistory)
	if err != nil {return "", fmt.Errorf(err.Error())}

	now := time.Now()
	ymdhms := now.Format("2001-01-01 15:01:05")

	workHistory.HistoryApprovalTime = ymdhms
	workHistory.Approved = true
	
	workHistoryAsBytes, _ = json.Marshal(workHistory)
	stub.PutState(args[0], workHistoryAsBytes)

	return string(args[0]), nil
}

func main() {
        if err := shim.Start(new(Albachain)); err != nil {fmt.Printf("Error creating new Albachain: %s", err)}
}