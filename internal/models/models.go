package models

import (
	"errors"
	"fmt"
)

type Player struct {
	Name          string
	Board         Board
	IsFirstPlayer bool
	Score         int
}

func (p Player) String() string {
	s := fmt.Sprintf("Name: %s - Score: %d", p.Name, p.Score)
	if p.IsFirstPlayer {
		s = fmt.Sprintf("%s - First Player", s)
	}

	return s
}

type Board struct {
	Score int
	//Wall
	//PatterLines
	//Floor
}

type FloorSpace struct {
	ScoreModifier int
	HasTile       bool
	Tile          Tile
}

type Factory struct {
	TileCollection
}

func NewFactory() *Factory {
	return &Factory{
		TileCollection: NewTileCollection(),
	}
}

func (f *Factory) DrawTiles(color TileColor) ([]Tile, error) {
	result, err := f.DrawAllTilesByColor(color)
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

type Bag struct {
	TileCollection
}

func NewBag() *Bag {
	return &Bag{
		TileCollection: NewTileCollection(),
	}
}

// DrawTile draws a random tile from the bag.  The tile is removed from the bag.
func (b *Bag) DrawTile() (Tile, error) {
	result, err := b.DrawRandomTile()
	if err != nil {
		if errors.Is(err, NoTilesError{}) {
			return result, InvalidActionError{Message: "The bag is empty"}
		} else {
			return result, err
		}
	}

	return result, err
}

type InvalidActionError struct {
	Message string
}

func (e InvalidActionError) Error() string {
	return fmt.Sprintf("Invalid action: %s", e.Message)
}

type NoTilesError struct{ error }
type NoTilesOfColorError struct{ error }
