package sdk

import (
	"github.com/algorand/go-algorand-sdk/v2/encoding/json"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// TransactionMsgpackToJson converts a msgpack-encoded Transaction to a
// json-encoded Transaction
func TransactionMsgpackToJson(msgpTxn []byte) (jsonTxn string, err error) {
	var txn types.Transaction
	err = msgpack.Decode(msgpTxn, &txn)
	if err == nil {
		jsonTxn = string(json.Encode(txn))
	}
	return
}

// TransactionJsonToMsgpack converts a json-encoded Transaction to a
// msgpack-encoded Transaction
func TransactionJsonToMsgpack(jsonTxn string) (msgpackTxn []byte, err error) {
	var txn types.Transaction
	err = json.Decode([]byte(jsonTxn), &txn)
	if err == nil {
		msgpackTxn = msgpack.Encode(txn)
	}
	return
}
