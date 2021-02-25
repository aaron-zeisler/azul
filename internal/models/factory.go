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
		TileCollection: NewTileCollection(WithErrorHandler(factoryErrorHandler{})),
	}
}

type factoryErrorHandler struct{}

func (eh factoryErrorHandler) HandleError(err error) error {
	if err != nil {
		if errors.As(err, &NoTilesError{}) {
			return InvalidActionError{Message: "This factory has no tiles"}
		} else if errors.As(err, &NoTilesOfColorError{}) {
			return InvalidActionError{Message: fmt.Sprintf("There are no %s tiles on this factory", err.(NoTilesOfColorError).Color)}
		} else {
			return fmt.Errorf("failed to draw the tiles: %w", err)
		}
	}
	return nil
}
