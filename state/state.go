package state

var Info info = info{
	Games: map[string]Game{},
}

type info struct {
	// The Key is always a Channel ID
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
	HeaderID string
	ThreadID string
	Players  map[string]Player
	Status   GameStatus
}

type Player struct {
	Status PlayerStatus
}

func (i *info) AddNewGame(pcID string, tcID string, hID string, pID string) {
	i.Games[pcID] = Game{
		HeaderID: hID,
		ThreadID: tcID,
		Players:  map[string]Player{pID: {Status: PlayerStatus(WAITING)}},
		Status:   GameStatus(DEALING),
	}
}

func (i *info) AddPlayer(pcID string, pID string) {
	i.Games[pcID].Players[pID] = Player{Status: PlayerStatus(WAITING)}
}

func (i *info) GameHasPlayer(pcID string, pID string) bool {
	_, has := i.Games[pcID].Players[pID]
	return has
}

func (i *info) ParentChannelHasGame(id string) bool {
	_, has := i.Games[id]
	return has
}

func (i *info) GetPlayers(id string) map[string]Player {
	return i.Games[id].Players
}

func (i *info) GetHeader(id string) string {
	return i.Games[id].HeaderID
}

func (i *info) GetThread(id string) string {
	return i.Games[id].ThreadID
}
