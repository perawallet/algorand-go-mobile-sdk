package sdk

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
)

func MakeKeyRegTxn(
	account string,
	note []byte,
	voteKey, selectionKey string,
	voteFirst, voteLast, voteKeyDilution *Uint64,
	suggestedParams *SuggestedParams,
) (txn []byte, err error) {
	paramsConverted, err := convertSuggestedParams(suggestedParams)
	if err != nil {
		return nil, fmt.Errorf("failed to convert suggested params: %v", err)
	}

	voteFirstDecoded, err := voteFirst.Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to extract voteFirst: %v", err)
	}

	voteLastDecoded, err := voteLast.Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to extract voteLast: %v", err)
	}

	voteKeyDilutionDecoded, err := voteKeyDilution.Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to extract voteKeyDilution: %v", err)
	}

	txnObj, err := transaction.MakeKeyRegTxn(
		account,
		note,
		paramsConverted,
		voteKey,
		selectionKey,
		voteFirstDecoded,
		voteLastDecoded,
		voteKeyDilutionDecoded,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to construct key reg txn: %v", err)
	}

	txn = msgpack.Encode(&txnObj)

	return txn, nil
}

func MakeKeyRegTxnWithStateProofKey(
	account string,
	note []byte,
	params *SuggestedParams,
	voteKey, selectionKey, stateProofPK string,
	voteFirst, voteLast, voteKeyDilution *Uint64,
	nonpart bool,
) (txn []byte, err error) {
	paramsConverted, err := convertSuggestedParams(params)
	if err != nil {
		return nil, fmt.Errorf("failed to convert suggested params: %v", err)
	}

	voteFirstDecoded, err := voteFirst.Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to extract voteFirst: %v", err)
	}

	voteLastDecoded, err := voteLast.Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to extract voteLast: %v", err)
	}

	voteKeyDilutionDecoded, err := voteKeyDilution.Extract()
	if err != nil {
		return nil, fmt.Errorf("failed to extract voteKeyDilution: %v", err)
	}

	txnObj, err := transaction.MakeKeyRegTxnWithStateProofKey(
		account,
		note,
		paramsConverted,
		voteKey,
		selectionKey,
		stateProofPK,
		voteFirstDecoded,
		voteLastDecoded,
		voteKeyDilutionDecoded,
		nonpart,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to construct key reg txn: %v", err)
	}

	txn = msgpack.Encode(&txnObj)

	return txn, nil
}
