package models

import "fmt"

type InvalidActionError struct {
	Message string
}

func (e InvalidActionError) Error() string {
	return fmt.Sprintf("Invalid action: %s", e.Message)
}

type NoTilesError struct{ error }
type NoTilesOfColorError struct{ error }
