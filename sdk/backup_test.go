package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateAndRecovery(t *testing.T) {
	key := GenerateBackupPrivateKey()
	cipher := GenerateBackupCipherKey("Algorand export 1.0", key)
	mnemonic, err := BackupMnemonicFromKey(key)
	require.NoError(t, err)
	recoveredKey, err := BackupMnemonicToKey(mnemonic)
	cipher2 := GenerateBackupCipherKey("Algorand export 1.0", recoveredKey)
	require.NoError(t, err)
	require.Equal(t, recoveredKey, key)
	require.Equal(t, cipher, cipher2)
}
