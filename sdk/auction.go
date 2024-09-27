package sdk

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/auction"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
)

func MakeBid(bidderAddress string, bidAmount, maxPrice, bidID int64, auctionAddress string, auctionID int64) (encodedBid []byte, err error) {
	// Sanity check for int64
	if bidAmount < 0 ||
		maxPrice < 0 ||
		bidID < 0 ||
		auctionID < 0 {
		err = fmt.Errorf("all numbers must not be negative")
		return
	}

	bid, err := auction.MakeBid(bidderAddress, uint64(bidAmount), uint64(maxPrice), uint64(bidID), auctionAddress, uint64(auctionID))
	if err != nil {
		return
	}

	encodedBid = msgpack.Encode(bid)
	return
}
