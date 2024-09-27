package sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

// checks json strings for equality
// inspired by https://gist.github.com/turtlemonvh/e4f7404e28387fadb8ad275a99596f67
func jsonEqual(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

type encodingTest struct {
	msgpack string
	json    string
}

func TestTransaction(t *testing.T) {
	t.Parallel()
	tests := []encodingTest{
		{
			msgpack: "iaRhcGFyiaJhbcQgZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5RmiiYW6sVGVzdCBBc3NldCAyomF1s2h0dHBzOi8vZXhhbXBsZS5jb22hY8QgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhZsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhbcQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhcsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhdM///////////6J1bqRUU1Qyo2ZlZc0D6KJmds4A3/ljo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToiomx2zgDf/Uukbm90ZcQOVGhpcyBpcyBhIG5vdGWjc25kxCC0kna9PsCXfquGoyHESerYAslsC9l8KVYTFRHS8R7r7KR0eXBlpGFjZmc=",
			json: `{"apar": {
				  "am": "ZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5Rmg=",
				  "an": "Test Asset 2",
				  "au": "https://example.com",
				  "c": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				  "f": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				  "m": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				  "r": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				  "t": 18446744073709551615,
				  "un": "TST2"
				},
				"fee": 1000,
				"fv": 14678371,
				"gen": "testnet-v1.0",
				"gh": "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
				"lv": 14679371,
				"note": "VGhpcyBpcyBhIG5vdGU=",
				"snd": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				"type": "acfg"
			  }`,
		},
		{
			msgpack: "i6RhcGFhkcQEdGVzdKRhcGF0ksQgACoyATtqMON+4ohUJO59fQVV6uCTn7aa/GvfndL+5/7EIAAHBAuPYqMysOAF8ALIwKUWNGgBCjFYJ8bPUnx4aXnnpGFwYniSgaFuxAVhbGljZYKhaQKhbsQDYm9ipGFwZmGSzRWzzRoKpGFwaWRko2ZlZc0E0qJmds0jKKJnaMQgMf0h6zjkEIEZPtNM3zsrg+iHQFS0fZxhgr7w35I464OibHbNIzKjc25kxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aR0eXBlpGFwcGw=",
			json: `{
				"apaa": [
				  "dGVzdA=="
				],
				"apat": [
				  "ACoyATtqMON+4ohUJO59fQVV6uCTn7aa/GvfndL+5/4=",
				  "AAcEC49iozKw4AXwAsjApRY0aAEKMVgnxs9SfHhpeec="
				],
				"apbx": [
				  {
					"n": "YWxpY2U="
				  },
				  {
					"i": 2,
					"n": "Ym9i"
				  }
				],
				"apfa": [
				  5555,
				  6666
				],
				"apid": 100,
				"fee": 1234,
				"fv": 9000,
				"gh": "Mf0h6zjkEIEZPtNM3zsrg+iHQFS0fZxhgr7w35I464M=",
				"lv": 9010,
				"snd": "CfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f0=",
				"type": "appl"
			  }`,
		},
		{
			msgpack: "i6RhcGFhkcQEdGVzdKRhcGF0ksQgACoyATtqMON+4ohUJO59fQVV6uCTn7aa/GvfndL+5/7EIAAHBAuPYqMysOAF8ALIwKUWNGgBCjFYJ8bPUnx4aXnnpGFwYniSgaFuxAdjaGFybGllgqFpAaFuxANkb2ekYXBmYZJkzRoKpGFwaWRko2ZlZc0E0qJmds0jKKJnaMQgMf0h6zjkEIEZPtNM3zsrg+iHQFS0fZxhgr7w35I464OibHbNIzKjc25kxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aR0eXBlpGFwcGw=",
			json: `{
				"apaa": [
				  "dGVzdA=="
				],
				"apat": [
				  "ACoyATtqMON+4ohUJO59fQVV6uCTn7aa/GvfndL+5/4=",
				  "AAcEC49iozKw4AXwAsjApRY0aAEKMVgnxs9SfHhpeec="
				],
				"apbx": [
				  {
					"n": "Y2hhcmxpZQ=="
				  },
				  {
					"i": 1,
					"n": "ZG9n"
				  }
				],
				"apfa": [
				  100,
				  6666
				],
				"apid": 100,
				"fee": 1234,
				"fv": 9000,
				"gh": "Mf0h6zjkEIEZPtNM3zsrg+iHQFS0fZxhgr7w35I464M=",
				"lv": 9010,
				"snd": "CfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f0=",
				"type": "appl"
			  }`,
		},
		{
			msgpack: "i6RhcGFhkcQEdGVzdKRhcGF0ksQgACoyATtqMON+4ohUJO59fQVV6uCTn7aa/GvfndL+5/7EIAAHBAuPYqMysOAF8ALIwKUWNGgBCjFYJ8bPUnx4aXnnpGFwYniSgaFuxAKM/4GhbsQCAACkYXBmYZJkzRoKpGFwaWRko2ZlZc0E0qJmds0jKKJnaMQgMf0h6zjkEIEZPtNM3zsrg+iHQFS0fZxhgr7w35I464OibHbNIzKjc25kxCAJ+9J2LAj4bFrmv23Xp6kB3mZ111Dgfoxcdphkfbbh/aR0eXBlpGFwcGw=",
			json: `{
				"apaa": [
				  "dGVzdA=="
				],
				"apat": [
				  "ACoyATtqMON+4ohUJO59fQVV6uCTn7aa/GvfndL+5/4=",
				  "AAcEC49iozKw4AXwAsjApRY0aAEKMVgnxs9SfHhpeec="
				],
				"apbx": [
				  {
					"n": "jP8="
				  },
				  {
					"n": "AAA="
				  }
				],
				"apfa": [
				  100,
				  6666
				],
				"apid": 100,
				"fee": 1234,
				"fv": 9000,
				"gh": "Mf0h6zjkEIEZPtNM3zsrg+iHQFS0fZxhgr7w35I464M=",
				"lv": 9010,
				"snd": "CfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f0=",
				"type": "appl"
			  }`,
		},
		{
			msgpack: "iaRhcGFhkcQEdGVzdKRhcGJ4kYCkYXBpZGSjZmVlzQTSomZ2zSMoomdoxCAx/SHrOOQQgRk+00zfOyuD6IdAVLR9nGGCvvDfkjjrg6Jsds0jMqNzbmTEIAn70nYsCPhsWua/bdenqQHeZnXXUOB+jFx2mGR9tuH9pHR5cGWkYXBwbA==",
			json: `{
				"apaa": [
				  "dGVzdA=="
				],
				"apbx": [
				  {}
				],
				"apid": 100,
				"fee": 1234,
				"fv": 9000,
				"gh": "Mf0h6zjkEIEZPtNM3zsrg+iHQFS0fZxhgr7w35I464M=",
				"lv": 9010,
				"snd": "CfvSdiwI+Gxa5r9t16epAd5mdddQ4H6MXHaYZH224f0=",
				"type": "appl"
			  }`,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("i=%d", i), func(t *testing.T) {
			expectedJson := test.json
			expectedMsgpack, err := base64.StdEncoding.DecodeString(test.msgpack)
			require.NoError(t, err)

			actualJson, err := TransactionMsgpackToJson(expectedMsgpack)
			require.NoError(t, err, "Could not convert transaction from msgpack to JSON")

			areJsonEqual, err := jsonEqual(expectedJson, actualJson)
			require.NoError(t, err)
			require.Truef(t, areJsonEqual, "Expected JSON does not match actual JSON.\nExpected:\n%s\n\nActual:\n%s", expectedJson, actualJson)

			actualMsgpack, err := TransactionJsonToMsgpack(expectedJson)
			require.NoError(t, err, "Could not convert transaction from JSON to msgpack")

			if !bytes.Equal(expectedMsgpack, actualMsgpack) {
				b64Expected := base64.StdEncoding.EncodeToString(expectedMsgpack)
				b64Actual := base64.StdEncoding.EncodeToString(actualMsgpack)
				require.Failf(t, "Expected msgpack does not match actual msgpack.\nExpected:\n%s\n\nActual:\n%s", b64Expected, b64Actual)
			}
		})
	}
}
