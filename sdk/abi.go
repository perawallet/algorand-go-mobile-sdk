package sdk

import (
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/abi"
	"github.com/algorand/go-algorand-sdk/v2/encoding/json"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// ABIType represents an ARC-4 ABI type.
type ABIType struct {
	value abi.Type
}

// ParseABIType parses a string into an ABIType.
func ParseABIType(typeString string) (*ABIType, error) {
	t, err := abi.TypeOf(typeString)
	if err != nil {
		return nil, err
	}
	return &ABIType{t}, nil
}

// String returns the string representation of the ABIType.
func (t *ABIType) String() string {
	return t.value.String()
}

// Encode takes an ABI value in JSON format and encodes it into a byte array.
//
// The JSON format for ABI values is the same as the format used by `goal app method --arg <argument>`.
//
// See `MarshalToJSON` and UnmarshalFromJSON from https://github.com/algorand/avm-abi/blob/3ac8977d88f2937721d1602027dde17f450e62dd/abi/json.go
// for more information about how the JSON values are handled.
func (t *ABIType) Encode(jsonValue string) ([]byte, error) {
	goValue, err := t.value.UnmarshalFromJSON([]byte(jsonValue))
	if err != nil {
		return nil, err
	}
	return t.value.Encode(goValue)
}

// Decode takes an encoded ABI value and decodes it into a JSON string.
//
// The JSON format for ABI values is the same as the format used by `goal app method --arg <argument>`.
//
// See `MarshalToJSON` and UnmarshalFromJSON from https://github.com/algorand/avm-abi/blob/3ac8977d88f2937721d1602027dde17f450e62dd/abi/json.go
// for more information about how the JSON values are handled.
func (t *ABIType) Decode(encodedValue []byte) (string, error) {
	goValue, err := t.value.Decode(encodedValue)
	if err != nil {
		return "", err
	}
	jsonValue, err := t.value.MarshalToJSON(goValue)
	return string(jsonValue), err
}

// ABIMethodJSONFromSignature takes a method signature and returns the JSON representation of the method.
func ABIMethodJSONFromSignature(signature string) (string, error) {
	method, err := abi.MethodFromSignature(signature)
	if err != nil {
		return "", err
	}
	return string(json.Encode(method)), nil
}

// GetABIMethodSignature takes a method JSON representation and returns the method signature.
func GetABIMethodSignature(methodJSON string) (string, error) {
	var method abi.Method
	err := json.Decode([]byte(methodJSON), &method)
	if err != nil {
		return "", err
	}
	return method.GetSignature(), nil
}

// AtomicTransactionComposer is a class for constructing and signing atomic transaction groups.
type AtomicTransactionComposer struct {
	value transaction.AtomicTransactionComposer
}

// NewAtomicTransactionComposer creates a new AtomicTransactionComposer.
func NewAtomicTransactionComposer() *AtomicTransactionComposer {
	return &AtomicTransactionComposer{}
}

// GetStatus returns the status of this composer's transaction group.
//
// The values that may be returned are:
// * 0 - BUILDING: The atomic group is still under construction.
// * 1 - BUILT: The atomic group has been finalized, but not yet signed.
// * 2 - SIGNED: The atomic group has been finalized and signed.
//
// Once a composer's status is at least BUILT, it may no longer be modified. A composer may advance
// to higher status levels, but it may never regress to lower status levels. If you wish to modify
// a composer that has already been BUILT, use `Clone()` to create a copy.
func (c *AtomicTransactionComposer) GetStatus() int {
	return c.value.GetStatus()
}

// Count returns the number of transactions currently in this atomic group.
func (c *AtomicTransactionComposer) Count() int {
	return c.value.Count()
}

// Clone creates a new composer with the same underlying transactions. The new composer's status will
// be BUILDING, so additional transactions may be added to it.
func (c *AtomicTransactionComposer) Clone() *AtomicTransactionComposer {
	return &AtomicTransactionComposer{c.value.Clone()}
}

// AddTransaction adds a transaction to this atomic group.
//
// An error will be thrown if the composer's status is not BUILDING, or if adding this transaction
// causes the current group to exceed MaxAtomicGroupSize (16).
func (c *AtomicTransactionComposer) AddTransaction(encodedTx []byte, signer TransactionSigner) error {
	var tx types.Transaction
	err := msgpack.Decode(encodedTx, &tx)
	if err != nil {
		return err
	}
	return c.value.AddTransaction(transaction.TransactionWithSigner{
		Txn:    tx,
		Signer: externalToInternalSigner{signer},
	})
}

// AddMethodCallParams contains the parameters for the method `AtomicTransactionComposer.AddMethodCall`
type AddMethodCallParams struct {
	value transaction.AddMethodCallParams
}

// NewAddMethodCallParams creates a new AddMethodCallParams object.
//
// onComplete is the OnCompletion action to take for this application call. The accepted values are:
// * 0 - NoOp
// * 1 - OptIn
// * 2 - CloseOut
// * 3 - ClearState
// * 4 - UpdateApplication
// * 5 - DeleteApplication
func NewAddMethodCallParams(
	appID int64,
	onComplete int,
	methodJson string,
	accounts *StringArray,
	foreignApps *Int64Array,
	foreignAssets *Int64Array,
	boxRefs *AppBoxRefArray,
	txnParams *SuggestedParams,
	note []byte,
	sender string,
	signer TransactionSigner,
) (*AddMethodCallParams, error) {
	if appID < 0 {
		return nil, errNegativeArgument
	}
	if onComplete < 0 || types.OnCompletion(onComplete) > types.DeleteApplicationOC {
		return nil, fmt.Errorf("invalid onComplete value: %d", onComplete)
	}

	var method abi.Method
	err := json.Decode([]byte(methodJson), &method)
	if err != nil {
		return nil, fmt.Errorf("could not decode method from JSON: %w", err)
	}

	internalForeignApps := make([]uint64, foreignApps.Length())
	for i := range internalForeignApps {
		value := foreignApps.Get(i)
		if value < 0 {
			return nil, errNegativeArgument
		}
		internalForeignApps[i] = uint64(value)
	}

	internalForeignAssets := make([]uint64, foreignAssets.Length())
	for i := range internalForeignAssets {
		value := foreignAssets.Get(i)
		if value < 0 {
			return nil, errNegativeArgument
		}
		internalForeignAssets[i] = uint64(value)
	}

	internalTxnParams, err := convertSuggestedParams(txnParams)
	if err != nil {
		return nil, err
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		return nil, err
	}

	params := transaction.AddMethodCallParams{
		AppID:           uint64(appID),
		Method:          method,
		OnComplete:      types.OnCompletion(onComplete),
		ForeignAccounts: accounts.Extract(),
		ForeignApps:     internalForeignApps,
		ForeignAssets:   internalForeignAssets,
		BoxReferences:   boxRefs.Extract(),
		SuggestedParams: internalTxnParams,
		Note:            note,
		Sender:          senderAddr,
		Signer:          externalToInternalSigner{signer},
	}
	return &AddMethodCallParams{params}, nil
}

// AddMethodArgument adds an ABI argument to the method call. This uses the same format as `ABIType.Encode()`
// for the argument value. This method can handle basic and reference ABI types, but not transaction
// argument types. See `AddMethodArgumentTransaction()` for transaction type support.
func (p *AddMethodCallParams) AddMethodArgument(valueJson string) error {
	numArgs := len(p.value.MethodArgs)
	if numArgs+1 > len(p.value.Method.Args) {
		return fmt.Errorf("too many arguments for method: '%s'", p.value.Method.Name)
	}
	argSpec := p.value.Method.Args[numArgs]
	var typeToDecode abi.Type
	if argSpec.IsTransactionArg() {
		return errors.New("cannot add a transaction argument using this method")
	}
	if argSpec.IsReferenceArg() {
		var proxyType string
		switch argSpec.Type {
		case abi.AccountReferenceType:
			proxyType = "address"
		case abi.AssetReferenceType, abi.ApplicationReferenceType:
			proxyType = "uint64"
		default:
			return fmt.Errorf("unsupported reference type: %s", argSpec.Type)
		}
		var err error
		typeToDecode, err = abi.TypeOf(proxyType)
		if err != nil {
			return fmt.Errorf("could not resolve reference type %s: %w", argSpec.Type, err)
		}
	} else {
		var err error
		typeToDecode, err = argSpec.GetTypeObject()
		if err != nil {
			return err
		}
	}
	goValue, err := typeToDecode.UnmarshalFromJSON([]byte(valueJson))
	if err != nil {
		return fmt.Errorf("cannot decode JSON value for argument type %s: %w", argSpec.Type, err)
	}
	p.value.MethodArgs = append(p.value.MethodArgs, goValue)
	return nil
}

// AddMethodArgumentTransaction adds a transaction ABI argument argument to the method call.
func (p *AddMethodCallParams) AddMethodArgumentTransaction(encodedTx []byte, signer TransactionSigner) error {
	numArgs := len(p.value.MethodArgs)
	if numArgs+1 > len(p.value.Method.Args) {
		return fmt.Errorf("too many arguments for method: '%s'", p.value.Method.Name)
	}
	argSpec := p.value.Method.Args[numArgs]
	if !argSpec.IsTransactionArg() {
		return fmt.Errorf("this method only accepts a transaction argument, got: '%s'", argSpec.Type)
	}
	var tx types.Transaction
	err := msgpack.Decode(encodedTx, &tx)
	if err != nil {
		return err
	}
	p.value.MethodArgs = append(p.value.MethodArgs, transaction.TransactionWithSigner{
		Txn:    tx,
		Signer: externalToInternalSigner{signer},
	})
	return nil
}

// AddPrograms adds the approval and clear state programs to the method call. Only needed for app
// creation or update.
func (p *AddMethodCallParams) AddPrograms(approvalProgram, clearStateProgram []byte) {
	p.value.ApprovalProgram = approvalProgram
	p.value.ClearProgram = clearStateProgram
}

// AddAppSchema adds global schema, local schema, and extra pages to the method call. Only needed for app creation.
func (p *AddMethodCallParams) AddAppSchema(globalSchemaUint, globalSchemaByteSlice, localSchemaUint, localSchemaByteSlice int64, extraPages int32) error {
	if globalSchemaUint < 0 || globalSchemaByteSlice < 0 || localSchemaUint < 0 || localSchemaByteSlice < 0 || extraPages < 0 {
		return errNegativeArgument
	}
	p.value.GlobalSchema = types.StateSchema{
		NumUint:      uint64(globalSchemaUint),
		NumByteSlice: uint64(globalSchemaByteSlice),
	}
	p.value.LocalSchema = types.StateSchema{
		NumUint:      uint64(localSchemaUint),
		NumByteSlice: uint64(localSchemaByteSlice),
	}
	p.value.ExtraPages = uint32(extraPages)
	return nil
}

// AddMethodCall adds a smart contract method call to this atomic group.
//
// An error will be thrown if the composer's status is not BUILDING, if adding this transaction
// causes the current group to exceed MaxAtomicGroupSize (16), or if the provided arguments are invalid
// for the given method.
func (c *AtomicTransactionComposer) AddMethodCall(params *AddMethodCallParams) error {
	return c.value.AddMethodCall(params.value)
}

// BuildGroup finalizes the transaction group and returns the finalized unsigned transactions.
//
// The composer's status will be at least BUILT after executing this method.
func (c *AtomicTransactionComposer) BuildGroup() (*BytesArray, error) {
	txnsWithSigners, err := c.value.BuildGroup()
	if err != nil {
		return nil, err
	}
	txnBytes := make([][]byte, len(txnsWithSigners))
	for i, txnWithSigner := range txnsWithSigners {
		txnBytes[i] = msgpack.Encode(&txnWithSigner.Txn)
	}
	return &BytesArray{txnBytes}, nil
}

// GatherSignatures obtains signatures for each transaction in this group. If signatures have
// already been obtained, this method will return cached versions of the signatures.
//
// The composer's status will be at least SIGNED after executing this method.
//
// An error will be thrown if signing any of the transactions fails. Otherwise, this will return an
// array of signed transactions.
func (c *AtomicTransactionComposer) GatherSignatures() (*BytesArray, error) {
	stxnBytes, err := c.value.GatherSignatures()
	if err != nil {
		return nil, err
	}
	return &BytesArray{stxnBytes}, nil
}
