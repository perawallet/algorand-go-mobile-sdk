package sdk

import (
	"encoding/base64"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
)

func mustDecodeAddress(t *testing.T, addr string) types.Address {
	t.Helper()
	decoded, err := types.DecodeAddress(addr)
	require.NoError(t, err)
	return decoded
}

func mustDecodeB64(t *testing.T, b64 string) []byte {
	decoded, err := base64.StdEncoding.DecodeString(b64)
	require.NoError(t, err)
	return decoded
}

func TestMultisigAccount(t *testing.T) {
	t.Parallel()
	addrs := []string{
		"2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY",
		"S64XU5HQEY2XLHVUSO6RI3JL6NHC32I4LJHM32ZOM5VC4QPON7BZZRCU2E",
		"W3KCADJF23RDTO3TMY63YQBKYDYFPHFBU75JQMX5QHOERRBOZ75L3B2J7Y",
	}

	account, err := MakeMultisigAccount(1, 2, &StringArray{addrs})
	require.NoError(t, err)

	expectedPks := make([]ed25519.PublicKey, len(addrs))
	for i, addr := range addrs {
		decoded := mustDecodeAddress(t, addr)
		expectedPks[i] = decoded[:]
	}
	expectedMA := crypto.MultisigAccount{
		Version:   1,
		Threshold: 2,
		Pks:       expectedPks,
	}
	require.Equal(t, expectedMA, account.value)

	t.Run("version", func(t *testing.T) {
		require.Equal(t, 1, account.Version())
	})

	t.Run("threshold", func(t *testing.T) {
		require.Equal(t, 2, account.Threshold())
	})

	t.Run("address", func(t *testing.T) {
		addr, err := account.Address()
		require.NoError(t, err)
		expectedAddr := "KXWFRJX453UVTXPRR2APX5EQQN3JQJVCW2E6AAGSXFHPINV353CETLQCYA"
		require.Equal(t, expectedAddr, addr)
	})

	t.Run("contributing addresses", func(t *testing.T) {
		contributingAddrs := account.ContributingAddresses()
		require.Equal(t, addrs, contributingAddrs.Extract())
	})
}

func TestExtractMultisigAccountFromSignedTransaction(t *testing.T) {
	t.Parallel()
	t.Run("valid", func(t *testing.T) {
		encodedStx := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAL+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgqJwa8QgttQgDSXW4jm7c2Y9vEAqwPBXnKGn+pgy/YHcSMQuz/qhc8RAiIsrT+BmBMjqCo/Hkuq1/NnmHcUzZiXYRyOtHsPxXi6dJo8m16N00iU//Pd65QbZYpqc6DMl+OsgBpMm0WydCqN0aHICoXYBpHNnbnLEIFXsWKb87ulZ3fGOgPv0kIN2mCaitongANK5TvQ2u+7Eo3R4bomjYW10zgAPQkCjZmVlzQPoomZ2AqNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds0D6qNyY3bEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/Do3NuZMQg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKkdHlwZaNwYXk=")
		account, err := ExtractMultisigAccountFromSignedTransaction(encodedStx)
		require.NoError(t, err)

		expectedAddrs := []string{
			"2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY",
			"S64XU5HQEY2XLHVUSO6RI3JL6NHC32I4LJHM32ZOM5VC4QPON7BZZRCU2E",
			"W3KCADJF23RDTO3TMY63YQBKYDYFPHFBU75JQMX5QHOERRBOZ75L3B2J7Y",
		}
		expectedPks := make([]ed25519.PublicKey, len(expectedAddrs))
		for i, addr := range expectedAddrs {
			decoded := mustDecodeAddress(t, addr)
			expectedPks[i] = decoded[:]
		}
		expectedMA := crypto.MultisigAccount{
			Version:   1,
			Threshold: 2,
			Pks:       expectedPks,
		}
		require.Equal(t, expectedMA, account.value)
	})

	t.Run("no msig", func(t *testing.T) {
		encodedStx := mustDecodeB64(t, "gqNzaWfEQC/nu0j+joowxkE3uMMx3SPZFzOHq8YdeTjqYVgV5r6xL/w4wRaYKhZXSVbeTnr2udAsbcxWUm7mhYTzH7AeEwejdHhuiaNhbXTOAA9CQKNmZWXNA+iiZnYCo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToiomx2zQPqo3JjdsQgl7l6dPAmNXWetJO9FG0r804t6RxaTs3rLmdqLkHub8Ojc25kxCDUYfSDPMXT2TQv+8ZvQtqEoK0G1xsTfW1wkIkR+AfA4qR0eXBlo3BheQ==")
		account, err := ExtractMultisigAccountFromSignedTransaction(encodedStx)
		require.NoError(t, err)
		require.Nil(t, account)
	})
}

