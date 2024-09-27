package sdk

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"golang.org/x/crypto/ed25519"
)

// MnemonicFromKey converts a 32-byte key into a 25 word mnemonic. The generated
// mnemonic includes a checksum. Each word in the mnemonic represents 11 bits
// of data, and the last 11 bits are reserved for the checksum.
func MnemonicFromKey(key []byte) (string, error) {
	return mnemonic.FromKey(key)
}

// MnemonicToKey converts a mnemonic generated using this library into the
// source key used to create it. It returns an error if the passed mnemonic has
// an incorrect checksum, if the number of words is unexpected, or if one of the
// passed words is not found in the words list.
func MnemonicToKey(mnemonicStr string) ([]byte, error) {
	return mnemonic.ToKey(mnemonicStr)
}

// MnemonicFromPrivateKey is a helper that converts an ed25519 private key to a
// human-readable mnemonic
func MnemonicFromPrivateKey(sk []byte) (string, error) {
	if len(sk) != ed25519.PrivateKeySize {
		return "", errWrongKeyLen
	}
	sk1 := ed25519.PrivateKey(sk)
	return mnemonic.FromPrivateKey(sk1)
}

// MnemonicToPrivateKey is a helper that converts a mnemonic directly to an
// ed25519 private key
func MnemonicToPrivateKey(mnemonicStr string) (sk []byte, err error) {
	return mnemonic.ToPrivateKey(mnemonicStr)
}

// MnemonicFromMasterDerivationKey is a helper that converts an MDK to a
// human-readable mnemonic
func MnemonicFromMasterDerivationKey(mdk []byte) (string, error) {
	var mdkType types.MasterDerivationKey
	if len(mdk) != len(mdkType) {
		return "", fmt.Errorf("Wrong length for master derivation key. Expected %d, got %d", len(mdkType), len(mdk))
	}
	copy(mdkType[:], mdk)
	return mnemonic.FromMasterDerivationKey(mdkType)
}

// MnemonicToMasterDerivationKey is a helper that converts a mnemonic directly
// to a master derivation key
func MnemonicToMasterDerivationKey(mnemonicStr string) (mdk []byte, err error) {
	mdkType, err := mnemonic.ToMasterDerivationKey(mnemonicStr)
	mdk = mdkType[:]
	return
}
