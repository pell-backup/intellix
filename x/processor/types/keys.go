package types

const (
	// ModuleName defines the module name
	ModuleName = "processor"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_processor"

	ProcessorKey = "Processor/value/"
	ProcessorCountKey = "Processor/count"
)

var (
	ParamsKey = []byte("p_processor")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
