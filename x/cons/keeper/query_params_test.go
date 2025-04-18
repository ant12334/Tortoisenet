package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/ant12334/Tortoisenet/testutil/keeper"
	"github.com/ant12334/Tortoisenet/x/cons/types"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := keepertest.ConsKeeper(t)
	params := types.DefaultParams()
	require.NoError(t, keeper.SetParams(ctx, params))

	response, err := keeper.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
