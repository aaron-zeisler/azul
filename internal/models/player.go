package models

import (
	"fmt"
)

type Player struct {
	Name          string
	Board         *Board
	IsFirstPlayer bool
}

func NewPlayer(name string, opts ...NewPlayerOption) Player {
	p := Player{
		Name:  name,
		Board: &Board{},
	}

	// Create the pattern lines on the board
	p.Board.ResetPatternLines()

	// Reset the floor
	p.Board.ResetFloor()

	for _, opt := range opts {
		p = opt(p)
	}

	return p
}

type NewPlayerOption func(p Player) Player

func FirstPlayer() NewPlayerOption {
	return func(p Player) Player {
		p.IsFirstPlayer = true
		return p
	}
}

func (p Player) String() string {
	s := fmt.Sprintf("Name: %s - Score: %d", p.Name, p.Board.Score)
	if p.IsFirstPlayer {
		s = fmt.Sprintf("%s - First Player", s)
	}

	return s
}

type Board struct {
	Score        int //TODO: Make this the score track, then create a Score() function that calculates ScoreTrack - Floor
	PatternLines map[int][]Tile
	Floor        []Tile //TODO: Add score modifiers to the floor spaces
	//Wall
}

func (b *Board) ResetPatternLines() {
	b.PatternLines = make(map[int][]Tile)
	for i := 0; i < NumPatternLines; i++ {
		b.PatternLines[i] = make([]Tile, 0, i+1)
	}
}

func (b *Board) PlaceTiles(patternLineNumber int, tiles []Tile) error {
	currentLine := b.PatternLines[patternLineNumber]
	maxTiles := cap(currentLine)
	currentTiles := len(currentLine)

	// If the line is already full, return an error
	if currentTiles >= maxTiles {
		return InvalidActionError{Message: "The line is already full of tiles, please choose another line"}
	}

	for _, tile := range tiles {
		if currentTiles >= maxTiles {
			if err := b.AddToFloor([]Tile{tile}); err != nil {
				return err
			}
		} else {
			b.PatternLines[patternLineNumber] = append(b.PatternLines[patternLineNumber], tile)
			currentTiles++
		}
	}

	return nil
}

func (b *Board) ResetFloor() {
	b.Floor = make([]Tile, 0, NumFloorSpaces)
}

func (b *Board) AddToFloor(tiles []Tile) error {
	maxTiles := cap(b.Floor)
	currentTiles := len(b.Floor)

	for _, tile := range tiles {
		if currentTiles >= maxTiles {
			// If there are more tiles than floor spaces, just throw out the extra tiles
			continue
		} else {
			b.Floor = append(b.Floor, tile)
			currentTiles++
		}
	}

	return nil
}

const NumFloorSpaces = 7
const NumPatternLines = 5

/*
type FloorSpace struct {
	ScoreModifier int
	HasTile       bool
	Tile          Tile
}
*/
