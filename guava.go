package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type GuavaChaincode struct {
}

var GuavaMapkey = "_guavamapkey"
var UserMapkey = "_usermapkey"

var accountcount int64 = 1
var transcount int64 = 1

var guavacount int64 = 1

type User struct {
	Username string `json:"username"`
	Owner    bool   `json:"owner"`
	Create   bool   `json:"create"`
	Approve  bool   `json:"approve"`
	Read     bool   `json:"read"`
}

type Transfer struct {
	From        int64   `json:"from"`        //account number who generated transfer
	To          int64   `json:"to"`          //account number receiving transfer
	Dec_value   float64 `json:"dec_value"`   //amount to decrease in from account
	Inc_value   float64 `json:"inc_value"`   //amount to increase in to account
	Fx_rate     float64 `json:"fx_rate"`     //fx_rate for the transfer
	Message     string  `json:"message"`     //description of desired transfer
	Status      string  `json:"status"`      //current status of transfer <accept,reject,pending>
	T_Type      string  `json:"type"`        //type of fund transfer <internal,external>
	Creator     string  `json:"creator"`     //the username of the user who created the transactions
	Approver    string  `json:"approver"`    //the username of the user who approved the payment
	Time        string  `json:"time"`        // time the transfer was created
	Transfer_id int64   `json:"transfer_id"` //unique identifier for transfer
}

// Transfers = make(map[String]Account[])

type Account struct {
	AccountName      string     `json:"name"`              // the name of the account
	AccountID        int64      `json:"id"`                //unique accountid
	Currency         string     `json:"currency"`          //currency representing the
	Country          string     `json:"country"`           //operational or savings acco
	Balance          float64    `json:"balance"`           //current account balance
	Type             string     `json:"type"`              //operational or savings acco
	IncomingTransfer []Transfer `json:"incoming_transfer"` //array of incoming transfers
	OutgoingTransfer []Transfer `json:"outgoing_transfer"` //array of outgoing transactions
}

