package sdk

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"

	"github.com/tyler-smith/go-bip39"
)

func GenerateBackupPrivateKey() ([]byte) {
	key := make([]byte, 16)
	rand.Read(key)
	return key
}

func GenerateBackupCipherKey(key string, input []byte) ([]byte) {
	cipher := hmac.New(sha256.New, []byte(key))
	cipher.Write(input)
	return cipher.Sum(nil) 
}

func BackupMnemonicFromKey(key []byte) (string, error) {
  	mnemonic, mnemonicError := bip39.NewMnemonic(key)
  	return mnemonic, mnemonicError
}

func BackupMnemonicToKey(mnemonic string) ([]byte, error) {
	privateKey, privateKeyError := bip39.EntropyFromMnemonic(mnemonic)
  	return privateKey, privateKeyError
}