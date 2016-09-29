package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type GuavaChaincode struct {
}

var OwnerAccountMapkey = "_accountmapkey"

var accountcount int64 = 1
var transcount int64 = 1

type AccountNumber struct {
	Id int64 `json:"id"`
}

type Transfer struct {
	From        AccountNumber `json:"from"`        //account number who generated transfer
	To          AccountNumber `json:"to"`          //account number receiving transfer
	Dec_value   float64       `json:"dec_value"`   //amount to decrease in from account
	Inc_value   float64       `json:"inc_value"`   //amount to increase in to account
	Fx_rate     float64       `json:"fx_rate"`     //fx_rate for the transfer
	Message     string        `json:"message"`     //description of desired transfer
	Status      string        `json:"status"`      //current status of transfer <accept,reject,pending>
	T_Type      string        `json:"type"`        //type of fund transfer <internal,external>
	Transfer_id int64         `json:"transfer_id"` //unique identifier for transfer
}

// Transfers = make(map[String]Account[])

type Account struct {
	Owner            string        `json:"owner"`             //owner of the account registered with the membership service
	AccountID        AccountNumber `json:"account_id"`        //unique accountid
	Currency         string        `json:"currency"`          //currency representing the
	Country          string        `json:"country"`           //operational or savings acco
	Balance          float64       `json:"balance"`           //current account balance
	Type             string        `json:"type"`              //operational or savings acco
	IncomingTransfer []Transfer    `json:"incoming_transfer"` //array of incoming transfers
	OutgoingTransfer []Transfer    `json:"outgoing_transfer"` //array of outgoing transactions
}

//need to clear the map so have to encapsulate in a struct
//type OwnerAccountMap struct{
var OwnerAccount = make(map[string][]int64)

