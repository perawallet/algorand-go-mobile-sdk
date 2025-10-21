package sdk

import (
	"bytes"
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/json"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// LogicSigAccount represents a LogicSig account
type LogicSigAccount struct {
	value crypto.LogicSigAccount
}

// ExtractLogicSigAccountFromSignedTransaction extracts a LogicSigAccount from a signed transaction.
// This will return nil if the transaction was not signed by a LogicSig account.
func ExtractLogicSigAccountFromSignedTransaction(encodedSignedTx []byte) (*LogicSigAccount, error) {
	var stx types.SignedTxn
	err := msgpack.Decode(encodedSignedTx, &stx)
	if err != nil {
		return nil, err
	}

	if stx.Lsig.Blank() {
		return nil, nil
	}

	var signerPublicKey *ed25519.PublicKey
	if stx.Lsig.Sig != (types.Signature{}) {
		var pk ed25519.PublicKey
		if !stx.AuthAddr.IsZero() {
			pk = stx.AuthAddr[:]
		} else {
			pk = stx.Txn.Sender[:]
		}
		signerPublicKey = &pk
	}

	account, err := crypto.LogicSigAccountFromLogicSig(stx.Lsig, signerPublicKey)
	if err != nil {
		return nil, err
	}

	return &LogicSigAccount{account}, nil
}

// DeserializeLogicSigAccountFromJSON deserializes a LogicSigAccount from a JSON string. See
// LogicSigAccount.ToJSON for serialization.
func DeserializeLogicSigAccountFromJSON(jsonStr string) (*LogicSigAccount, error) {
	var account crypto.LogicSigAccount
	err := json.Decode([]byte(jsonStr), &account)
	if err != nil {
		return nil, err
	}
	return &LogicSigAccount{account}, nil
}

// MakeLogicSigAccountEscrow creates a new escrow LogicSigAccount.
// The address of this account will be a hash of its program.
func MakeLogicSigAccountEscrow(program []byte, args *BytesArray) (*LogicSigAccount, error) {
	var extractedArgs [][]byte
	if args != nil {
		extractedArgs = args.Extract()
	}

	account, err := crypto.MakeLogicSigAccountEscrowChecked(program, extractedArgs)
	if err != nil {
		return nil, err
	}

	return &LogicSigAccount{account}, nil
}

// MakeLogicSigAccountDelegatedSign creates a new delegated LogicSigAccount. This
// type of LogicSig has the authority to sign transactions on behalf of another
// account, called the delegating account. If the delegating account is a
// multisig account, use MakeLogicSigAccountDelegatedMsig instead.
//
// This version of the function takes the private key of the delegating account, `signerSk`, and
// signs the program with that key. If you instead wish to provide the signature, use
// MakeLogicSigAccountDelegatedAttachSig.
func MakeLogicSigAccountDelegatedSign(program []byte, args *BytesArray, signerSk []byte) (*LogicSigAccount, error) {
	var extractedArgs [][]byte
	if args != nil {
		extractedArgs = args.Extract()
	}

	account, err := crypto.MakeLogicSigAccountDelegated(program, extractedArgs, signerSk)
	if err != nil {
		return nil, err
	}

	return &LogicSigAccount{account}, nil
}

// MakeLogicSigAccountDelegatedSign creates a new delegated LogicSigAccount. This
// type of LogicSig has the authority to sign transactions on behalf of another
// account, called the delegating account. If the delegating account is a
// multisig account, use MakeLogicSigAccountDelegatedMsig instead.
//
// This version of the function takes the signer address, `signer`, and its signature over the
// program, `signature`. See LogicSigProgramForSigning to calculate the bytes that must
// be signed. If you instead wish to provide the private key and sign directly, use
// MakeLogicSigAccountDelegatedSign.
func MakeLogicSigAccountDelegatedAttachSig(program []byte, args *BytesArray, signer string, signature []byte) (*LogicSigAccount, error) {
	var extractedArgs [][]byte
	if args != nil {
		extractedArgs = args.Extract()
	}

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

	account, err := crypto.MakeLogicSigAccountEscrowChecked(program, extractedArgs)
	if err != nil {
		return nil, err
	}

	account.Lsig.Sig = s
	account.SigningKey = signerAddr[:]

	if !crypto.VerifyLogicSig(account.Lsig, signerAddr) {
		return nil, errors.New("invalid signature provided")
	}

	return &LogicSigAccount{account}, nil
}

// MakeLogicSigAccountDelegatedMsig creates a new delegated LogicSigAccount.
// This type of LogicSig has the authority to sign transactions on behalf of
// another account, called the delegating account. Use this function if the
// delegating account is a multisig account, otherwise use of one MakeLogicSigAccountDelegatedSign
// or MakeLogicSigAccountDelegatedAttachSig.
//
// The parameter msigAccount is the delegating multisig account.
//
// You must use the methods AppendSignMultisigSignature or AppendAttachMultisigSignature on the
// returned LogicSigAccount to add signatures from the members of the multisig account. The multisig
// account's threshold for signatures must be met for this to be a valid delegated LogicSig.
func MakeLogicSigAccountDelegatedMsig(program []byte, args *BytesArray, msigAccount *MultisigAccount) (*LogicSigAccount, error) {
	var extractedArgs [][]byte
	if args != nil {
		extractedArgs = args.Extract()
	}

	account, err := crypto.MakeLogicSigAccountEscrowChecked(program, extractedArgs)
	if err != nil {
		return nil, err
	}

	// Construct the MultisigSig
	msig := types.MultisigSig{
		Version:   msigAccount.value.Version,
		Threshold: msigAccount.value.Threshold,
		Subsigs:   make([]types.MultisigSubsig, len(msigAccount.value.Pks)),
	}
	for i, pk := range msigAccount.value.Pks {
		c := make([]byte, len(pk))
		copy(c, pk)
		msig.Subsigs[i].Key = c
	}
	account.Lsig.Msig = msig

	return &LogicSigAccount{account}, nil
}

// AppendSignMultisigSignature adds an additional signature from a member of the
// delegating multisig account. This version of the function uses the passed in private key to
// calculate the signature; if you instead wish to provide the signature, see AppendAttachMultisigSignature.
//
// The LogicSigAccount must represent a delegated LogicSig backed by a multisig
// account.
func (lsa *LogicSigAccount) AppendSignMultisigSignature(signerSk []byte) error {
	if len(signerSk) != ed25519.PrivateKeySize {
		return fmt.Errorf("Incorrect privateKey length expected %d, got %d", ed25519.PrivateKeySize, len(signerSk))
	}

	if lsa.value.Lsig.Msig.Blank() {
		return errors.New("empty multisig in logicsig")
	}

	// Sign the program
	programData := LogicSigProgramForSigning(lsa.value.Lsig.Logic)
	signature := ed25519.Sign(ed25519.PrivateKey(signerSk), programData)

	// Get the public key from the private key
	publicKey := ed25519.PrivateKey(signerSk).Public().(ed25519.PublicKey)
	var signerAddr types.Address
	copy(signerAddr[:], publicKey)

	// Find the signer's index in the multisig
	signerIndex := -1
	for i, subsig := range lsa.value.Lsig.Msig.Subsigs {
		var pkAddr types.Address
		copy(pkAddr[:], subsig.Key[:])
		if pkAddr == signerAddr {
			signerIndex = i
			break
		}
	}
	if signerIndex == -1 {
		return errors.New("signer address does not match any of the addresses in the multisig account")
	}

	// Attach the signature
	var s types.Signature
	copy(s[:], signature)
	lsa.value.Lsig.Msig.Subsigs[signerIndex].Sig = s

	// Note: We don't verify here because multisig may not have enough signatures yet
	// Verification happens when the LogicSig is actually used to sign a transaction

	return nil
}

// AppendAttachMultisigSignature adds an additional signature from a member of the
// delegating multisig account. This version of the function requires you to provide the signer's
// address, `signer`, and its signature over the program, `signature`. See LogicSigProgramForSigning
// to calculate the bytes that must be signed. If you instead wish to provide the private key directly,
// see AppendSignMultisigSignature.
//
// The LogicSigAccount must represent a delegated LogicSig backed by a multisig
// account.
func (lsa *LogicSigAccount) AppendAttachMultisigSignature(signer string, signature []byte) error {
	if len(signature) != ed25519.SignatureSize {
		return fmt.Errorf("incorrect signature length expected %d, got %d", ed25519.SignatureSize, len(signature))
	}
	// Copy signature into a Signature, and check that it's the expected length
	var s types.Signature
	n := copy(s[:], signature)
	if n != len(s) {
		return errInvalidSignatureReturned
	}

	signerAddr, err := types.DecodeAddress(signer)
	if err != nil {
		return err
	}

	if lsa.value.Lsig.Msig.Blank() {
		return errors.New("empty multisig in logicsig")
	}

	signerIndex := -1
	for i, subsig := range lsa.value.Lsig.Msig.Subsigs {
		var pkAddr types.Address
		copy(pkAddr[:], subsig.Key[:])
		if pkAddr == signerAddr {
			signerIndex = i
			break
		}
	}
	if signerIndex == -1 {
		return errors.New("signer address does not match any of the addresses in the multisig account")
	}

	lsa.value.Lsig.Msig.Subsigs[signerIndex].Sig = s

	// Note: We don't verify here because multisig may not have enough signatures yet
	// Verification happens when the LogicSig is actually used to sign a transaction

	return nil
}

// IsDelegated returns true if and only if the LogicSigAccount has been delegated to
// another account with a signature.
func (lsa *LogicSigAccount) IsDelegated() bool {
	return lsa.value.IsDelegated()
}

// Address returns the address over which this LogicSigAccount has authority.
//
// If the LogicSig is delegated to another account, this will return the address
// of that account.
//
// If the LogicSig is not delegated to another account, this will return an
// escrow address that is the hash of the LogicSig's program code.
func (lsa *LogicSigAccount) Address() (string, error) {
	addr, err := lsa.value.Address()
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}

// ToJSON returns a JSON serialization of a LogicSigAccount. See DeserializeLogicSigAccountFromJSON for
// deserialization.
func (lsa *LogicSigAccount) ToJSON() string {
	return string(json.Encode(&lsa.value))
}

// SignLogicSigTransaction signs a transaction with a LogicSigAccount.
//
// Note: any type of transaction can be signed by a LogicSig, but the network will reject the
// transaction if the LogicSig's program declines the transaction.
func SignLogicSigTransaction(account *LogicSigAccount, encodedTx []byte) ([]byte, error) {
	var tx types.Transaction
	err := msgpack.Decode(encodedTx, &tx)
	if err != nil {
		return nil, err
	}

	_, stxBytes, err := crypto.SignLogicSigAccountTransaction(account.value, tx)
	return stxBytes, err
}

// LogicSigProgramForSigning returns the bytes that should be signed for a delegated LogicSig.
func LogicSigProgramForSigning(program []byte) []byte {
	return bytes.Join([][]byte{[]byte("Program"), program}, []byte{})
}
