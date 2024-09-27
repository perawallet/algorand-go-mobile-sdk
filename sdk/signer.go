package sdk

import (
	"bytes"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// TransactionSigner is an interface which can sign transactions from an atomic transaction group.
//
// SignTransactions(txGroup,indexesToSign) signs the transactions in txGroup at the indexes
// specified in indexesToSign.
//
// Equals(other) returns true if the two signers are equivalent.
type TransactionSigner interface {
	SignTransactions(txGroup *BytesArray, indexesToSign *Int64Array) (*BytesArray, error)
	Equals(other TransactionSigner) bool
}

type internalToExternalSigner struct {
	internalSigner transaction.TransactionSigner
}

func (s internalToExternalSigner) SignTransactions(txGroup *BytesArray, indexesToSign *Int64Array) (*BytesArray, error) {
	txns := make([]types.Transaction, txGroup.Length())
	for i, txBytes := range txGroup.Extract() {
		var tx types.Transaction
		err := msgpack.Decode(txBytes, &tx)
		if err != nil {
			return nil, err
		}
		txns[i] = tx
	}

	indexes := make([]int, indexesToSign.Length())
	for i, index := range indexesToSign.Extract() {
		indexes[i] = int(index)
	}

	stxBytes, err := s.internalSigner.SignTransactions(txns, indexes)
	if err != nil {
		return nil, err
	}
	return &BytesArray{stxBytes}, nil
}

func (s internalToExternalSigner) Equals(other TransactionSigner) bool {
	if casted, ok := other.(internalToExternalSigner); ok {
		return s.internalSigner.Equals(casted.internalSigner)
	}
	return false
}

type externalToInternalSigner struct {
	externalSigner TransactionSigner
}

func (s externalToInternalSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	txBytes := make([][]byte, len(txGroup))
	for i, tx := range txGroup {
		txBytes[i] = msgpack.Encode(&tx)
	}
	int64Indexes := make([]int64, len(indexesToSign))
	for i, index := range indexesToSign {
		int64Indexes[i] = int64(index)
	}
	stxBytes, err := s.externalSigner.SignTransactions(&BytesArray{txBytes}, &Int64Array{int64Indexes})
	if err != nil {
		return nil, err
	}
	return stxBytes.Extract(), nil
}

func (s externalToInternalSigner) Equals(other transaction.TransactionSigner) bool {
	if castedSigner, ok := other.(externalToInternalSigner); ok {
		return s.externalSigner.Equals(castedSigner.externalSigner)
	}
	return false
}

// MakeBasicAccountSigner creates a TransactionSigner for a basic account from a private key.
func MakeBasicAccountSigner(sk []byte) (TransactionSigner, error) {
	account, err := crypto.AccountFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}
	return internalToExternalSigner{transaction.BasicAccountTransactionSigner{account}}, nil
}

// MakeLogicSigAccountSigner creates a TransactionSigner for a LogicSigAccount.
func MakeLogicSigAccountSigner(ls *LogicSigAccount) TransactionSigner {
	return internalToExternalSigner{transaction.LogicSigAccountTransactionSigner{ls.value}}
}

// MakeMultiSigAccountTransactionSigner creates a TransactionSigner for a MultisigAccount with the
// given component account private keys.
//
// There must be enough private keys to meet the multisig threshold.
func MakeMultiSigAccountTransactionSigner(msig *MultisigAccount, sks *BytesArray) (TransactionSigner, error) {
	privateKeys := sks.Extract()
	seenPkIndexes := make(map[int]bool)
	for _, sk := range privateKeys {
		account, err := crypto.AccountFromPrivateKey(sk)
		if err != nil {
			return nil, err
		}
		pkIndex := -1
		for i, pk := range msig.value.Pks {
			if bytes.Equal(pk, account.PublicKey) {
				pkIndex = i
				break
			}
		}
		if seenPkIndexes[pkIndex] {
			return nil, fmt.Errorf("duplicate private key for public key %s", account.Address)
		}
		seenPkIndexes[pkIndex] = true
	}
	if len(seenPkIndexes) < int(msig.value.Threshold) {
		return nil, fmt.Errorf("not enough private keys to meet multisig threshold. Have %d, need %d", len(seenPkIndexes), msig.value.Threshold)
	}
	return internalToExternalSigner{transaction.MultiSigAccountTransactionSigner{msig.value, privateKeys}}, nil
}
