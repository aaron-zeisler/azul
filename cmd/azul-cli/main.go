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
	response, err := interactions.PromptToDrawFactoryTiles()
	if err != nil {
		panic(err)
	}

	var drawSource models.DrawSource
	var drawDisplay string
	if response.DrawSourceType == models.DrawSourceCenter {
		drawSource = game.CenterOfTheTable
		drawDisplay = "the center of the table"
	} else {
		drawSource = game.Factories[response.FactoryNumber]
		drawDisplay = fmt.Sprintf("factory #%d", response.FactoryNumber)
	}
	fmt.Printf("Drawing %s tiles from %s\n", response.TileColor, drawDisplay)
	drawSource.DrawAllTilesByColor(response.TileColor)

	// Put the rest of the factory's tiles in the center of the table
	if response.DrawSourceType == models.DrawSourceFactory {
		leftovers := game.Factories[response.FactoryNumber].DrawAllTiles()
		for _, tile := range leftovers {
			game.CenterOfTheTable.AddTile(tile)
		}
	}

	fmt.Println("After drawing tiles:")
	interactions.DisplayGameState(game)

	//TODO: Put the drawn tiles onto the player's game board (!!!)
}
