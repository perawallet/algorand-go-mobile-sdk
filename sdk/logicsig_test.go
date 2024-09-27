package sdk

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
)

func makeTestMultisigAccount(t *testing.T) (*MultisigAccount, crypto.Account, crypto.Account, crypto.Account) {
	t.Helper()
	// DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA
	mn1 := "auction inquiry lava second expand liberty glass involve ginger illness length room item discover ahead table doctor term tackle cement bonus profit right above catch"
	sk1, err := mnemonic.ToPrivateKey(mn1)
	require.NoError(t, err)
	acct1, err := crypto.AccountFromPrivateKey(sk1)
	require.NoError(t, err)

	// BFRTECKTOOE7A5LHCF3TTEOH2A7BW46IYT2SX5VP6ANKEXHZYJY77SJTVM
	mn2 := "since during average anxiety protect cherry club long lawsuit loan expand embark forum theory winter park twenty ball kangaroo cram burst board host ability left"
	sk2, err := mnemonic.ToPrivateKey(mn2)
	require.NoError(t, err)
	acct2, err := crypto.AccountFromPrivateKey(sk2)
	require.NoError(t, err)

	// 47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU
	mn3 := "advice pudding treat near rule blouse same whisper inner electric quit surface sunny dismiss leader blood seat clown cost exist hospital century reform able sponsor"
	sk3, err := mnemonic.ToPrivateKey(mn3)
	require.NoError(t, err)
	acct3, err := crypto.AccountFromPrivateKey(sk3)
	require.NoError(t, err)

	ma, err := MakeMultisigAccount(1, 2, &StringArray{[]string{
		acct1.Address.String(),
		acct2.Address.String(),
		acct3.Address.String(),
	}})
	require.NoError(t, err)

	return ma, acct1, acct2, acct3
}

func TestMakeLogicSigAccountEscrow(t *testing.T) {
	t.Parallel()
	program := []byte{0x1, 0x20, 0x1, 0x1, 0x22}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}

	lsigAccount, err := MakeLogicSigAccountEscrow(program, &BytesArray{args})
	require.NoError(t, err)

	require.Equal(t, program, lsigAccount.value.Lsig.Logic)
	require.Equal(t, args, lsigAccount.value.Lsig.Args)
	require.Equal(t, types.Signature{}, lsigAccount.value.Lsig.Sig)
	require.True(t, lsigAccount.value.Lsig.Msig.Blank())
	require.Equal(t, ed25519.PublicKey(nil), lsigAccount.value.SigningKey)

	require.False(t, lsigAccount.IsDelegated())

	actualAddr, err := lsigAccount.Address()
	require.NoError(t, err)
	expectedAddr := "6Z3C3LDVWGMX23BMSYMANACQOSINPFIRF77H7N3AWJZYV6OH6GWTJKVMXY"
	require.Equal(t, expectedAddr, actualAddr)
}

