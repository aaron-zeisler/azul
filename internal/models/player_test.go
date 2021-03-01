package models

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

func TestBoard_ScoreTile(t *testing.T) {
	type state struct {
		filledWallTiles []WallCoordinate
		scoringTile     WallCoordinate
	}
	type expected struct {
		result WallScore
	}
	testCases := map[string]struct {
		state    state
		expected expected
	}{
		"Empty wall, no adjacent tiles, should be 1 point": {
			state{
				scoringTile: WallCoordinate{Row: 2, Col: 3},
			},
			expected{
				WallScore{
					Score: 1,
					Tiles: []WallCoordinate{{Row: 2, Col: 3}},
				},
			},
		},
		"There are other tiles, but no adjacent tiles, should be 1 point": {
			state{
				filledWallTiles: []WallCoordinate{
					{Row: 0, Col: 0}, {Row: 0, Col: 4}, {Row: 4, Col: 0}, {Row: 4, Col: 4},
				},
				scoringTile: WallCoordinate{Row: 2, Col: 3},
			},
			expected{
				WallScore{
					Score: 1,
					Tiles: []WallCoordinate{{Row: 2, Col: 3}},
				},
			},
		},
		"One adjacent tile to the left": {
			state{
				filledWallTiles: []WallCoordinate{
					{Row: 2, Col: 2},
				},
				scoringTile: WallCoordinate{Row: 2, Col: 3},
			},
			expected{
				WallScore{
					Score: 2,
					Tiles: []WallCoordinate{{Row: 2, Col: 2}, {Row: 2, Col: 3}},
				},
			},
		},
		"One adjacent tile to the right": {
			state{
				filledWallTiles: []WallCoordinate{
					{Row: 2, Col: 4},
				},
				scoringTile: WallCoordinate{Row: 2, Col: 3},
			},
			expected{
				WallScore{
					Score: 2,
					Tiles: []WallCoordinate{{Row: 2, Col: 4}, {Row: 2, Col: 3}},
				},
			},
		},
		"One adjacent tile above": {
			state{
				filledWallTiles: []WallCoordinate{
					{Row: 1, Col: 3},
				},
				scoringTile: WallCoordinate{Row: 2, Col: 3},
			},
			expected{
				WallScore{
					Score: 2,
					Tiles: []WallCoordinate{{Row: 1, Col: 3}, {Row: 2, Col: 3}},
				},
			},
		},
		"One adjacent tile below": {
			state{
				filledWallTiles: []WallCoordinate{
					{Row: 3, Col: 3},
				},
				scoringTile: WallCoordinate{Row: 2, Col: 3},
			},
			expected{
				WallScore{
					Score: 2,
					Tiles: []WallCoordinate{{Row: 3, Col: 3}, {Row: 2, Col: 3}},
				},
			},
		},
		"Three in a row horizontally": {
			state{
				filledWallTiles: []WallCoordinate{
					{Row: 2, Col: 2}, {Row: 2, Col: 4},
				},
				scoringTile: WallCoordinate{Row: 2, Col: 3},
			},
			expected{
				WallScore{
					Score: 3,
					Tiles: []WallCoordinate{{Row: 2, Col: 2}, {Row: 2, Col: 4}, {Row: 2, Col: 3}},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assertions.New(t)

			b := &Board{}
			b.ResetWall()
			for _, tile := range tc.state.filledWallTiles {
				b.Wall[tile.Row][tile.Col].HasTile = true
			}

			result := b.ScoreTile(tc.state.scoringTile)

			assert.So(result.Score, should.Resemble, tc.expected.result.Score)
			//TODO: Assert the result.Tiles slice
		})
	}
}
