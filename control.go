package netorcai

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

// Game state
const (
	GAME_NOT_RUNNING = iota
	GAME_RUNNING     = iota
	GAME_FINISHED    = iota
)

// Client state
const (
	CLIENT_UNLOGGED = iota
	CLIENT_LOGGED   = iota
	CLIENT_READY    = iota
	CLIENT_THINKING = iota
	CLIENT_KICKED   = iota
)

type GlobalState struct {
	Mutex sync.Mutex

	Listener net.Listener

	GameState int

	GameLogic []*GameLogicClient
	Players   []*PlayerOrVisuClient
	Visus     []*PlayerOrVisuClient

	NbPlayersMax                int
	NbVisusMax                  int
	NbTurnsMax                  int
	Autostart                   bool
	Fast                        bool
	MillisecondsBeforeFirstTurn float64
	MillisecondsBetweenTurns    float64
}

// Debugging helpers
const (
	debugGlobalStateMutex = false
)

func LockGlobalStateMutex(gs *GlobalState, reason, who string) {
	if debugGlobalStateMutex {
		log.WithFields(log.Fields{
			"reason": reason,
			"who":    who,
		}).Debug("Desire global state mutex")
	}
	gs.Mutex.Lock()
	if debugGlobalStateMutex {
		log.WithFields(log.Fields{
			"reason": reason,
			"who":    who,
		}).Debug("Got global state mutex")
	}
}

func UnlockGlobalStateMutex(gs *GlobalState, reason, who string) {
	if debugGlobalStateMutex {
		log.WithFields(log.Fields{
			"reason": reason,
			"who":    who,
		}).Debug("Release global state mutex")
	}
	gs.Mutex.Unlock()
}

func areAllExpectedClientsConnected(gs *GlobalState) bool {
	return (len(gs.Players) == gs.NbPlayersMax) &&
		(len(gs.Visus) == gs.NbVisusMax) &&
		(len(gs.GameLogic) == 1)
}

func autostart(gs *GlobalState) {
	if gs.Autostart && areAllExpectedClientsConnected(gs) {
		log.Info("Automatic starting conditions are met")
		gs.GameState = GAME_RUNNING
		gs.GameLogic[0].start <- 1
	}
}

