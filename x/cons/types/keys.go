package types

const (
	// ModuleName defines the module name
	ModuleName = "cons"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_cons"
)

var (
	ParamsKey = []byte("p_cons")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
