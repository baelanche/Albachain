package main
// main 패키지는 main 함수를 반드시 가져야한다. 시작함수가 되는 패키지임
import (
	"encoding/json"
        "fmt" // 기본
	"strconv"

        "github.com/hyperledger/fabric/core/chaincode/shim"
        "github.com/hyperledger/fabric/protos/peer" // 하이퍼레저 패브릭 라이브러리
)

type MyAccount struct {
	AccountNum string `json:"key"`
	Value int `json:"value"`
	Owner string `json:"owner"`
}

type SimpleAsset struct {
}

func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

/*
func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {
	args := stub.GetStringArgs()
	if len(args) != 2 {
		return shim.Error("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1])) // key, value(문자열을 그대로 저장할 수 없어서 byte 로 바꿈) 순서로 put
	if err != nil { // 에러가 없지 않으면
		return shim.Error(fmt.Sprintf("Failed to create asset: %s", args[0]))
	}
	return shim.Success(nil)
}
*/

func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = t.setMyAccount(stub, args)
	} else if fn == "get" {
		result, err = t.getMyAccount(stub, args)
	} else {
		return shim.Error("Not supported chaincode function !!")
	}
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

/*
func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	err := stub.PutState(args[0], []byte(args[1]))
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return args[1], nil
}
*/

func (t *SimpleAsset) getMyAccount(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s", err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

func (t *SimpleAsset) setMyAccount(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect arguments !!")
	}

	value, _ := strconv.Atoi(args[1])

	var data = MyAccount{AccountNum: args[0], Value: value, Owner: args[2]}
	dataAsBytes, _ := json.Marshal(data)
	err := stub.PutState(args[0], dataAsBytes)

	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", err)
	}
	return string(dataAsBytes), nil
}

func main() {
        if err := shim.Start(new(SimpleAsset)); err != nil {
                fmt.Printf("Error creating new Smart Contract: %s", err)
        }
}

