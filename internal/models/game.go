package models

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
	CenterOfTheTable *TileCollection
	Bag              *Bag
	//DiscardPile      []Tile
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

	g.ResetBag()
	g.ResetFactories()

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
			g.Players[i] = players[i]
		}
	}
}

func (g *Game) ResetFactories() {
	numFactories := g.Config.PlayersToFactoriesMap[len(g.Players)]
	g.Factories = make(map[int]*Factory, numFactories)

	for i := 0; i < numFactories; i++ {
		factory := NewFactory()
		for t := 0; t < g.Config.TilesPerFactory; t++ {
			tile, err := g.Bag.DrawRandomTile()
			if err != nil {
				panic(err)
			}
			factory.AddTile(tile)
		}

		g.Factories[i] = factory
	}
}

func (g *Game) ResetBag() {
	g.Bag = NewBag()

	var tileCounter int
	for _, color := range g.Config.TileColors {
		for i := 0; i < g.Config.TilesPerColor; i++ {
			g.Bag.AddTile(Tile{Color: color})
			tileCounter++
		}
	}
}
