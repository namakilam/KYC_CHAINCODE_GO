package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

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

type SimpleChainCode struct {
}

func main()  {
	err := shim.Start(new(SimpleChainCode))
	if err != nil {
		fmt.Printf("Error starting simple chaincode due to :%s", err)
	}
}

func (t *SimpleChainCode) Init(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {
	fmt.Print("Initializing ChainCode.....")

	return []byte("INITIALIZATION O.K"), nil
}

func (t *SimpleChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "insert":
		return t.insertDataIntoLedger(stub, args)
	case "update":
	    	return t.updateDataIntoLedger(stub, args)
	}

	return nil, errors.New("Unknown Function Invocation")
}

func (t *SimpleChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	switch function {
	case "retrieve":
		return t.readDataFromLedger(stub, args)
	}

	return nil, errors.New("Unknown Function Invocation")
}

func (t *SimpleChainCode) updateDataIntoLedger(stub shim.ChaincodeStubInterface, args[] string) ([]byte, error) {
    if len(args) != 2 {
	return nil, errors.New("Incorrect Number of Arguments. Required : 2")
    }

    key := args[0]
    customerInfo, err := stub.GetState(key)

    if err != nil {
	return nil, errors.New("Entry for given key not found!. Please insert into the ledger first.")
    }
    
    var customer, customer1 Customer
    err = json.Unmarshal(customerInfo, &customer)
    err = json.Unmarshal([]byte(args[1]), &customer1) 
   
    if err != nil {
	return nil, errors.New("Unable to parse Customer String. Please ensure a valid JSON.")
    }

    if customer.Aadhar == customer1.Aadhar && customer.PAN == customer1.PAN {
	value, err := json.Marshal(customer1)
	if err != nil {
	    return nil, err
	}

	err = stub.PutState(key, value)
	if err != nil {
	    return nil, err
	}
	return []byte("Update Successful"), nil
    } else {
	return nil, errors.New("Cannot Update Immutable Fields (AADHAR NUMBER, PAN)")
    }
}


func (t *SimpleChainCode) readDataFromLedger(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
	    return nil, errors.New("Incorrect Number of Arguments. Required : 1")
	}

	customerInfo, err := stub.GetState(args[0])
	if err != nil {
	    return nil, err
	}

	return customerInfo, nil
}

func (t * SimpleChainCode) insertDataIntoLedger(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    if len(args) != 1 {
	return nil, errors.New("Incorrect Number of Arguments. Required : 1") 
    }

    var customer Customer
    err := json.Unmarshal([]byte(args[0]), &customer)
    if err != nil {
	return nil, err
    }

    key := customer.Aadhar
    value, err := json.Marshal(customer)
    if err != nil {
	return nil, err
    }
    err = stub.PutState(key, value) 
    if err != nil {
	return nil, err
    }

    fmt.Println("Ledger state successfully updated")

    return []byte("Insert Success"), nil
}
