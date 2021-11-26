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
	} else if fn == "addWorkHistory" {
		result, err = t.addWorkHistory(stub, args)
	} else if fn == "workHistoryApproval" {
		result, err = t.workHistoryApproval(stub, args)
	} else if fn == "addWorkplaceByEmployer" {
		result, err = t.addWorkplaceByEmployer(stub, args)
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
	ymd := now.Format("20010101")

	var workplaceList []Workplace
	var value = Employer{EmployerId: args[0], EmployerName: args[1], WorkplaceList: workplaceList, JoinDate: ymd}
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
고용주 근무지 추가
param : EmployerId, EmployerName, WorkplaceNumber, WorkplaceName, Wage
*/
func (t *Albachain) addWorkplaceByEmployer(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 5 {return "", fmt.Errorf("Call addWorkplaceByEmployer failed")}
	
	employerAsBytes, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf(err.Error())}
	if employerAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	employer := Employer{}
	err = json.Unmarshal(employerAsBytes, &employer)
	if err != nil {return "", fmt.Errorf(err.Error())}

	now := time.Now()
	ymd := now.Format("20010101")

	var workerList []string
	var value = Workplace{WorkplaceNumber: args[2], WorkplaceName: args[3], EmployerId: args[0], EmployerName: args[1], WorkerList: workerList, DefaultWage: args[4], JoinDate: ymd, RetireDate: ""}
	workplaceAsBytes, _ := json.Marshal(value)
	stub.PutState(args[2], workplaceAsBytes)

	employer.WorkplaceList = append(employer.WorkplaceList, value)
	
	employerAsBytes, _ = json.Marshal(employer)
	stub.PutState(args[0], employerAsBytes)

	return string(args[2]), nil
}

/*
노동자 근무지 추가 승인
param : WorkerId, WorkplaceNumber, Wage
*/
func (t *Albachain) workplaceApproval(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {return "", fmt.Errorf("Call workplaceApproval failed")}
	
	workerAsBytes, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf(err.Error())}
	if workerAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	worker := Worker{}
	err = json.Unmarshal(workerAsBytes, &worker)
	if err != nil {return "", fmt.Errorf(err.Error())}

	now := time.Now()
	ymd := now.Format("20010101")

	worker.WorkJoinDate = ymd
	worker.Wage = args[2]
	worker.Approved = true
	
	workerAsBytes, _ = json.Marshal(worker)
	stub.PutState(args[0], workerAsBytes)

	workplaceAsBytes, err2 := stub.GetState(args[1])
	if err2 != nil {return "", fmt.Errorf(err.Error())}
	if workplaceAsBytes == nil {return "", fmt.Errorf("The wrong approach")}

	workplace := Workplace{}
	err2 = json.Unmarshal(workplaceAsBytes, &workplace)
	if err2 != nil {return "", fmt.Errorf(err.Error())}

	workplace.WorkerList = append(workplace.WorkerList, args[0])
	workplaceAsBytes, _ = json.Marshal(workplace)
	stub.PutState(args[1], workplaceAsBytes)

	return string(args[0]), nil
}

/*
근무 기록 추가
param : WorkHistoryNumber, WorkerId, WorkerName, WorkplaceNumber, WorkplaceName, WorkStartTime, WorkFinishTime, Wage
*/
func (t *Albachain) addWorkHistory(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 8 {return "", fmt.Errorf("Failed to call addWorkHistory")}

	/* duplicate check */
	historyNumber, err := stub.GetState(args[0])
	if err != nil {return "", fmt.Errorf("Failed to get workhistory: %s", err)}
	if historyNumber != nil {return "", fmt.Errorf("This history id already exists")}

	now := time.Now()
	ymdhms := now.Format("2001-01-01 15:01:05")

	var value = WorkHistory{WorkHistoryNumber: args[0], WorkerId: args[1], WorkerName: args[2], WorkplaceNumber: args[3], WorkplaceName: args[4], WorkStartTime: args[5], WorkFinishTime: args[6], Wage: args[7], HistoryCreateTime: ymdhms, HistoryApprovalTime: "", Approved: false}
	valueAsBytes, _ := json.Marshal(value)
	err2 := stub.PutState(args[0], valueAsBytes)

	if err2 != nil {return "", fmt.Errorf("Error during addWorkHistory function")}
	return string(valueAsBytes), nil
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

/*
func (t *SimpleChaincode) transferMarblesBasedOnColor(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0       1
	// "color", "bob"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	color := args[0]
	newOwner := strings.ToLower(args[1])
	fmt.Println("- start transferMarblesBasedOnColor ", color, newOwner)

	// Query the color~name index by color
	// This will execute a key range query on all keys starting with 'color'
	coloredMarbleResultsIterator, err := stub.GetStateByPartialCompositeKey("color~name", []string{color})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer coloredMarbleResultsIterator.Close()

	// Iterate through result set and for each marble found, transfer to newOwner
	var i int
	for i = 0; coloredMarbleResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble name from the composite key
		responseRange, err := coloredMarbleResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// get the color and name from color~name composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedColor := compositeKeyParts[0]
		returnedMarbleName := compositeKeyParts[1]
		fmt.Printf("- found a marble from index:%s color:%s name:%s\n", objectType, returnedColor, returnedMarbleName)

		// Now call the transfer function for the found marble.
		// Re-use the same function that is used to transfer individual marbles
		response := t.transferMarble(stub, []string{returnedMarbleName, newOwner})
		// if the transfer failed break out of loop and return error
		if response.Status != shim.OK {
			return shim.Error("Transfer failed: " + response.Message)
		}
	}

	responsePayload := fmt.Sprintf("Transferred %d %s marbles to %s", i, color, newOwner)
	fmt.Println("- end transferMarblesBasedOnColor: " + responsePayload)
	return shim.Success([]byte(responsePayload))
}
*/

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