package sdk

import (
	"encoding/base64"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
)

// MakeAndSignKeyRegTxnWithFlatFee creates and signs a key registration transaction
// using a flat fee. It returns the signed transaction bytes and/or an error.
//
// Parameters:
//   - sender: the account address (checksummed, human-readable) for which we register participation.
//   - fee: the flat fee.
//   - firstRound, lastRound: first and last rounds in which the transaction is valid.
//   - note: optional note as []byte.
//   - genesisID: the network's genesis ID.
//   - genesisHashB64: base64-encoded genesis hash.
//   - voteKeyB64: base64-encoded string corresponding to the root participation public key.
//   - selectionKeyB64: base64-encoded string corresponding to the VRF public key.
//   - voteFirst, voteLast: the first/last round for which this participation key is valid.
//   - voteKeyDilution: the dilution for the 2-level participation key.
//   - sk: the sender's secret key (ed25519.PrivateKey) for signing.
//
// Returns:
//   - signedTxnBytes: The bytes of the signed transaction.
//   - err: Error, if any occurred.
func MakeKeyRegTxn(
    sender string,
    note []byte,
    voteKeyB64, selectionKeyB64 string,
    voteFirst, voteLast, voteKeyDilution *Uint64,
	suggestedParams *SuggestedParams,
) (txn []byte, err error) {
	internalTxnParams, err := convertSuggestedParams(suggestedParams)
	
	decodedVoteFirst, _ := voteFirst.Extract()
	decodedVoteLast, _ := voteLast.Extract()
	decodedVoteKeyDilution, _ := voteKeyDilution.Extract()



    voteKey, err := base64.StdEncoding.DecodeString(voteKeyB64)

	selectionKey, err := base64.StdEncoding.DecodeString(selectionKeyB64)


    txnObj, err := transaction.MakeKeyRegTxn(
        sender,
        note,
		internalTxnParams,
        string(voteKey),
        string(selectionKey),
        decodedVoteFirst,
        decodedVoteLast,
        decodedVoteKeyDilution,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to construct key reg txn: %v", err)
    }

	txn = msgpack.Encode(&txnObj)

    return txn, nil
}