package types

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

func PriceFeedVoteKey(requestId []byte, operatorId string) []byte {
	return append(
		append(KeyPrefix(PriceFeedVoteKeyPrefix), requestId...),
		[]byte(operatorId)...,
	)
}

func FinalizedRequestPrice(requestId []byte) []byte {
	return append(KeyPrefix(FinalizedPriceFeedKeyPrefix), requestId...)
}