func TestMakeLogicSigAccountDelegated(t *testing.T) {
	t.Parallel()
	program := []byte{0x1, 0x20, 0x1, 0x1, 0x22}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}

	account, err := crypto.AccountFromPrivateKey(ed25519.PrivateKey{0xd2, 0xdc, 0x4c, 0xcc, 0xe9, 0x98, 0x62, 0xff, 0xcf, 0x8c, 0xeb, 0x93, 0x6, 0xc4, 0x8d, 0xa6, 0x80, 0x50, 0x82, 0xa, 0xbb, 0x29, 0x95, 0x7a, 0xac, 0x82, 0x68, 0x9a, 0x8c, 0x49, 0x5a, 0x38, 0x5e, 0x67, 0x4f, 0x1c, 0xa, 0xee, 0xec, 0x37, 0x71, 0x89, 0x8f, 0x61, 0xc7, 0x6f, 0xf5, 0xd2, 0x4a, 0x19, 0x79, 0x3e, 0x2c, 0x91, 0xfa, 0x8, 0x51, 0x62, 0x63, 0xe3, 0x85, 0x73, 0xea, 0x42})
	require.NoError(t, err)
	signature := types.Signature{0x3e, 0x5, 0x3d, 0x39, 0x4d, 0xfb, 0x12, 0xbc, 0x65, 0x79, 0x9f, 0xea, 0x31, 0x8a, 0x7b, 0x8e, 0xa2, 0x51, 0x8b, 0x55, 0x2c, 0x8a, 0xbe, 0x6c, 0xd7, 0xa7, 0x65, 0x2d, 0xd8, 0xb0, 0x18, 0x7e, 0x21, 0x5, 0x2d, 0xb9, 0x24, 0x62, 0x89, 0x16, 0xe5, 0x61, 0x74, 0xcd, 0xf, 0x19, 0xac, 0xb9, 0x6c, 0x45, 0xa4, 0x29, 0x91, 0x99, 0x11, 0x1d, 0xe4, 0x7c, 0xe4, 0xfc, 0x12, 0xec, 0xce, 0x2}

	t.Run("provide sk", func(t *testing.T) {
		lsigAccount, err := MakeLogicSigAccountDelegatedSign(program, &BytesArray{args}, account.PrivateKey)
		require.NoError(t, err)

		require.Equal(t, program, lsigAccount.value.Lsig.Logic)
		require.Equal(t, args, lsigAccount.value.Lsig.Args)
		require.Equal(t, signature, lsigAccount.value.Lsig.Sig)
		require.True(t, lsigAccount.value.Lsig.Msig.Blank())
		require.Equal(t, account.PublicKey, lsigAccount.value.SigningKey)

		require.True(t, lsigAccount.IsDelegated())

		actualAddr, err := lsigAccount.Address()
		require.NoError(t, err)
		require.Equal(t, account.Address.String(), actualAddr)
	})

	t.Run("provide sig", func(t *testing.T) {
		lsigAccount, err := MakeLogicSigAccountDelegatedAttachSig(program, &BytesArray{args}, account.Address.String(), signature[:])
		require.NoError(t, err)

		require.Equal(t, program, lsigAccount.value.Lsig.Logic)
		require.Equal(t, args, lsigAccount.value.Lsig.Args)
		require.Equal(t, signature, lsigAccount.value.Lsig.Sig)
		require.True(t, lsigAccount.value.Lsig.Msig.Blank())
		require.Equal(t, account.PublicKey, lsigAccount.value.SigningKey)

		require.True(t, lsigAccount.IsDelegated())

		actualAddr, err := lsigAccount.Address()
		require.NoError(t, err)
		require.Equal(t, account.Address.String(), actualAddr)
	})
}

func TestMakeLogicSigAccountDelegatedMsig(t *testing.T) {
	t.Parallel()
	program := []byte{0x1, 0x20, 0x1, 0x1, 0x22}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}

	ma, acct1, acct2, _ := makeTestMultisigAccount(t)

	lsigAccount, err := MakeLogicSigAccountDelegatedMsig(program, &BytesArray{args}, ma)
	require.NoError(t, err)

	expectedMsig := types.MultisigSig{
		Version:   ma.value.Version,
		Threshold: ma.value.Threshold,
		Subsigs: []types.MultisigSubsig{
			{
				Key: ma.value.Pks[0],
			},
			{
				Key: ma.value.Pks[1],
			},
			{
				Key: ma.value.Pks[2],
			},
		},
	}
	expectedSig1 := types.Signature{0x49, 0x13, 0xb8, 0x5, 0xd1, 0x9e, 0x7f, 0x2c, 0x10, 0x80, 0xf6, 0x33, 0x7e, 0x18, 0x54, 0xa7, 0xce, 0xea, 0xee, 0x10, 0xdd, 0xbd, 0x13, 0x65, 0x84, 0xbf, 0x93, 0xb7, 0x5f, 0x30, 0x63, 0x15, 0x91, 0xca, 0x23, 0xc, 0xed, 0xef, 0x23, 0xd1, 0x74, 0x1b, 0x52, 0x9d, 0xb0, 0xff, 0xef, 0x37, 0x54, 0xd6, 0x46, 0xf4, 0xb5, 0x61, 0xfc, 0x8b, 0xbc, 0x2d, 0x7b, 0x4e, 0x63, 0x5c, 0xbd, 0x2}
	expectedSig2 := types.Signature{0x64, 0xbc, 0x55, 0xdb, 0xed, 0x91, 0xa2, 0x41, 0xd4, 0x2a, 0xb6, 0x60, 0xf7, 0xe1, 0x4a, 0xb9, 0x99, 0x9a, 0x52, 0xb3, 0xb1, 0x71, 0x58, 0xce, 0xfc, 0x3f, 0x4f, 0xe7, 0xcb, 0x22, 0x41, 0x14, 0xad, 0xa9, 0x3d, 0x5e, 0x84, 0x5, 0x2, 0xa, 0x17, 0xa6, 0x69, 0x83, 0x3, 0x22, 0x4e, 0x86, 0xa3, 0x8b, 0x6a, 0x36, 0xc5, 0x54, 0xbe, 0x20, 0x50, 0xff, 0xd3, 0xee, 0xa8, 0xb3, 0x4, 0x9}

	require.Equal(t, program, lsigAccount.value.Lsig.Logic)
	require.Equal(t, args, lsigAccount.value.Lsig.Args)
	require.Equal(t, types.Signature{}, lsigAccount.value.Lsig.Sig)
	require.Equal(t, expectedMsig, lsigAccount.value.Lsig.Msig)
	require.Equal(t, ed25519.PublicKey(nil), lsigAccount.value.SigningKey)

	require.True(t, lsigAccount.IsDelegated())

	actualAddr, err := lsigAccount.Address()
	require.NoError(t, err)
	expectedAddr, err := ma.value.Address()
	require.NoError(t, err)
	require.Equal(t, expectedAddr.String(), actualAddr)

	// attach first signature with AppendSignMultisigSignature
	err = lsigAccount.AppendSignMultisigSignature(acct1.PrivateKey)
	require.NoError(t, err)

	expectedMsig.Subsigs[0].Sig = expectedSig1
	require.Equal(t, expectedMsig, lsigAccount.value.Lsig.Msig)

	// attach second signature with AppendAttachMultisigSignature
	err = lsigAccount.AppendAttachMultisigSignature(acct2.Address.String(), expectedSig2[:])
	require.NoError(t, err)

	expectedMsig.Subsigs[1].Sig = expectedSig2
	require.Equal(t, expectedMsig, lsigAccount.value.Lsig.Msig)
}

