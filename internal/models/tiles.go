package models

import (
	"fmt"
	"math/rand"
	"time"
)

type TileCollection struct {
	Tiles []Tile
}

func NewTileCollection() TileCollection {
	rand.Seed(time.Now().UnixNano())

	return TileCollection{
		Tiles: make([]Tile, 0),
	}
}

func (tc *TileCollection) DrawRandomTile() (Tile, error) {
	tileCount := len(tc.Tiles)
	if tileCount == 0 {
		return Tile{}, InvalidActionError{Message: "There are no tiles"}
	}

	// Choose a random tile from the slice
	selectedIndex := rand.Intn(tileCount)
	selectedTile := tc.Tiles[selectedIndex]

	// Remove that tile from the slice
	removeTileFromSlice(tc.Tiles, selectedIndex)

	return selectedTile, nil
}

func (tc *TileCollection) DrawAllTiles() []Tile {
	// Grab all the tiles. This will be returned
	tiles := tc.Tiles

	// Reset the list of Tiles to an empty slice
	tc.Tiles = make([]Tile, 0)

	return tiles
}

func (tc *TileCollection) DrawAllTilesByColor(color TileColor) ([]Tile, error) {
	// validations:
	// is the collection empty?
	// does this collection have any tiles of the specified color?
	result := make([]Tile, 0)

	if len(tc.Tiles) == 0 {
		return result, NoTilesError{}
		//return result, InvalidActionError{Message: "There are no tiles"}
	}

	// Start at the end of the slice and work backwards
	for i := len(tc.Tiles) - 1; i >= 0; i-- {
		if tc.Tiles[i].Color == color {
			// Add the tile to the function's result
			result = append(result, tc.Tiles[i])
			// Remove the tile from the fatory
			tc.Tiles = removeTileFromSlice(tc.Tiles, i)
		}
	}

	if len(result) == 0 {
		return result, NoTilesOfColorError{}
		//return result, InvalidActionError{Message: fmt.Sprintf("There are no %s tiles", color)}
	}

	return result, nil
}

func (tc *TileCollection) AddTile(t Tile) {
	tc.Tiles = append(tc.Tiles, t)
}

func (tc *TileCollection) TileCount() int {
	return len(tc.Tiles)
}

func (tc *TileCollection) String() string {
	return fmt.Sprintf("%v", tc.Tiles)
}

// removeTileFromSlice removes the tile at the specified index, and returns a new slice with
// the tile removed.  This function swaps the tile at the end of the slice with the one that's
// being removed, so the order of tiles is slightly adjusted.
func removeTileFromSlice(tiles []Tile, tileIndex int) []Tile {
	tileCount := len(tiles)

	tiles[tileCount-1], tiles[tileIndex] = tiles[tileIndex], tiles[tileCount-1]
	tiles = tiles[:tileCount-1]

	return tiles
}

type Tile struct {
	Color TileColor
}

type TileColor string

const (
	Orange TileColor = "orange"
	Blue   TileColor = "blue"
	White  TileColor = "white"
	Black  TileColor = "black"
	Red    TileColor = "red"
)

func (t TileColor) String() string {
	return string(t)
}