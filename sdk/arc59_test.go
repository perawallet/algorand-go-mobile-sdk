package sdk

import (
	"encoding/base64"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/stretchr/testify/require"
)

func TestMakeAndSignARC59OptInTxn(t *testing.T) {
	t.Parallel()
	// corresponds to SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4
	sk, err := mnemonic.ToPrivateKey("ocean tank film evil fresh ability capital huge ensure chat small dentist garlic slam decide extra fly train cross rib dog federal monitor about thought")
	require.NoError(t, err)

	gh, err := base64.StdEncoding.DecodeString("SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=")
	require.NoError(t, err)

	suggested_params := SuggestedParams{
		Fee:             0,
		GenesisID:       "testnet-v1.0",
		GenesisHash:     gh,
		FirstRoundValid: 40432872,
		LastRoundValid:  40433872,
		FlatFee:         false,
	}

	// https://testnet.explorer.perawallet.app/tx-group/UGbP1Hz6KLznhcbB6+6qCS2UiUWzLVsc+nO2WJvi7uY=/
	txnsByteArray, err := MakeAndSignARC59OptInTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"MEKFJGDJTHSBCAUMH5UFV7BGICQ3UCGUVR5CD6GURFUBYHUYSWQDLEGVXU",
		655494101,
		655977010,
		&suggested_params,
		sk,
	)
	require.NoError(t, err)
	txns := txnsByteArray.Extract()

	expectedStxnBytes, err := base64.StdEncoding.DecodeString("gqNzaWfEQFeOrmXMibgrrFvt/LpN9V4T3VAUqIywyVqODPrY/00syeYaP83zNHBvND2HvriQQo8NyRaEvb30bYnozcR8pACjdHhuiqNhbXTOAAGGoKNmZWXNA+iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIHnJL9z3H4kzJsIJNZOQflcpkswjKbuFC6ij79fN5zVqomx2zgJo+NCjcmN2xCBhFFSYaZnkEQKMP2ha/CZAoboI1Kx6IfjUiWgcHpiVoKNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWjcGF5")
	require.NoError(t, err)
	require.Equal(t, expectedStxnBytes, txns[0])

	expectedStxnBytes_2, err := base64.StdEncoding.DecodeString("gqNzaWfEQE8zvhg+o8qoawSeZByUfpoNU0Y0CJ0/Y5WQta0lZQgXcqOsNq0hOSo9TZVWqK/A35QxJBRddAHw98ZeXNiB3AmjdHhujKRhcGFhksQE6FQIEMQIAAAAACcZajKkYXBhc5HOJxlqMqRhcGF0kcQgYRRUmGmZ5BECjD9oWvwmQKG6CNSseiH41IloHB6YlaCkYXBpZM4nEgvVo2ZlZc0H0KJmds4CaPToo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQgeckv3PcfiTMmwgk1k5B+VymSzCMpu4ULqKPv183nNWqibHbOAmj40KNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWkYXBwbA==")
	require.NoError(t, err)
	require.Equal(t, expectedStxnBytes_2, txns[1])

}

