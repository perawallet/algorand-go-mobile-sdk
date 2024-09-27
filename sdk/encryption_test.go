package sdk

import (
	"testing"
	"github.com/stretchr/testify/require"
)

const (
	firstValidSecretKey = "0123456789ABCDEFGHIJKLMNOPQRSTUV"
	firstInvalidSecretKey = "0123456789AB"
	secondValidSecretKey = "9876456789ABCDEFGAIJKKMNDPQSSTCV"
)


func TestSecretboxEncryption(t *testing.T) {
	testString := "testdata"
	testData := []byte(testString)
	secretKey := firstValidSecretKey
	secretData := []byte(secretKey)

	encryptedContent := Encrypt(testData, secretData)

	require.NotEmpty(t, encryptedContent)
	require.NotEmpty(t, encryptedContent.EncryptedData)
	require.Empty(t, encryptedContent.DecryptedData)
	require.Equal(t, encryptedContent.ErrorCode, 0)

	decryptedContent := Decrypt(encryptedContent.EncryptedData, secretData)

	require.NotEmpty(t, decryptedContent)
	require.NotEmpty(t, decryptedContent.DecryptedData)
	require.NotEmpty(t, decryptedContent.EncryptedData)
	require.Equal(t, decryptedContent.ErrorCode, 0)
	require.Equal(t, decryptedContent.DecryptedData, testData)
	require.Equal(t, decryptedContent.EncryptedData, encryptedContent.EncryptedData)
}

func TestSecretboxInvalidSecretKey(t *testing.T) {
	testString := "testdata"
	testData := []byte(testString)
	secretKey := firstInvalidSecretKey
	secretData := []byte(secretKey)

	encryptedContent := Encrypt(testData, secretData)

	require.NotEmpty(t, encryptedContent)
	require.Equal(t, encryptedContent.ErrorCode, 1)
	require.Empty(t, encryptedContent.EncryptedData)
	require.Empty(t, encryptedContent.DecryptedData)

	decryptedContent := Decrypt(encryptedContent.EncryptedData, secretData)

	require.NotEmpty(t, decryptedContent)
	require.Empty(t, encryptedContent.EncryptedData)
	require.Empty(t, encryptedContent.DecryptedData)
	require.Equal(t, encryptedContent.ErrorCode, 1)
}

func TestSecretboxSameDataWithDifferentEncryptions(t *testing.T) {
	testString := "testdata"
	testData := []byte(testString)
	secretKey := firstValidSecretKey
	secretData := []byte(secretKey)

	firstEncryptedContent := Encrypt(testData, secretData)

	require.NotEmpty(t, firstEncryptedContent)
	require.NotEmpty(t, firstEncryptedContent.EncryptedData)
	require.Empty(t, firstEncryptedContent.DecryptedData)
	require.Equal(t, firstEncryptedContent.ErrorCode, 0)

	secondSecretKey := secondValidSecretKey
	secondSecretData := []byte(secondSecretKey)

	secondEncryptedContent := Encrypt(testData, secondSecretData)

	require.NotEmpty(t, secondEncryptedContent)
	require.NotEmpty(t, secondEncryptedContent.EncryptedData)
	require.Empty(t, secondEncryptedContent.DecryptedData)
	require.Equal(t, secondEncryptedContent.ErrorCode, 0)

	firstDecryptedContent := Decrypt(firstEncryptedContent.EncryptedData, secretData)

	require.NotEmpty(t, firstDecryptedContent)
	require.NotEmpty(t, firstDecryptedContent.DecryptedData)
	require.NotEmpty(t, firstDecryptedContent.EncryptedData)
	require.Equal(t, firstDecryptedContent.ErrorCode, 0)
	require.Equal(t, firstDecryptedContent.DecryptedData, testData)
	require.Equal(t, firstDecryptedContent.EncryptedData, firstEncryptedContent.EncryptedData)

	secondDecryptedContent := Decrypt(secondEncryptedContent.EncryptedData, secondSecretData)

	require.NotEmpty(t, secondDecryptedContent)
	require.NotEmpty(t, secondDecryptedContent.DecryptedData)
	require.NotEmpty(t, secondDecryptedContent.EncryptedData)
	require.Equal(t, secondDecryptedContent.ErrorCode, 0)
	require.Equal(t, secondDecryptedContent.DecryptedData, testData)
	require.Equal(t, secondDecryptedContent.DecryptedData, firstDecryptedContent.DecryptedData)
	require.Equal(t, secondDecryptedContent.EncryptedData, secondDecryptedContent.EncryptedData)

	secondDecryptedContentWithFirstSecret := Decrypt(secondEncryptedContent.EncryptedData, secretData)

	require.NotEmpty(t, secondDecryptedContentWithFirstSecret)
	require.Equal(t, secondDecryptedContentWithFirstSecret.ErrorCode, 4)
	require.NotEmpty(t, secondDecryptedContentWithFirstSecret.EncryptedData)
	require.Empty(t, secondDecryptedContentWithFirstSecret.DecryptedData)

	var extractedContent []byte
	copy(extractedContent[:], secondEncryptedContent.EncryptedData[:24])

	invalidDecryptedContent := Decrypt(extractedContent, secretData)

	require.NotEmpty(t, invalidDecryptedContent)
	require.Empty(t, invalidDecryptedContent.DecryptedData)
	require.Empty(t, invalidDecryptedContent.EncryptedData)
	require.NotEqual(t, invalidDecryptedContent.EncryptedData, secondEncryptedContent.EncryptedData)
	require.Equal(t, invalidDecryptedContent.ErrorCode, 3)
}