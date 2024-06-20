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
	Players map[string]Player
	Status  GameStatus
}

type Player struct {
	Status   PlayerStatus
	HeaderID string
	ThreadID string
}

func (i *info) AddNewGame(pcID string, tcID string, hID string, pID string) {
	i.Games[pcID] = Game{
		Players: map[string]Player{pID: {Status: PlayerStatus(WAITING),
			HeaderID: hID,
			ThreadID: tcID,
		},
		},
		Status: GameStatus(DEALING),
	}
}

func (i *info) AddPlayer(pcID string, pID string, tID string, hID string) {
	i.Games[pcID].Players[pID] = Player{Status: PlayerStatus(WAITING), HeaderID: hID, ThreadID: tID}
}

func (i *info) GameHasPlayer(pcID string, pID string) bool {
	_, has := i.Games[pcID].Players[pID]
	return has
}

func (i *info) ParentChannelHasGame(ID string) bool {
	_, has := i.Games[ID]
	return has
}

func (i *info) GetPlayers(ID string) map[string]Player {
	return i.Games[ID].Players
}

func (i *info) GetHeader(gID string, pID string) string {
	return i.Games[gID].Players[pID].HeaderID
}

func (i *info) GetThread(gID string, pID string) string {
	return i.Games[gID].Players[pID].ThreadID
}
