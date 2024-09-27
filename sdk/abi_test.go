package sdk

import (
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/abi"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/json"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/stretchr/testify/require"
)

func TestABIType(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		typeString   string
		jsonValue    string
		encodedValue []byte
	}{
		{
			typeString:   "uint64",
			jsonValue:    "123",
			encodedValue: []byte{0, 0, 0, 0, 0, 0, 0, 123},
		},
		{
			typeString: "string",
			jsonValue:  `"hello world"`,
			encodedValue: []byte{
				0, 11, // length
				104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, // "hello world"
			},
		},
		{
			typeString: "address",
			jsonValue:  `"WWYNX3TKQYVEREVSW6QQP3SXSFOCE3SKUSEIVJ7YAGUPEACNI5UGI4DZCE"`,
			encodedValue: []byte{
				181, 176, 219, 238, 106, 134, 42, 72,
				146, 178, 183, 161, 7, 238, 87, 145,
				92, 34, 110, 74, 164, 136, 138, 167,
				248, 1, 168, 242, 0, 77, 71, 104,
			},
		},
		{
			typeString: "byte[]",
			jsonValue:  `"aGVsbG8gd29ybGQ="`,
			encodedValue: []byte{
				0, 11, // length
				104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, // "hello world"
			},
		},
	}

	for _, tc := range testCases {
		name := fmt.Sprintf("%s:%s", tc.typeString, tc.jsonValue)
		t.Run(name, func(t *testing.T) {
			abiType, err := ParseABIType(tc.typeString)
			require.NoError(t, err)

			require.Equal(t, tc.typeString, abiType.String())

			encoded, err := abiType.Encode(tc.jsonValue)
			require.NoError(t, err)
			require.Equal(t, tc.encodedValue, encoded)

			decoded, err := abiType.Decode(encoded)
			require.NoError(t, err)
			require.Equal(t, tc.jsonValue, decoded)
		})
	}
}

func TestABIMethodFunctions(t *testing.T) {
	t.Parallel()
	methodJSONWithArgNames := `{"name":"add","args":[{"name":"a","type":"uint8"},{"name":"b","type":"uint16"}],"returns":{"type":"uint32"}}`
	methodJSONWithoutArgNames := `{"name":"add","args":[{"type":"uint8"},{"type":"uint16"}],"returns":{"type":"uint32"}}`
	signature := "add(uint8,uint16)uint32"

	methodFromSignature, err := ABIMethodJSONFromSignature(signature)
	require.NoError(t, err)
	ok, err := jsonEqual(methodJSONWithoutArgNames, methodFromSignature)
	require.True(t, ok, "expected: %s, actual %s", methodJSONWithoutArgNames, methodFromSignature)

	for _, method := range []string{methodJSONWithArgNames, methodJSONWithoutArgNames} {
		signatureFromJSON, err := GetABIMethodSignature(method)
		require.NoError(t, err)
		require.Equal(t, signature, signatureFromJSON)
	}
}

func TestNewAtomicTransactionComposer(t *testing.T) {
	t.Parallel()
	atc := NewAtomicTransactionComposer()
	require.Equal(t, atc.GetStatus(), transaction.BUILDING)
	require.Equal(t, atc.Count(), 0)
	copyAtc := atc.Clone()
	require.Equal(t, atc, copyAtc)
}

func TestATCAddTransaction(t *testing.T) {
	t.Parallel()
	atc := NewAtomicTransactionComposer()

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     addr,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			Note:       []byte{180, 81, 121, 57, 252, 250, 210, 113},
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: addr,
			Amount:   5000,
		},
	}

	encodedTxn := msgpack.Encode(&tx)

	err = atc.AddTransaction(encodedTxn, nil)
	require.NoError(t, err)

	require.Equal(t, atc.GetStatus(), transaction.BUILDING)
	require.Equal(t, atc.Count(), 1)
}