var GuavaMap = make(map[string][]int64)
var UserMap = make(map[string][]User)

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

	//this is a test entry into the worldstate
	err := stub.PutState("hello", []byte(args[0]))
	if err != nil {
		return nil, err
	}

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
	} else if function == "create_user" {

		return t.create_user(stub, args)
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
// create_account - create new account expected arguments <account_name, guava_id, currency, country, acctype, username>
// ============================================================================================================================
func (t *GuavaChaincode) create_account(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var account_name, currency, country, acctype, guava_id string // Entities
	var account_number int64
	//this is just for testing
	var initialbalance float64 = 100
	var err error

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5 arguments <account_name, guava_id, currency, country, acctype>")
	}

	account_name = args[0]
	guava_id = args[1]
	//guava_id, err = strconv.ParseInt(args[1], 10, 64)
	currency = args[2]
	country = args[3]
	acctype = args[4]

	account_number = accountcount
	accountcount = accountcount + 1

	if strings.Compare(guava_id, "-1") == 0 {

		guava_id = strconv.FormatInt(guavacount, 10)
		guavacount = guavacount + 1
	}

	incoming_t := make([]Transfer, 0)
	outgoing_t := make([]Transfer, 0)

	new_Account := &Account{
		AccountName:      account_name,
		AccountID:        account_number,
		Currency:         currency,
		Country:          country,
		Balance:          initialbalance,
		Type:             acctype,
		IncomingTransfer: incoming_t,
		OutgoingTransfer: outgoing_t}

	new_Account_m, _ := json.Marshal(new_Account)

	new_Account_string := string(new_Account_m)

	err = stub.PutState(strconv.FormatInt(account_number, 10), []byte(new_Account_string)) //store the account
	if err != nil {
		return nil, err
	}

	//add the account number to the OwnerAccountMap
	GuavaMap[guava_id] = append(GuavaMap[guava_id], account_number)

	//var accounts AccountStruct
	jsonAsBytes, _ := json.Marshal(GuavaMap)
	ownwermap_string := string(jsonAsBytes)

	err = stub.PutState(GuavaMapkey, []byte(ownwermap_string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// create_transfer - create new account expected arguments <message, fx_rate, value, from_id, to_id, trans_type>
// ============================================================================================================================

func (t *GuavaChaincode) create_transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var trans_type, message, status, time, creator, approver string // Entities
	var from_id, to_id string
	var fx_rate string
	var value_inc, value_dec string
	var err error

	var from_id_int, to_id_int int64

	if len(args) != 9 {
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
	time = args[7]
	creator = args[8]

	if strings.Compare(trans_type, "internal") == 0 {
		approver = args[8]
		status = "approved"
	} else {
		approver = "pending"
	}

	dec_float, err := strconv.ParseFloat(value_dec, 64)
	inc_float, err := strconv.ParseFloat(value_inc, 64)
	fx_rate_float, err := strconv.ParseFloat(fx_rate, 64)

	from_id_int, err = strconv.ParseInt(from_id, 10, 64)
	to_id_int, err = strconv.ParseInt(to_id, 10, 64)

	//create transfer

	new_transfer := &Transfer{
		From:        from_id_int,
		To:          to_id_int,
		Dec_value:   dec_float,
		Inc_value:   inc_float,
		Fx_rate:     fx_rate_float,
		Message:     message,
		Status:      status,
		T_Type:      trans_type,
		Creator:     creator,
		Approver:    approver,
		Time:        time,
		Transfer_id: trans_id}

	//find the account entry for from_id
	fromAccountAsBytes, err := stub.GetState(from_id)
	if err != nil {
		return nil, errors.New("Could not find the account that is sending funds " + from_id)
	}

	from_acc := Account{}
	json.Unmarshal(fromAccountAsBytes, &from_acc)

	//check that account has enough funds, decrement if internal otherwise set status as pending

	if from_acc.Balance < new_transfer.Dec_value {
		return nil, errors.New("from account does not have enough funds " + from_id)
	} else if strings.Compare(new_transfer.T_Type, "internal") == 0 {

		from_acc.Balance = from_acc.Balance - new_transfer.Dec_value
	} else {
		new_transfer.Status = "pending"
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

	//increment this value
	if strings.Compare(new_transfer.T_Type, "internal") == 0 {
		to_acc.Balance = to_acc.Balance + new_transfer.Inc_value
		to_acc.IncomingTransfer = append(to_acc.IncomingTransfer, *new_transfer)

	}
	//update the account states

	toaccAsBytes, _ := json.Marshal(to_acc)

	to_acc_string := string(toaccAsBytes)
	err = stub.PutState(to_id, []byte(to_acc_string))
	if err != nil {
		return nil, err
	}

	fromaccAsBytes, _ := json.Marshal(from_acc)
	from_acc_string := string(fromaccAsBytes)
	err = stub.PutState(from_id, []byte(from_acc_string))
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// ============================================================================================================================
// increment_value - increase balance in account <account_id, value>
// ============================================================================================================================

func (t *GuavaChaincode) increment_value(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var account_id string

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments.")
	}

	account_id = args[0]
	inc_val, err := strconv.ParseFloat(args[1], 64)

	incAccountAsBytes, err := stub.GetState(account_id)
	if err != nil {
		return nil, errors.New("Could not find the account to increment" + account_id)
	}

	inc_acc := Account{}
	json.Unmarshal(incAccountAsBytes, &inc_acc)

	inc_acc.Balance = inc_acc.Balance + inc_val
	newAccountAsBytes, _ := json.Marshal(inc_acc)
	newacc_string := string(newAccountAsBytes)
	err = stub.PutState(account_id, []byte(newacc_string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// decrement_value - decrease balance in account <account_id, value>
// ============================================================================================================================

func (t *GuavaChaincode) decrement_value(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var account_id string

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments.")
	}

	account_id = args[0]
	dec_val, err := strconv.ParseFloat(args[1], 64)

	decAccountAsBytes, err := stub.GetState(account_id)
	if err != nil {
		return nil, errors.New("Could not find the account to increment" + account_id)
	}

	dec_acc := Account{}
	json.Unmarshal(decAccountAsBytes, &dec_acc)

	dec_acc.Balance = dec_acc.Balance - dec_val
	newAccountAsBytes, _ := json.Marshal(dec_acc)
	newacc_string := string(newAccountAsBytes)
	err = stub.PutState(account_id, []byte(newacc_string))
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// ============================================================================================================================
// accept_transfer - accept the transfer from the incoming queue <account_id, transaction_id>
// ============================================================================================================================

func (t *GuavaChaincode) accept_transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var receiving_id, sending_id, transfer_id, approver string

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments.")
	}

	receiving_id = args[0]
	sending_id = args[1]
	transfer_id = args[2]
	dec_value, err := strconv.ParseFloat(args[3], 64)
	inc_value, err := strconv.ParseFloat(args[4], 64)
	approver = args[5]
	//get the account from the passed in receiving_id (should be account who accepted)

	//rec_id_int, err := strconv.ParseInt(receiving_id, 10, 64)
	//sen_id_int, err := strconv.ParseInt(sending_id, 10, 64)
	tran_id_int, err := strconv.ParseInt(transfer_id, 10, 64)

	recAccountAsBytes, err := stub.GetState(receiving_id)
	if err != nil {
		return nil, errors.New("Could not find the account that is receiving funds " + receiving_id)
	}

	receiving_acc := Account{}
	json.Unmarshal(recAccountAsBytes, &receiving_acc)

	// find the account that is sending the transaction from the transaction
	sendAccountAsBytes, err := stub.GetState(sending_id)
	if err != nil {
		return nil, errors.New("Could not find the account that is sending funds " + sending_id)
	}

	sending_acc := Account{}
	json.Unmarshal(sendAccountAsBytes, &sending_acc)

	// decrement sending account
	// increcment receiving account

	if sending_acc.Balance < dec_value {
		return nil, errors.New("sending account does not have enough funds " + sending_id)
	} else {

		sending_acc.Balance = sending_acc.Balance - dec_value
		receiving_acc.Balance = receiving_acc.Balance + inc_value
	}

	trans_list_o := sending_acc.OutgoingTransfer

	for i := 0; i < len(trans_list_o); i++ {
		transl := &trans_list_o[i]
		if transl.Transfer_id == tran_id_int {
			transl.Status = "approved"
			transl.Approver = approver
			receiving_acc.IncomingTransfer = append(receiving_acc.IncomingTransfer, *transl)
		}
	}

	//update the account states

	newrecAccountasBytes, _ := json.Marshal(receiving_acc)
	rec_acc_string := string(newrecAccountasBytes)
	err = stub.PutState(receiving_id, []byte(rec_acc_string))
	if err != nil {
		return nil, err
	}

	newsendAccountAsBytes, _ := json.Marshal(sending_acc)
	send_acc_string := string(newsendAccountAsBytes)
	err = stub.PutState(sending_id, []byte(send_acc_string))
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// ============================================================================================================================
// reject_transfer - reject the transfer from the incoming queue <account_id, transaction_id>
// ============================================================================================================================

func (t *GuavaChaincode) reject_transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var sending_id, transfer_id, approver string

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments.")
	}

	sending_id = args[0]
	transfer_id = args[1]
	approver = args[2]
	//get the account from the passed in receiving_id (should be account who accepted)

	//	sen_id_int, err := strconv.ParseInt(sending_id, 10, 64)
	tran_id_int, err := strconv.ParseInt(transfer_id, 10, 64)

	// find the account that is sending the transfer
	sendAccountAsBytes, err := stub.GetState(sending_id)
	if err != nil {
		return nil, errors.New("Could not find the account that is sending funds " + sending_id)
	}

	sending_acc := Account{}
	json.Unmarshal(sendAccountAsBytes, &sending_acc)

	trans_list_o := sending_acc.OutgoingTransfer

	for i := 0; i < len(trans_list_o); i++ {
		transl := &trans_list_o[i]
		if transl.Transfer_id == tran_id_int {
			transl.Status = "rejected"
			transl.Approver = approver
		}
	}

	newsendAccountAsBytes, _ := json.Marshal(sending_acc)
	send_acc_string := string(newsendAccountAsBytes)
	err = stub.PutState(sending_id, []byte(send_acc_string))
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// ============================================================================================================================
// create_user - create a new user with the specific access rights <username, owner, create, approve, read>
// ============================================================================================================================

func (t *GuavaChaincode) create_user(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var username string

	var owner, create, approve, read bool
	var guava_id string

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments.")
	}

	username = args[0]
	owner, err := strconv.ParseBool(args[1])
	create, err = strconv.ParseBool(args[2])
	approve, err = strconv.ParseBool(args[3])
	read, err = strconv.ParseBool(args[4])
	//guava_id, err = strconv.ParseInt(args[5], 10, 64)

	guava_id = args[5]
	// create User struct

	guava_id_int, err := strconv.ParseInt(guava_id, 10, 64)

	new_user := &User{
		Username: username,
		Owner:    owner,
		Create:   create,
		Approve:  approve,
		Read:     read}

	//find guava_id in user map
	if guava_id_int <= guavacount-1 {
		//add the account number to the OwnerAccountMap
		// add user struct to the array

		UserMap[guava_id] = append(UserMap[guava_id], *new_user)

		//var accounts AccountStruct
		jsonAsBytes, _ := json.Marshal(UserMap)
		usermap_string := string(jsonAsBytes)

		// add the new map to the world state

		err = stub.PutState(UserMapkey, []byte(usermap_string))
		if err != nil {
			return nil, err
		}

	} else {
		return nil, errors.New("Guava id does not exist")

	}

	return nil, nil

}
