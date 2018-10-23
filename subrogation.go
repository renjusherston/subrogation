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
	Claimref             string `json:"claimref"`
	Insuredname          string `json:"insuredname"`
	Policynumber         string `json:"policynumber"`
	Claimnumber          string `json:"claimnumber"`
	Tortcarriername      string `json:"tortcarriername"`
	Tortcarrieraddress   string `json:"tortcarrieraddress"`
	Tortcarrieremail     string `json:"tortcarrieremail"`
}

type Priliminary struct {
	Claimref1             string `json:"claimref1"`
	Insuredname1         string `json:"insuredname1"`
	Policynumber1         string `json:"policynumber1"`
	Claimnumber1          string `json:"claimnumber1"`
	Tortcarriername1      string `json:"tortcarriername1"`
	Tortcarrieraddress1   string `json:"tortcarrieraddress1"`
	Tortcarrieremail1     string `json:"tortcarrieremail1"`
	Dateofaccident1       string `json:"dateofaccident1"`
	Tortdefendentname1    string `json:"tortdefendentname1"`
	Accidentstreet1       string `json:"accidentstreet1"`
	Accidenttown1         string `json:"accidenttown1"`
	Accidentcounty1       string `json:"accidentcounty1"`
	Accidentstate1        string `json:"accidentstate1"`
	Propertydamageamount1 string `json:"propertydamageamount1"`
	Claimamount1          string `json:"claimamount1"`
	Attorneyname1         string `json:"attorneyname1"`
	Attorneyid1           string `json:"attorneyid1"`
	Releaserep1           string `json:"releaserep1"`
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

		keys = append(keys, klaim)
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

		if strings.ToLower(klaim.Claimref) == klaimref && klaim.Insuredname != "" {
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

		if priliminary.Insuredname1 != "" {
			if strings.ToLower(priliminary.Claimref1) == klaimref {
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

	claimref := args[0]
	insuredname := args[1]
	policynumber := args[2]
	claimnumber := args[3]
	tortcarriername := args[4]
	tortcarrieraddress := args[5]
	tortcarrieremail := args[6]

	if claimref != "" {
		str := `{"claimref": "` + claimref + `", "insuredname": "` + insuredname + `", "policynumber": "` + policynumber + `", "claimnumber": "` + claimnumber + `", "tortcarriername": "` + tortcarriername + `", "tortcarrieraddress": "` + tortcarrieraddress + `", "tortcarrieremail": "` + tortcarrieremail + `"}`

		err = stub.PutState(strconv.FormatInt(ctime, 10), []byte(str))  //store cert with user name as key
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
	claimref := args[0]
	insuredname := args[1]
	policynumber := args[2]
	claimnumber := args[3]
	tortcarriername := args[4]
	tortcarrieraddress := args[5]
	tortcarrieremail := args[6]
	dateofaccident := args[7]
	tortdefendentname := args[8]
	accidentstreet := args[9]
	accidenttown := args[10]
	accidentcounty := args[11]
	accidentstate := args[12]
	propertydamageamount := args[13]
	claimamount := args[14]
	attorneyname := args[15]
	attorneyid := args[16]
	releaserep := args[17]

	if insuredname != "" {
		//build the cert json string manually
		str := `{"claimref1": "` + claimref + `", "insuredname1": "` + insuredname + `", "policynumber1": "` + policynumber + `", "claimnumber1": "` + claimnumber + `", "tortcarriername1": "` + tortcarriername + `", "tortcarrieraddress1": "` + tortcarrieraddress + `", "tortcarrieremail1": "` + tortcarrieremail + `", "dateofaccident1": "` + dateofaccident + `", , "tortdefendentname1": "` + tortdefendentname + `", , "accidentstreet1": "` + accidentstreet + `", , "accidenttown1": "` + accidenttown + `", "accidentcounty1": "` + accidentcounty + `", "accidentstate1": "` + accidentstate + `", "propertydamageamount1": "` + propertydamageamount + `", "claimamount1": "` + claimamount + `", "attorneyname1": "` + attorneyname + `", "attorneyid1": "` + attorneyid + `", "releaserep1": "` + releaserep + `"}`

		fmt.Printf("String: %s", str)

		err = stub.PutState(strconv.FormatInt(ctime, 10), []byte(str)) //store cert with user name as key
	}
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}

	return nil, nil
}

// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
	return time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}
