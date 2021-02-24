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

	game.DisplayState()
	fmt.Printf("Number of tiles left in the bag: %d\n", game.Bag.TileCount())

	// Draw tiles from a factory
	response, err := interactions.PromptToDrawFactoryTiles()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Drawing %s from factory #%d\n", response.TileColor, response.FactoryNumber)
	game.Factories[response.FactoryNumber].DrawTiles(response.TileColor)

	// Put the rest of the factory's tiles in the center of the table
	leftovers := game.Factories[response.FactoryNumber].DrawAllTiles()
	for _, tile := range leftovers {
		game.CenterOfTheTable.AddTile(tile)
	}

	fmt.Println("After drawing tiles:")
	game.DisplayState()

	//TODO: Put the drawn tiles onto the player's game board (!!!)
}
