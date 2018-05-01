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
	CLIENT_THINKING = iota
	CLIENT_FINISHED = iota
	CLIENT_KICKED   = iota
)

type PlayerClient struct {
	client   *Client
	newTurn  chan MessageTurn
	gameEnds chan MessageGameEnds
}

type GameLogicClient struct {
	client *Client
}

type VisuClient struct {
	client  *Client
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

func handleClient(client *Client, globalState *GlobalState) {
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
					newTurn: make(chan MessageTurn),
				}

				globalState.players = append(globalState.players, playerClient)

				log.WithFields(log.Fields{
					"nickname":       client.nickname,
					"remote address": client.conn.RemoteAddr(),
					"player count":   len(globalState.players),
				}).Info("New player accepted")

				globalState.mutex.Unlock()

				// Player behavior is handled in dedicated function.
				handlePlayer(&playerClient, globalState)
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

				log.WithFields(log.Fields{
					"nickname":       client.nickname,
					"remote address": client.conn.RemoteAddr(),
					"visu count":     len(globalState.visus),
				}).Info("New visualization accepted")

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

				log.WithFields(log.Fields{
					"nickname":       client.nickname,
					"remote address": client.conn.RemoteAddr(),
				}).Info("Game logic accepted")

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

func handlePlayer(playerClient *PlayerClient, globalState *GlobalState) {
	turnBuffer := make([]MessageTurn, 1)
	lastTurnNumberSent := -1

	for {
		select {
		case turn := <-playerClient.newTurn:
			// A new turn has been received.
			if playerClient.client.state == CLIENT_READY {
				// If no turn is buffered and if the player is waiting for
				// a new turn, directly send the turn to the player.
				lastTurnNumberSent = turn.TurnNumber
				err := sendTurn(playerClient.client, turn)
				if err != nil {
					kickLoggedPlayer(playerClient.client, globalState,
						fmt.Sprintf("Cannot send TURN. %v", err.Error()))
					return
				}
				playerClient.client.state = CLIENT_THINKING
			} else if len(turnBuffer) > 0 {
				// The client is not ready, and a message is already buffered.
				// Update the turn buffer with the new message.
				turnBuffer[0] = turn
			} else {
				// The client is not ready, and the turn buffer is empty.
				// Put the new message into the turn buffer.
				turnBuffer = append(turnBuffer, turn)
			}
		case msg := <-playerClient.client.incomingMessages:
			// A new message has been received from the player socket.
			if msg.err != nil {
				kickLoggedPlayer(playerClient.client, globalState,
					fmt.Sprintf("Cannot read TURN_ACK. %v", msg.err.Error()))
				return
			}
			_, err := readTurnACKMessage(msg.content, lastTurnNumberSent)
			if err != nil {
				kickLoggedPlayer(playerClient.client, globalState,
					fmt.Sprintf("Invalid TURN_ACK received. %v",
						err.Error()))
				return
			}

			// TODO: transmit message to game logic

			if len(turnBuffer) > 0 {
				// If a TURN is buffered, send it right now.
				err := sendTurn(playerClient.client, turnBuffer[0])
				if err != nil {
					kickLoggedPlayer(playerClient.client, globalState,
						fmt.Sprintf("Cannot send TURN. %v", err.Error()))
					return
				}

				// Empty turn buffer
				turnBuffer = turnBuffer[0:]
				playerClient.client.state = CLIENT_THINKING
			} else {
				playerClient.client.state = CLIENT_READY
			}
		}
	}
}

func kick(client *Client, reason string) {
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

func kickLoggedPlayer(client *Client, gs *GlobalState, reason string) {
	// Remove the player from the global state
	gs.mutex.Lock()

	// Locate the player in the array
	playerIndex := -1
	for index, value := range gs.players {
		if value.client == client {
			playerIndex = index
			break
		}
	}

	if playerIndex == -1 {
		log.Error("Could not remove player: Did not find it")
	} else {
		// Remove the player by placing it at the end of the slice,
		// then reducing the slice length
		gs.players[len(gs.players)-1], gs.players[playerIndex] = gs.players[playerIndex], gs.players[len(gs.players)-1]
		gs.players = gs.players[:len(gs.players)-1]
	}

	gs.mutex.Unlock()

	// Kick the player
	kick(client, reason)
}

func sendLoginACK(client *Client) error {
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

func sendTurn(client *Client, msg MessageTurn) error {
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
