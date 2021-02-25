package models

import "errors"

type Bag struct {
	*TileCollection
}

func NewBag() *Bag {
	return &Bag{
		TileCollection: NewTileCollection(WithErrorHandler(bagErrorHandler{})),
	}
}

type bagErrorHandler struct{}

func (eh bagErrorHandler) HandleError(err error) error {
	if err != nil {
		if errors.Is(err, NoTilesError{}) {
			return InvalidActionError{Message: "The bag is empty"}
		}
	}

	return err
}
