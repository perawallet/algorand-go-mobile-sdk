package sdk

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/stretchr/testify/require"
)

func TestMakeBasicAccountSigner(t *testing.T) {
	t.Parallel()
	account := crypto.GenerateAccount()
	signer, err := MakeBasicAccountSigner(account.PrivateKey)
	require.NoError(t, err)

	require.True(t, signer.Equals(signer))

	params := types.SuggestedParams{
		Fee:             0,
		GenesisID:       "testnet-v1.0",
		GenesisHash:     mustDecodeB64(t, "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="),
		FirstRoundValid: 2,
		LastRoundValid:  1002,
	}
	txn, err := transaction.MakePaymentTxn("2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY", "S64XU5HQEY2XLHVUSO6RI3JL6NHC32I4LJHM32ZOM5VC4QPON7BZZRCU2E", 1_000_000, nil, "", params)
	require.NoError(t, err)
	encodedTxn := msgpack.Encode(&txn)
	indexesToSign := []int64{0}

	result, err := signer.SignTransactions(&BytesArray{[][]byte{encodedTxn}}, &Int64Array{indexesToSign})
	require.NoError(t, err)
	require.Equal(t, 1, result.Length())

	_, expectedStxnBytes, err := crypto.SignTransaction(account.PrivateKey, txn)
	require.NoError(t, err)
	require.Equal(t, expectedStxnBytes, result.Get(0))
}

func TestMakeLogicSigAccountSigner(t *testing.T) {
	t.Parallel()
	program := []byte{0x1, 0x20, 0x1, 0x1, 0x22}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}
	lsigAccount, err := crypto.MakeLogicSigAccountEscrowChecked(program, args)
	require.NoError(t, err)
	signer := MakeLogicSigAccountSigner(&LogicSigAccount{lsigAccount})

	require.True(t, signer.Equals(signer))

	params := types.SuggestedParams{
		Fee:             0,
		GenesisID:       "testnet-v1.0",
		GenesisHash:     mustDecodeB64(t, "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="),
		FirstRoundValid: 2,
		LastRoundValid:  1002,
	}
	txn, err := transaction.MakePaymentTxn("2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY", "S64XU5HQEY2XLHVUSO6RI3JL6NHC32I4LJHM32ZOM5VC4QPON7BZZRCU2E", 1_000_000, nil, "", params)
	require.NoError(t, err)
	encodedTxn := msgpack.Encode(&txn)
	indexesToSign := []int64{0}

	result, err := signer.SignTransactions(&BytesArray{[][]byte{encodedTxn}}, &Int64Array{indexesToSign})
	require.NoError(t, err)
	require.Equal(t, 1, result.Length())

	_, expectedStxnBytes, err := crypto.SignLogicSigAccountTransaction(lsigAccount, txn)
	require.NoError(t, err)
	require.Equal(t, expectedStxnBytes, result.Get(0))
}

func TestMakeMultiSigAccountTransactionSigner(t *testing.T) {
	t.Parallel()
	ma, acct1, _, acct3 := makeTestMultisigAccount(t)
	signer, err := MakeMultiSigAccountTransactionSigner(ma, &BytesArray{[][]byte{acct1.PrivateKey, acct3.PrivateKey}})
	require.NoError(t, err)

	require.True(t, signer.Equals(signer))

	params := types.SuggestedParams{
		Fee:             0,
		GenesisID:       "testnet-v1.0",
		GenesisHash:     mustDecodeB64(t, "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="),
		FirstRoundValid: 2,
		LastRoundValid:  1002,
	}
	txn, err := transaction.MakePaymentTxn("2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY", "S64XU5HQEY2XLHVUSO6RI3JL6NHC32I4LJHM32ZOM5VC4QPON7BZZRCU2E", 1_000_000, nil, "", params)
	require.NoError(t, err)
	encodedTxn := msgpack.Encode(&txn)
	indexesToSign := []int64{0}

	result, err := signer.SignTransactions(&BytesArray{[][]byte{encodedTxn}}, &Int64Array{indexesToSign})
	require.NoError(t, err)
	require.Equal(t, 1, result.Length())

	_, expectedStxnBytes, err := crypto.SignMultisigTransaction(acct1.PrivateKey, ma.value, txn)
	require.NoError(t, err)
	_, expectedStxnBytes, err = crypto.AppendMultisigTransaction(acct3.PrivateKey, ma.value, expectedStxnBytes)
	require.NoError(t, err)
	require.Equal(t, expectedStxnBytes, result.Get(0))
}