func TestMakeAndSignARC59SendTxn(t *testing.T) {
	t.Parallel()
	// corresponds to SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4
	sk, err := mnemonic.ToPrivateKey("ocean tank film evil fresh ability capital huge ensure chat small dentist garlic slam decide extra fly train cross rib dog federal monitor about thought")
	require.NoError(t, err)

	gh, err := base64.StdEncoding.DecodeString("SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=")
	require.NoError(t, err)

	suggested_params := SuggestedParams{
		Fee:             0,
		GenesisID:       "testnet-v1.0",
		GenesisHash:     gh,
		FirstRoundValid: 40432872,
		LastRoundValid:  40433872,
		FlatFee:         false,
	}

	amount := MakeUint64(10)
	min_balance_requirement := MakeUint64(228100)
	algoAmount := MakeUint64(0)

	txnsByteArray, err := MakeAndSignARC59SendTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"MKKKFL5JBJTOCEMEZUAJKTWD5FYAI2FOLW5BP5N5YR37ZG5FHLTUYCFC6U",
		"MEKFJGDJTHSBCAUMH5UFV7BGICQ3UCGUVR5CD6GURFUBYHUYSWQDLEGVXU",
		"6LD3JUWPR72DX5JNGPHH2QEM2IKRA3SXGFSRR4P7TB5JMVPQTIYCXLUMCQ",
		&amount,
		&min_balance_requirement,
		5,
		655494101,
		655977010,
		&suggested_params,
		&algoAmount,
		sk,
	)
	require.NoError(t, err)
	txns := txnsByteArray.Extract()

	require.Equal(
		t,
		"gqNzaWfEQAJI+XbWaTbRFOipX1jvScgbbuW+QW9Hd2+jXcTU91516zWztA7jZewq78vq3mITew0HyPfmJSCAX7FaYMauZgGjdHhuiqNhbXTOAAN7BKNmZWXNA+iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIHtvbPII4jrvH+i4COcUIOLk3LLQSSF2X5IuB8aVvhSHomx2zgJo+NCjcmN2xCBhFFSYaZnkEQKMP2ha/CZAoboI1Kx6IfjUiWgcHpiVoKNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWjcGF5",
		base64.StdEncoding.EncodeToString(txns[0]),
	)

	require.Equal(
		t,
		"gqNzaWfEQBjnzGDGxeQsJlRavC0m2/iXdjOCNgiKhDqnDC3f9MBMSJKoSbSTTaaXTmjCRxgcRLbZ/v9t24ccOyUaF042zwejdHhui6RhYW10CqRhcmN2xCBhFFSYaZnkEQKMP2ha/CZAoboI1Kx6IfjUiWgcHpiVoKNmZWXNA+iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIHtvbPII4jrvH+i4COcUIOLk3LLQSSF2X5IuB8aVvhSHomx2zgJo+NCjc25kxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6R0eXBlpWF4ZmVypHhhaWTOJxlqMg==",
		base64.StdEncoding.EncodeToString(txns[1]),
	)

	require.Equal(
		t,
		"gqNzaWfEQC+328x1ZoqDprBNBz2Zc2ZvC2CTVR/IS9hvpilakRKyY7UwaJVAHA+SYjEbSnsdp5rIM34omz68YmDMplEWeQejdHhujaRhcGFhk8QECFMe18QgYpSir6kKZuERhM0AlU7D6XAEaK5duhf1vcR3/JulOufECAAAAAAAAAAApGFwYXORzicZajKkYXBhdJLEIGKUoq+pCmbhEYTNAJVOw+lwBGiuXboX9b3Ed/ybpTrnxCDyx7TSz4/0O/UtM859QIzSFRBuVzFlGPH/mHqWVfCaMKRhcGJ4kYGhbsQgYpSir6kKZuERhM0AlU7D6XAEaK5duhf1vcR3/JulOuekYXBpZM4nEgvVo2ZlZc0XcKJmds4CaPToo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQge29s8gjiOu8f6LgI5xQg4uTcstBJIXZfki4HxpW+FIeibHbOAmj40KNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWkYXBwbA==",
		base64.StdEncoding.EncodeToString(txns[2]),
	)

}

func TestMakeAndSignARC59ClaimTxn(t *testing.T) {
	t.Parallel()
	// corresponds to SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4
	sk, err := mnemonic.ToPrivateKey("ocean tank film evil fresh ability capital huge ensure chat small dentist garlic slam decide extra fly train cross rib dog federal monitor about thought")
	require.NoError(t, err)

	gh, err := base64.StdEncoding.DecodeString("SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=")
	require.NoError(t, err)

	suggested_params := SuggestedParams{
		Fee:             0,
		GenesisID:       "testnet-v1.0",
		GenesisHash:     gh,
		FirstRoundValid: 40432872,
		LastRoundValid:  40433872,
		FlatFee:         false,
	}

	// test with opt-in true

	txnsByteArrayOptInTrue, err := MakeAndSignARC59ClaimTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"6LD3JUWPR72DX5JNGPHH2QEM2IKRA3SXGFSRR4P7TB5JMVPQTIYCXLUMCQ",
		655494101,
		655977010,
		&suggested_params,
		true,
		false,
		sk,
	)
	require.NoError(t, err)
	txns_opt_in_true := txnsByteArrayOptInTrue.Extract()

	require.Equal(
		t,
		"gqNzaWfEQDGOkcKzLGrBcOeZmgxLI4z0VLbkz+1zZ8LPEKleCYmL4IvQS23cefFDxzpaObW/df79rrsiUQgbPfjmqEnCVQijdHhujaRhcGFhksQEv5AuPMQIAAAAACcZajKkYXBhc5HOJxlqMqRhcGF0kcQg8se00s+P9Dv1LTPOfUCM0hUQblcxZRjx/5h6llXwmjCkYXBieJGBoW7EIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpGFwaWTOJxIL1aNmZWXNC7iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIHc1uwt20dMbyuv/VyzNDHyhm69fQvuDhfq19k03YxvBomx2zgJo+NCjc25kxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6R0eXBlpGFwcGw=",
		base64.StdEncoding.EncodeToString(txns_opt_in_true[0]),
	)
	require.NoError(t, err)

	// test with opt-in false

	txnsByteArrayOptInFalse, err := MakeAndSignARC59ClaimTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"6LD3JUWPR72DX5JNGPHH2QEM2IKRA3SXGFSRR4P7TB5JMVPQTIYCXLUMCQ",
		655494101,
		655977010,
		&suggested_params,
		false,
		false,
		sk,
	)
	require.NoError(t, err)
	txns_opt_in_false := txnsByteArrayOptInFalse.Extract()

	require.Equal(
		t,
		"gqNzaWfEQAqRStLhafCrzrJzJQ5LQaRkHW09jT+Z61ZhG6i1HZ1RGaxtLX1wu7xsVxlVbqUvRGuH8jKKrnKAUBjGQpSIkgCjdHhuiqRhcmN2xCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6NmZWXNA+iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIHnfsWb5RyTDo2Jl2WK5A8NxnXGZJad9oOK87LeY5p0somx2zgJo+NCjc25kxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6R0eXBlpWF4ZmVypHhhaWTOJxlqMg==",
		base64.StdEncoding.EncodeToString(txns_opt_in_false[0]),
	)
	require.Equal(
		t,
		"gqNzaWfEQDLs1J4n4YEoG8rtQdfj/yjNFLlBhXotYUpLAgXxC/Ev8CA3ha4+CMga9NkZy/LzlkO9pH+bO5ASaUaSeZUCCgejdHhujaRhcGFhksQEv5AuPMQIAAAAACcZajKkYXBhc5HOJxlqMqRhcGF0kcQg8se00s+P9Dv1LTPOfUCM0hUQblcxZRjx/5h6llXwmjCkYXBieJGBoW7EIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpGFwaWTOJxIL1aNmZWXNC7iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIHnfsWb5RyTDo2Jl2WK5A8NxnXGZJad9oOK87LeY5p0somx2zgJo+NCjc25kxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6R0eXBlpGFwcGw=",
		base64.StdEncoding.EncodeToString(txns_opt_in_false[1]),
	)
	require.NoError(t, err)
}

