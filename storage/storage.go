package storage

import "github.com/0xfyer/discord-bot/game"

type State struct {
	// The Key is always a Guild ID
	Games map[string]Game
}

type GameStatus int
type PlayerStatus int

const (
	WAITING  GameStatus   = 0
	DEALING  GameStatus   = 1
	DECIDING PlayerStatus = 0
	HIT      PlayerStatus = 1
	STAND    PlayerStatus = 2
)

type Game struct {
	Players map[string]Player
	Status  GameStatus
	Game    game.Game
}

type Player struct {
	Status   PlayerStatus
	HeaderID string
	ThreadID string
}

func New() *State {
	return &State{
		Games: map[string]Game{},
	}

}

func (i *State) AddNewGame(gID string, game game.Game) {
	i.Games[gID] = Game{
		Game:    game,
		Players: map[string]Player{},
		Status:  GameStatus(DEALING),
	}
}

func (i *State) AddPlayer(gID string, pID string, tID string, hID string) {
	i.Games[gID].Players[pID] = Player{Status: PlayerStatus(WAITING), HeaderID: hID, ThreadID: tID}
}

func (i *State) GameHasPlayer(gID string, pID string) bool {
	_, has := i.Games[gID].Players[pID]
	return has
}

func (i *State) GuildHasGame(ID string, game game.Game) bool {
	_, has := i.Games[ID]
	return has
}

func (i *State) GetPlayers(ID string) map[string]Player {
	return i.Games[ID].Players
}

func (i *State) GetHeader(gID string, pID string) string {
	return i.Games[gID].Players[pID].HeaderID
}

func (i *State) GetThread(gID string, pID string) string {
	return i.Games[gID].Players[pID].ThreadID
}
