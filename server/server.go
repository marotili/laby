package main

import (
	"encoding/gob"
	"laby/game"
	"log"
	"net"
	"sync"
	"time"
)

// var game Game = game.NewGame()
var gameState *GameState
var initOnce sync.Once

func InitGame() {
	initOnce.Do(func() {
		g, err := game.NewGame()
		if err != nil {
			log.Fatal("Failed to initialize game")
		}
		gameState = &GameState{
			dataLock:    sync.Mutex{},
			playerData:  make(map[*Player]*PerPlayerState, 0),
			game:        g,
			gameStarted: false,
		}
	})
}

type GameState struct {
	dataLock    sync.Mutex
	playerData  map[*Player]*PerPlayerState
	game        *game.Game
	gameStarted bool
}

type Player struct {
	conn       net.Conn
	gamePlayer game.Player
}

func NewPlayer(conn net.Conn, gamePlayer game.Player) *Player {
	return &Player{
		conn:       conn,
		gamePlayer: gamePlayer,
	}
}

type PerPlayerState struct {
	newActions []game.ActionType
	isReady    bool
}

func NewPlayerState() *PerPlayerState {
	return &PerPlayerState{
		newActions: make([]game.ActionType, 0),
		isReady:    false,
	}
}

func PlayerIsSynchronized(player *Player) bool {
	gameState.dataLock.Lock()
	defer gameState.dataLock.Unlock()
	if len(gameState.playerData[player].newActions) == 0 {
		return true
	}

	return false
}

func SetPlayerSynchronized(player *Player) {
	gameState.dataLock.Lock()
	defer gameState.dataLock.Unlock()

	gameState.playerData[player].newActions = make([]game.ActionType, 0)
}

func OtherPlayers(player *Player) []*Player {
	gameState.dataLock.Lock()
	defer gameState.dataLock.Unlock()
	otherPlayers := make([]*Player, 0, len(gameState.playerData)-1)

	for otherPlayer, _ := range gameState.playerData {
		otherPlayers = append(otherPlayers, otherPlayer)
	}

	return otherPlayers
}

func CompileData(player *Player) map[game.Player][]game.ActionType {
	gameState.dataLock.Lock()
	defer gameState.dataLock.Unlock()

	data := make(map[game.Player][]game.ActionType, 0)
	for otherPlayer, state := range gameState.playerData {
		copyData := make([]game.ActionType, len(state.newActions))
		copy(copyData, state.newActions)
		data[otherPlayer.gamePlayer] = copyData
	}
	return data
}

func AddNewActions(player *Player, actions []game.ActionType) {
	gameState.dataLock.Lock()
	defer gameState.dataLock.Unlock()
	gameState.playerData[player].newActions = actions
}

func SetPlayerReady(player *Player) {
	gameState.dataLock.Lock()
	defer gameState.dataLock.Unlock()
	gameState.playerData[player].isReady = true
}

func AllPlayersReady() bool {
	gameState.dataLock.Lock()
	defer gameState.dataLock.Unlock()

	and := true
	for _, state := range gameState.playerData {
		and = and && state.isReady // if one player not ready ret false
	}

	return and
}

func StartGame() {
	gameState.dataLock.Lock()
	defer gameState.dataLock.Unlock()

	gameState.gameStarted = true
}

func GameAddPlayer(conn net.Conn) *Player {
	gameState.dataLock.Lock()
	defer gameState.dataLock.Unlock()

	playerId := len(gameState.playerData)
	gamePlayer := gameState.game.NewPlayer(playerId)
	newPlayer := NewPlayer(conn, gamePlayer)

	gameState.playerData[newPlayer] = NewPlayerState()
	return newPlayer
}

func handleConnection(conn net.Conn) {
	player := GameAddPlayer(conn)
	defer func() {
	}()

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	var req game.ClientRequest
	for {
		log.Println("Wait for request")
		err := dec.Decode(&req)
		if err != nil {
			log.Fatal("Failed to decode client req")
		}

		log.Println("Got request", req)
		switch req {
		case game.ClientReqSendAction:
			var actions []game.ActionType = make([]game.ActionType, 0, 100)
			var numActions int

			dec.Decode(&numActions)
			if err != nil {
				log.Fatal("Failed to decode action length")
			}

			for i := 0; i < numActions; i++ {
				var action game.ActionType
				err = dec.Decode(&action)
				if err != nil {
					log.Fatal("Failed to decode actions")
				}

				actions = append(actions, action)
			}

			actionDenied := false
			if PlayerIsSynchronized(player) {
				log.Println(actions)
				for _, action := range actions {
					if action == game.ActionPlayerReady {
						if gameState.gameStarted {
							actionDenied = true
						} else {
							SetPlayerReady(player)
						}
					}
				}

				if actionDenied {
					log.Println("Action denied from player", player)
					enc.Encode(game.ServerActionDenied)
				} else {
					// Test game state
					copyActions := make([]game.ActionType, len(actions))
					copy(copyActions, actions)
					AddNewActions(player, copyActions)
					log.Println("Action ok from player", player)
					enc.Encode(game.ServerActionOk)
				}
			} else {
				log.Println("Player not synchronized")
				enc.Encode(game.ServerActionWait)
				// ignore
			}

		case game.ClientReqUpdate:
			var data map[game.Player][]game.ActionType = CompileData(player)
			enc.Encode(len(data))
			for player, actions := range data {
				enc.Encode(player)
				enc.Encode(len(actions))
				for _, action := range actions {
					enc.Encode(action)
				}
			}

			SetPlayerSynchronized(player)
		default:
			log.Println("Unknown command from client")
		}

		// time.Sleep(time.Millisecond * 50)
	}
	time.Sleep(100)
}

// var (
// 	game Game
// 	mu   sync.Mutex
// )

func main() {
	var err error

	InitGame()

	// if game, err = NewGame(); err != nil {
	// 	log.Fatal(err)
	// }

	// go func() {
	listen, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
	// }()
}