func handleClient(client *Client, globalState *GlobalState,
	gameLogicExit chan int) {
	log.WithFields(log.Fields{
		"remote address": client.Conn.RemoteAddr(),
	}).Debug("New connection")

	defer client.Conn.Close()

	go readClientMessages(client)

	msg := <-client.incomingMessages
	if msg.err != nil {
		log.WithFields(log.Fields{
			"err":            msg.err,
			"remote address": client.Conn.RemoteAddr(),
		}).Debug("Cannot receive client first message")
		Kick(client, fmt.Sprintf("Invalid first message: %v", msg.err.Error()))
		return
	}

	loginMessage, err := readLoginMessage(msg.content)
	if err != nil {
		log.WithFields(log.Fields{
			"err":            err,
			"remote address": client.Conn.RemoteAddr(),
		}).Debug("Cannot read LOGIN message")
		Kick(client, fmt.Sprintf("Invalid first message: %v", err.Error()))
		return
	}
	client.nickname = loginMessage.nickname

	LockGlobalStateMutex(globalState, "New client", "Login manager")
	switch loginMessage.role {
	case "player":
		if globalState.GameState != GAME_NOT_RUNNING {
			UnlockGlobalStateMutex(globalState, "New client", "Login manager")
			Kick(client, "LOGIN denied: Game has been started")
		} else if len(globalState.Players) >= globalState.NbPlayersMax {
			UnlockGlobalStateMutex(globalState, "New client", "Login manager")
			Kick(client, "LOGIN denied: Maximum number of players reached")
		} else {
			err = sendLoginACK(client)
			if err != nil {
				UnlockGlobalStateMutex(globalState, "New client", "Login manager")
				Kick(client, "LOGIN denied: Could not send LOGIN_ACK")
			} else {
				pvClient := &PlayerOrVisuClient{
					client:     client,
					playerID:   -1,
					isPlayer:   true,
					gameStarts: make(chan MessageGameStarts),
					newTurn:    make(chan MessageTurn),
					gameEnds:   make(chan MessageGameEnds),
					playerInfo: nil,
				}

				globalState.Players = append(globalState.Players, pvClient)

				log.WithFields(log.Fields{
					"nickname":       client.nickname,
					"remote address": client.Conn.RemoteAddr(),
					"player count":   len(globalState.Players),
				}).Info("New player accepted")

				UnlockGlobalStateMutex(globalState, "New client", "Login manager")

				// Automatically start the game if conditions are met
				autostart(globalState)

				// Player behavior is handled in dedicated function.
				handlePlayerOrVisu(pvClient, globalState)
			}
		}
	case "visualization":
		if len(globalState.Visus) >= globalState.NbVisusMax {
			UnlockGlobalStateMutex(globalState, "New client", "Login manager")
			Kick(client, "LOGIN denied: Maximum number of visus reached")
		} else {
			err = sendLoginACK(client)
			if err != nil {
				UnlockGlobalStateMutex(globalState, "New client", "Login manager")
				Kick(client, "LOGIN denied: Could not send LOGIN_ACK")
			} else {
				pvClient := &PlayerOrVisuClient{
					client:     client,
					playerID:   -1,
					isPlayer:   false,
					gameStarts: make(chan MessageGameStarts),
					newTurn:    make(chan MessageTurn),
					gameEnds:   make(chan MessageGameEnds),
				}

				globalState.Visus = append(globalState.Visus, pvClient)

				log.WithFields(log.Fields{
					"nickname":       client.nickname,
					"remote address": client.Conn.RemoteAddr(),
					"visu count":     len(globalState.Visus),
				}).Info("New visualization accepted")

				UnlockGlobalStateMutex(globalState, "New client", "Login manager")

				// Automatically start the game if conditions are met
				autostart(globalState)

				// Visu behavior is handled in dedicated function.
				handlePlayerOrVisu(pvClient, globalState)
			}
		}
	case "game logic":
		if globalState.GameState != GAME_NOT_RUNNING {
			UnlockGlobalStateMutex(globalState, "New client", "Login manager")
			Kick(client, "LOGIN denied: Game has been started")
		} else if len(globalState.GameLogic) >= 1 {
			UnlockGlobalStateMutex(globalState, "New client", "Login manager")
			Kick(client, "LOGIN denied: A game logic is already logged in")
		} else {
			err = sendLoginACK(client)
			if err != nil {
				UnlockGlobalStateMutex(globalState, "New client", "Login manager")
				Kick(client, "LOGIN denied: Could not send LOGIN_ACK")
			} else {
				glClient := &GameLogicClient{
					client:             client,
					playerAction:       make(chan MessageDoTurnPlayerAction, 1),
					playerDisconnected: make(chan int, 1),
					start:              make(chan int, 1),
				}

				globalState.GameLogic = append(globalState.GameLogic, glClient)

				log.WithFields(log.Fields{
					"nickname":       client.nickname,
					"remote address": client.Conn.RemoteAddr(),
				}).Info("Game logic accepted")

				UnlockGlobalStateMutex(globalState, "New client", "Login manager")

				// Automatically start the game if conditions are met
				autostart(globalState)

				// Game logic behavior is handled in dedicated function
				handleGameLogic(glClient, globalState, gameLogicExit)
			}
		}
	}
}

func Kick(client *Client, reason string) {
	if client.state == CLIENT_KICKED {
		return
	}

	client.state = CLIENT_KICKED
	log.WithFields(log.Fields{
		"remote address": client.Conn.RemoteAddr(),
		"nickname":       client.nickname,
		"reason":         reason,
	}).Warn("Kicking client")

	msg := MessageKick{
		MessageType: "KICK",
		KickReason:  reason,
	}

	content, err := json.Marshal(msg)
	if err == nil {
		_ = sendMessage(client, content)
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

func sendLoginACK(client *Client) error {
	msg := MessageLoginAck{
		MessageType: "LOGIN_ACK",
	}

	content, err := json.Marshal(msg)
	if err == nil {
		err = sendMessage(client, content)
	}
	return err
}

func Cleanup() {
	LockGlobalStateMutex(globalGS, "Cleanup", "Main")
	log.Warn("Closing listening socket.")
	globalGS.Listener.Close()

	nbClients := len(globalGS.Players) + len(globalGS.Visus) +
		len(globalGS.GameLogic)
	if nbClients > 0 {
		log.Warn("Sending KICK messages to clients")
		kickChan := make(chan int)
		for _, client := range append(globalGS.Players, globalGS.Visus...) {
			go func(c *Client) {
				Kick(c, "netorcai abort")
				kickChan <- 0
			}(client.client)
		}
		for _, client := range globalGS.GameLogic {
			go func(c *Client) {
				Kick(c, "netorcai abort")
				kickChan <- 0
			}(client.client)
		}

		for i := 0; i < nbClients; i++ {
			<-kickChan
		}

		log.Warn("Closing client sockets")
		for _, client := range append(globalGS.Players, globalGS.Visus...) {
			client.client.Conn.Close()
		}
		for _, client := range globalGS.GameLogic {
			client.client.Conn.Close()
		}
	}

	UnlockGlobalStateMutex(globalGS, "Cleanup", "Main")
}
