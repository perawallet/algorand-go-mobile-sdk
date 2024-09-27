package sdk

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

type SuggestedParams struct {
	// Fee is the suggested transaction fee
	// Fee is in units of micro-Algos per byte.
	// Fee may fall to zero but transactions must still have a fee of
	// at least MinTxnFee for the current network protocol.
	Fee int64

	// Genesis ID
	GenesisID string

	// Genesis hash
	GenesisHash []byte

	// FirstRoundValid is the first protocol round on which the txn is valid
	FirstRoundValid int64

	// LastRoundValid is the final protocol round on which the txn may be committed
	LastRoundValid int64

	// FlatFee indicates whether the passed fee is per-byte or per-transaction
	FlatFee bool
}

func convertSuggestedParams(params *SuggestedParams) (internalParams types.SuggestedParams, err error) {
	if params.Fee < 0 || params.FirstRoundValid < 0 || params.LastRoundValid < 0 {
		err = fmt.Errorf("Could not convert suggested params: %v", errNegativeArgument)
		return
	}

	internalParams = types.SuggestedParams{
		Fee:             types.MicroAlgos(params.Fee),
		GenesisID:       params.GenesisID,
		GenesisHash:     params.GenesisHash,
		FirstRoundValid: types.Round(params.FirstRoundValid),
		LastRoundValid:  types.Round(params.LastRoundValid),
		FlatFee:         params.FlatFee,
	}

	return
}

