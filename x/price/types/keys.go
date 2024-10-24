package types

import "fmt"

const (
	// ModuleName defines the module name
	ModuleName = "price"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_price"

	PriceFeedVoteKeyPrefix        = "price_feed_vote"
	EventTypeVoteRequestPriceFeed = "vote_request_price_feed"
	AttributeKeyTaskIndex         = "task_index"
	AttributeKeyOperatorId        = "operator_id"
	AttributeKeyRequestId         = "request_id"

	FinalizedPriceFeedKeyPrefix        = "finalized_price_feed"
	EventTypeVoteFinalizedRequestPrice = "vote_finalized_request_price"
)

var (
	ParamsKey = []byte("p_price")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

func PriceFeedVoteKey(taskIndex uint32, operatorId string) []byte {
	return append(
		append(KeyPrefix(PriceFeedVoteKeyPrefix), []byte(fmt.Sprintf("%d", taskIndex))...),
		[]byte(operatorId)...,
	)
}

func FinalizedRequestPrice(taskIndex uint32) []byte {
	return append(KeyPrefix(FinalizedPriceFeedKeyPrefix), []byte(fmt.Sprintf("%d", taskIndex))...)
}
