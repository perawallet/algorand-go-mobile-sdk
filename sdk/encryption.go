package sdk

import (
	"io"
	"crypto/rand"
	"golang.org/x/crypto/nacl/secretbox"
)


type Encryption struct {
	EncryptedData []byte
	DecryptedData []byte
	ErrorCode int

	// ErrorCode Descriptions
	// 0 => No Error
	// 1 => Invalid SecretKey
	// 2 => Random Generator Error
	// 3 => Invalid encrypted data length
	// 4 => Decryption error
}

func Encrypt(data []byte, sk []byte) *Encryption  {
	var secretKey [32]byte

	if len(sk) != len(secretKey) {
		return &Encryption{
			ErrorCode: 1,
		}
	}

	copy(secretKey[:], sk)
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return &Encryption{
			ErrorCode: 2,
		}
	}

	encrypted := secretbox.Seal(nonce[:], data, &nonce, &secretKey)

	encryptedData := &Encryption{
		EncryptedData: encrypted,
		ErrorCode: 0,
	}

	return encryptedData
}

func Decrypt(data []byte, sk []byte) *Encryption {
	var secretKey [32]byte
	if len(sk) != len(secretKey) {
		return &Encryption{
			EncryptedData: data,
			ErrorCode: 1,
		}
	}

	copy(secretKey[:], sk)

	if len(data) < 24 {
		return &Encryption{
			EncryptedData: data,
			ErrorCode: 3,
		}
	}

	var decryptNonce [24]byte
	copy(decryptNonce[:], data[:24])

	decrypted, ok := secretbox.Open(nil, data[24:], &decryptNonce, &secretKey)

	if !ok {
		return &Encryption{
			EncryptedData: data,
			ErrorCode: 4,
		}
	}

	return &Encryption{
		EncryptedData: data,
		DecryptedData: decrypted,
		ErrorCode: 0,
	}
}