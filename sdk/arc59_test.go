package sdk

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/stretchr/testify/require"
)

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
	algoAmount := MakeUint64(20)

	// test with arc 59 false
	txnsByteArray, err := MakeAndSignARC59SendTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"MKKKFL5JBJTOCEMEZUAJKTWD5FYAI2FOLW5BP5N5YR37ZG5FHLTUYCFC6U",
		"YIIC6GF4DUJYZTYTZ5UEOAXONUUKZRDFOTV4EKSGD5E7BYE6EE3IVPYEDQ",
		"6LD3JUWPR72DX5JNGPHH2QEM2IKRA3SXGFSRR4P7TB5JMVPQTIYCXLUMCQ",
		&amount,
		&min_balance_requirement,
		5,
		643020148,
		655977010,
		&suggested_params,
		false,
		&algoAmount,
		sk,
	)
	require.NoError(t, err)
	txns := txnsByteArray.Extract()
	require.Equal(t, len(txns), 4)
	require.Equal(
		t,
		"gqNzaWfEQH6SBWU7qw/2IImisTh6tdiJvvBbfNiwDmw4KXzowkMhB1naf/MVee7D1dAO6FfkBRVZTPUqBvZt1ENOAZoX+gyjdHhuiqNhbXTOAAN7GKNmZWXNA+iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIEH2wbyJwlTR8QvNTIBsrUopXCnUTeA/FO0UKNHAX057omx2zgJo+NCjcmN2xCDCEC8YvB0TjM8Tz2hHAu5tKKzEZXTrwipGH0nw4J4hNqNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWjcGF5",
		base64.StdEncoding.EncodeToString(txns[0]),
	)
	require.Equal(
		t,
		"gqNzaWfEQC8Navcft57gM4gOSofCc8/3MrY6eo/SIOrmzWyHDvrAvyrt1EAuzxwgagDPSaYRkxulL7fAUt65kxEqbT7uzgqjdHhujKRhcGFhksQE6FQIEMQIAAAAACcZajKkYXBhc5HOJxlqMqRhcGF0kcQgwhAvGLwdE4zPE89oRwLubSisxGV068IqRh9J8OCeITakYXBpZM4mU7V0o2ZlZc0H0KJmds4CaPToo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQgQfbBvInCVNHxC81MgGytSilcKdRN4D8U7RQo0cBfTnuibHbOAmj40KNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWkYXBwbA==",
		base64.StdEncoding.EncodeToString(txns[1]),
	)
	require.Equal(
		t,
		"gqNzaWfEQDnGRnANVgGFmQxtVb0ymPI1JTgFiFUvJxQ9C8Wu+jRvQMB4B+tJp6kKK+M5Gay8AGcZdI7qNWC94PYgc4R2wAqjdHhui6RhYW10CqRhcmN2xCDCEC8YvB0TjM8Tz2hHAu5tKKzEZXTrwipGH0nw4J4hNqNmZWXNA+iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIEH2wbyJwlTR8QvNTIBsrUopXCnUTeA/FO0UKNHAX057omx2zgJo+NCjc25kxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6R0eXBlpWF4ZmVypHhhaWTOJxlqMg==",
		base64.StdEncoding.EncodeToString(txns[2]),
	)
	require.Equal(
		t,
		"gqNzaWfEQKE+AQk9falzkBIpzhQ4yycSXnIq2yzETYnuJPJ6UqFooLCRL4qPBFdwEKvAQEA3ZEYawURdCRB9eYtEKwdx5gqjdHhujaRhcGFhk8QECFMe18QgYpSir6kKZuERhM0AlU7D6XAEaK5duhf1vcR3/JulOufECAAAAAAAAAAUpGFwYXORzicZajKkYXBhdJLEIGKUoq+pCmbhEYTNAJVOw+lwBGiuXboX9b3Ed/ybpTrnxCDyx7TSz4/0O/UtM859QIzSFRBuVzFlGPH/mHqWVfCaMKRhcGJ4kYGhbsQgYpSir6kKZuERhM0AlU7D6XAEaK5duhf1vcR3/JulOuekYXBpZM4mU7V0o2ZlZc0bWKJmds4CaPToo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQgQfbBvInCVNHxC81MgGytSilcKdRN4D8U7RQo0cBfTnuibHbOAmj40KNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWkYXBwbA==",
		base64.StdEncoding.EncodeToString(txns[3]),
	)

	// test with arc 59 true
	txnsByteArray, err = MakeAndSignARC59SendTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"MKKKFL5JBJTOCEMEZUAJKTWD5FYAI2FOLW5BP5N5YR37ZG5FHLTUYCFC6U",
		"YIIC6GF4DUJYZTYTZ5UEOAXONUUKZRDFOTV4EKSGD5E7BYE6EE3IVPYEDQ",
		"6LD3JUWPR72DX5JNGPHH2QEM2IKRA3SXGFSRR4P7TB5JMVPQTIYCXLUMCQ",
		&amount,
		&min_balance_requirement,
		5,
		643020148,
		655977010,
		&suggested_params,
		true,
		&algoAmount,
		sk,
	)
	require.NoError(t, err)
	txns = txnsByteArray.Extract()
	require.Equal(t, len(txns), 3)
	require.Equal(
		t,
		"gqNzaWfEQFWCZAqrBkIsumkX/ZETn4cq/ECZFyW/658Jhkwu7vmMKE3X7MWYk+uIJIF9+3J9fGvH+slHnjR+vTcqXMzAOwGjdHhuiqNhbXTOAAN7GKNmZWXNA+iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIOnT3buWnGrbSyyOJnIQQ9K2XCELvgTYUubOZbDoSEJOomx2zgJo+NCjcmN2xCDCEC8YvB0TjM8Tz2hHAu5tKKzEZXTrwipGH0nw4J4hNqNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWjcGF5",
		base64.StdEncoding.EncodeToString(txns[0]),
	)
	require.Equal(
		t,
		"gqNzaWfEQP6D1SAZX4981f+ylVf/VUKNcRv6423AVRppPBKDw/AagAm7EyXjM/83iK+suvkA3RHXWf+XuzytIVcUlnPfagKjdHhui6RhYW10CqRhcmN2xCDCEC8YvB0TjM8Tz2hHAu5tKKzEZXTrwipGH0nw4J4hNqNmZWXNA+iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIOnT3buWnGrbSyyOJnIQQ9K2XCELvgTYUubOZbDoSEJOomx2zgJo+NCjc25kxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6R0eXBlpWF4ZmVypHhhaWTOJxlqMg==",
		base64.StdEncoding.EncodeToString(txns[1]),
	)
	require.Equal(
		t,
		"gqNzaWfEQLHMQiQreD8rXvu0fJJ64KUIY/RYzDv/07EVEGiuhfyz2nFI6O2U9QS4s4xEGAXS5FwdxfcOf0alZzZS4z7guwSjdHhujaRhcGFhk8QECFMe18QgYpSir6kKZuERhM0AlU7D6XAEaK5duhf1vcR3/JulOufECAAAAAAAAAAUpGFwYXORzicZajKkYXBhdJLEIGKUoq+pCmbhEYTNAJVOw+lwBGiuXboX9b3Ed/ybpTrnxCDyx7TSz4/0O/UtM859QIzSFRBuVzFlGPH/mHqWVfCaMKRhcGJ4kYGhbsQgYpSir6kKZuERhM0AlU7D6XAEaK5duhf1vcR3/JulOuekYXBpZM4mU7V0o2ZlZc0bWKJmds4CaPToo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQg6dPdu5acattLLI4mchBD0rZcIQu+BNhS5s5lsOhIQk6ibHbOAmj40KNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWkYXBwbA==",
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

	generated_txs, err := MakeAndSignARC59ClaimTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"6LD3JUWPR72DX5JNGPHH2QEM2IKRA3SXGFSRR4P7TB5JMVPQTIYCXLUMCQ",
		643020148,
		655977010,
		&suggested_params,
		true,
		false,
		sk,
	)
	require.NoError(t, err)
	txs := generated_txs.Extract()
	require.Equal(t, len(txs), 1)
	require.Equal(
		t,
		"gqNzaWfEQIj+/2PrjZ6UBzM8C4SKcDxx0LZ5W7ks7Y83XgJkRxdGrAajL3q1fbAbE40bxc1s0b+mTaqb+0h1bqbri9bIBQyjdHhujaRhcGFhksQEv5AuPMQIAAAAACcZajKkYXBhc5HOJxlqMqRhcGF0kcQg8se00s+P9Dv1LTPOfUCM0hUQblcxZRjx/5h6llXwmjCkYXBieJGBoW7EIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpGFwaWTOJlO1dKNmZWXNC7iiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEID5DmNnDlKg3YPonhwOMQMEL1Anae17u80wk6i8UoFk/omx2zgJo+NCjc25kxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6R0eXBlpGFwcGw=",
		base64.StdEncoding.EncodeToString(txs[0]),
	)
	require.NoError(t, err)

	// test with opt-in false

	generated_txns_2, err := MakeAndSignARC59ClaimTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"6LD3JUWPR72DX5JNGPHH2QEM2IKRA3SXGFSRR4P7TB5JMVPQTIYCXLUMCQ",
		643020148,
		655977010,
		&suggested_params,
		false,
		false,
		sk,
	)
	require.NoError(t, err)
	txs_2 := generated_txns_2.Extract()
	require.Equal(t, len(txs_2), 2)

	require.Equal(
		t,
		"gqNzaWfEQEDGFTJf/se2k41PE5BydfzRqKgNd4JRufCme20jlKhTH+z/Sxc2q/9+jS+lMfI+6SgC/GjEEQ6apsa2tLEcEA6jdHhuiaRhcmN2xCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6Jmds4CaPToo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQgqbROyTlncTEZsLLw33xsoreRE2i2O2cA3NhZOFl5OoeibHbOAmj40KNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWlYXhmZXKkeGFpZM4nGWoy",
		base64.StdEncoding.EncodeToString(txs_2[0]),
	)
	require.Equal(
		t,
		"gqNzaWfEQJVgdG7cF32LMiZzgMu3eaV1ExRKWu3bHz+zuXCSZGAVKREDekesgNOqORoSSGhJaHm68SM4QiUZ6gdonzzagg+jdHhujaRhcGFhksQEv5AuPMQIAAAAACcZajKkYXBhc5HOJxlqMqRhcGF0kcQg8se00s+P9Dv1LTPOfUCM0hUQblcxZRjx/5h6llXwmjCkYXBieJGBoW7EIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpGFwaWTOJlO1dKNmZWXND6CiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIKm0Tsk5Z3ExGbCy8N98bKK3kRNotjtnANzYWThZeTqHomx2zgJo+NCjc25kxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6R0eXBlpGFwcGw=",
		base64.StdEncoding.EncodeToString(txs_2[1]),
	)

	// test with claim true

	generated_txns_3, err := MakeAndSignARC59ClaimTxn(
		"SENDSCOFWLP5OZVFWWU5BXSRLVVETTU5IVDRTALPQTIZTAK44IF2SJ57P4",
		"6LD3JUWPR72DX5JNGPHH2QEM2IKRA3SXGFSRR4P7TB5JMVPQTIYCXLUMCQ",
		643020148,
		655977010,
		&suggested_params,
		false,
		true,
		sk,
	)
	require.NoError(t, err)
	txs_3 := generated_txns_3.Extract()
	require.Equal(t, len(txs_3), 3)

	fmt.Println(base64.StdEncoding.EncodeToString(txs_3[0]))
	fmt.Println(base64.StdEncoding.EncodeToString(txs_3[1]))
	fmt.Println(base64.StdEncoding.EncodeToString(txs_3[2]))

	require.Equal(
		t,
		"gqNzaWfEQEoLwKPyWIBtOQ5wXZ8WJlQ9/wMD9rixH3V5NwuvFk2v95VPvAJcO77D7yfdXXsK6emBZWJq5RSXNZ9IMb53ygSjdHhui6RhcGFhkcQENi3K16RhcGF0kcQg8se00s+P9Dv1LTPOfUCM0hUQblcxZRjx/5h6llXwmjCkYXBieJGBoW7EIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpGFwaWTOJlO1dKJmds4CaPToo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQgdjURnaVRpNL/JxD3ExEpJgm6hVcpN5yLr55KcGNk/c+ibHbOAmj40KNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWkYXBwbA==",
		base64.StdEncoding.EncodeToString(txs_3[0]),
	)
	require.Equal(
		t,
		"gqNzaWfEQPXWP4OUs/IEpADudMkP8vIjFE3RZxeA2d/9Z64vK5CoojSTa76GRbYZXp0P+E7TA294nJcCaHlibTMU1v5pRg6jdHhuiaRhcmN2xCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6Jmds4CaPToo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQgdjURnaVRpNL/JxD3ExEpJgm6hVcpN5yLr55KcGNk/c+ibHbOAmj40KNzbmTEIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpHR5cGWlYXhmZXKkeGFpZM4nGWoy",
		base64.StdEncoding.EncodeToString(txs_3[1]),
	)
	require.Equal(
		t,
		"gqNzaWfEQPl0cDSIcbF0hhdocLiPuafwxdTLBYCoBdSCY43mwFG8zLtrYz1Js52KFE5aUVMgMuhQn4a4jDTDkDKaj69J2gKjdHhujaRhcGFhksQEv5AuPMQIAAAAACcZajKkYXBhc5HOJxlqMqRhcGF0kcQg8se00s+P9Dv1LTPOfUCM0hUQblcxZRjx/5h6llXwmjCkYXBieJGBoW7EIJEaOQnFst/XZqW1qdDeUV1qSc6dRUcZgW+E0ZmBXOILpGFwaWTOJlO1dKNmZWXNF3CiZnbOAmj06KNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIHY1EZ2lUaTS/ycQ9xMRKSYJuoVXKTeci6+eSnBjZP3Pomx2zgJo+NCjc25kxCCRGjkJxbLf12altanQ3lFdaknOnUVHGYFvhNGZgVziC6R0eXBlpGFwcGw=",
		base64.StdEncoding.EncodeToString(txs_3[2]),
	)
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
