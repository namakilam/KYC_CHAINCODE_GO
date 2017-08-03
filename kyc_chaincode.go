opyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Address struct {
	Address_Line string `json:"address_line"`
	City string `json:"city"`
}

type Customer struct {
	Name string `json:"name"`
	Gender string `json:"gender"`
	DOB string `json:"dob"`
	Aadhar string `json:"aadhar_no"`
	Address Address `json:"address"`
	PAN string `json:"pan_no"`
	Cibil_Score int32 `json:"cibil_score"`
	Marital_Status string `json:"marital_status"`
	Education map[string]string `json:"education"`
	Employement map[string]string `json:"employement"`
	Health map[string]string `json:"health"`
	Possesions map[string]string `json:"possesions"`
}

// Init initializes the chaincode state
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### example_cc Init ###########")
	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)

}

// Invoke makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### example_cc Invoke ###########")
	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if args[0] == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}

	if args[0] == "query" {
		// queries an entity state
		return t.query(stub, args)
	}
	if args[0] == "move" {
		// Deletes an entity from its state
		return t.move(stub, args)
	}
	if args[0] == "insert" {
		// Deletes an entity from its state
		return t.insertDataIntoLedger(stub, args)
	}
	if args[0] == "update" {
		// Deletes an entity from its state
		return t.updateDataIntoLedger(stub, args)
	}
	if args[0] == "retrieve" {
		// Deletes an entity from its state
		return t.readDataFromLedger(stub, args)
	}
	if args[0] == "history" {
		// Deletes an entity from its state
		return t.readHistoryFromLedger(stub, args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}

func (t *SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var X int          // Transaction value
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	A = args[1]
	B = args[2]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[1]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var A string // Entities
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[1]

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

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}


func (t *SimpleChaincode) updateDataIntoLedger(stub shim.ChaincodeStubInterface, args[] string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect Number of Arguments. Required : 2")
	}
	args = args[1:]
	key := args[0]
	customerInfo, err := stub.GetState(key)

	if err != nil {
		return shim.Error("Entry for given key not found!. Please insert into the ledger first.")
	}

	var customer, customer1 Customer
	err = json.Unmarshal(customerInfo, &customer)
	err = json.Unmarshal([]byte(args[1]), &customer1)

	if err != nil {
		return shim.Error("Unable to parse Customer String. Please ensure a valid JSON.")
	}

	if customer.Aadhar == customer1.Aadhar && customer.PAN == customer1.PAN {
		value, err := json.Marshal(customer1)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.PutState(key, value)
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte("Update Successful"))
	} else {
		return shim.Error("Cannot Update Immutable Fields (AADHAR NUMBER, PAN)")
	}
}

func (t *SimpleChaincode) readHistoryFromLedger(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect Number of Arguments. Required : 2")
	}
	args = args[1:]

	key := args[0]
	historyItr, err := stub.GetHistoryForKey(key)

	if err != nil {
		return shim.Error(err.Error())
	}

	var history []string

	for historyItr.HasNext() {
		alters, err := historyItr.Next()
		if err != nil {
			fmt.Println(err)
		} else {
			history =  append(history, string(alters.Value))
		}
	}

	val, err := json.Marshal(history)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(val)
}

func (t *SimpleChaincode) readDataFromLedger(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect Number of Arguments. Required : 1")
	}
	args = args[1:]

	customerInfo, err := stub.GetState(args[0])
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(customerInfo)
}

func (t * SimpleChaincode) insertDataIntoLedger(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect Number of Arguments. Required : 1")
	}
	args = args[1:]

	var customer Customer
	err := json.Unmarshal([]byte(args[0]), &customer)
	if err != nil {
		return shim.Error(err.Error())
	}

	key := customer.Aadhar
	value, err := json.Marshal(customer)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(key, value)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Ledger state successfully updated")

	return shim.Success([]byte("Insert Success"))
}



func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

