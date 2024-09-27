package sdk

import (
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// MakeARC59OptInTxn creates the app call transaction for opting into ARC59 protocol.
func MakeARC59OptInTxn(
	sender,
	appAddress string,
	appID,
	assetID int64,
	suggestedParams *SuggestedParams,
) (assignedTxns *BytesArray, err error) {
	appCallSuggestedParams := *suggestedParams
	appCallSuggestedParams.FlatFee = true
	if appCallSuggestedParams.Fee == 0 {
		appCallSuggestedParams.Fee = 2 * 1000
	} else {
		appCallSuggestedParams.Fee *= 2
	}

	amount := MakeUint64(100000)
	paymentTxn, paymentTxnError := MakePaymentTxn(
		sender,
		appAddress,
		&amount,
		nil,
		"",
		suggestedParams,
	)
	if paymentTxnError != nil {
		err = paymentTxnError
		return
	}

	appArgumentsByteArray, appArgumentsError := MakeAppArgumentsByteArrayWithAsset(
		"arc59_optRouterIn(uint64)void",
		assetID,
	)
	if appArgumentsError != nil {
		err = appArgumentsError
		return
	}

	accountStringArray := StringArray{values: []string{appAddress}}
	assetInt64Array := Int64Array{values: []int64{assetID}}

	appCallTxn, appCallTxnError := MakeApplicationNoOpTx(
		appID,
		&appArgumentsByteArray,
		&accountStringArray,
		&Int64Array{values: []int64{}}, // empty array
		&assetInt64Array,
		&AppBoxRefArray{value: []types.AppBoxReference{}}, // empty ref
		&appCallSuggestedParams,
		sender,
		nil,
	)
	if appCallTxnError != nil {
		err = appCallTxnError
		return
	}

	bytesArrayTxns := BytesArray{values: [][]byte{paymentTxn, appCallTxn}}
	assignedTxns, err = AssignGroupID(&bytesArrayTxns)
	return
}

// MakeAndSignARC59OptInTxn creates the app call transaction for opting into ARC59 protocol and returns the signed transaction.
func MakeAndSignARC59OptInTxn(
	sender,
	appAddress string,
	appID,
	assetID int64,
	suggestedParams *SuggestedParams,
	sk []byte,
) (signedTxns *BytesArray, err error) {
	groupTxns, groupTxnError := MakeARC59OptInTxn(
		sender,
		appAddress,
		appID,
		assetID,
		suggestedParams,
	)
	if groupTxnError != nil {
		err = groupTxnError
		return
	}

	// Sign each transaction individually
	txns := make([][]byte, len(groupTxns.values))
	for i, txn := range groupTxns.values {
		signedTxn, signError := SignTransaction(sk, txn)
		if signError != nil {
			err = signError
			return
		}
		txns[i] = signedTxn
	}


	signedTxns = &BytesArray{values: txns}
	return
}

// MakeARC59SendTxn creates the payment, asset transfer and app call transactions for sending an asset with the ARC59 protocol.
func MakeARC59SendTxn(
	sender,
	receiver,
	appAddress,
	inboxAccountAddressOrEmptyString string,
	amount,
	minimumBalanceRequirement *Uint64,
	innerTxCount,
	appID,
	assetID int64,
	suggestedParams *SuggestedParams,
	extraAlgoAmount *Uint64,
) (assignedTxns *BytesArray, err error) {
	appCallSuggestedParams := *suggestedParams
	appCallSuggestedParams.FlatFee = true
	if appCallSuggestedParams.Fee == 0 {
		appCallSuggestedParams.Fee = (innerTxCount + 1) * 1000
	} else {
		appCallSuggestedParams.Fee *= (innerTxCount + 1)
	}

	decodedAlgoAmount, err := extraAlgoAmount.Extract()
	decodedTxnAmount, err := minimumBalanceRequirement.Extract()
	totalTxnAmount := decodedAlgoAmount + decodedTxnAmount
	txnPaymentAmount := MakeUint64(totalTxnAmount)

	if decodedAlgoAmount > 0 {
		appCallSuggestedParams.Fee += 1000
	}

	paymentTxn, paymentTxnError := MakePaymentTxn(
		sender,
		appAddress,
		&txnPaymentAmount,
		nil,
		"",
		suggestedParams,
	)
	if paymentTxnError != nil {
		err = paymentTxnError
		return
	}

	assetTxn, assetTxnError := MakeAssetTransferTxn(
		sender,
		appAddress,
		"",
		amount,
		nil,
		suggestedParams,
		assetID,
	)
	if assetTxnError != nil {
		err = assetTxnError
		return
	}

	receiverDecoded, err := DecodeAddress(receiver)

	appArgumentsByteArray, appArgumentsError := MakeAppArgumentsByteArrayWithAddressAndAmount(
		"arc59_sendAsset(axfer,address,uint64)address",
		receiverDecoded,
		decodedAlgoAmount,
	)
	if appArgumentsError != nil {
		err = appArgumentsError
		return
	}

	foreignAccounts := []string{receiver}
	if inboxAccountAddressOrEmptyString != "" {
		foreignAccounts = append(foreignAccounts, inboxAccountAddressOrEmptyString)
	}

	accountStringArray := StringArray{values: foreignAccounts}
	assetInt64Array := Int64Array{values: []int64{assetID}}
	boxRefArray := MakeAppBoxRefArray(uint64(appID), receiverDecoded)

	appCallTxn, appCallTxnError := MakeApplicationNoOpTx(
		appID,
		&appArgumentsByteArray,
		&accountStringArray,
		&Int64Array{values: []int64{}}, // empty array
		&assetInt64Array,
		&boxRefArray,
		&appCallSuggestedParams,
		sender,
		nil,
	)
	if appCallTxnError != nil {
		err = appCallTxnError
		return
	}

	bytesArrayTxns := BytesArray{values: [][]byte{paymentTxn, assetTxn, appCallTxn}}

	assignedTxns, err = AssignGroupID(&bytesArrayTxns)
	return
}

// MakeAndSignARC59SendTxn creates the payment, asset transfer and app call transactions for sending an asset with
// the ARC59 protocol and signs these transactions.
func MakeAndSignARC59SendTxn(
	sender,
	receiver,
	appAddress,
	inboxAccountAddressOrEmptyString string,
	amount,
	minimumBalanceRequirement *Uint64,
	innerTxCount,
	appID,
	assetID int64,
	suggestedParams *SuggestedParams,
	extraAlgoAmount *Uint64,
	sk []byte,
) (signedTxns *BytesArray, err error) {
	groupTxns, groupTxnError := MakeARC59SendTxn(
		sender,
		receiver,
		appAddress,
		inboxAccountAddressOrEmptyString,
		amount,
		minimumBalanceRequirement,
		innerTxCount,
		appID,
		assetID,
		suggestedParams,
		extraAlgoAmount,
	)
	if groupTxnError != nil {
		err = groupTxnError
		return
	}

	// Sign each transaction individually
	txns := make([][]byte, len(groupTxns.values))
	for i, txn := range groupTxns.values {
		signedTxn, signError := SignTransaction(sk, txn)
		if signError != nil {
			err = signError
			return
		}
		txns[i] = signedTxn
	}

	signedTxns = &BytesArray{values: txns}
	return
}

// MakeARC59ClaimTxn creates the app call transaction, and opt in transaction if needed, to claim the asset from the ARC59 protocol.
func MakeARC59ClaimTxn(
	receiver,
	inboxAccountAddress string,
	appID,
	assetID int64,
	suggestedParams *SuggestedParams,
	isOptedInToAsset,
	isClaimingAlgo bool,
) (assignedTxns *BytesArray, err error) {
	appCallSuggestedParams := *suggestedParams
	appCallSuggestedParams.FlatFee = true
	if appCallSuggestedParams.Fee == 0 {
		appCallSuggestedParams.Fee = 3 * 1000
	} else {
		appCallSuggestedParams.Fee *= 3
	}

	if isClaimingAlgo {
		appCallSuggestedParams.Fee += 2 * 1000
	}

	var inboxAccountStringArray StringArray
	if inboxAccountAddress == "" {
		inboxAccountStringArray = StringArray{values: []string{}}
	} else {
		inboxAccountStringArray = StringArray{values: []string{inboxAccountAddress}}
	}

	decodedReceiver, err := DecodeAddress(receiver)

	claimAlgoTxnSuggestedParams := *suggestedParams
	claimAlgoTxnSuggestedParams.FlatFee = true
	claimAlgoTxnSuggestedParams.Fee = 0

	methodNameHex := MethodName("arc59_claimAlgo()void")
	methodNameBytes, err := hex.DecodeString(methodNameHex)
	methodNameAppArgs := [][]byte{methodNameBytes}
	methodNameBytesArray := BytesArray{values: methodNameAppArgs}

	claimAlgoTxnBoxReferenceItem := types.AppBoxReference{
		Name:  decodedReceiver[:],
	}
	claimAlgoTxnBoxReferenceArray := []types.AppBoxReference{claimAlgoTxnBoxReferenceItem}
	claimAlgoAppBoxRefArray := AppBoxRefArray{value: claimAlgoTxnBoxReferenceArray}

	claimAlgoAppCallTxn, claimAlgoAppCallTxnError := MakeApplicationNoOpTx(
		appID,
		&methodNameBytesArray,
		&inboxAccountStringArray,
		&Int64Array{values: []int64{}}, // empty array
		&Int64Array{values: []int64{}}, // empty array
		&claimAlgoAppBoxRefArray,
		&claimAlgoTxnSuggestedParams,
		receiver,
		nil,
	)
	if claimAlgoAppCallTxnError != nil {
		err = claimAlgoAppCallTxnError
		return
	}

	appArgumentsByteArray, appArgumentsError := MakeAppArgumentsByteArrayWithAsset(
		"arc59_claim(uint64)void",
		assetID,
	)
	if appArgumentsError != nil {
		err = appArgumentsError
		return
	}

	assetInt64Array := Int64Array{values: []int64{assetID}}
	boxRefArray := MakeAppBoxRefArray(uint64(appID), decodedReceiver)

	appCallTxn, appCallTxnError := MakeApplicationNoOpTx(
		appID,
		&appArgumentsByteArray,
		&inboxAccountStringArray,
		&Int64Array{values: []int64{}}, // empty array
		&assetInt64Array,
		&boxRefArray,
		&appCallSuggestedParams,
		receiver,
		nil,
	)
	if appCallTxnError != nil {
		err = appCallTxnError
		return
	}

	bytesArrayTxns := BytesArray{values: [][]byte{}}
	if !isOptedInToAsset {
		optInAmount := MakeUint64(0)
		optInTxn, assetTxnError := MakeAssetTransferTxn(
			receiver,
			receiver,
			"",
			&optInAmount,
			nil,
			suggestedParams,
			assetID,
		)
		if assetTxnError != nil {
			err = assetTxnError
			return
		}

		if isClaimingAlgo {
			bytesArrayTxns = BytesArray{values: [][]byte{claimAlgoAppCallTxn, optInTxn, appCallTxn}}
		} else {
			bytesArrayTxns = BytesArray{values: [][]byte{optInTxn, appCallTxn}}
		}
	} else {
		if isClaimingAlgo {
			bytesArrayTxns = BytesArray{values: [][]byte{claimAlgoAppCallTxn, appCallTxn}}
		} else {
			bytesArrayTxns = BytesArray{values: [][]byte{appCallTxn}}
		}
	}

	assignedTxns, err = AssignGroupID(&bytesArrayTxns)
	return
}

// MakeAndSignARC59ClaimTxn creates the app call transaction, and opt in transaction if needed, to claim the asset from
// the ARC59 protocol and signs these transactions.
func MakeAndSignARC59ClaimTxn(
	receiver,
	inboxAccountAddress string,
	appID,
	assetID int64,
	suggestedParams *SuggestedParams,
	isOptedInToAsset,
	isClaimingAlgo bool,
	sk []byte,
) (signedTxns *BytesArray, err error) {
	groupTxns, groupTxnError := MakeARC59ClaimTxn(
		receiver,
		inboxAccountAddress,
		appID,
		assetID,
		suggestedParams,
		isOptedInToAsset,
		isClaimingAlgo,
	)
	if groupTxnError != nil {
		err = groupTxnError
		return
	}

	// Sign each transaction individually
	txns := make([][]byte, len(groupTxns.values))
	for i, txn := range groupTxns.values {
		signedTxn, signError := SignTransaction(sk, txn)
		if signError != nil {
			err = signError
			return
		}
		txns[i] = signedTxn
	}

	signedTxns = &BytesArray{values: txns}
	return
}

// MakeARC59RejectTxn creates the app call transaction to reject the asset from the ARC59 protocol.
func MakeARC59RejectTxn(
	receiver,
	inboxAccountAddress,
	creatorAccountAddress string,
	appID,
	assetID int64,
	suggestedParams *SuggestedParams,
	isClaimingAlgo bool,
) (assignedTxns *BytesArray, err error) {
	appCallSuggestedParams := *suggestedParams
	appCallSuggestedParams.FlatFee = true
	if appCallSuggestedParams.Fee == 0 {
		appCallSuggestedParams.Fee = 3 * 1000
	} else {
		appCallSuggestedParams.Fee *= 3
	}

	if isClaimingAlgo {
		appCallSuggestedParams.Fee += 2 * 1000
	}

	var inboxAccountStringArray StringArray
	if inboxAccountAddress == "" {
		inboxAccountStringArray = StringArray{values: []string{}}
	} else {
		inboxAccountStringArray = StringArray{values: []string{inboxAccountAddress}}
	}

	decodedReceiver, err := DecodeAddress(receiver)

	claimAlgoTxnSuggestedParams := *suggestedParams
	claimAlgoTxnSuggestedParams.FlatFee = true
	claimAlgoTxnSuggestedParams.Fee = 0

	methodNameHex := MethodName("arc59_claimAlgo()void")
	methodNameBytes, err := hex.DecodeString(methodNameHex)
	methodNameAppArgs := [][]byte{methodNameBytes}
	methodNameBytesArray := BytesArray{values: methodNameAppArgs}

	claimAlgoTxnBoxReferenceItem := types.AppBoxReference{
		Name:  decodedReceiver[:],
	}
	claimAlgoTxnBoxReferenceArray := []types.AppBoxReference{claimAlgoTxnBoxReferenceItem}
	claimAlgoAppBoxRefArray := AppBoxRefArray{value: claimAlgoTxnBoxReferenceArray}

	claimAlgoAppCallTxn, claimAlgoAppCallTxnError := MakeApplicationNoOpTx(
		appID,
		&methodNameBytesArray,
		&inboxAccountStringArray,
		&Int64Array{values: []int64{}}, // empty array
		&Int64Array{values: []int64{}}, // empty array
		&claimAlgoAppBoxRefArray,
		&claimAlgoTxnSuggestedParams,
		receiver,
		nil,
	)
	if claimAlgoAppCallTxnError != nil {
		err = claimAlgoAppCallTxnError
		return
	}

	appArgumentsByteArray, appArgumentsError := MakeAppArgumentsByteArrayWithAsset(
		"arc59_reject(uint64)void",
		assetID,
	)
	if appArgumentsError != nil {
		err = appArgumentsError
		return
	}

	assetInt64Array := Int64Array{values: []int64{assetID}}
	boxRefArray := MakeAppBoxRefArray(uint64(appID), decodedReceiver)

	var accountStringArray StringArray
	if inboxAccountAddress == "" {
		accountStringArray = StringArray{values: []string{creatorAccountAddress}}
	} else {
		accountStringArray = StringArray{values: []string{inboxAccountAddress, creatorAccountAddress}}
	}

	appCallTxn, err := MakeApplicationNoOpTx(
		appID,
		&appArgumentsByteArray,
		&accountStringArray,
		&Int64Array{values: []int64{}}, // empty array
		&assetInt64Array,
		&boxRefArray,
		&appCallSuggestedParams,
		receiver,
		nil,
	)

	bytesArrayTxns := BytesArray{values: [][]byte{}}
	if isClaimingAlgo {
		bytesArrayTxns = BytesArray{values: [][]byte{claimAlgoAppCallTxn, appCallTxn}}
	} else {
		bytesArrayTxns = BytesArray{values: [][]byte{appCallTxn}}
	}

	assignedTxns, err = AssignGroupID(&bytesArrayTxns)
	return
}

// MakeAndSignARC59RejectTxn creates the app call transaction to reject the asset from the ARC59 protocol and signs the transaction.
func MakeAndSignARC59RejectTxn(
	receiver,
	inboxAccountAddress,
	creatorAccountAddress string,
	appID,
	assetID int64,
	suggestedParams *SuggestedParams,
	isClaimingAlgo bool,
	sk []byte,
) (signedTxns *BytesArray, err error) {
	groupTxns, groupTxnError := MakeARC59RejectTxn(
		receiver,
		inboxAccountAddress,
		creatorAccountAddress,
		appID,
		assetID,
		suggestedParams,
		isClaimingAlgo,
	)
	if groupTxnError != nil {
		err = groupTxnError
		return
	}

	// Sign each transaction individually
	txns := make([][]byte, len(groupTxns.values))
	for i, txn := range groupTxns.values {
		signedTxn, signError := SignTransaction(sk, txn)
		if signError != nil {
			err = signError
			return
		}
		txns[i] = signedTxn
	}

	signedTxns = &BytesArray{values: txns}
	return
}

// Helper Functions

// MethodName converts the text to encoded hex method name
func MethodName(text string) string {
	hash := sha512.New512_256()
	hash.Write([]byte(text))
	checksum := hash.Sum(nil)
	truncated := checksum[:32]
	hexDigs := hex.EncodeToString(truncated)
	return hexDigs[:8]
}

// DecodeAddress converts the address as decoded.
func DecodeAddress(address string) (types.Address, error) {
	return types.DecodeAddress(address)
}

// MakeAppArgumentsByteArrayWithAddressAndAmount creates the app arguments for app call transaction with address and algo amount, then converts to the BytesArray.
func MakeAppArgumentsByteArrayWithAddressAndAmount(method string, decodedReceiver types.Address, amount uint64) (bytes BytesArray, err error) {
	methodNameHex := MethodName(method)
	methodNameBytes, err := hex.DecodeString(methodNameHex)

	appArgs := [][]byte{
		methodNameBytes,
		decodedReceiver[:],
		EncodeUIntAsBytes(amount),
	}
	bytes = BytesArray{values: appArgs}
	return
}

func EncodeIntAsBytes(e int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(e))
	return b
}

func EncodeUIntAsBytes(e uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, e)
	return b
}

// MakeAppArgumentsByteArrayWithAsset creates the app arguments for app call transaction with asset ID and converts to the BytesArray.
func MakeAppArgumentsByteArrayWithAsset(method string, assetID int64) (bytes BytesArray, err error) {
	methodNameHex := MethodName(method)

	methodNameBytes, err := hex.DecodeString(methodNameHex)

	appArgs := [][]byte{
		methodNameBytes,
		EncodeIntAsBytes(assetID),
	}
	bytes = BytesArray{values: appArgs}
	return
}

// MakeAppBoxRefArray converts the type of app box reference to the AppBoxRefArray.
func MakeAppBoxRefArray(appID uint64, decodedReceiver types.Address) AppBoxRefArray {
	boxReferenceItem := types.AppBoxReference{
		AppID: appID,
		Name:  decodedReceiver[:],
	}
	boxReferenceArray := []types.AppBoxReference{boxReferenceItem}
	return AppBoxRefArray{value: boxReferenceArray}
}
