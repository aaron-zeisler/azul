package models

import (
	"errors"
	"fmt"
)

type Factory struct {
	*TileCollection
}

func NewFactory() *Factory {
	return &Factory{
		TileCollection: NewTileCollection(),
	}
}

func (f *Factory) DrawAllTilesByColor(color TileColor) ([]Tile, error) {
	result, err := f.TileCollection.DrawAllTilesByColor(color)
	if err != nil {
		if errors.Is(err, NoTilesError{}) {
			return result, InvalidActionError{Message: "This factory has no tiles"}
		} else if errors.Is(err, NoTilesOfColorError{}) {
			return result, InvalidActionError{Message: fmt.Sprintf("There are no %s tiles on this factory", color)}
		} else {
			return result, fmt.Errorf("failed to draw the %s tiles: %w", color, err)
		}
	}

	return result, nil
}
