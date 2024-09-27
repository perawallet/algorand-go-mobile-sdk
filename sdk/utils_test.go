package sdk

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/stretchr/testify/require"
)

func TestUint64(t *testing.T) {
	t.Parallel()
	tests := []uint64{
		0,
		1,
		math.MaxUint32,
		math.MaxUint32 + 1,
		math.MaxUint64,
	}

	for _, test := range tests {
		t.Run(fmt.Sprint(test), func(t *testing.T) {
			value := MakeUint64(test)

			extracted, err := value.Extract()
			if err != nil {
				t.Fatal(err)
			}

			if test != extracted {
				t.Errorf("Wrong exacted value. Expected %d, got %d", test, extracted)
			}
		})
	}
}

func randomBytes(s []byte) {
	_, err := rand.Read(s)
	if err != nil {
		panic(err)
	}
}

func TestEncodeDecode(t *testing.T) {
	t.Parallel()
	a := types.Address{}
	for i := 0; i < 1000; i++ {
		randomBytes(a[:])
		addr := a.String()
		b, err := types.DecodeAddress(addr)
		require.NoError(t, err)
		require.Equal(t, a, b)
		require.True(t, IsValidAddress(a.String()))

		require.False(t, IsValidAddress("SONNPE7I3TYWE7VQQA7VCGZK54WXEICYROMWQYXCOB5A5RHM46LVNGZLRU"))

	}
}
