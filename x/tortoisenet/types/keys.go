package types

const (
	// ModuleName defines the module name
	ModuleName = "tortoisenet"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_tortoisenet"
)

var (
	ParamsKey = []byte("p_tortoisenet")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
