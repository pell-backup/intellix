package types

const (
	// ModuleName defines the module name
	ModuleName = "intellix"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_intellix"
)

var (
	ParamsKey = []byte("p_intellix")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
