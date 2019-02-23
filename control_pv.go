package netorcai

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type PlayerOrVisuClient struct {
	client          *Client
	playerID        int
	isPlayer        bool
	isSpecialPlayer bool
	gameStarts      chan MessageGameStarts
	newTurn         chan MessageTurn
	gameEnds        chan MessageGameEnds
	playerInfo      *PlayerInformation
}

func waitPlayerOrVisuFinition(pvClient *PlayerOrVisuClient) {
	for {
		select {
		case kickReason := <-pvClient.client.canTerminate:
			Kick(pvClient.client, kickReason)
			return
		case <-pvClient.client.incomingMessages:
		}
	}
}

func handlePlayerOrVisu(pvClient *PlayerOrVisuClient,
	globalState *GlobalState) {
	turnBuffer := make([]MessageTurn, 0)
	lastTurnNumberSent := -1
	var glClient *GameLogicClient

	for {
		select {
		case kickReason := <-pvClient.client.canTerminate:
			Kick(pvClient.client, kickReason)
			return
		case gameStarts := <-pvClient.gameStarts:
			// A game start has been received.
			err := sendGameStarts(pvClient.client, gameStarts)
			if err != nil {
				KickLoggedPlayerOrVisu(pvClient, globalState,
					fmt.Sprintf("Cannot send GAME_STARTS. %v", err.Error()))
				return
			}
			pvClient.client.state = CLIENT_READY

			// Set glClient from the global state now
			LockGlobalStateMutex(globalState, "Local copy of GL pointer", "client")
			glClient = globalState.GameLogic[0]
			UnlockGlobalStateMutex(globalState, "Local copy of GL pointer", "client")
		case gameEnds := <-pvClient.gameEnds:
			// A game end has been received.
			err := sendGameEnds(pvClient.client, gameEnds)
			if err != nil {
				KickLoggedPlayerOrVisu(pvClient, globalState,
					fmt.Sprintf("Cannot send GAME_ENDS. %v", err.Error()))
				return
			}

			// Leave the client
			Kick(pvClient.client, "Game is finished")
			waitPlayerOrVisuFinition(pvClient)
			return
		case turn := <-pvClient.newTurn:
			// A new turn has been received.
			log.WithFields(log.Fields{
				"playerID": pvClient.playerID,
			}).Debug("Client received a new TURN (from GL goroutine)")

			if pvClient.client.state == CLIENT_READY {
				// The client is ready, the message can be sent right now.
				lastTurnNumberSent = turn.TurnNumber
				err := sendTurn(pvClient.client, turn)
				if err != nil {
					KickLoggedPlayerOrVisu(pvClient, globalState,
						fmt.Sprintf("Cannot send TURN. %v", err.Error()))
					return
				}
				pvClient.client.state = CLIENT_THINKING
			} else if pvClient.client.state == CLIENT_THINKING {
				// The client is still computing something (its decisions for
				// a player, or just updating its display for a visualization).
				// The turn message is therefore buffered.
				if len(turnBuffer) > 0 {
					// Update the turn buffer with the new message.
					turnBuffer[0] = turn
				} else {
					// Put the new message into the turn buffer.
					turnBuffer = append(turnBuffer, turn)
				}
			}
		case msg := <-pvClient.client.incomingMessages:
			// A new message has been received from the player socket.
			if msg.err != nil {
				KickLoggedPlayerOrVisu(pvClient, globalState,
					fmt.Sprintf("Cannot read TURN_ACK. %v", msg.err.Error()))
				return
			}
			turnAckMsg, err := readTurnAckMessage(msg.content,
				lastTurnNumberSent)
			if err != nil {
				KickLoggedPlayerOrVisu(pvClient, globalState,
					fmt.Sprintf("Invalid TURN_ACK received. %v",
						err.Error()))
				return
			}

			log.WithFields(log.Fields{
				"playerID": pvClient.playerID,
			}).Debug("Client received a TURN_ACK (from socket)")

			// Check client state
			if pvClient.client.state != CLIENT_THINKING {
				KickLoggedPlayerOrVisu(pvClient, globalState,
					"Received a TURN_ACK but the client state is not THINKING")
				return
			}

			if pvClient.isPlayer {
				// Forward the player actions to the game logic
				glClient.playerAction <- MessageDoTurnPlayerAction{
					PlayerID:   pvClient.playerID,
					TurnNumber: turnAckMsg.turnNumber,
					Actions:    turnAckMsg.actions,
				}
			}

			// If a TURN is buffered, send it right now.
			if len(turnBuffer) > 0 {
				lastTurnNumberSent = turnBuffer[0].TurnNumber
				err := sendTurn(pvClient.client, turnBuffer[0])
				if err != nil {
					KickLoggedPlayerOrVisu(pvClient, globalState,
						fmt.Sprintf("Cannot send TURN. %v", err.Error()))
					return
				}

				// Empty turn buffer
				turnBuffer = turnBuffer[:0]
				pvClient.client.state = CLIENT_THINKING
			} else {
				pvClient.client.state = CLIENT_READY
			}
		}
	}
}

