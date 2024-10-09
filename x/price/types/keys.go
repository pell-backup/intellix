package types

const (
	// ModuleName defines the module name
	ModuleName = "price"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_price"
)

var (
	ParamsKey = []byte("p_price")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
