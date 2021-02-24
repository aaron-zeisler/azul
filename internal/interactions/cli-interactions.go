package interactions

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aaron-zeisler/azul/internal/models"
)

type DrawFactoryTilesResponse struct {
	FactoryNumber int
	TileColor     models.TileColor
}

func PromptToDrawFactoryTiles() (DrawFactoryTilesResponse, error) {
	response := DrawFactoryTilesResponse{}

	//TODO: Allow the player to draw from the center of the table

	factoryNumber, err := PromptForInt("Which factory would you like to draw from?")
	if err != nil {
		return response, err
	}

	tileColor, err := PromptForString("Which color would you like to draw?")
	if err != nil {
		return response, err
	}

	return DrawFactoryTilesResponse{
		FactoryNumber: factoryNumber,
		TileColor:     models.TileColor(tileColor),
	}, nil
}

type NewPlayersResponse struct {
	Players map[int]models.Player
}

//TODO: Separate prompting from validation
func PromptForNewPlayers(config models.GameConfig) (NewPlayersResponse, error) {
	response := NewPlayersResponse{
		Players: make(map[int]models.Player),
	}

	var numPlayers int
	var err error
	for numPlayers == 0 {
		numPlayers, err = PromptForInt(fmt.Sprintf("How many players are playing the game? (min: %d; max: %d)", config.MinNumberOfPlayers, config.MaxNumberOfPlayers))
		if err != nil {
			return response, err
		}

		if numPlayers < config.MinNumberOfPlayers || numPlayers > config.MaxNumberOfPlayers {
			fmt.Println("Invalid number of players")
		}
	}

	for i := 0; i < numPlayers; i++ {
		playerName, err := PromptForString(fmt.Sprintf("What is player #%d's name?", i))
		if err != nil {
			return response, err
		}

		player := models.Player{
			Name:          playerName,
			IsFirstPlayer: i == 0,
		}
		response.Players[i] = player
	}

	return response, nil
}

func PromptForString(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(prompt)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to parse the response for '%s': %w", prompt, err)
	}

	answer = strings.Trim(answer, "\n")

	return answer, nil
}

func PromptForInt(prompt string) (int, error) {
	answerStr, err := PromptForString(prompt)
	if err != nil {
		return 0, err
	}

	answer, err := strconv.Atoi(answerStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert the answer to an integer: %w", err)
	}

	return answer, nil
}
