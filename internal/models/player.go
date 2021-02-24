package models

import (
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
