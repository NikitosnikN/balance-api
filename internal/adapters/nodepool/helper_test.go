package nodepool

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDetermineTrueBalance(t *testing.T) {
	t.Parallel()

	t.Run("no elements in map", func(t *testing.T) {
		map_ := map[string]int{}

		result, err := determineTrueBalance(map_)

		require.Error(t, err)
		require.Equal(t, "", result)
	})
	t.Run("map has correct value", func(t *testing.T) {
		map_ := map[string]int{
			"x": 3,
			"y": 1,
			"z": 1,
			"a": 0,
		}

		result, err := determineTrueBalance(map_)

		require.NoError(t, err)
		require.Equal(t, "x", result)
	})
	t.Run("map has duplicated value", func(t *testing.T) {
		map_ := map[string]int{
			"y": 2,
			"z": 2,
			"a": 0,
		}

		_, err := determineTrueBalance(map_)

		require.Error(t, err)
	})
}
