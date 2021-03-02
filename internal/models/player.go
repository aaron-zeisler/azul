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

	for _, opt := range opts {
		p = opt(p)
	}

	// Create the pattern lines on the board
	p.Board.ResetPatternLines()

	// Reset the floor
	p.Board.ResetFloor()

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
	Score        int
	PatternLines map[int][]Tile
	Floor        []FloorSpace
	Wall         [][]WallSpace
}

func (b *Board) ResetPatternLines() {
	b.PatternLines = make(map[int][]Tile)
	for i := 0; i < NumPatternLines; i++ {
		b.ResetPatternLine(i)
	}
}

func (b *Board) ResetPatternLine(rowNumber int) {
	b.PatternLines[rowNumber] = make([]Tile, 0, rowNumber+1)
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
		// If the pattern line is full, place the tile on the Floor instead
		// Also, if this is the 1st player tile, place it on the floor
		if currentTiles >= maxTiles || tile.Color == FirstPlayerTile {
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
	b.Floor = make([]FloorSpace, 0, NumFloorSpaces)
}

func (b *Board) AddToFloor(tiles []Tile) error {
	maxTiles := cap(b.Floor)
	currentTiles := len(b.Floor)

	for _, tile := range tiles {
		if currentTiles >= maxTiles {
			// If there are more tiles than floor spaces, just throw out the extra tiles
			continue
		} else {
			b.Floor = append(b.Floor, FloorSpace{Tile: tile, ScoreModifier: FloorScoreModifiers[currentTiles]})
			currentTiles++
		}
	}

	return nil
}

const NumFloorSpaces = 7
const NumPatternLines = 5

type FloorSpace struct {
	Tile
	ScoreModifier int
}

var FloorScoreModifiers = map[int]int{
	0: -1,
	1: -1,
	2: -2,
	3: -2,
	4: -2,
	5: -3,
	6: -3,
}

func (f FloorSpace) String() string {
	return fmt.Sprintf("%v(%d)", f.Tile, f.ScoreModifier)
}

type WallSpace struct {
	Tile
	HasTile bool
}

func (w WallSpace) String() string {
	if w.HasTile {
		return fmt.Sprintf("%s", w.Tile)
	} else {
		return "{empty}"
	}
}

var wallLayout = [][]Tile{
	{{Color: Blue}, {Color: Orange}, {Color: Red}, {Color: Black}, {Color: White}},
	{{Color: White}, {Color: Blue}, {Color: Orange}, {Color: Red}, {Color: Black}},
	{{Color: Black}, {Color: White}, {Color: Blue}, {Color: Orange}, {Color: Red}},
	{{Color: Red}, {Color: Black}, {Color: White}, {Color: Blue}, {Color: Orange}},
	{{Color: Orange}, {Color: Red}, {Color: Black}, {Color: White}, {Color: Blue}},
}

func (b *Board) ResetWall() {
	b.Wall = make([][]WallSpace, 5)

	for i := 0; i < len(wallLayout); i++ {
		b.Wall[i] = make([]WallSpace, 5)
		for j := 0; j < len(wallLayout[i]); j++ {
			b.Wall[i][j].Tile = wallLayout[i][j]
		}
	}
}

func (b *Board) MoveTileToWall(tile Tile, rowNumber int) WallCoordinate {
	// Return the coordinate where the tile was placed on the wall
	coord := WallCoordinate{Row: rowNumber}

	// Find the coordinate
	for i := 0; i < len(b.Wall[rowNumber]); i++ {
		if b.Wall[rowNumber][i].Color == tile.Color {
			coord.Col = i
			break
		}
	}

	// Set the wall tile's HasTile property to true
	b.Wall[coord.Row][coord.Col].HasTile = true
	return coord
}

// Move tiles from the pattern lines to the wall, score the tiles, and reset the pattern lines
func (b *Board) ScorePatternLines() {
	// Iterate through each pattern line
	for i := 0; i < NumPatternLines; i++ {
		// If the pattern line if full...
		if len(b.PatternLines[i]) == i+1 {
			// Move a tile of that color to the wall
			wallCoord := b.MoveTileToWall(b.PatternLines[i][i], i)
			// Score the new tile in the wall
			wallScore := b.ScoreTile(wallCoord)
			b.Score += wallScore.Score
			// Discard all the other tiles and reset the pattern line
			b.ResetPatternLine(i)
		}
	}
}

// Score and reset the floor
func (b *Board) ScoreFloor() {
	for _, tile := range b.Floor {
		b.Score -= tile.ScoreModifier
	}
	b.ResetFloor()
}

type WallScore struct {
	Score int
	Tiles []WallCoordinate
}

type WallCoordinate struct {
	Row int
	Col int
}

// Calculate the score for a single tile placed into the wall
// This function returns the total score, along with a list of all the
// sets of linked tiles that contributed to the score.
func (b *Board) ScoreTile(coord WallCoordinate) WallScore {
	result := WallScore{
		Tiles: make([]WallCoordinate, 0),
	}

	rows := len(wallLayout)
	cols := len(wallLayout[0])

	// Add the scoring tile to the result
	result.Tiles = append(result.Tiles, coord)

	// Calculate the score for horizontal adjacent tiles
	horizontalTiles := make([]WallCoordinate, 0)

	// Count the adjacent tiles to the right
	if coord.Col < cols-1 {
		for i := coord.Col + 1; i < cols; i++ {
			fmt.Printf("Looking at {%d,%d}: HasTile: %t\n", coord.Row, i, b.Wall[coord.Row][i].HasTile)
			if b.Wall[coord.Row][i].HasTile {
				horizontalTiles = append(horizontalTiles, WallCoordinate{Row: coord.Row, Col: i})
			} else {
				break
			}
		}
	}

	// Count the adjacent tiles to the left
	if coord.Col > 0 {
		for i := coord.Col - 1; i >= 0; i-- {
			fmt.Printf("Looking at {%d,%d}: HasTile: %t\n", coord.Row, i, b.Wall[coord.Row][i].HasTile)
			if b.Wall[coord.Row][i].HasTile {
				horizontalTiles = append(horizontalTiles, WallCoordinate{Row: coord.Row, Col: i})
			} else {
				break
			}
		}
	}

	// If there were adjacent tiles, add the tiles and score to the result
	if len(horizontalTiles) > 0 {
		result.Tiles = append(result.Tiles, horizontalTiles...)
		result.Score += len(horizontalTiles) + 1
	}

	// Calculate the score for vertical adjacent tiles
	verticalTiles := make([]WallCoordinate, 0)

	// Count the adjacent tiles below
	if coord.Row < rows-1 {
		for i := coord.Row + 1; i < rows; i++ {
			fmt.Printf("Looking at {%d,%d}: HasTile: %t\n", i, coord.Col, b.Wall[i][coord.Col].HasTile)
			if b.Wall[i][coord.Col].HasTile {
				verticalTiles = append(verticalTiles, WallCoordinate{Row: i, Col: coord.Col})
			} else {
				break
			}
		}
	}

	// Count the adjacent tiles above
	if coord.Row > 0 {
		for i := coord.Row - 1; i >= 0; i-- {
			fmt.Printf("Looking at {%d,%d}: HasTile: %t\n", i, coord.Col, b.Wall[i][coord.Col].HasTile)
			if b.Wall[i][coord.Col].HasTile {
				verticalTiles = append(verticalTiles, WallCoordinate{Row: i, Col: coord.Col})
			} else {
				break
			}
		}
	}

	// If there were adjacent tiles, add the specified tile
	if len(verticalTiles) > 0 {
		result.Tiles = append(result.Tiles, verticalTiles...)
		result.Score += len(verticalTiles) + 1
	}

	// If the tile doesn't have any adjacent tiles, it's worth 1 point
	if len(result.Tiles) == 1 {
		result.Score = 1
	}

	return result
}
