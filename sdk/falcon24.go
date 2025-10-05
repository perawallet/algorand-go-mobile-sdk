package sdk

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/algorandfoundation/falcon-signatures/algorand"
	"github.com/algorandfoundation/falcon-signatures/falcongo"
	"golang.org/x/crypto/pbkdf2"
)

type AlgorandKeyInfo struct {
	AlgorandAddress string
	PublicKey       string
	PrivateKey      string
}

func DeriveFromBIP39(mnemonic string) (*AlgorandKeyInfo, error) {
	mnemonic = strings.TrimSpace(mnemonic)

	seed := deriveSeedFromMnemonic(mnemonic)

	kp, err := falcongo.GenerateKeyPair(seed)
	if err != nil {
		return nil, err
	}

	address, err := algorand.GetAddressFromPublicKey(kp.PublicKey)
	if err != nil {
		return nil, err
	}

	return &AlgorandKeyInfo{
		AlgorandAddress: string(address),
		PublicKey:       strings.ToLower(hex.EncodeToString(kp.PublicKey[:])),
		PrivateKey:      strings.ToLower(hex.EncodeToString(kp.PrivateKey[:])),
	}, nil
}

func GetAlgorandAddress(mnemonic string) (string, error) {
	info, err := DeriveFromBIP39(mnemonic)
	if err != nil {
		return "", err
	}
	return info.AlgorandAddress, nil
}

func GetPublicKey(mnemonic string) (string, error) {
	info, err := DeriveFromBIP39(mnemonic)
	if err != nil {
		return "", err
	}
	return info.PublicKey, nil
}

func (ki *AlgorandKeyInfo) ToJSON() (string, error) {
	data, err := json.MarshalIndent(ki, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

const (
	kdfIterations = 100000
	kdfKeyLen     = 48
	kdfSaltStr    = "falcon-cli-seed-v1"
)

func deriveSeedFromMnemonic(mnemonic string) []byte {
	return pbkdf2.Key([]byte(mnemonic), []byte(kdfSaltStr), kdfIterations, kdfKeyLen, sha512.New)
}