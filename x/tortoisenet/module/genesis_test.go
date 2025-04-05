package tortoisenet_test

import (
	"testing"

	keepertest "github.com/ant12334/Tortoisenet/testutil/keeper"
	"github.com/ant12334/Tortoisenet/testutil/nullify"
	tortoisenet "github.com/ant12334/Tortoisenet/x/tortoisenet/module"
	"github.com/ant12334/Tortoisenet/x/tortoisenet/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.TortoisenetKeeper(t)
	tortoisenet.InitGenesis(ctx, k, genesisState)
	got := tortoisenet.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
