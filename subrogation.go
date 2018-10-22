//Author: renju vm
//Project: Subrogation


package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"strings"
	"time"
)

// Subrogationcode klaim registration chaincode
type Subrogationcode struct {
}

var certIndexStr = "_certindex" //name for the key/value that will store a list of all known certs
var opentransStr = "_opentrans" //name for the key/value that will store all klaims

type Claim struct {
	Claimref     string `json:"claimref"` //the fieldtags are needed to keep track klaim
	Tortname     string `json:"tortname"`
	Tortinsuarer     string `json:"tortinsuarer"`
	Accidentdate        string `json:"accidentdate"`
	Policynumber        string `json:"policynumber"`
	Carrieremail    string `json:"carrieremail"`
	Claimamount string `json:"claimamount"`
}

type Priliminary struct {
	Claimref string `json:"claimref"`
	Insurer  string `json:"insuarer"`
	Adjustername  string `json:"adjustername"`
	Insuredname  string `json:"insuredname"`
	Referance  string `json:"referance"`
	Policylimits  string `json:"policylimits"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(Subrogationcode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *Subrogationcode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("start", []byte(strconv.Itoa(Aval))) //making a test var "start", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(certIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations
// ============================================================================================================================
func (t *Subrogationcode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *Subrogationcode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		res, err := t.Init(stub, "init", args)
		return res, err
	} else if function == "reg_claim" { //create a new  klaim
		res, err := t.reg_claim(stub, args)
		return res, err
	} else if function == "reg_priliminaries" { //create a new  klaim
		res, err := t.reg_priliminaries(stub, args)
		return res, err
	}

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Quering klaims
// ============================================================================================================================
func (t *Subrogationcode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "getAllclaims" {
		return t.getAllclaims(stub, args)
	} else if function == "getClaim" { //read a variable
		return t.getClaim(stub, args)
	} else if function == "getPriliminaries" { //read a variable
		return t.getPriliminaries(stub, args)
	}

	return nil, errors.New("Received unknown function query")
}



// ============================================================================================================================
// Read all - read all matching variable from chaincode state
// ============================================================================================================================
func (t *Subrogationcode) getAllclaims(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	keysIter, err := stub.RangeQueryState("", "")
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}
	defer keysIter.Close()

	var keys []Claim
	for keysIter.HasNext() {
		key, _, iterErr := keysIter.Next()
		if iterErr != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}
		vals, err := stub.GetState(key)
		if err != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}

		var klaim Claim
		json.Unmarshal(vals, &klaim)

if klaim.Tortname != "" {
		keys = append(keys, klaim)
	}
	}

	jsonKeys, err := json.Marshal(keys)
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error marshaling JSON: %s", err)
	}

	return jsonKeys, nil

}

// ============================================================================================================================
// Validate - validate a variable from chaincode state
// ============================================================================================================================
func (t *Subrogationcode) getClaim(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var klaimref string

	klaimref = strings.ToLower(args[0])

	keysIter, err := stub.RangeQueryState("", "")
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}
	defer keysIter.Close()

	var keys []Claim

	for keysIter.HasNext() {
		key, _, iterErr := keysIter.Next()
		if iterErr != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}
		vals, err := stub.GetState(key)
		if err != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}

		var klaim Claim
		json.Unmarshal(vals, &klaim)

		if klaim.Claimref == klaimref && klaim.Tortname != "" {
			keys = append(keys, klaim)
		}
	}

	jsonKeys, err := json.Marshal(keys)
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error marshaling JSON: %s", err)
	}

	return jsonKeys, nil

}

// ============================================================================================================================
// Validate - validate invoice from chaincode state
// ============================================================================================================================
func (t *Subrogationcode) getPriliminaries(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var klaimref string

	klaimref = strings.ToLower(args[0])

	if klaimref == "" {
		return nil, errors.New("Referance is missing")
	}

	keysIter, err := stub.RangeQueryState("", "")
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}
	defer keysIter.Close()

	var keys []Priliminary

	for keysIter.HasNext() {
		key, _, iterErr := keysIter.Next()
		if iterErr != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}
		vals, err := stub.GetState(key)
		if err != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}

		var priliminary Priliminary
		json.Unmarshal(vals, &priliminary)

		if priliminary.Insurer != "" {
			if priliminary.Claimref == klaimref {
				keys = append(keys, priliminary)
			}
		}

	}

	jsonKeys, err := json.Marshal(keys)
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error marshaling JSON: %s", err)
	}

	return jsonKeys, nil

}

// ============================================================================================================================
// Init claim - create a new claim, store into chaincode state
// ============================================================================================================================
func (t *Subrogationcode) reg_claim(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	ctime := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	claimref := strings.ToLower(args[0])
	tortname := strings.ToLower(args[1])
	tortinsuarer := strings.ToLower(args[2])
	accidentdate := strings.ToLower(args[3])
	policynumber := strings.ToLower(args[4])
	carrieremail := strings.ToLower(args[5])
	claimamount := strings.ToLower(args[6])

if(tortname != ""){
	//build the cert json string manually
	str := `{"claimref": "` + claimref + `", "tortname": "` + tortname + `", "tortinsuarer": "` + tortinsuarer + `", "accidentdate": "` + accidentdate + `", "policynumber": "` + policynumber + `", "carrieremail": "` + carrieremail + `", "claimamount": "` + claimamount + `"}`

	err = stub.PutState(strconv.FormatInt(ctime, 10), []byte(str)) //store cert with user name as key
}

	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}

	return nil, nil
}

// ============================================================================================================================
// Init priliminaries - create a priliminary entry, store into chaincode state
// ============================================================================================================================
func (t *Subrogationcode) reg_priliminaries(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	ctime := time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))

	claimref := strings.ToLower(args[0])
	insuarer := strings.ToLower(args[1])
	adjustername := strings.ToLower(args[2])
	insuredname := strings.ToLower(args[3])
	referance := strings.ToLower(args[4])
	policylimits := strings.ToLower(args[5])

if insuarer != "" {
	//build the cert json string manually
	str := `{"claimref": "` + claimref + `", "insuarer": "` + insuarer + `", "adjustername": "` + adjustername + `", "insuredname": "` + insuredname + `", "referance": "` + referance + `", "policylimits": "` + policylimits + `"}`
	err = stub.PutState(strconv.FormatInt(ctime, 10), []byte(str)) //store cert with user name as key
}
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