// MakePaymentTxn constructs a payment transaction using the passed parameters.
// `from` and `to` addresses should be checksummed, human-readable addresses
func MakePaymentTxn(from, to string, amount *Uint64, note []byte, closeRemainderTo string, params *SuggestedParams) (encoded []byte, err error) {
	internalAmount, err := amount.Extract()
	if err != nil {
		err = fmt.Errorf("Could not decode transaction amount: %v", err)
		return
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	tx, err := transaction.MakePaymentTxn(from, to, internalAmount, note, closeRemainderTo, internalParams)
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeRekeyTxn constructs a rekey transaction using the passed parameters.
func MakeRekeyTxn(from, rekeyTo string, params *SuggestedParams) (encoded []byte, err error) {
	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	tx, err := transaction.MakePaymentTxn(from, from, 0, nil, "", internalParams)
	tx.Rekey(rekeyTo)
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeAssetCreateTxn constructs an asset creation transaction using the passed parameters.
// - account is a checksummed, human-readable address which will send the transaction.
// - note is a byte array
func MakeAssetCreateTxn(account string, note []byte, params *SuggestedParams, total *Uint64, decimals int32, defaultFrozen bool, manager, reserve, freeze, clawback, unitName, assetName, url string, metadataHash []byte) (encoded []byte, err error) {
	if decimals < 0 {
		err = errNegativeArgument
		return
	}

	internalTotal, err := total.Extract()
	if err != nil {
		err = fmt.Errorf("Could not extract asset total: %v", err)
		return
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	tx, err := transaction.MakeAssetCreateTxn(account, note, internalParams, internalTotal, uint32(decimals), defaultFrozen, manager, reserve, freeze, clawback, unitName, assetName, url, string(metadataHash))
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeAssetConfigTxn creates a tx template for changing the
// keys for an asset. An empty string means a zero key (which
// cannot be changed after becoming zero); to keep a key
// unchanged, you must specify that key.
// - account is a checksummed, human-readable address for which we register the given participation key.
func MakeAssetConfigTxn(account string, note []byte, params *SuggestedParams, index int64, newManager, newReserve, newFreeze, newClawback string) (encoded []byte, err error) {
	if index < 0 {
		err = errNegativeArgument
		return
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	tx, err := transaction.MakeAssetConfigTxn(account, note, internalParams, uint64(index), newManager, newReserve, newFreeze, newClawback, false)
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeAssetTransferTxn creates a tx for sending some asset from an asset holder to another user
// the recipient address must have previously issued an asset acceptance transaction for this asset
// - account is a checksummed, human-readable address that will send the transaction and assets
// - recipient is a checksummed, human-readable address what will receive the assets
// - closeAssetsTo is a checksummed, human-readable address that behaves as a close-to address for the asset transaction; the remaining assets not sent to recipient will be sent to closeAssetsTo. Leave blank for no close-to behavior.
// - amount is the number of assets to send
// - note is an arbitrary byte array
// - creator is the address of the asset creator
// - index is the asset index
func MakeAssetTransferTxn(account, recipient, closeAssetsTo string, amount *Uint64, note []byte, params *SuggestedParams, index int64) (encoded []byte, err error) {
	if index < 0 {
		err = errNegativeArgument
		return
	}

	internalAmount, err := amount.Extract()
	if err != nil {
		err = fmt.Errorf("Could not decode transaction amount: %v", err)
		return
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	tx, err := transaction.MakeAssetTransferTxn(account, recipient, internalAmount, note, internalParams, closeAssetsTo, uint64(index))
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeAssetAcceptanceTxn creates a tx for marking an account as willing to accept the given asset
// - account is a checksummed, human-readable address that will send the transaction and begin accepting the asset
// - note is an arbitrary byte array
// - index is the asset index
func MakeAssetAcceptanceTxn(account string, note []byte, params *SuggestedParams, index int64) (encoded []byte, err error) {
	if index < 0 {
		err = errNegativeArgument
		return
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	tx, err := transaction.MakeAssetAcceptanceTxn(account, note, internalParams, uint64(index))
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeAssetRevocationTxn creates a tx for revoking an asset from an account and sending it to another
// - account is a checksummed, human-readable address; it must be the revocation manager / clawback address from the asset's parameters
// - target is a checksummed, human-readable address; it is the account whose assets will be revoked
// - recipient is a checksummed, human-readable address; it will receive the revoked assets
// - amount defines the number of assets to clawback
// - index is the asset index
func MakeAssetRevocationTxn(account, target string, amount *Uint64, recipient string, note []byte, params *SuggestedParams, index int64) (encoded []byte, err error) {
	if index < 0 {
		err = errNegativeArgument
		return
	}

	internalAmount, err := amount.Extract()
	if err != nil {
		err = fmt.Errorf("Could not decode transaction amount: %v", err)
		return
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	tx, err := transaction.MakeAssetRevocationTxn(account, target, internalAmount, recipient, note, internalParams, uint64(index))
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeAssetDestroyTxn creates a tx template for destroying an asset, removing it from the record.
// All outstanding asset amount must be held by the creator, and this transaction must be issued by the asset manager.
// - account is a checksummed, human-readable address that will send the transaction; it also must be the asset manager
// - index is the asset index
func MakeAssetDestroyTxn(account string, note []byte, params *SuggestedParams, index int64) (encoded []byte, err error) {
	if index < 0 {
		err = errNegativeArgument
		return
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	tx, err := transaction.MakeAssetDestroyTxn(account, note, internalParams, uint64(index))
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeAssetFreezeTxn constructs a transaction that freezes or unfreezes an account's asset holdings
// It must be issued by the freeze address for the asset
// - account is a checksummed, human-readable address which will send the transaction.
// - note is an optional arbitrary byte array
// - assetIndex is the index for tracking the asset
// - target is the account to be frozen or unfrozen
// - newFreezeSetting is the new state of the target account
func MakeAssetFreezeTxn(account string, note []byte, params *SuggestedParams, assetIndex int64, target string, newFreezeSetting bool) (encoded []byte, err error) {
	if assetIndex < 0 {
		err = errNegativeArgument
		return
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	tx, err := transaction.MakeAssetFreezeTxn(account, note, internalParams, uint64(assetIndex), target, newFreezeSetting)
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeApplicationCreateTx makes a transaction for creating an application (see below for args desc.)
//
// - optIn:        true for opting in on complete, false for no-op.
//
//   - accounts      lists the accounts (in addition to the sender) that may be accessed
//     from the application logic.
//
//   - appArgs       ApplicationArgs lists some transaction-specific arguments accessible
//     from application logic.
//
//   - appIdx        ApplicationID is the application being interacted with, or 0 if
//     creating a new application.
//
//   - approvalProg  ApprovalProgram determines whether or not this ApplicationCall
//     transaction will be approved or not.
//
//   - clearProg     ClearStateProgram executes when a clear state ApplicationCall
//     transaction is executed. This program may not reject the
//     transaction, only update state.
//
//   - foreignApps   lists the applications (in addition to txn.ApplicationID) whose global
//     states may be accessed by this application. The access is read-only.
//
// - foreignAssets lists the assets whose global state may be accessed by this application. The access is read-only.
//
// - boxRefs	   lists the boxes which be accessed by this application call.
//
//   - globalSchema  GlobalStateSchema sets limits on the number of strings and
//     integers that may be stored in the GlobalState. The larger these
//     limits are, the larger minimum balance must be maintained inside
//     the creator's account (in order to 'pay' for the state that can
//     be used). The GlobalStateSchema is immutable.
//
//   - localSchema   LocalStateSchema sets limits on the number of strings and integers
//     that may be stored in an account's LocalState for this application.
//     The larger these limits are, the larger minimum balance must be
//     maintained inside the account of any users who opt into this
//     application. The LocalStateSchema is immutable.
//
//   - extraPages    ExtraProgramPages specifies the additional app program size requested in pages.
//     A page is 1024 bytes. This field enables execution of app programs
//     larger than the default maximum program size.
//
//   - onComplete    This is the faux application type used to distinguish different
//     application actions. Specifically, OnCompletion specifies what
//     side effects this transaction will have if it successfully makes
//     it into a block.
func MakeApplicationCreateTx(
	optIn bool,
	approvalProg []byte,
	clearProg []byte,
	globalSchemaUint int64,
	globalSchemaByteSlice int64,
	localSchemaUint int64,
	localSchemaByteSlice int64,
	extraPages int32,
	appArgs *BytesArray,
	accounts *StringArray,
	foreignApps *Int64Array,
	foreignAssets *Int64Array,
	boxRefs *AppBoxRefArray,
	params *SuggestedParams,
	sender string,
	note []byte,
) (encoded []byte, err error) {
	if globalSchemaUint < 0 || globalSchemaByteSlice < 0 || localSchemaUint < 0 || localSchemaByteSlice < 0 || extraPages < 0 {
		err = errNegativeArgument
		return
	}

	internalForeignApps := make([]uint64, foreignApps.Length())
	for i := range internalForeignApps {
		value := foreignApps.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignApps[i] = uint64(value)
	}

	internalForeignAssets := make([]uint64, foreignAssets.Length())
	for i := range internalForeignAssets {
		value := foreignAssets.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignAssets[i] = uint64(value)
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	globalSchema := types.StateSchema{
		NumUint:      uint64(globalSchemaUint),
		NumByteSlice: uint64(globalSchemaByteSlice),
	}

	localSchema := types.StateSchema{
		NumUint:      uint64(localSchemaUint),
		NumByteSlice: uint64(localSchemaByteSlice),
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		err = fmt.Errorf("Could not decode sender address: %v", err)
		return
	}

	tx, err := transaction.MakeApplicationCreateTxWithBoxes(optIn, approvalProg, clearProg, globalSchema, localSchema, uint32(extraPages), appArgs.Extract(), accounts.Extract(), internalForeignApps, internalForeignAssets, boxRefs.Extract(), internalParams, senderAddr, note, types.Digest{}, [32]byte{}, types.Address{})
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeApplicationUpdateTx makes a transaction for updating an application's programs (see above for args desc.)
func MakeApplicationUpdateTx(
	appIdx int64,
	appArgs *BytesArray,
	accounts *StringArray,
	foreignApps *Int64Array,
	foreignAssets *Int64Array,
	boxRefs *AppBoxRefArray,
	approvalProg []byte,
	clearProg []byte,
	params *SuggestedParams,
	sender string,
	note []byte,
) (encoded []byte, err error) {
	if appIdx < 0 {
		err = errNegativeArgument
		return
	}

	internalForeignApps := make([]uint64, foreignApps.Length())
	for i := range internalForeignApps {
		value := foreignApps.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignApps[i] = uint64(value)
	}

	internalForeignAssets := make([]uint64, foreignAssets.Length())
	for i := range internalForeignAssets {
		value := foreignAssets.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignAssets[i] = uint64(value)
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		err = fmt.Errorf("Could not decode sender address: %v", err)
		return
	}

	tx, err := transaction.MakeApplicationUpdateTxWithBoxes(uint64(appIdx), appArgs.Extract(), accounts.Extract(), internalForeignApps, internalForeignAssets, boxRefs.Extract(), approvalProg, clearProg, internalParams, senderAddr, note, types.Digest{}, [32]byte{}, types.Address{})
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeApplicationDeleteTx makes a transaction for deleting an application (see above for args desc.)
func MakeApplicationDeleteTx(
	appIdx int64,
	appArgs *BytesArray,
	accounts *StringArray,
	foreignApps *Int64Array,
	foreignAssets *Int64Array,
	boxRefs *AppBoxRefArray,
	params *SuggestedParams,
	sender string,
	note []byte,
) (encoded []byte, err error) {
	if appIdx < 0 {
		err = errNegativeArgument
		return
	}

	internalForeignApps := make([]uint64, foreignApps.Length())
	for i := range internalForeignApps {
		value := foreignApps.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignApps[i] = uint64(value)
	}

	internalForeignAssets := make([]uint64, foreignAssets.Length())
	for i := range internalForeignAssets {
		value := foreignAssets.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignAssets[i] = uint64(value)
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		err = fmt.Errorf("Could not decode sender address: %v", err)
		return
	}

	tx, err := transaction.MakeApplicationDeleteTxWithBoxes(uint64(appIdx), appArgs.Extract(), accounts.Extract(), internalForeignApps, internalForeignAssets, boxRefs.Extract(), internalParams, senderAddr, note, types.Digest{}, [32]byte{}, types.Address{})
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeApplicationOptInTx makes a transaction for opting in to (allocating
// some account-specific state for) an application (see above for args desc.)
func MakeApplicationOptInTx(
	appIdx int64,
	appArgs *BytesArray,
	accounts *StringArray,
	foreignApps *Int64Array,
	foreignAssets *Int64Array,
	boxRefs *AppBoxRefArray,
	params *SuggestedParams,
	sender string,
	note []byte,
) (encoded []byte, err error) {
	if appIdx < 0 {
		err = errNegativeArgument
		return
	}

	internalForeignApps := make([]uint64, foreignApps.Length())
	for i := range internalForeignApps {
		value := foreignApps.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignApps[i] = uint64(value)
	}

	internalForeignAssets := make([]uint64, foreignAssets.Length())
	for i := range internalForeignAssets {
		value := foreignAssets.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignAssets[i] = uint64(value)
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		err = fmt.Errorf("Could not decode sender address: %v", err)
		return
	}

	tx, err := transaction.MakeApplicationOptInTxWithBoxes(uint64(appIdx), appArgs.Extract(), accounts.Extract(), internalForeignApps, internalForeignAssets, boxRefs.Extract(), internalParams, senderAddr, note, types.Digest{}, [32]byte{}, types.Address{})
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeApplicationCloseOutTx makes a transaction for closing out of
// (deallocating all account-specific state for) an application (see above for args desc.)
func MakeApplicationCloseOutTx(
	appIdx int64,
	appArgs *BytesArray,
	accounts *StringArray,
	foreignApps *Int64Array,
	foreignAssets *Int64Array,
	boxRefs *AppBoxRefArray,
	params *SuggestedParams,
	sender string,
	note []byte,
) (encoded []byte, err error) {
	if appIdx < 0 {
		err = errNegativeArgument
		return
	}

	internalForeignApps := make([]uint64, foreignApps.Length())
	for i := range internalForeignApps {
		value := foreignApps.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignApps[i] = uint64(value)
	}

	internalForeignAssets := make([]uint64, foreignAssets.Length())
	for i := range internalForeignAssets {
		value := foreignAssets.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignAssets[i] = uint64(value)
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		err = fmt.Errorf("Could not decode sender address: %v", err)
		return
	}

	tx, err := transaction.MakeApplicationCloseOutTxWithBoxes(uint64(appIdx), appArgs.Extract(), accounts.Extract(), internalForeignApps, internalForeignAssets, boxRefs.Extract(), internalParams, senderAddr, note, types.Digest{}, [32]byte{}, types.Address{})
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeApplicationClearStateTx makes a transaction for clearing out all
// account-specific state for an application. It may not be rejected by the
// application's logic. (see above for args desc.)
func MakeApplicationClearStateTx(
	appIdx int64,
	appArgs *BytesArray,
	accounts *StringArray,
	foreignApps *Int64Array,
	foreignAssets *Int64Array,
	boxRefs *AppBoxRefArray,
	params *SuggestedParams,
	sender string,
	note []byte,
) (encoded []byte, err error) {
	if appIdx < 0 {
		err = errNegativeArgument
		return
	}

	internalForeignApps := make([]uint64, foreignApps.Length())
	for i := range internalForeignApps {
		value := foreignApps.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignApps[i] = uint64(value)
	}

	internalForeignAssets := make([]uint64, foreignAssets.Length())
	for i := range internalForeignAssets {
		value := foreignAssets.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignAssets[i] = uint64(value)
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		err = fmt.Errorf("Could not decode sender address: %v", err)
		return
	}

	tx, err := transaction.MakeApplicationClearStateTxWithBoxes(uint64(appIdx), appArgs.Extract(), accounts.Extract(), internalForeignApps, internalForeignAssets, boxRefs.Extract(), internalParams, senderAddr, note, types.Digest{}, [32]byte{}, types.Address{})
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeApplicationNoOpTx makes a transaction for interacting with an existing
// application, potentially updating any account-specific local state and
// global state associated with it. (see above for args desc.)
func MakeApplicationNoOpTx(
	appIdx int64,
	appArgs *BytesArray,
	accounts *StringArray,
	foreignApps *Int64Array,
	foreignAssets *Int64Array,
	boxRefs *AppBoxRefArray,
	params *SuggestedParams,
	sender string,
	note []byte,
) (encoded []byte, err error) {
	if appIdx < 0 {
		err = errNegativeArgument
		return
	}

	internalForeignApps := make([]uint64, foreignApps.Length())
	for i := range internalForeignApps {
		value := foreignApps.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignApps[i] = uint64(value)
	}

	internalForeignAssets := make([]uint64, foreignAssets.Length())
	for i := range internalForeignAssets {
		value := foreignAssets.Get(i)
		if value < 0 {
			err = errNegativeArgument
			return
		}
		internalForeignAssets[i] = uint64(value)
	}

	internalParams, err := convertSuggestedParams(params)
	if err != nil {
		return
	}

	senderAddr, err := types.DecodeAddress(sender)
	if err != nil {
		err = fmt.Errorf("Could not decode sender address: %v", err)
		return
	}

	tx, err := transaction.MakeApplicationNoOpTxWithBoxes(uint64(appIdx), appArgs.Extract(), accounts.Extract(), internalForeignApps, internalForeignAssets, boxRefs.Extract(), internalParams, senderAddr, note, types.Digest{}, [32]byte{}, types.Address{})
	if err == nil {
		encoded = msgpack.Encode(tx)
	}

	return
}

// MakeOptInAndAssetTransferTxns makes transactions for opting in to one asset
// and sending the asset to the opted-in receiver.
func MakeOptInAndAssetTransferTxns(
	sender,
	receiver string,
	transactionAmount,
	senderAlgoAmount,
	senderMinBalanceAmount,
	receiverAlgoAmount,
	receiverMinBalanceAmount *Uint64,
	note []byte,
	closeRemainderTo string,
	assetID int64,
	params *SuggestedParams,
) (transactions *TransactionSignerArray, err error) {
	var ASSET_OPT_IN_MBR uint64 = 100000
	var FLAT_FEE uint64 = 1000
	var ACCOUNT_MBR uint64 = 100000

	if params.Fee > 0 {
		FLAT_FEE = uint64(params.Fee)
	}

	senderExtractedAmount, _ := senderAlgoAmount.Extract()
	senderExtractedMinBalance, _ := senderMinBalanceAmount.Extract()
	receiverExtractedAmount, _ := receiverAlgoAmount.Extract()
	receiverExtractedMinBalance, _ := receiverMinBalanceAmount.Extract()

	// If receiver has enough algo to opt-in to the asset
	if receiverExtractedAmount-receiverExtractedMinBalance >= ASSET_OPT_IN_MBR+FLAT_FEE {
		optInAmount := MakeUint64(0)
		receiverOptInTxn, receiverOptInTxnError := MakeAssetTransferTxn(
			receiver,
			receiver,
			"",
			&optInAmount,
			nil,
			params,
			assetID,
		)
		if receiverOptInTxnError != nil {
			err = receiverOptInTxnError
			return
		}

		senderAssetTxn, senderAssetTxnError := MakeAssetTransferTxn(
			sender,
			receiver,
			"",
			transactionAmount,
			nil,
			params,
			assetID,
		)
		if senderAssetTxnError != nil {
			err = senderAssetTxnError
			return
		}

		bytesArrayTxns := BytesArray{values: [][]byte{receiverOptInTxn, senderAssetTxn}}
		assignedTxns, groupIDError := AssignGroupID(&bytesArrayTxns)
		if groupIDError != nil {
			err = groupIDError
			return
		}

		transactionSignerItem1 := TransactionSignerItem{
			signer:      receiver,
			transaction: assignedTxns.Get(0),
		}
		transactionSignerItem2 := TransactionSignerItem{
			signer:      sender,
			transaction: assignedTxns.Get(1),
		}
		signerArray := []TransactionSignerItem{transactionSignerItem1, transactionSignerItem2}

		transactionList := TransactionSignerArray{
			signerItems:  signerArray,
			transactions: assignedTxns,
		}
		transactions = &transactionList
		return
	} else {
		// If receiver does not have enough algo to opt-in to the asset

		var receiverExtraAlgoAmount uint64 = 0
		if receiverExtractedAmount == 0 {
			receiverExtraAlgoAmount = ACCOUNT_MBR + ASSET_OPT_IN_MBR + FLAT_FEE
		} else {
			algoAmount := (ASSET_OPT_IN_MBR + FLAT_FEE) - (receiverExtractedAmount - receiverExtractedMinBalance)
			if algoAmount > 0 {
				receiverExtraAlgoAmount = algoAmount
			}
		}

		if senderExtractedAmount-senderExtractedMinBalance < receiverExtractedAmount+FLAT_FEE+FLAT_FEE {
			err = fmt.Errorf("sender does not have enough algo to cover recivers needs")
			return
		}

		paymentAMount := MakeUint64(receiverExtraAlgoAmount)
		paymentTxn, paymentTxnError := MakePaymentTxn(
			sender,
			receiver,
			&paymentAMount,
			nil,
			"",
			params,
		)
		if paymentTxnError != nil {
			err = paymentTxnError
			return
		}

		optInAmount := MakeUint64(0)
		receiverOptInTxn, receiverOptInTxnError := MakeAssetTransferTxn(
			receiver,
			receiver,
			"",
			&optInAmount,
			nil,
			params,
			assetID,
		)
		if receiverOptInTxnError != nil {
			err = receiverOptInTxnError
			return
		}

		senderAssetTxn, senderAssetTxnError := MakeAssetTransferTxn(
			sender,
			receiver,
			"",
			transactionAmount,
			nil,
			params,
			assetID,
		)
		if senderAssetTxnError != nil {
			err = senderAssetTxnError
			return
		}

		bytesArrayTxns := BytesArray{values: [][]byte{paymentTxn, receiverOptInTxn, senderAssetTxn}}
		assignedTxns, groupIDError := AssignGroupID(&bytesArrayTxns)
		if groupIDError != nil {
			err = groupIDError
			return
		}
		transactionSignerItem1 := TransactionSignerItem{
			signer:      sender,
			transaction: assignedTxns.Get(0),
		}
		transactionSignerItem2 := TransactionSignerItem{
			signer:      receiver,
			transaction: assignedTxns.Get(1),
		}
		transactionSignerItem3 := TransactionSignerItem{
			signer:      sender,
			transaction: assignedTxns.Get(2),
		}
		signerArray := []TransactionSignerItem{transactionSignerItem1, transactionSignerItem2, transactionSignerItem3}

		transactionList := TransactionSignerArray{
			signerItems:  signerArray,
			transactions: assignedTxns,
		}
		transactions = &transactionList
		return
	}
}
