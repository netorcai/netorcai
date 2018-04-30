package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	CLIENT_FINISHED = iota
	CLIENT_KICKED   = iota
)

type PlayerClient struct {
	client  Client
	newTurn chan string
}

type GameLogicClient struct {
	client Client
}

type VisuClient struct {
	client  Client
	newTurn chan string
}

type GlobalState struct {
	mutex sync.Mutex

	gameState int

	gameLogic []GameLogicClient
	players   []PlayerClient
	visus     []VisuClient

	nbPlayersMax                int
	nbVisusMax                  int
	nbTurnsMax                  int
	millisecondsBeforeFirstTurn float64
}

func handleClient(client Client, globalState *GlobalState) {
	log.WithFields(log.Fields{
		"remote address": client.conn.RemoteAddr(),
	}).Debug("New connection")

	defer client.conn.Close()

	go readClientMessages(client)

	msg := <-client.incomingMessages
	if msg.err != nil {
		log.WithFields(log.Fields{
			"err":            msg.err,
			"remote address": client.conn.RemoteAddr(),
		}).Debug("Cannot receive client first message")
		kick(client, fmt.Sprintf("Invalid first message: %v", msg.err.Error()))
		return
	}

	loginMessage, err := readLoginMessage(msg.content)
	if err != nil {
		log.WithFields(log.Fields{
			"err":            err,
			"remote address": client.conn.RemoteAddr(),
		}).Debug("Cannot read LOGIN message")
		kick(client, fmt.Sprintf("Invalid first message: %v", err.Error()))
		return
	}
	client.nickname = loginMessage.nickname

	globalState.mutex.Lock()
	switch loginMessage.role {
	case "player":
		if globalState.gameState != GAME_NOT_RUNNING {
			globalState.mutex.Unlock()
			kick(client, "LOGIN denied: Game is not running")
		} else if len(globalState.players) >= globalState.nbPlayersMax {
			globalState.mutex.Unlock()
			kick(client, "LOGIN denied: Maximum number of players reached")
		} else {
			err = sendLoginACK(client)
			if err != nil {
				globalState.mutex.Unlock()
				kick(client, "LOGIN denied: Could not send LOGIN_ACK")
			} else {
				playerClient := PlayerClient{
					client:  client,
					newTurn: make(chan string),
				}

				globalState.players = append(globalState.players, playerClient)
				globalState.mutex.Unlock()

				// TODO: call handlePlayer
			}
		}
	case "visualization":
		if len(globalState.visus) >= globalState.nbVisusMax {
			globalState.mutex.Unlock()
			kick(client, "LOGIN denied: Maximum number of visus reached")
		} else {
			err = sendLoginACK(client)
			if err != nil {
				globalState.mutex.Unlock()
				kick(client, "LOGIN denied: Could not send LOGIN_ACK")
			} else {
				visuClient := VisuClient{
					client:  client,
					newTurn: make(chan string),
				}

				globalState.visus = append(globalState.visus, visuClient)
				globalState.mutex.Unlock()

				// TODO: call handleVisu
			}
		}
	case "game logic":
		if len(globalState.gameLogic) >= 1 {
			globalState.mutex.Unlock()
			kick(client, "LOGIN denied: A game logic is already logged in")
		} else {
			err = sendLoginACK(client)
			if err != nil {
				globalState.mutex.Unlock()
				kick(client, "LOGIN denied: Could not send LOGIN_ACK")
			} else {
				glClient := GameLogicClient{
					client: client,
				}

				globalState.gameLogic = append(globalState.gameLogic, glClient)
				globalState.mutex.Unlock()

				// TODO: call handleGameLogic
			}
		}
	default:
		globalState.mutex.Unlock()
		kick(client, fmt.Sprintf("LOGIN denied: Unknown role '%v'",
			loginMessage.role))
	}

}

func kick(client Client, reason string) {
	client.state = CLIENT_KICKED
	log.WithFields(log.Fields{
		"remote address": client.conn.RemoteAddr(),
		"nickname":       client.nickname,
		"reason":         reason,
	}).Warn("Kicking client")

	msg := MessageKick{
		MessageType: "KICK",
		KickReason:  reason,
	}

	content, err := json.Marshal(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot marshal JSON message")
	} else {
		_ = sendMessage(client, content)
		time.Sleep(500 * time.Millisecond)
	}
}

func sendLoginACK(client Client) error {
	msg := MessageLoginAck{
		MessageType: "LOGIN_ACK",
	}

	content, err := json.Marshal(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot marshal JSON message")
		return err
	} else {
		err = sendMessage(client, content)
		return err
	}
}
