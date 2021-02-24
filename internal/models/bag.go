package models

import "errors"

type Bag struct {
	*TileCollection
}

func NewBag() *Bag {
	return &Bag{
		TileCollection: NewTileCollection(),
	}
}

// DrawTile draws a random tile from the bag.  The tile is removed from the bag.
func (b *Bag) DrawRandomTile() (Tile, error) {
	result, err := b.TileCollection.DrawRandomTile()
	if err != nil {
		if errors.Is(err, NoTilesError{}) {
			return result, InvalidActionError{Message: "The bag is empty"}
		} else {
			return result, err
		}
	}

	return result, err
}