func TestMakeAndSignARC59RejectTxn(t *testing.T) {
	t.Parallel()
	// corresponds to SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4
	sk, err := mnemonic.ToPrivateKey("ocean tank film evil fresh ability capital huge ensure chat small dentist garlic slam decide extra fly train cross rib dog federal monitor about thought")
	require.NoError(t, err)

	gh, err := base64.StdEncoding.DecodeString("SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=")
	require.NoError(t, err)

	suggested_params := SuggestedParams{
		Fee:             0,
		GenesisID:       "testnet-v1.0",
		GenesisHash:     gh,
		FirstRoundValid: 40432872,
		LastRoundValid:  40433872,
		FlatFee:         false,
	}

	txn, err := MakeAndSignARC59RejectTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"6LD3JUWPR72DX5JNGPHH2QEM2IKRA3SXGFSRR4P7TB5JMVPQTIYCXLUMCQ",
		"CRTRSWA2Y242PCGITJHCF3WYXAUKCF5VH5IMUT4DCMEV7FQGBMUHPMMWJ4",
		655494101,
		655977010,
		&suggested_params,
		false,
		sk,
	)
	require.NoError(t, err)

	require.Equal(
		t,
		"gqNzaWfEQCwZMrgN4U0DHGHzp53UWg7+sbVMxcBWZUfdpOPytBIPN+Z081TkRlPXUg9Kf0dvc6KLBfWyfHqRltz2oVh09w+jdHhujaRhcGFhksQEibPJzcQIAAAAACcZajKkYXBhc5HOJxlqMqRhcGF0ksQg8se00s+P9Dv1LTPOfUCM0hUQblcxZRjx/5h6llXwmjDEIBRnGVgaxrmniMiaTiLu2LgooRe1P1DKT4MTCV+WBgsopGFwYniRgaFuxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6RhcGlkzicSC9WjZmVlzQu4omZ2zgJo9OijZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKjZ3JwxCBFmFWCaRDjmg0D6GrBgrBn8VcQXwVucbeOIYUBDgy1qKJsds4CaPjQo3NuZMQgkRo5CcWy39dmpbWp0N5RXWpJzp1FRxmBb4TRmYFc4gukdHlwZaRhcHBs",
		base64.StdEncoding.EncodeToString(txn.Extract()[0]),
	)

}

func TestMethodName(t *testing.T) {
	require.Equal(
		t,
		MethodName("arc59_optRouterIn(uint64)void"),
		"e8540810",
	)

	require.Equal(
		t,
		MethodName("arc59_sendAsset(axfer,address)address"),
		"2bea37bb",
	)
}
