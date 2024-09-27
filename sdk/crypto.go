package sdk

import (
	"fmt"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

func GenerateSK() []byte {
	account := crypto.GenerateAccount()
	return account.PrivateKey
}

func GenerateAddressFromSK(sk []byte) (string, error) {
	addr, err := crypto.GenerateAddressFromSK(sk)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

func GenerateAddressFromPublicKey(pk []byte) (string, error) {
	var a types.Address
	n := copy(a[:], pk)
	if n != ed25519.PublicKeySize {
		return "", fmt.Errorf("given public key has the wrong size, expected %d, got %d", ed25519.PublicKeySize, n)
	}
	return a.String(), nil
}

// SignTransaction accepts a private key and a transaction, and returns the
// bytes of a signed txn.
func SignTransaction(sk []byte, encodedTx []byte) (stxBytes []byte, err error) {
	if len(sk) != ed25519.PrivateKeySize {
		err = fmt.Errorf("Incorrect privateKey length expected %d, got %d", ed25519.PrivateKeySize, len(sk))
		return
	}

	var tx types.Transaction
	err = msgpack.Decode(encodedTx, &tx)
	if err != nil {
		return
	}

	_, stxBytes, err = crypto.SignTransaction(sk, tx)
	return
}

// AttachSignature accepts a signature and a transaction, and returns the bytes of a the signed transaction
func AttachSignature(signature, encodedTx []byte) (stxBytes []byte, err error) {
	if len(signature) != ed25519.SignatureSize {
		err = fmt.Errorf("incorrect signature length expected %d, got %d", ed25519.SignatureSize, len(signature))
		return
	}

	// Copy signature into a Signature, and check that it's the expected length
	var s types.Signature
	n := copy(s[:], signature)
	if n != len(s) {
		err = errInvalidSignatureReturned
		return
	}

	var tx types.Transaction
	err = msgpack.Decode(encodedTx, &tx)

	if err != nil {
		return nil, err
	}

	// Construct the SignedTxn
	stx := types.SignedTxn{
		Sig: s,
		Txn: tx,
	}

	// Encode the SignedTxn
	stxBytes = msgpack.Encode(stx)
	return
}

// AttachSignatureWithSigner accepts a signature, a transaction, and a signer address and returns the bytes of a the signed transaction
func AttachSignatureWithSigner(signature, encodedTx []byte, signer string) (stxBytes []byte, err error) {
	if len(signature) != ed25519.SignatureSize {
		err = fmt.Errorf("incorrect signature length expected %d, got %d", ed25519.SignatureSize, len(signature))
		return
	}

	// Copy signature into a Signature, and check that it's the expected length
	var s types.Signature
	n := copy(s[:], signature)
	if n != len(s) {
		err = errInvalidSignatureReturned
		return
	}

	var tx types.Transaction
	err = msgpack.Decode(encodedTx, &tx)

	if err != nil {
		return nil, err
	}

	signerAddr, err := types.DecodeAddress(signer)
	if err != nil {
		return nil, err
	}

	// Construct the SignedTxn
	stx := types.SignedTxn{
		Sig: s,
		Txn: tx,
	}
	if tx.Sender != signerAddr {
		stx.AuthAddr = signerAddr
	}

	// Encode the SignedTxn
	stxBytes = msgpack.Encode(stx)
	return
}

func SignBid(sk []byte, encodedBid []byte) (sBid []byte, err error) {
	if len(sk) != ed25519.PrivateKeySize {
		err = fmt.Errorf("Incorrect privateKey length expected %d, got %d", ed25519.PrivateKeySize, len(sk))
		return
	}

	var bid types.Bid
	err = msgpack.Decode(encodedBid, &bid)
	if err != nil {
		return
	}

	sBid, err = crypto.SignBid(sk, bid)
	return
}

// GetTxID takes an encoded txn and return the txid as string
func GetTxID(encodedTxn []byte) string {
	var tx types.Transaction
	err := msgpack.Decode(encodedTxn, &tx)
	if err != nil {
		panic("Could not decode transaction")
	}

	return crypto.TransactionIDString(tx)
}

// GetProgramAddress returns the escrow account address for a program, derived from its TEAL bytecode.
func AddressFromProgram(program []byte) string {
	return crypto.AddressFromProgram(program).String()
}

// AssignGroupID computes and return list of encoded transactions with Group field set.
func AssignGroupID(txns *BytesArray) (assignedTxns *BytesArray, err error) {
	txgroup := make([]types.Transaction, txns.Length())
	for i, encodedTxn := range txns.Extract() {
		err = msgpack.Decode(encodedTxn, &txgroup[i])
		if err != nil {
			err = fmt.Errorf("Could not decode transaction at index %d: %v", i, err)
			return
		}
	}

	txgroup, err = transaction.AssignGroupID(txgroup, "")
	if err == nil {
		assignedTxns = &BytesArray{
			values: make([][]byte, len(txgroup)),
		}

		for i := range txgroup {
			assignedTxns.values[i] = msgpack.Encode(&txgroup[i])
		}
	}

	return
}

func decodeTxns(txns *BytesArray) (decoded []types.Transaction, err error) {
	decoded = make([]types.Transaction, txns.Length())
	for i, encodedTxn := range txns.Extract() {
		err = msgpack.Decode(encodedTxn, &decoded[i])
		if err != nil {
			err = fmt.Errorf("Could not decode transaction at index %d: %v", i, err)
			return
		}
	}

	return
}

func verifyTxnsGroupID(txgroup []types.Transaction) (valid bool, err error) {
	if len(txgroup) == 0 {
		err = fmt.Errorf("Input transaction group has 0 elements")
		return
	}

	emptyGroup := types.Digest{}

	// a group of size 1 may have an empty group ID
	if len(txgroup) == 1 && txgroup[0].Group == emptyGroup {
		valid = true
		return
	}

	var inputGroup types.Digest
	for i := range txgroup {
		if i == 0 {
			inputGroup = txgroup[i].Group
		} else {
			// ensure all txns in the input have the same group ID
			if txgroup[i].Group != inputGroup {
				valid = false
				return
			}
		}
		// zero out group IDs so we can compute it again
		txgroup[i].Group = emptyGroup
	}

	gid, err := crypto.ComputeGroupID(txgroup)
	if err == nil {
		valid = gid == inputGroup
	}

	return
}

// VerifyGroupID verifies that an atomic group of transactions all contain the correct group ID
func VerifyGroupID(txns *BytesArray) (valid bool, err error) {
	txgroup, err := decodeTxns(txns)
	if err == nil {
		valid, err = verifyTxnsGroupID(txgroup)
	}
	return
}

// FindAndVerifyTxnGroups takes an array of encoded transactions and finds and verifies consecutive
// transactions which claim to be an atomic group. It returns an array of integers which is the same
// length as the array of input transactions. If two indexes in this array have the same value, then
// the transactions at the same indexes in the input array are in the same atomic group.
func FindAndVerifyTxnGroups(txns *BytesArray) (groups *Int64Array, err error) {
	allTxns, err := decodeTxns(txns)
	if err != nil {
		return
	}

	if len(allTxns) == 0 {
		err = fmt.Errorf("Input transaction group has 0 elements")
		return
	}

	groupAssignment := make([]int64, len(allTxns))
	var groupStarts []int

	emptyGroup := types.Digest{}
	prevGroup := types.Digest{}

	// find groups
	var groupCount int64
	for i, txn := range allTxns {
		if txn.Group != prevGroup || txn.Group == emptyGroup {
			groupCount++
			groupStarts = append(groupStarts, i)
		}

		groupAssignment[i] = groupCount - 1

		prevGroup = txn.Group
	}

	// verify groups
	for i, indexGroupStart := range groupStarts {
		var indexGroupEnd int
		if i+1 < len(groupStarts) {
			indexGroupEnd = groupStarts[i+1]
		} else {
			indexGroupEnd = len(allTxns)
		}

		group := allTxns[indexGroupStart:indexGroupEnd]
		if len(group) == 0 {
			err = fmt.Errorf("Zero length group error")
			return
		}

		var valid bool
		valid, err = verifyTxnsGroupID(group)
		if err != nil {
			err = fmt.Errorf("Error when verifying group: %v", err)
			return
		}

		if !valid {
			err = fmt.Errorf("The transactions in range [%d:%d] form an invalid group", indexGroupStart, indexGroupEnd)
			return
		}
	}

	groups = &Int64Array{
		values: groupAssignment,
	}

	return
}


// SignBytes signs the bytes and returns the signature
func SignBytes(sk []byte, bytesToSign []byte) (signature []byte, err error) {
	if len(sk) != ed25519.PrivateKeySize {
		err = fmt.Errorf("Incorrect privateKey length expected %d, got %d", ed25519.PrivateKeySize, len(sk))
		return
	}

	// sign the bytes
	signature, err = crypto.SignBytes(sk, bytesToSign)
	return
}