//}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(GuavaChaincode))
	if err != nil {
		fmt.Printf("Error starting Guava chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *GuavaChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//var Aval int
	//	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// what does this do?
	//	Aval, err = strconv.Atoi(args[0])
	//	if err != nil {
	//		return nil, errors.New("Expecting integer value for asset holding")
	//	}

	//clear the accounts to
	/*	var accounthashkey OwnerAccountMap
		jsonAsBytes, _ = json.Marshal(accounts)								//clear the open trade struct
		err = stub.PutState(AccountsHashkey, jsonAsBytes)
		if err != nil {
			return nil, err
		}*/

	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *GuavaChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *GuavaChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "create_account" { //create a new account

		return t.create_account(stub, args)
	} else if function == "create_transfer" { //create a new transfer

		return t.create_transfer(stub, args)
	} else if function == "increment_value" { //increment balance for an account

		return t.increment_value(stub, args)
	} else if function == "decrement_value" { //decrement balance for an account

		return t.decrement_value(stub, args)
	} else if function == "accept_transfer" { //accept a pending open transfer

		return t.accept_transfer(stub, args)
	} else if function == "reject_transfer" { //reject a pending open transfer

		return t.reject_transfer(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error invoke function not found

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *GuavaChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a account_num
// ============================================================================================================================
func (t *GuavaChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var account_num, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	account_num = args[0]
	valAsbytes, err := stub.GetState(account_num) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + account_num + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}

// ============================================================================================================================
// create_account - create new account expected arguments <account_owner, currency, country, acctype>
// ============================================================================================================================
func (t *GuavaChaincode) create_account(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var account_owner, currency, country, acctype string // Entities
	var account_number int64
	var initialbalance float64 = 0
	var err error
	fmt.Println("running write()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. Account owner and currency type")
	}

	account_owner = args[0]
	currency = args[1]
	country = args[2]
	acctype = args[3]

	account_number = accountcount
	accountcount = accountcount + 1

	//check if marble already exists

	accountAsBytes, err := stub.GetState(strconv.FormatInt(account_number, 10))
	found_acc := Account{}
	json.Unmarshal(accountAsBytes, &found_acc)
	if err == nil {
		return nil, errors.New("this Account arleady exists under" + found_acc.Owner)
	}

	incoming_t := make([]Transfer, 30)
	outgoing_t := make([]Transfer, 30)

	Acc_numstruct := AccountNumber{account_number}
	//bal_float := strconv.FormatFloat(initialbalance, 'f', -1, 64)

	new_Account := &Account{
		Owner:            account_owner,
		AccountID:        Acc_numstruct,
		Currency:         currency,
		Country:          country,
		Balance:          initialbalance,
		Type:             acctype,
		IncomingTransfer: incoming_t,
		OutgoingTransfer: outgoing_t}

	new_Account_m, _ := json.Marshal(new_Account)

	/*str := `{"owner": "` + account_owner + `", "account_id": "` + account_number +  `", "currency": "` + currency + `", "country": "` + country + `, "balance": "` + strconv.FormatFloat(initialbalance, 'f', 4 ,64) +  `", "type": "` + acctype + `", "country": "` + country + `", "incoming_transfer": "` + nil + `", "outgoing_transfer": "` + nil + `"}`*/

	err = stub.PutState(strconv.FormatInt(account_number, 10), []byte(new_Account_m)) //store the account
	if err != nil {
		return nil, err
	}

	//add the account number to the OwnerAccountMap
	OwnerAccount["account_owner"] = append(OwnerAccount["account_owner"], account_number)

	//var accounts AccountStruct
	jsonAsBytes, _ := json.Marshal(OwnerAccount)
	err = stub.PutState(OwnerAccountMapkey, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// create_transfer - create new account expected arguments <message, fx_rate, value, from_id, to_id, trans_type>
// ============================================================================================================================

func (t *GuavaChaincode) create_transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var trans_type, message, status string // Entities
	var from_id, to_id string
	var fx_rate string
	var value_inc, value_dec string
	var err error

	var from_id_int, to_id_int int64

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments.")
	}

	var trans_id = transcount
	transcount = transcount + 1

	message = args[0]
	fx_rate = args[1]
	value_inc = args[2]
	value_dec = args[3]
	from_id = args[4]
	to_id = args[5]
	trans_type = args[6]
	status = args[7]

	dec_float, err := strconv.ParseFloat(value_dec, 64)
	inc_float, err := strconv.ParseFloat(value_inc, 64)
	fx_rate_float, err := strconv.ParseFloat(fx_rate, 64)

	from_id_int, err = strconv.ParseInt(from_id, 10, 64)
	to_id_int, err = strconv.ParseInt(to_id, 10, 64)

	fromacc_struct := AccountNumber{from_id_int}
	toacc_struct := AccountNumber{to_id_int}

	//create transfer

	new_transfer := &Transfer{
		From:        fromacc_struct,
		To:          toacc_struct,
		Dec_value:   dec_float,
		Inc_value:   inc_float,
		Fx_rate:     fx_rate_float,
		Message:     message,
		Status:      status,
		T_Type:      trans_type,
		Transfer_id: trans_id}

	//new_Transfer_m, _ := json.Marshal(new_transfer)

	//find the account entry for from_id
	fromAccountAsBytes, err := stub.GetState(from_id)
	if err != nil {
		return nil, errors.New("Could not find the account that is sending funds " + from_id)
	}

	from_acc := Account{}
	json.Unmarshal(fromAccountAsBytes, &from_acc)

	//check that account has enough funds, decrement if internal
	//TODO: handle internal/external
	if from_acc.Balance < new_transfer.Dec_value {
		return nil, errors.New("from account does not have enough funds " + from_id)
	} else {

		from_acc.Balance = from_acc.Balance - new_transfer.Dec_value
	}

	//add transfer to outgoing transfer
	from_acc.OutgoingTransfer = append(from_acc.OutgoingTransfer, *new_transfer)

	//find account entry for to_id
	toAccountAsBytes, err := stub.GetState(to_id)
	if err != nil {
		return nil, errors.New("Could not find this account that is receiving funds")
	}
	to_acc := Account{}
	json.Unmarshal(toAccountAsBytes, &to_acc)

	//add transaction to incoming transfer
	to_acc.IncomingTransfer = append(to_acc.IncomingTransfer, *new_transfer)

	//increment this value TODO:handle internal,external transfer
	to_acc.Balance = to_acc.Balance - new_transfer.Inc_value

	//update the account states

	toaccAsBytes, _ := json.Marshal(to_acc)
	err = stub.PutState(to_id, toaccAsBytes)
	if err != nil {
		return nil, err
	}

	fromaccAsBytes, _ := json.Marshal(from_acc)
	err = stub.PutState(from_id, fromaccAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// ============================================================================================================================
// increment_value - increase balance in account <account_id, value>
// ============================================================================================================================

func (t *GuavaChaincode) increment_value(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	return nil, nil
}

// ============================================================================================================================
// decrement_value - decrease balance in account <account_id, value>
// ============================================================================================================================

func (t *GuavaChaincode) decrement_value(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil

}

// ============================================================================================================================
// accept_transfer - accept the transfer from the incoming queue <account_id, transaction_id>
// ============================================================================================================================

func (t *GuavaChaincode) accept_transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	// transfer the value
	return nil, nil

}

// ============================================================================================================================
// reject_transfer - reject the transfer from the incoming queue <account_id, transaction_id>
// ============================================================================================================================

func (t *GuavaChaincode) reject_transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	return nil, nil

}
