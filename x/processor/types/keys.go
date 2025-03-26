package types

import "fmt"

const (
	// ModuleName defines the module name
	ModuleName = "processor"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_processor"

	ProcessorKey      = "Processor/value/"
	ProcessorCountKey = "Processor/count"

	EventTypeVoteRequestProcessor  = "vote_request_processor"
	EventTypeVoteResponseProcessor = "vote_response_processor"
	AttributeKeyTaskIndex          = "task_index"
	AttributeKeyOperatorId         = "operator_id"
	AttributeKeyRequestId          = "request_id"
)

var (
	ParamsKey = []byte("p_processor")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

// MsgVoteRequestProcessorKey returns the key for the vote request processor
func MsgVoteRequestProcessorKey(requestId []byte, operatorId string) []byte {
	return append(
		append(
			KeyPrefix(ProcessorKey+EventTypeVoteRequestProcessor),
			[]byte(operatorId)...,
		),
		requestId...,
	)
}

func MsgVoteResponseProcessorKey(taskIndex uint32, requestId []byte) []byte {
	return append(
		append(KeyPrefix(ProcessorCountKey+EventTypeVoteResponseProcessor),
			[]byte(fmt.Sprintf("%d", taskIndex))...),
		requestId...,
	)
}
