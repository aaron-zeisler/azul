package interactions

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aaron-zeisler/azul/internal/models"
)

type PlaceFactoryTilesResponse struct {
	PatternLineNumber int
}

func PromptToPlaceFactoryTiles() (PlaceFactoryTilesResponse, error) {
	response := PlaceFactoryTilesResponse{}

	patternLineNumber, err := PromptForInt("Which line would you like to place the tiles on?")
	if err != nil {
		return response, err
	}
	response.PatternLineNumber = patternLineNumber

	return response, nil
}

type DrawFactoryTilesResponse struct {
	DrawSourceType models.DrawSourceType
	FactoryNumber  int
	TileColor      models.TileColor
}

func PromptToDrawFactoryTiles() (DrawFactoryTilesResponse, error) {
	response := DrawFactoryTilesResponse{}

	drawSourceTypeStr, err := PromptForString("Would you like to draw from a factory or from the center of the table (type 'factory' or 'center')?")
	if err != nil {
		return response, err
	}
	//TODO: Validate this
	response.DrawSourceType = models.DrawSourceType(drawSourceTypeStr)

	if response.DrawSourceType == models.DrawSourceFactory {
		factoryNumber, err := PromptForInt("Which factory would you like to draw from?")
		if err != nil {
			return response, err
		}

		response.FactoryNumber = factoryNumber
	}

	tileColor, err := PromptForString("Which color would you like to draw?")
	if err != nil {
		return response, err
	}
	response.TileColor = models.TileColor(tileColor)

	return response, nil
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

		opts := make([]models.NewPlayerOption, 0)
		if i == 0 {
			opts = append(opts, models.FirstPlayer())
		}
		player := models.NewPlayer(playerName, opts...)

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

func DisplayGameState(game *models.Game) {
	// Print out the players and their boards
	fmt.Println("PLAYERS:")
	for i := 0; i < len(game.Players); i++ {
		fmt.Printf("Player #%d: %s\n", i, game.Players[i])

		fmt.Println("Game Board:")

		// Print the pattern lines
		printPatternLines(game.Players[i].Board)

		// Print the floor
		printFloor(game.Players[i].Board.Floor)
	}
	fmt.Println()

	// Print out the factories and their tiles
	fmt.Println("FACTORIES:")
	for i := 0; i < len(game.Factories); i++ {
		fmt.Printf("Factory #%d: %s\n", i, game.Factories[i])
	}
	// Print the tiles in center of the table
	fmt.Printf("Center of the Table: %s\n", game.CenterOfTheTable.Tiles)
	fmt.Println()
}

func printPatternLines(board *models.Board) {
	fmt.Println("Pattern Lines:")
	//fmt.Printf("%v\n", board.PatternLines)

	for line := 0; line < models.NumPatternLines; line++ {
		lineString := fmt.Sprintf("%d", line)
		for tile := 0; tile <= line; tile++ {
			if tile >= len(board.PatternLines[line]) {
				lineString = fmt.Sprintf("%s {empty}", lineString)
			} else {
				lineString = fmt.Sprintf("%s %s", lineString, board.PatternLines[line][tile])
			}
		}
		fmt.Println(lineString)
	}
}

func printFloor(floor []models.FloorSpace) {
	lineString := "Floor:"

	for tile := 0; tile < models.NumFloorSpaces; tile++ {
		if tile >= len(floor) {
			lineString = fmt.Sprintf("%s {empty}(%d)", lineString, models.ScoreModifiers[tile])
		} else {
			lineString = fmt.Sprintf("%s %s", lineString, floor[tile])
		}
	}

	fmt.Println(lineString)
}
