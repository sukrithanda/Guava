# Guava

Guava chaincode written in GO

// If you dont have a guava id you will have to pass guava_id as -1 and a new guava_id will be created

create_account - create new account expected arguments <account_name, guava_id, currency, country, acctype(OPR, SAVINGS)>



create_transfer - create new account expected arguments <message, fx_rate, value_inc, value_dec, from_id, to_id, tans_type(internal, payment), time, creator>

increment_value - increase balance in account <account_id, value>

decrement_value - decrease balance in account <account_id, value>

accept_transfer - accept the transfer from the outgoing array<to_id, from_id, transfer_id, dec_value, inc_value, approver>

reject_transfer - reject the transfer int the outgoing array <from_id, trans_id, approver>

create_user - create a new user with the specific access rights and add it to the User map <username, owner, create, approve, read>
