package models

import "fmt"

type GameConfig struct {
	TileColors         []TileColor
	TilesPerColor      int
	TilesPerFactory    int
	MinNumberOfPlayers int
	MaxNumberOfPlayers int

	// PlayersToFactoriesMap specifies how many factories to include in the game depending on
	// how many players are playing. In this map, the key is number of players, and the value
	// is number of factories.
	PlayersToFactoriesMap map[int]int
}

var DefaultGameConfig = GameConfig{
	TileColors:         []TileColor{Orange, Blue, White, Black, Red},
	TilesPerColor:      20,
	TilesPerFactory:    4,
	MinNumberOfPlayers: 2,
	MaxNumberOfPlayers: 4,
	PlayersToFactoriesMap: map[int]int{
		2: 5,
		3: 7,
		4: 9,
	},
}

type Game struct {
	Config           GameConfig
	Players          map[int]Player
	Factories        map[int]*Factory
	CenterOfTheTable TileCollection
	Bag              *Bag
}

func NewGame(opts ...NewGameOption) *Game {
	g := &Game{
		Config:           DefaultGameConfig,
		Players:          make(map[int]Player),
		CenterOfTheTable: NewTileCollection(),
	}

	for _, opt := range opts {
		opt(g)
	}

	g.InitBag()
	g.InitFactories()

	return g
}

type NewGameOption func(g *Game)

func WithConfig(config GameConfig) NewGameOption {
	return func(g *Game) {
		g.Config = config
	}
}

func WithPlayers(players map[int]Player) NewGameOption {
	return func(g *Game) {
		for i := 0; i < len(players); i++ {
			g.Players[i] = Player{
				Name:          players[i].Name,
				IsFirstPlayer: players[i].IsFirstPlayer,
			}
		}
	}
}

func (g *Game) InitFactories() {
	numFactories := g.Config.PlayersToFactoriesMap[len(g.Players)]
	g.Factories = make(map[int]*Factory, numFactories)

	for i := 0; i < numFactories; i++ {
		factory := NewFactory()
		for t := 0; t < g.Config.TilesPerFactory; t++ {
			tile, _ := g.Bag.DrawTile()
			factory.AddTile(tile)
		}

		g.Factories[i] = factory
	}
}

func (g *Game) InitBag() {
	g.Bag = NewBag()

	var tileCounter int
	for _, color := range g.Config.TileColors {
		for i := 0; i < g.Config.TilesPerColor; i++ {
			g.Bag.AddTile(Tile{Color: color})
			tileCounter += 1
		}
	}
}

func (g *Game) DisplayState() {
	// Print out the players and their boards
	fmt.Println("PLAYERS:")
	for i := 0; i < len(g.Players); i++ {
		fmt.Printf("Player #%d: %s\n", i, g.Players[i])
	}
	fmt.Println()

	// Print out the factories and their tiles
	fmt.Println("FACTORIES:")
	for i := 0; i < len(g.Factories); i++ {
		fmt.Printf("Factory #%d: %s\n", i, g.Factories[i])
	}
	// Print the tiles in center of the table
	fmt.Printf("Center of the Table: %s\n", g.CenterOfTheTable.Tiles)
	fmt.Println()
}