func TestSignMultisigTransaction(t *testing.T) {
	t.Parallel()
	// corresponds to 2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY
	sk, err := mnemonic.ToPrivateKey("carbon another pair valley ride lumber exhibit chunk forget select nerve topic refuse ball bomb draw chunk toward motor detect process smile envelope abstract rule")
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

	multisigAccount, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
		mustDecodeAddress(t, "2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY"),
		mustDecodeAddress(t, "S64XU5HQEY2XLHVUSO6RI3JL6NHC32I4LJHM32ZOM5VC4QPON7BZZRCU2E"),
		mustDecodeAddress(t, "W3KCADJF23RDTO3TMY63YQBKYDYFPHFBU75JQMX5QHOERRBOZ75L3B2J7Y"),
	})
	require.NoError(t, err)

	stxBytes, err := SignMultisigTransaction(sk, &MultisigAccount{multisigAccount}, encodedTxn)
	require.NoError(t, err)

	expectedStxBytes := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAL+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgaJwa8QgttQgDSXW4jm7c2Y9vEAqwPBXnKGn+pgy/YHcSMQuz/qjdGhyAqF2AaRzZ25yxCBV7Fim/O7pWd3xjoD79JCDdpgmoraJ4ADSuU70NrvuxKN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
	require.Equal(t, expectedStxBytes, stxBytes)
}

func TestAttachMultisigSignature(t *testing.T) {
	t.Parallel()
	signerAddress := "2RQ7JAZ4YXJ5SNBP7PDG6QW2QSQK2BWXDMJX23LQSCERD6AHYDRH4N4MXY"
	signerSig := mustDecodeB64(t, "L+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TBw==")

	params := types.SuggestedParams{
		Fee:             0,
		GenesisID:       "testnet-v1.0",
		GenesisHash:     mustDecodeB64(t, "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI="),
		FirstRoundValid: 2,
		LastRoundValid:  1002,
	}
	txn, err := transaction.MakePaymentTxn(signerAddress, "S64XU5HQEY2XLHVUSO6RI3JL6NHC32I4LJHM32ZOM5VC4QPON7BZZRCU2E", 1_000_000, nil, "", params)
	require.NoError(t, err)
	encodedTxn := msgpack.Encode(&txn)

	multisigAccount, err := crypto.MultisigAccountWithParams(1, 2, []types.Address{
		mustDecodeAddress(t, signerAddress),
		mustDecodeAddress(t, "S64XU5HQEY2XLHVUSO6RI3JL6NHC32I4LJHM32ZOM5VC4QPON7BZZRCU2E"),
		mustDecodeAddress(t, "W3KCADJF23RDTO3TMY63YQBKYDYFPHFBU75JQMX5QHOERRBOZ75L3B2J7Y"),
	})
	require.NoError(t, err)

	stxBytes, err := AttachMultisigSignature(signerAddress, signerSig, &MultisigAccount{multisigAccount}, encodedTxn)
	require.NoError(t, err)

	expectedStxBytes := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAL+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgaJwa8QgttQgDSXW4jm7c2Y9vEAqwPBXnKGn+pgy/YHcSMQuz/qjdGhyAqF2AaRzZ25yxCBV7Fim/O7pWd3xjoD79JCDdpgmoraJ4ADSuU70NrvuxKN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
	require.Equal(t, expectedStxBytes, stxBytes)
}