func TestATCAddMethodCallParams(t *testing.T) {
	t.Parallel()

	uint64ArrayToInt64Array := func(uint64Array []uint64) []int64 {
		int64Array := make([]int64, len(uint64Array))
		for i, v := range uint64Array {
			int64Array[i] = int64(v)
		}
		return int64Array
	}

	convertAppBoxReferences := func(refs []types.AppBoxReference) *AppBoxRefArray {
		var array AppBoxRefArray
		for _, ref := range refs {
			array.Append(int64(ref.AppID), ref.Name)
		}
		return &array
	}

	convertSuggestedParams := func(sp types.SuggestedParams) *SuggestedParams {
		return &SuggestedParams{
			Fee:             int64(sp.Fee),
			FlatFee:         sp.FlatFee,
			FirstRoundValid: int64(sp.FirstRoundValid),
			LastRoundValid:  int64(sp.LastRoundValid),
			GenesisID:       sp.GenesisID,
			GenesisHash:     sp.GenesisHash,
		}
	}

	appID := uint64(10)
	onComplete := types.OptInOC
	method := abi.Method{
		Name: "example",
		Args: []abi.Arg{
			{
				Type: "uint8",
			},
			{
				Type: "uint16",
			},
			{
				Type: "account",
			},
			{
				Type: "pay",
			},
		},
		Returns: abi.Return{
			Type: "uint32",
		},
	}
	foreignAccounts := []string{"DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA"}
	foreignApps := []uint64{10, 20, 30}
	foreignAssets := []uint64{5, 7, 9}
	boxRefs := []types.AppBoxReference{
		{AppID: 0, Name: []byte("box1")},
		{AppID: 20, Name: []byte("box2")},
	}
	txnParams := types.SuggestedParams{
		GenesisID:       "testnet-v1.0",
		GenesisHash:     mustDecodeB64(t, "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="),
		FirstRoundValid: 2,
		LastRoundValid:  1002,
		Fee:             1000,
		FlatFee:         true,
	}
	note := []byte("note value")
	sender, err := types.DecodeAddress("2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY")
	require.NoError(t, err)
	var signer TransactionSigner

	params, err := NewAddMethodCallParams(
		int64(appID),
		int(onComplete),
		string(json.Encode(&method)),
		&StringArray{foreignAccounts},
		&Int64Array{uint64ArrayToInt64Array(foreignApps)},
		&Int64Array{uint64ArrayToInt64Array(foreignAssets)},
		convertAppBoxReferences(boxRefs),
		convertSuggestedParams(txnParams),
		note,
		sender.String(),
		signer,
	)
	require.NoError(t, err)

	expectedParams := transaction.AddMethodCallParams{
		AppID:           appID,
		OnComplete:      onComplete,
		Method:          method,
		ForeignAccounts: foreignAccounts,
		ForeignAssets:   foreignAssets,
		ForeignApps:     foreignApps,
		BoxReferences:   boxRefs,
		SuggestedParams: txnParams,
		Note:            note,
		Sender:          sender,
		Signer:          externalToInternalSigner{signer},
	}
	require.Equal(t, expectedParams, params.value)

	err = params.AddMethodArgument("100")
	require.NoError(t, err)

	expectedParams.MethodArgs = []interface{}{uint8(100)}
	require.Equal(t, expectedParams, params.value)

	err = params.AddMethodArgument("222")
	require.NoError(t, err)

	expectedParams.MethodArgs = append(expectedParams.MethodArgs, uint16(222))
	require.Equal(t, expectedParams, params.value)

	accountArg, err := types.DecodeAddress("W3KCADJF23RDTO3TMY63YQBKYDYFPHFBU75JQMX5QHOERRBOZ75L3B2J7Y")
	require.NoError(t, err)

	err = params.AddMethodArgument(fmt.Sprintf(`"%s"`, accountArg))
	require.NoError(t, err)

	expectedParams.MethodArgs = append(expectedParams.MethodArgs, accountArg[:])
	require.Equal(t, expectedParams, params.value)

	var paymentSigner TransactionSigner
	paymentArg, err := transaction.MakePaymentTxn(
		"S64XU5HQEY2XLHVUSO6RI3JL6NHC32I4LJHM32ZOM5VC4QPON7BZZRCU2E",
		"2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY",
		1234567,
		[]byte("payment txn note"),
		"W3KCADJF23RDTO3TMY63YQBKYDYFPHFBU75JQMX5QHOERRBOZ75L3B2J7Y",
		txnParams,
	)
	require.NoError(t, err)

	err = params.AddMethodArgumentTransaction(msgpack.Encode(&paymentArg), paymentSigner)
	require.NoError(t, err)

	expectedParams.MethodArgs = append(expectedParams.MethodArgs, transaction.TransactionWithSigner{
		Txn:    paymentArg,
		Signer: externalToInternalSigner{paymentSigner},
	})
	require.Equal(t, expectedParams, params.value)

	approval := []byte("approval program")
	clear := []byte("clear program")
	params.AddPrograms(approval, clear)

	expectedParams.ApprovalProgram = approval
	expectedParams.ClearProgram = clear
	require.Equal(t, expectedParams, params.value)

	gSchema := types.StateSchema{
		NumByteSlice: 1,
		NumUint:      2,
	}
	lSchema := types.StateSchema{
		NumByteSlice: 3,
		NumUint:      4,
	}
	extraPages := uint32(5)

	err = params.AddAppSchema(int64(gSchema.NumUint), int64(gSchema.NumByteSlice), int64(lSchema.NumUint), int64(lSchema.NumByteSlice), int32(extraPages))
	require.NoError(t, err)

	expectedParams.GlobalSchema = gSchema
	expectedParams.LocalSchema = lSchema
	expectedParams.ExtraPages = extraPages
	require.Equal(t, expectedParams, params.value)
}

