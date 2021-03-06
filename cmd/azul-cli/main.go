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

	// Find the first player and set them as the current player
	var firstPlayerKey int
	for i, player := range game.Players {
		if player.IsFirstPlayer {
			firstPlayerKey = i
			break
		}
	}
	currentPlayerKey := firstPlayerKey
	currentPlayer := game.Players[currentPlayerKey]

	// Set the current round, and the turn number within the round
	//currentRound := 1
	currentTurnNumber := 1

	// This is the beginning of the game loop
	for {
		// Display which player's turn it is
		fmt.Println()
		fmt.Printf("CURRENT PLAYER: %s\n", currentPlayer.Name)

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

		//Put the drawn tiles onto the player's game board (and/or floor)
		placeResponse, err := interactions.PromptToPlaceFactoryTiles()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Placing tiles on pattern line #%d\n", placeResponse.PatternLineNumber)
		err = currentPlayer.Board.PlaceTiles(placeResponse.PatternLineNumber, drawnTiles)
		if err != nil {
			panic(err)
		}

		fmt.Println("After placing tiles:")
		interactions.DisplayGameState(game)

		// Move the pointer to the next player
		currentPlayerKey++
		if currentPlayerKey >= len(game.Players) {
			currentPlayerKey = 0
			currentTurnNumber++
		}
		currentPlayer = game.Players[currentPlayerKey]

		// Check for the end of the round
		if isRoundOver(game) {
			// Score the round and show the game state
			game.ScoreRound()

			fmt.Println("TESTING -- GAME OVER")
			fmt.Println("Here's the game state after scoring:")
			interactions.DisplayGameState(game)
			return
		}
	}
}

func isRoundOver(game *models.Game) bool {
	return factoriesAreEmpty(game) && centerOfTableIsEmpty(game)
}

func factoriesAreEmpty(game *models.Game) bool {
	for _, factory := range game.Factories {
		if factory.HasTiles() {
			return false
		}
	}

	return true
}

func centerOfTableIsEmpty(game *models.Game) bool {
	return !game.CenterOfTheTable.HasTiles()
}
