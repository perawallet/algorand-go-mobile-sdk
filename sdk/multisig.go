package sdk

import (
	"errors"
	"fmt"
	"math"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
	"golang.org/x/crypto/ed25519"
)

// MultiSigAccount represents a MultiSig account
type MultisigAccount struct {
	value crypto.MultisigAccount
}

// MakeMultisigAccount creates a new instance of a MultiSig account. The order of the addresses matters.
func MakeMultisigAccount(version int, threshold int, addrs *StringArray) (*MultisigAccount, error) {
	addresses := make([]types.Address, addrs.Length())
	for i, addrStr := range addrs.Extract() {
		addr, err := types.DecodeAddress(addrStr)
		if err != nil {
			return nil, fmt.Errorf("could not decode address '%s': %w", addrStr, err)
		}
		addresses[i] = addr
	}
	if version < 0 || version > math.MaxUint8 {
		return nil, fmt.Errorf("version %d out of range", version)
	}
	if threshold < 0 || threshold > math.MaxUint8 {
		return nil, fmt.Errorf("threshold %d out of range", version)
	}

	ma, err := crypto.MultisigAccountWithParams(uint8(version), uint8(threshold), addresses)
	if err != nil {
		return nil, err
	}

	return &MultisigAccount{ma}, nil
}

// ExtractMultisigAccountFromSignedTransaction extracts a MultisigAccount from a signed transaction.
// This will return nil if the transaction was not signed by a multisig account.
func ExtractMultisigAccountFromSignedTransaction(encodedSignedTx []byte) (*MultisigAccount, error) {
	var stx types.SignedTxn
	err := msgpack.Decode(encodedSignedTx, &stx)
	if err != nil {
		return nil, err
	}

	if stx.Msig.Blank() {
		return nil, nil
	}

	ma, err := crypto.MultisigAccountFromSig(stx.Msig)
	if err != nil {
		return nil, err
	}

	return &MultisigAccount{ma}, nil
}

// Address generates the corresponding account address for this MultisigAccount
func (ma *MultisigAccount) Address() (string, error) {
	addr, err := ma.value.Address()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

// Version returns the version of this MultisigAccount
func (ma *MultisigAccount) Version() int {
	return int(ma.value.Version)
}

// Threshold returns the threshold of this MultisigAccount
func (ma *MultisigAccount) Threshold() int {
	return int(ma.value.Threshold)
}

// ContributingAddresses returns the individual addresses that make up this MultisigAccount
func (ma *MultisigAccount) ContributingAddresses() *StringArray {
	addrs := make([]string, len(ma.value.Pks))
	for i, pk := range ma.value.Pks {
		var addr types.Address
		copy(addr[:], pk[:])
		addrs[i] = addr.String()
	}
	return &StringArray{values: addrs}
}

// SignMultisigTransaction signs and contributes a single signature to a transaction from a
// MultisigAccount. The transaction will only be properly signed once the MultiSigAccount threshold
// has been reached. MergeMultisigTransactions can be used to merge multiple partially-signed
// transactions into a single transaction.
//
// The argument `sk` must be the private key of one of the contributing addresses of the MultisigAccount.
func SignMultisigTransaction(sk []byte, account *MultisigAccount, encodedTx []byte) ([]byte, error) {
	if len(sk) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("Incorrect privateKey length expected %d, got %d", ed25519.PrivateKeySize, len(sk))
	}

	var tx types.Transaction
	err := msgpack.Decode(encodedTx, &tx)
	if err != nil {
		return nil, err
	}

	_, stxBytes, err := crypto.SignMultisigTransaction(sk, account.value, tx)
	return stxBytes, err
}

// AttachMultisigSignature attaches a single signature to a transaction from a MultisigAccount. The
// transaction will only be properly signed once the MultiSigAccount threshold has been reached.
// MergeMultisigTransactions can be used to merge multiple partially-signed transactions into a
// single transaction.
//
// The argument `signer` must be one of the contributing addresses of the MultisigAccount, and
// `signature` must be a signature of the transaction from that address.
func AttachMultisigSignature(signer string, signature []byte, account *MultisigAccount, encodedTx []byte) ([]byte, error) {
	if len(signature) != ed25519.SignatureSize {
		return nil, fmt.Errorf("incorrect signature length expected %d, got %d", ed25519.SignatureSize, len(signature))
	}

	// Copy signature into a Signature, and check that it's the expected length
	var s types.Signature
	n := copy(s[:], signature)
	if n != len(s) {
		return nil, errInvalidSignatureReturned
	}

	signerAddr, err := types.DecodeAddress(signer)
	if err != nil {
		return nil, err
	}

	var tx types.Transaction
	err = msgpack.Decode(encodedTx, &tx)
	if err != nil {
		return nil, err
	}

	signerIndex := -1
	for i, pk := range account.value.Pks {
		var pkAddr types.Address
		copy(pkAddr[:], pk[:])
		if pkAddr == signerAddr {
			signerIndex = i
			break
		}
	}
	if signerIndex == -1 {
		return nil, errors.New("signer address does not match any of the addresses in the multisig account")
	}

	// Construct the MultisigSig
	msig := types.MultisigSig{
		Version:   account.value.Version,
		Threshold: account.value.Threshold,
		Subsigs:   make([]types.MultisigSubsig, len(account.value.Pks)),
	}
	for i, pk := range account.value.Pks {
		c := make([]byte, len(pk))
		copy(c, pk)
		msig.Subsigs[i].Key = c
	}
	msig.Subsigs[signerIndex].Sig = s

	// Construct the SignedTxn
	stx := types.SignedTxn{
		Msig: msig,
		Txn:  tx,
	}
	msigAddr, err := account.value.Address()
	if err != nil {
		return nil, err
	}
	if tx.Sender != msigAddr {
		stx.AuthAddr = msigAddr
	}

	// Encode the SignedTxn
	stxBytes := msgpack.Encode(stx)
	return stxBytes, nil
}

// MergeMultisigTransactions merges multiple partially-signed transactions into a single transaction.
// The transactions to be merged must be signed by the same MultisigAccount. See
// SignMultisigTransaction and AttachMultisigSignature for creating partially-signed transactions.
func MergeMultisigTransactions(encodedSignedTx1, encodedSignedTx2 []byte) ([]byte, error) {
	_, stxnBytes, err := crypto.MergeMultisigTransactions(encodedSignedTx1, encodedSignedTx2)
	return stxnBytes, err
}
