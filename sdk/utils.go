package sdk

// https://github.com/algorand/go-algorand-sdk/compare/a140151ac15136f234b9cf094e58cb69cb37e2c0...MobileCompatible

import (
	"fmt"
	"math"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

type Uint64 struct {
	Upper int64
	Lower int64
}

func MakeUint64(value uint64) Uint64 {
	return Uint64{
		Upper: int64(value >> 32),
		Lower: int64(math.MaxUint32 & value),
	}
}

func (i Uint64) Extract() (value uint64, err error) {
	if i.Upper < 0 || i.Upper > int64(math.MaxUint32) {
		err = fmt.Errorf("Upper value of Uint64 not in correct range. Expected value between 0 and %d, got %d", int64(math.MaxUint32), i.Upper)
		return
	}

	if i.Lower < 0 || i.Lower > int64(math.MaxUint32) {
		err = fmt.Errorf("Lower value of Uint64 not in correct range. Expected value between 0 and %d, got %d", int64(math.MaxUint32), i.Lower)
		return
	}

	value = uint64(i.Upper)<<32 | uint64(i.Lower)

	return
}

type TransactionSignerArray struct {
	signerItems  []TransactionSignerItem
	transactions *BytesArray
}

func (tsa *TransactionSignerArray) Length() int {
	return len(tsa.signerItems)
}

func (tsa *TransactionSignerArray) GetSigner(index int) string {
	return tsa.signerItems[index].signer
}

func (tsa *TransactionSignerArray) GetTxnFromSigner(index int) []byte {
	return tsa.signerItems[index].transaction
}

func (tsa *TransactionSignerArray) GetAssignedFlattenTxns() []byte {
	return tsa.transactions.Flatten()
}

func (tsa *TransactionSignerArray) ExtractAssignedFlattenTxns() [][]byte {
	return tsa.transactions.Extract()
}

func (tsa *TransactionSignerArray) GetTxn(index int) []byte {
	return tsa.transactions.Get(index)
}

type TransactionSignerItem struct {
	signer      string
	transaction []byte
}

type StringArray struct {
	values []string
}

func (sa *StringArray) Length() int {
	return len(sa.values)
}

func (sa *StringArray) Append(value string) {
	sa.values = append(sa.values, string([]byte(value))) // deep copy the string
}

func (sa *StringArray) Get(index int) string {
	return sa.values[index]
}

func (sa *StringArray) Set(index int, value string) {
	sa.values[index] = string([]byte(value)) // deep copy the string
}

func (sa *StringArray) Extract() []string {
	return sa.values[:]
}

type BytesArray struct {
	values [][]byte
}

func (ba *BytesArray) Length() int {
	return len(ba.values)
}

func (ba *BytesArray) Append(value []byte) {
	cp := make([]byte, len(value))
	copy(cp, value)
	ba.values = append(ba.values, cp)
}

func (ba *BytesArray) Get(index int) []byte {
	return ba.values[index]
}

func (ba *BytesArray) Set(index int, value []byte) {
	cp := make([]byte, len(value))
	copy(cp, value)
	ba.values[index] = cp
}

// Flatten returns a single byte array containing all the contained byte arrays, in order.
func (ba *BytesArray) Flatten() []byte {
	newLength := 0
	for _, value := range ba.values {
		newLength += len(value)
	}
	result := make([]byte, 0, newLength)
	for _, value := range ba.values {
		result = append(result, value...)
	}
	return result
}

func (ba *BytesArray) Extract() [][]byte {
	return ba.values[:]
}

type Int64Array struct {
	values []int64
}

func (ia *Int64Array) Length() int {
	return len(ia.values)
}

func (ia *Int64Array) Append(value int64) {
	ia.values = append(ia.values, value)
}

func (ia *Int64Array) Get(index int) int64 {
	return ia.values[index]
}

func (ia *Int64Array) Set(index int, value int64) {
	ia.values[index] = value
}

func (ia *Int64Array) Extract() []int64 {
	return ia.values[:]
}

type AppBoxRefArray struct {
	value []types.AppBoxReference
}

func (ba *AppBoxRefArray) Length() int {
	return len(ba.value)
}

func (ba *AppBoxRefArray) Append(appID int64, boxName []byte) error {
	if appID < 0 {
		return fmt.Errorf("appID must be positive: %d", appID)
	}
	ba.value = append(ba.value, types.AppBoxReference{AppID: uint64(appID), Name: boxName})
	return nil
}

func (ba *AppBoxRefArray) GetAppID(index int) int64 {
	return int64(ba.value[index].AppID)
}

func (ba *AppBoxRefArray) GetBoxName(index int) []byte {
	return ba.value[index].Name
}

func (ba *AppBoxRefArray) Set(index int, appID int64, boxName []byte) error {
	if appID < 0 {
		return fmt.Errorf("appID must be positive: %d", appID)
	}
	ba.value[index] = types.AppBoxReference{AppID: uint64(appID), Name: boxName}
	return nil
}

func (ba *AppBoxRefArray) Extract() []types.AppBoxReference {
	return ba.value[:]
}

func IsValidAddress(addr string) bool {
	_, err := types.DecodeAddress(addr)
	return err == nil
}
