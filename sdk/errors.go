package sdk

import "errors"

// crypto
var errInvalidSignatureReturned = errors.New("ed25519 library returned an invalid signature")
var errFailedToCopyPK = errors.New("failed to copy the public key")

// mnemonic
var errWrongKeyLen = errors.New("wrong key length") // TODO: check the actual error text for this

// transaction
var errNegativeArgument = errors.New("all integer arguments must be >= 0")