func TestATCAddMethodCall(t *testing.T) {
	t.Parallel()
	atc := NewAtomicTransactionComposer()

	methodSig := "add()uint32"
	method, err := abi.MethodFromSignature(methodSig)
	require.NoError(t, err)

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	err = atc.AddMethodCall(
		&AddMethodCallParams{
			transaction.AddMethodCallParams{
				AppID:  4,
				Method: method,
				Sender: addr,
				Signer: externalToInternalSigner{},
			},
		})
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), transaction.BUILDING)
	require.Equal(t, atc.Count(), 1)
}

func TestATCAddMethodCallWithManualForeignArgs(t *testing.T) {
	t.Parallel()
	atc := NewAtomicTransactionComposer()

	methodSig := "add(application)uint32"
	method, err := abi.MethodFromSignature(methodSig)
	require.NoError(t, err)

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	arg_addr_str := "E4VCHISDQPLIZWMALIGNPK2B2TERPDMR64MZJXE3UL75MUDXZMADX5OWXM"
	arg_addr, err := types.DecodeAddress(arg_addr_str)
	require.NoError(t, err)

	params := transaction.AddMethodCallParams{
		AppID:           4,
		Method:          method,
		Sender:          addr,
		Signer:          externalToInternalSigner{},
		MethodArgs:      []interface{}{2},
		ForeignApps:     []uint64{1},
		ForeignAssets:   []uint64{5},
		ForeignAccounts: []string{arg_addr_str},
	}
	err = atc.AddMethodCall(&AddMethodCallParams{params})
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), transaction.BUILDING)
	require.Equal(t, atc.Count(), 1)
	txnBytes, err := atc.BuildGroup()
	require.NoError(t, err)

	var txns []types.Transaction
	for _, txBytes := range txnBytes.Extract() {
		var tx types.Transaction
		err = msgpack.Decode(txBytes, &tx)
		require.NoError(t, err)
		txns = append(txns, tx)
	}

	require.Equal(t, len(txns[0].ForeignApps), 2)
	require.Equal(t, txns[0].ForeignApps[0], types.AppIndex(1))
	require.Equal(t, txns[0].ForeignApps[1], types.AppIndex(2))
	// verify original params object hasn't changed.
	require.Equal(t, params.ForeignApps, []uint64{1})

	require.Equal(t, len(txns[0].ForeignAssets), 1)
	require.Equal(t, txns[0].ForeignAssets[0], types.AssetIndex(5))

	require.Equal(t, len(txns[0].Accounts), 1)
	require.Equal(t, txns[0].Accounts[0], arg_addr)
}

func TestATCGatherSignatures(t *testing.T) {
	t.Parallel()
	atc := NewAtomicTransactionComposer()
	account := crypto.GenerateAccount()
	txSigner := internalToExternalSigner{transaction.BasicAccountTransactionSigner{Account: account}}

	addr, err := types.DecodeAddress("DN7MBMCL5JQ3PFUQS7TMX5AH4EEKOBJVDUF4TCV6WERATKFLQF4MQUPZTA")
	require.NoError(t, err)

	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     addr,
			Fee:        217000,
			FirstValid: 972508,
			LastValid:  973508,
			Note:       []byte{180, 81, 121, 57, 252, 250, 210, 113},
			GenesisID:  "testnet-v31.0",
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: addr,
			Amount:   5000,
		},
	}

	err = atc.AddTransaction(msgpack.Encode(&tx), txSigner)
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), transaction.BUILDING)
	require.Equal(t, atc.Count(), 1)

	sigs, err := atc.GatherSignatures()
	require.NoError(t, err)
	require.Equal(t, atc.GetStatus(), transaction.SIGNED)
	require.Equal(t, 1, sigs.Length())

	txBytes, _ := atc.BuildGroup()
	var firstTx types.Transaction
	err = msgpack.Decode(txBytes.Get(0), &firstTx)
	require.NoError(t, err)
	require.Equal(t, types.Digest{}, firstTx.Group)

	_, expectedSig, err := crypto.SignTransaction(account.PrivateKey, tx)
	require.NoError(t, err)
	require.Equal(t, len(expectedSig), len(sigs.Get(0)))
	require.Equal(t, expectedSig, sigs.Get(0))
}
