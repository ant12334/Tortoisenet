package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/ant12334/Tortoisenet/testutil/keeper"
	"github.com/ant12334/Tortoisenet/x/cons/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.ConsKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