func TestMergeMultisigTransactions(t *testing.T) {
	t.Parallel()
	t.Run("no overlap", func(t *testing.T) {
		encodedStxSig1 := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAL+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgaJwa8QgttQgDSXW4jm7c2Y9vEAqwPBXnKGn+pgy/YHcSMQuz/qjdGhyAqF2AaRzZ25yxCBV7Fim/O7pWd3xjoD79JCDdpgmoraJ4ADSuU70NrvuxKN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
		encodedStxSig3 := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgaJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKBonBrxCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw4KicGvEILbUIA0l1uI5u3NmPbxAKsDwV5yhp/qYMv2B3EjELs/6oXPEQIiLK0/gZgTI6gqPx5LqtfzZ5h3FM2Yl2EcjrR7D8V4unSaPJtejdNIlP/z3euUG2WKanOgzJfjrIAaTJtFsnQqjdGhyAqF2AaRzZ25yxCBV7Fim/O7pWd3xjoD79JCDdpgmoraJ4ADSuU70NrvuxKN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
		mergedStx, err := MergeMultisigTransactions(encodedStxSig1, encodedStxSig3)
		require.NoError(t, err)
		encodedStxSig1And3 := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAL+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgqJwa8QgttQgDSXW4jm7c2Y9vEAqwPBXnKGn+pgy/YHcSMQuz/qhc8RAiIsrT+BmBMjqCo/Hkuq1/NnmHcUzZiXYRyOtHsPxXi6dJo8m16N00iU//Pd65QbZYpqc6DMl+OsgBpMm0WydCqN0aHICoXYBpHNnbnLEIFXsWKb87ulZ3fGOgPv0kIN2mCaitongANK5TvQ2u+7Eo3R4bomjYW10zgAPQkCjZmVlzQPoomZ2AqNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds0D6qNyY3bEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/Do3NuZMQg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKkdHlwZaNwYXk=")
		require.Equal(t, encodedStxSig1And3, mergedStx)
	})
	t.Run("overlap", func(t *testing.T) {
		encodedStxSig1 := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAL+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgaJwa8QgttQgDSXW4jm7c2Y9vEAqwPBXnKGn+pgy/YHcSMQuz/qjdGhyAqF2AaRzZ25yxCBV7Fim/O7pWd3xjoD79JCDdpgmoraJ4ADSuU70NrvuxKN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
		encodedStxSig1And3 := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAL+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgqJwa8QgttQgDSXW4jm7c2Y9vEAqwPBXnKGn+pgy/YHcSMQuz/qhc8RAiIsrT+BmBMjqCo/Hkuq1/NnmHcUzZiXYRyOtHsPxXi6dJo8m16N00iU//Pd65QbZYpqc6DMl+OsgBpMm0WydCqN0aHICoXYBpHNnbnLEIFXsWKb87ulZ3fGOgPv0kIN2mCaitongANK5TvQ2u+7Eo3R4bomjYW10zgAPQkCjZmVlzQPoomZ2AqNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds0D6qNyY3bEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/Do3NuZMQg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKkdHlwZaNwYXk=")
		mergedStx, err := MergeMultisigTransactions(encodedStxSig1, encodedStxSig1And3)
		require.NoError(t, err)
		require.Equal(t, encodedStxSig1And3, mergedStx)
	})
	t.Run("overlap with different signatures", func(t *testing.T) {
		// sig here differs from the one in encodedStxSig1And3
		encodedStxSig1 := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAQ+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgaJwa8QgttQgDSXW4jm7c2Y9vEAqwPBXnKGn+pgy/YHcSMQuz/qjdGhyAqF2AaRzZ25yxCBV7Fim/O7pWd3xjoD79JCDdpgmoraJ4ADSuU70NrvuxKN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
		encodedStxSig1And3 := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAL+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgqJwa8QgttQgDSXW4jm7c2Y9vEAqwPBXnKGn+pgy/YHcSMQuz/qhc8RAiIsrT+BmBMjqCo/Hkuq1/NnmHcUzZiXYRyOtHsPxXi6dJo8m16N00iU//Pd65QbZYpqc6DMl+OsgBpMm0WydCqN0aHICoXYBpHNnbnLEIFXsWKb87ulZ3fGOgPv0kIN2mCaitongANK5TvQ2u+7Eo3R4bomjYW10zgAPQkCjZmVlzQPoomZ2AqNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds0D6qNyY3bEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/Do3NuZMQg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKkdHlwZaNwYXk=")
		_, err := MergeMultisigTransactions(encodedStxSig1, encodedStxSig1And3)
		require.ErrorContains(t, err, "mismatched duplicate signatures")
	})
	t.Run("different msigs", func(t *testing.T) {
		encodedStxMsig1 := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgaJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKBonBrxCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw4KicGvEILbUIA0l1uI5u3NmPbxAKsDwV5yhp/qYMv2B3EjELs/6oXPEQIiLK0/gZgTI6gqPx5LqtfzZ5h3FM2Yl2EcjrR7D8V4unSaPJtejdNIlP/z3euUG2WKanOgzJfjrIAaTJtFsnQqjdGhyAqF2AaRzZ25yxCBV7Fim/O7pWd3xjoD79JCDdpgmoraJ4ADSuU70NrvuxKN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
		encodedStxMsig2 := mustDecodeB64(t, "g6Rtc2lng6ZzdWJzaWeTgqJwa8Qg1GH0gzzF09k0L/vGb0LahKCtBtcbE31tcJCJEfgHwOKhc8RAL+e7SP6OijDGQTe4wzHdI9kXM4erxh15OOphWBXmvrEv/DjBFpgqFldJVt5Oeva50CxtzFZSbuaFhPMfsB4TB4GicGvEIJe5enTwJjV1nrSTvRRtK/NOLekcWk7N6y5nai5B7m/DgaJwa8Qg7C0NjeEq5un7i63DFepDX8kkj1acaRia6umO9DnJgIijdGhyAqF2AaRzZ25yxCAG82xYho8S8l9dBd1U/FY5osUnPuM36f6uJztT62jZ2aN0eG6Jo2FtdM4AD0JAo2ZlZc0D6KJmdgKjZ2VurHRlc3RuZXQtdjEuMKJnaMQgSGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiKibHbNA+qjcmN2xCCXuXp08CY1dZ60k70UbSvzTi3pHFpOzesuZ2ouQe5vw6NzbmTEINRh9IM8xdPZNC/7xm9C2oSgrQbXGxN9bXCQiRH4B8DipHR5cGWjcGF5")
		_, err := MergeMultisigTransactions(encodedStxMsig1, encodedStxMsig2)
		require.ErrorContains(t, err, "multisig parameters do not match")
	})
}