func TestLogicSigJSON(t *testing.T) {
	t.Parallel()
	program := []byte{0x1, 0x20, 0x1, 0x1, 0x22}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}

	ma, acct1, _, _ := makeTestMultisigAccount(t)

	lsigAccount1, err := MakeLogicSigAccountDelegatedMsig(program, &BytesArray{args}, ma)
	require.NoError(t, err)
	expectedJson1 := `{
	"lsig": {
		"arg": [
			"AQ==",
			"AgM="
		],
		"l": "ASABASI=",
		"msig": {
			"subsig": [
				{
					"pk": "G37AsEvqYbeWkJfmy/QH4QinBTUdC8mKvrEiCairgXg="
				},
				{
					"pk": "CWMyCVNzifB1ZxF3OZHH0D4bc8jE9Sv2r/Aaolz5wnE="
				},
				{
					"pk": "5/D4TQaBHfnzHI2HixFV9GcdUaGFwgCQhmf0SVhwaKE="
				}
			],
			"thr": 2,
			"v": 1
		}
	}
}`

	lsigAccount2, err := MakeLogicSigAccountDelegatedSign(program, nil, acct1.PrivateKey)
	require.NoError(t, err)
	expectedJson2 := `{
	"sigkey": "G37AsEvqYbeWkJfmy/QH4QinBTUdC8mKvrEiCairgXg=",
	"lsig": {
		"l": "ASABASI=",
		"sig": "SRO4BdGefywQgPYzfhhUp87q7hDdvRNlhL+Tt18wYxWRyiMM7e8j0XQbUp2w/+83VNZG9LVh/Iu8LXtOY1y9Ag=="
	}
}`

	for _, testcase := range []struct {
		lsigAccount  *LogicSigAccount
		expectedJson string
	}{
		{
			lsigAccount:  lsigAccount1,
			expectedJson: expectedJson1,
		},
		{
			lsigAccount:  lsigAccount2,
			expectedJson: expectedJson2,
		},
	} {
		actualJson := testcase.lsigAccount.ToJSON()
		ok, err := jsonEqual(testcase.expectedJson, actualJson)
		require.NoError(t, err)
		require.True(t, ok, "expected: %s, actual: %s", testcase.expectedJson, actualJson)

		roundTripAccount, err := DeserializeLogicSigAccountFromJSON(actualJson)
		require.NoError(t, err)
		require.Equal(t, testcase.lsigAccount, roundTripAccount)
	}
}

