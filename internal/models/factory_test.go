package models

import (
	"testing"

	"github.com/smartystreets/assertions"

	"github.com/aaron-zeisler/azul/internal/testutils"
)

func TestFactory_DrawAllTilesByColor(t *testing.T) {
	type state struct {
		tiles []Tile
		color TileColor
	}
	type expected struct {
		result         []Tile
		err            error
		remainingTiles []Tile
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"Error case: the factory is empty": {
			state{
				tiles: []Tile{},
				color: Orange,
			},
			expected{
				err:            InvalidActionError{Message: "This factory has no tiles"},
				remainingTiles: []Tile{},
			},
		},
		"Error case: the factory doesn't have any tiles of the specified color": {
			state{
				tiles: []Tile{{Color: Black}, {Color: Black}, {Color: Blue}},
				color: Orange,
			},
			expected{
				err:            InvalidActionError{"There are no orange tiles on this factory"},
				remainingTiles: []Tile{{Color: Black}, {Color: Black}, {Color: Blue}},
			},
		},
		"The factory has one of the specified tiles": {
			state{
				tiles: []Tile{{Color: Blue}, {Color: Orange}, {Color: White}},
				color: Orange,
			},
			expected{
				result:         []Tile{{Color: Orange}},
				remainingTiles: []Tile{{Color: Blue}, {Color: White}},
			},
		},
		"The factory has two of the specified tiles": {
			state{
				tiles: []Tile{{Color: Orange}, {Color: Blue}, {Color: Orange}, {Color: White}},
				color: Orange,
			},
			expected{
				result:         []Tile{{Color: Orange}, {Color: Orange}},
				remainingTiles: []Tile{{Color: Blue}, {Color: White}},
			},
		},
		"The factory is full of the specified tiles": {
			state{
				tiles: []Tile{{Color: Orange}, {Color: Orange}, {Color: Orange}},
				color: Orange,
			},
			expected{
				result:         []Tile{{Color: Orange}, {Color: Orange}, {Color: Orange}},
				remainingTiles: []Tile{},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			f := NewFactory()
			for _, tile := range tc.state.tiles {
				f.AddTile(tile)
			}

			result, err := f.DrawAllTilesByColor(tc.state.color)

			assert.So(err, testutils.ShouldEqualError, tc.expected.err)
			if tc.expected.err == nil {
				assert.So(result, shouldEqualTileSlice, tc.expected.result)
			}
			assert.So(f.Tiles, shouldEqualTileSlice, tc.expected.remainingTiles)
		})
	}
}

func shouldEqualTileSlice(actual interface{}, expected ...interface{}) string {
	if lenResult := assertions.ShouldEqual(len(actual.([]Tile)), len(expected[0].([]Tile))); lenResult != "" {
		return lenResult
	}
	for _, tile := range actual.([]Tile) {
		if containsResult := assertions.ShouldBeIn(tile, expected[0].([]Tile)); containsResult != "" {
			return containsResult
		}
	}
	return ""
}
