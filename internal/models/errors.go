package models

import "fmt"

type InvalidActionError struct {
	Message string
}

func (e InvalidActionError) Error() string {
	return fmt.Sprintf("Invalid action: %s", e.Message)
}

type NoTilesError struct{}

func (e NoTilesError) Error() string {
	return "There are no tiles"
}

type NoTilesOfColorError struct {
	Color TileColor
}

func (e NoTilesOfColorError) Error() string {
	return fmt.Sprintf("There are no %s tiles", e.Color)
}
