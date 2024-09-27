package sdk

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/stretchr/testify/require"
)

func TestMakePaymentTxn(t *testing.T) {
	t.Parallel()
	fromAddress := "47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU"
	toAddress := "PNWOET7LLOWMBMLE4KOCELCX6X3D3Q4H2Q4QJASYIEOF7YIPPQBG3YQ5YI"
	params := SuggestedParams{
		Fee:             4,
		FirstRoundValid: 12466,
		LastRoundValid:  13466,
		GenesisID:       "devnet-v33.0",
		GenesisHash:     mustDecodeB64(t, "JgsgCaCTqIaLeVhyL6XlRu3n7Rfk2FxMeK+wRSaQ7dI="),
	}
	amount := MakeUint64(1000)
	encodedTx, err := MakePaymentTxn(fromAddress, toAddress, &amount, mustDecodeB64(t, "6gAVR0Nsv5Y="), "IDUTJEUIEVSMXTU4LGTJWZ2UE2E6TIODUKU6UW3FU3UKIQQ77RLUBBBFLA", &params)
	require.NoError(t, err)

	expectedEncodedTx := mustDecodeB64(t, "i6NhbXTNA+ilY2xvc2XEIEDpNJKIJWTLzpxZpptnVCaJ6aHDoqnqW2Wm6KRCH/xXo2ZlZc0EmKJmds0wsqNnZW6sZGV2bmV0LXYzMy4womdoxCAmCyAJoJOohot5WHIvpeVG7eftF+TYXEx4r7BFJpDt0qJsds00mqRub3RlxAjqABVHQ2y/lqNyY3bEIHts4k/rW6zAsWTinCIsV/X2PcOH1DkEglhBHF/hD3wCo3NuZMQg5/D4TQaBHfnzHI2HixFV9GcdUaGFwgCQhmf0SVhwaKGkdHlwZaNwYXk=")
	require.Equal(t, expectedEncodedTx, encodedTx)
}

func TestMakeRekeyTxn(t *testing.T) {
	t.Parallel()
	fromAddress := "47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU"
	rekeyToAddress := "PNWOET7LLOWMBMLE4KOCELCX6X3D3Q4H2Q4QJASYIEOF7YIPPQBG3YQ5YI"
	params := SuggestedParams{
		Fee:             4,
		FirstRoundValid: 12466,
		LastRoundValid:  13466,
		GenesisID:       "devnet-v33.0",
		GenesisHash:     mustDecodeB64(t, "JgsgCaCTqIaLeVhyL6XlRu3n7Rfk2FxMeK+wRSaQ7dI="),
	}
	encodedTx, err := MakeRekeyTxn(fromAddress, rekeyToAddress, &params)
	require.NoError(t, err)

	expectedEncodedTx := mustDecodeB64(t, "iaNmZWXNA+iiZnbNMLKjZ2VurGRldm5ldC12MzMuMKJnaMQgJgsgCaCTqIaLeVhyL6XlRu3n7Rfk2FxMeK+wRSaQ7dKibHbNNJqjcmN2xCDn8PhNBoEd+fMcjYeLEVX0Zx1RoYXCAJCGZ/RJWHBooaVyZWtlecQge2ziT+tbrMCxZOKcIixX9fY9w4fUOQSCWEEcX+EPfAKjc25kxCDn8PhNBoEd+fMcjYeLEVX0Zx1RoYXCAJCGZ/RJWHBooaR0eXBlo3BheQ==")
	require.Equal(t, expectedEncodedTx, encodedTx)
}

func TestMakeApplicationCreateTx(t *testing.T) {
	t.Parallel()
	params := SuggestedParams{
		FlatFee:         true,
		Fee:             1000,
		FirstRoundValid: 2063137,
		LastRoundValid:  2064137,
		GenesisID:       "devnet-v1.0",
		GenesisHash:     mustDecodeB64(t, "sC3P7e2SdbqKJK0tbiCdK9tdSpbe6XeCGKdoNzmlj0E="),
	}
	note := mustDecodeB64(t, "8xMCTuLQ810=")
	program := []byte{1, 32, 1, 1, 34}
	args := make([][]byte, 2)
	args[0] = []byte("123")
	args[1] = []byte("456")
	foreignApps := make([]int64, 1)
	foreignApps[0] = 10
	foreignAssets := foreignApps
	gSchema := types.StateSchema{NumUint: uint64(1), NumByteSlice: uint64(1)}
	lSchema := types.StateSchema{NumUint: uint64(1), NumByteSlice: uint64(1)}
	extraPages := int32(2)
	addr := make([]string, 1)
	addr[0] = "47YPQTIGQEO7T4Y4RWDYWEKV6RTR2UNBQXBABEEGM72ESWDQNCQ52OPASU"
	boxReferences := make([]types.AppBoxReference, 3)
	boxReferences[0] = types.AppBoxReference{AppID: 0, Name: []byte("box_name")}
	boxReferences[1] = types.AppBoxReference{AppID: 10, Name: []byte("box_name")}
	boxReferences[2] = types.AppBoxReference{AppID: 10, Name: []byte("box_name2")}

	encodedTx, err := MakeApplicationCreateTx(
		false,
		program,
		program,
		int64(gSchema.NumUint),
		int64(gSchema.NumByteSlice),
		int64(lSchema.NumUint),
		int64(lSchema.NumByteSlice),
		extraPages,
		&BytesArray{args},
		&StringArray{addr},
		&Int64Array{foreignApps},
		&Int64Array{foreignAssets},
		&AppBoxRefArray{boxReferences},
		&params,
		types.ZeroAddress.String(),
		note,
	)
	require.NoError(t, err)

	expectedEncodedTx := mustDecodeB64(t, "3gARpGFwYWGSxAMxMjPEAzQ1NqRhcGFwxAUBIAEBIqRhcGFzkQqkYXBhdJHEIOfw+E0GgR358xyNh4sRVfRnHVGhhcIAkIZn9ElYcGihpGFwYniTgaFuxAhib3hfbmFtZYKhaQGhbsQIYm94X25hbWWCoWkBoW7ECWJveF9uYW1lMqRhcGVwAqRhcGZhkQqkYXBnc4KjbmJzAaNudWkBpGFwbHOCo25icwGjbnVpAaRhcHN1xAUBIAEBIqNmZWXNA+iiZnbOAB97IaNnZW6rZGV2bmV0LXYxLjCiZ2jEILAtz+3tknW6iiStLW4gnSvbXUqW3ul3ghinaDc5pY9Bomx2zgAffwmkbm90ZcQI8xMCTuLQ812kdHlwZaRhcHBs")
	require.Equal(t, expectedEncodedTx, encodedTx)
}