func TestExtractLogicSigAccountFromSignedTransaction(t *testing.T) {
	t.Parallel()
	t.Run("valid escrow", func(t *testing.T) {
		encodedStx := mustDecodeB64(t, "g6Rsc2lngqNhcmeSxAEBxAICA6FsxAUBIAEBIqRzZ25yxCD2di2sdbGZfWwslhgGgFB0kNeVES/+f7dgsnOK+cfxraN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
		account, err := ExtractLogicSigAccountFromSignedTransaction(encodedStx)
		require.NoError(t, err)

		expectedLsig := crypto.LogicSigAccount{
			Lsig: types.LogicSig{
				Logic: []byte{0x1, 0x20, 0x1, 0x1, 0x22},
				Args: [][]byte{
					{0x01},
					{0x02, 0x03},
				},
			},
		}
		require.Equal(t, expectedLsig, account.value)
	})

	t.Run("valid delegated", func(t *testing.T) {
		encodedStx := mustDecodeB64(t, "g6Rsc2lng6NhcmeSxAEBxAICA6FsxAUBIAEBIqNzaWfEQEkTuAXRnn8sEID2M34YVKfO6u4Q3b0TZYS/k7dfMGMVkcojDO3vI9F0G1KdsP/vN1TWRvS1YfyLvC17TmNcvQKkc2ducsQgG37AsEvqYbeWkJfmy/QH4QinBTUdC8mKvrEiCairgXijdHhuiaNhbXTOAA9CQKNmZWXNA+iiZnYCo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToiomx2zQPqo3JjdsQgl7l6dPAmNXWetJO9FG0r804t6RxaTs3rLmdqLkHub8Ojc25kxCDUYfSDPMXT2TQv+8ZvQtqEoK0G1xsTfW1wkIkR+AfA4qR0eXBlo3BheQ==")
		account, err := ExtractLogicSigAccountFromSignedTransaction(encodedStx)
		require.NoError(t, err)

		expectedSigKey := mustDecodeAddress(t, "DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
		var expectedSig types.Signature
		copy(expectedSig[:], mustDecodeB64(t, "SRO4BdGefywQgPYzfhhUp87q7hDdvRNlhL+Tt18wYxWRyiMM7e8j0XQbUp2w/+83VNZG9LVh/Iu8LXtOY1y9Ag=="))
		expectedLsig := crypto.LogicSigAccount{
			SigningKey: expectedSigKey[:],
			Lsig: types.LogicSig{
				Logic: []byte{0x1, 0x20, 0x1, 0x1, 0x22},
				Args: [][]byte{
					{0x01},
					{0x02, 0x03},
				},
				Sig: expectedSig,
			},
		}
		require.Equal(t, expectedLsig, account.value)
	})

	t.Run("no lsig", func(t *testing.T) {
		encodedStx := mustDecodeB64(t, "gqNzaWfEQC/nu0j+joowxkE3uMMx3SPZFzOHq8YdeTjqYVgV5r6xL/w4wRaYKhZXSVbeTnr2udAsbcxWUm7mhYTzH7AeEwejdHhuiaNhbXTOAA9CQKNmZWXNA+iiZnYCo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToiomx2zQPqo3JjdsQgl7l6dPAmNXWetJO9FG0r804t6RxaTs3rLmdqLkHub8Ojc25kxCDUYfSDPMXT2TQv+8ZvQtqEoK0G1xsTfW1wkIkR+AfA4qR0eXBlo3BheQ==")
		account, err := ExtractLogicSigAccountFromSignedTransaction(encodedStx)
		require.NoError(t, err)
		require.Nil(t, account)
	})
}

func TestSignLogicSigTransaction(t *testing.T) {
	t.Parallel()
	program := []byte{0x1, 0x20, 0x1, 0x1, 0x22}
	args := [][]byte{
		{0x01},
		{0x02, 0x03},
	}
	lsigAccount, err := MakeLogicSigAccountEscrow(program, &BytesArray{args})
	require.NoError(t, err)

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

	stxBytes, err := SignLogicSigTransaction(lsigAccount, encodedTxn)
	require.NoError(t, err)

	expectedStxBytes := mustDecodeB64(t, "g6Rsc2lngqNhcmeSxAEBxAICA6FsxAUBIAEBIqRzZ25yxCD2di2sdbGZfWwslhgGgFB0kNeVES/+f7dgsnOK+cfxraN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
	require.Equal(t, expectedStxBytes, stxBytes)
}

func TestLogicSigProgramForSigning(t *testing.T) {
	t.Parallel()
	program := []byte{0x1, 0x20, 0x1, 0x1, 0x22}
	actual := LogicSigProgramForSigning(program)
	expected := []byte{0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x1, 0x20, 0x1, 0x1, 0x22}
	require.Equal(t, expected, actual)
}
