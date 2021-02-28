package main

import (
	"fmt"

	"github.com/aaron-zeisler/azul/internal/interactions"
	"github.com/aaron-zeisler/azul/internal/models"
)

func main() {
	fmt.Println("AZUL STARTING ...")
	fmt.Println()

	config := models.DefaultGameConfig

	// Prompt for the number of players
	playerSetup, err := interactions.PromptForNewPlayers(config)
	if err != nil {
		panic(err)
	}

	// Initialize the game
	game := models.NewGame(
		models.WithConfig(config),
		models.WithPlayers(playerSetup.Players))

	interactions.DisplayGameState(game)
	fmt.Printf("Number of tiles left in the bag: %d\n", game.Bag.TileCount())

	// Draw tiles from a factory or the center of the table
	drawResponse, err := interactions.PromptToDrawFactoryTiles()
	if err != nil {
		panic(err)
	}

	var drawSource models.DrawSource
	var drawDisplay string
	if drawResponse.DrawSourceType == models.DrawSourceCenter {
		drawSource = game.CenterOfTheTable
		drawDisplay = "the center of the table"
	} else {
		drawSource = game.Factories[drawResponse.FactoryNumber]
		drawDisplay = fmt.Sprintf("factory #%d", drawResponse.FactoryNumber)
	}
	fmt.Printf("Drawing %s tiles from %s\n", drawResponse.TileColor, drawDisplay)
	drawnTiles, err := drawSource.DrawAllTilesByColor(drawResponse.TileColor)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Put the rest of the factory's tiles in the center of the table
	if drawResponse.DrawSourceType == models.DrawSourceFactory {
		leftovers := game.Factories[drawResponse.FactoryNumber].DrawAllTiles()
		for _, tile := range leftovers {
			game.CenterOfTheTable.AddTile(tile)
		}
	}

	//Put the drawn tiles onto the player's game board
	placeResponse, err := interactions.PromptToPlaceFactoryTiles()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Placing tiles on pattern line #%d\n", placeResponse.PatternLineNumber)
	err = game.Players[0].Board.PlaceTiles(placeResponse.PatternLineNumber, drawnTiles)
	if err != nil {
		panic(err)
	}

	fmt.Println("After placing tiles:")
	interactions.DisplayGameState(game)
}