func KickLoggedPlayerOrVisu(pvClient *PlayerOrVisuClient,
	gs *GlobalState, reason string) {
	// Remove the client from the global state
	LockGlobalStateMutex(gs, "Kick player or visu", "player/visu")

	if pvClient.isPlayer {
		// Mark the player as disconnected
		if pvClient.playerInfo != nil {
			pvClient.playerInfo.IsConnected = false
		}

		if pvClient.isSpecialPlayer {
			// Locate the player in the array
			playerIndex := -1
			for index, player := range gs.SpecialPlayers {
				if player.client == pvClient.client {
					playerIndex = index
					break
				}
			}

			if gs.GameState == GAME_RUNNING && gs.Fast {
				gs.GameLogic[0].playerDisconnected <- pvClient.playerID
			}

			if playerIndex != -1 {
				// Remove the player by placing it at the end of the slice,
				// then reducing the slice length
				gs.SpecialPlayers[len(gs.SpecialPlayers)-1], gs.SpecialPlayers[playerIndex] = gs.SpecialPlayers[playerIndex], gs.SpecialPlayers[len(gs.SpecialPlayers)-1]
				gs.SpecialPlayers = gs.SpecialPlayers[:len(gs.SpecialPlayers)-1]
			}
		} else {
			// Locate the player in the array
			playerIndex := -1
			for index, player := range gs.Players {
				if player.client == pvClient.client {
					playerIndex = index
					break
				}
			}

			if gs.GameState == GAME_RUNNING && gs.Fast {
				gs.GameLogic[0].playerDisconnected <- pvClient.playerID
			}

			if playerIndex != -1 {
				// Remove the player by placing it at the end of the slice,
				// then reducing the slice length
				gs.Players[len(gs.Players)-1], gs.Players[playerIndex] = gs.Players[playerIndex], gs.Players[len(gs.Players)-1]
				gs.Players = gs.Players[:len(gs.Players)-1]
			}
		}
	} else {
		// Locate the visu in the array
		visuIndex := -1
		for index, visu := range gs.Visus {
			if visu.client == pvClient.client {
				visuIndex = index
				break
			}
		}

		if visuIndex != -1 {
			// Remove the visu by placing it at the end of the slice,
			// then reducing the slice length
			gs.Visus[len(gs.Visus)-1], gs.Visus[visuIndex] = gs.Visus[visuIndex], gs.Visus[len(gs.Visus)-1]
			gs.Visus = gs.Visus[:len(gs.Visus)-1]
		}
	}

	UnlockGlobalStateMutex(gs, "Kick player or visu", "player/visu")

	// Kick the client
	Kick(pvClient.client, reason)
}

func sendGameStarts(client *Client, msg MessageGameStarts) error {
	content, err := json.Marshal(msg)
	if err == nil {
		log.WithFields(log.Fields{
			"nickname":       client.nickname,
			"remote address": client.Conn.RemoteAddr(),
			"content":        string(content),
		}).Debug("Sending GAME_STARTS to client")
		err = sendMessage(client, content)
	}
	return err
}

func sendTurn(client *Client, msg MessageTurn) error {
	content, err := json.Marshal(msg)
	if err == nil {
		log.WithFields(log.Fields{
			"nickname":       client.nickname,
			"remote address": client.Conn.RemoteAddr(),
			"content":        string(content),
		}).Debug("Sending TURN to client")
		err = sendMessage(client, content)
	}
	return err
}

func sendGameEnds(client *Client, msg MessageGameEnds) error {
	content, err := json.Marshal(msg)
	if err == nil {
		log.WithFields(log.Fields{
			"nickname":       client.nickname,
			"remote address": client.Conn.RemoteAddr(),
			"content":        string(content),
		}).Debug("Sending GAME_ENDS to client")
		err = sendMessage(client, content)
	}
	return err
}
