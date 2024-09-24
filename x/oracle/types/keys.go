package types

const (
	// ModuleName defines the module name
	ModuleName = "oracle"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_oracle"
)

var (
	ParamsKey = []byte("p_oracle")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
