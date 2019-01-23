package netorcai

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type GameLogicClient struct {
	client *Client
	// Messages to aggregate from player clients
	playerAction chan MessageDoTurnPlayerAction
	// Control messages
	start chan int
}

func handleGameLogic(glClient *GameLogicClient, globalState *GlobalState,
	onexit chan int) {
	// Wait for the game to start
	select {
	case <-glClient.start:
		log.Info("Starting game")
	case msg := <-glClient.client.incomingMessages:
		globalState.Mutex.Lock()
		if msg.err == nil {
			Kick(glClient.client, "Received a game logic message but "+
				"the game has not started")
		} else {
			Kick(glClient.client, fmt.Sprintf("Game logic error. %v",
				msg.err.Error()))
		}
		globalState.GameLogic = globalState.GameLogic[:0]
		globalState.Mutex.Unlock()
		onexit <- 1
		return
	}

	// Generate randomized player identifiers
	globalState.Mutex.Lock()
	initialNbPlayers := len(globalState.Players)
	playerIDs := rand.Perm(len(globalState.Players))
	for playerIndex, player := range globalState.Players {
		player.playerID = playerIDs[playerIndex]
	}

	// Generate player information
	playersInfo := []*PlayerInformation{}
	for _, player := range globalState.Players {
		info := &PlayerInformation{
			PlayerID:      player.playerID,
			Nickname:      player.client.nickname,
			RemoteAddress: player.client.Conn.RemoteAddr().String(),
			IsConnected:   true,
		}
		player.playerInfo = info
		playersInfo = append(playersInfo, info)
	}

	// Sort player information by player_id
	sort.Slice(playersInfo, func(i, j int) bool {
		return playersInfo[i].PlayerID < playersInfo[j].PlayerID
	})

	// Send DO_INIT
	err := sendDoInit(glClient, len(globalState.Players),
		globalState.NbTurnsMax)
	globalState.Mutex.Unlock()

	if err != nil {
		Kick(glClient.client, fmt.Sprintf("Cannot send DO_INIT. %v",
			err.Error()))
		onexit <- 1
		return
	}

	// Wait for first turn (DO_INIT_ACK)
	var msg ClientMessage
	select {
	case msg = <-glClient.client.incomingMessages:
		if msg.err != nil {
			Kick(glClient.client,
				fmt.Sprintf("Cannot read DO_INIT_ACK. %v", msg.err.Error()))
			onexit <- 1
			return
		}
	case <-time.After(3 * time.Second):
		Kick(glClient.client, "Did not receive DO_INIT_ACK after 3 seconds.")
		onexit <- 1
		return
	}

	doTurnAckMsg, err := readDoInitAckMessage(msg.content)
	if err != nil {
		Kick(glClient.client,
			fmt.Sprintf("Invalid DO_INIT_ACK message. %v", err.Error()))
		onexit <- 1
		return
	}

	// Send GAME_STARTS to all clients
	globalState.Mutex.Lock()
	for _, player := range globalState.Players {
		player.gameStarts <- MessageGameStarts{
			MessageType:      "GAME_STARTS",
			PlayerID:         player.playerID,
			PlayersInfo:      []*PlayerInformation{},
			NbPlayers:        initialNbPlayers,
			NbTurnsMax:       globalState.NbTurnsMax,
			DelayFirstTurn:   globalState.MillisecondsBeforeFirstTurn,
			DelayTurns:       globalState.MillisecondsBetweenTurns,
			InitialGameState: doTurnAckMsg.InitialGameState,
		}
	}

	for _, visu := range globalState.Visus {
		visu.gameStarts <- MessageGameStarts{
			MessageType:      "GAME_STARTS",
			PlayerID:         visu.playerID,
			PlayersInfo:      playersInfo,
			NbPlayers:        initialNbPlayers,
			NbTurnsMax:       globalState.NbTurnsMax,
			DelayFirstTurn:   globalState.MillisecondsBeforeFirstTurn,
			DelayTurns:       globalState.MillisecondsBetweenTurns,
			InitialGameState: doTurnAckMsg.InitialGameState,
		}
	}
	globalState.Mutex.Unlock()

	// Wait before really starting the game
	log.WithFields(log.Fields{
		"duration (ms)": globalState.MillisecondsBeforeFirstTurn,
	}).Debug("Sleeping before first turn")
	time.Sleep(time.Duration(globalState.MillisecondsBeforeFirstTurn) *
		time.Millisecond)

	// Order the game logic to compute a TURN (without any action)
	turnNumber := 0
	playerActions := make([]MessageDoTurnPlayerAction, 0)
	var playerActionsMutex sync.Mutex
	sendDoTurn(glClient, playerActions)

	for {
		select {
		case action := <-glClient.playerAction:
			// A client sent its actions.
			// Replace the current message from this player if it exists,
			// and place it at the end of the array.
			// This may happen if the client was late in a previous turn but
			// catched up in current turn by sending two TURN_ACK.
			playerActionsMutex.Lock()
			actionFound := false
			for actionIndex, act := range playerActions {
				if act.PlayerID == action.PlayerID {
					playerActions[len(playerActions)-1], playerActions[actionIndex] = playerActions[actionIndex], playerActions[len(playerActions)-1]
					playerActions[len(playerActions)-1] = action
					actionFound = true
					break
				}
			}

			if !actionFound {
				// Append the action into the actions array
				playerActions = append(playerActions, action)
			}

			playerActionsMutex.Unlock()

		case msg := <-glClient.client.incomingMessages:
			// New message received from the game logic
			if msg.err != nil {
				Kick(glClient.client,
					fmt.Sprintf("Cannot read DO_TURN_ACK. %v",
						msg.err.Error()))
				onexit <- 1
				return
			}

			doTurnAckMsg, err := readDoTurnAckMessage(msg.content,
				initialNbPlayers)
			if err != nil {
				Kick(glClient.client,
					fmt.Sprintf("Invalid DO_TURN_ACK message. %v",
						err.Error()))
				onexit <- 1
				return
			}

			turnNumber = turnNumber + 1
			if turnNumber < globalState.NbTurnsMax {
				// Forward the TURN to the clients
				globalState.Mutex.Lock()
				for _, player := range globalState.Players {
					player.newTurn <- MessageTurn{
						MessageType: "TURN",
						TurnNumber:  turnNumber - 1,
						GameState:   doTurnAckMsg.GameState,
						PlayersInfo: []*PlayerInformation{},
					}
				}
				for _, visu := range globalState.Visus {
					visu.newTurn <- MessageTurn{
						MessageType: "TURN",
						TurnNumber:  turnNumber - 1,
						GameState:   doTurnAckMsg.GameState,
						PlayersInfo: playersInfo,
					}
				}
				globalState.Mutex.Unlock()

				// Trigger a new DO_TURN in some time
				go func() {
					log.WithFields(log.Fields{
						"duration (ms)": globalState.MillisecondsBetweenTurns,
					}).Debug("Sleeping before next turn")
					time.Sleep(time.Duration(
						globalState.MillisecondsBetweenTurns) *
						time.Millisecond)

					playerActionsMutex.Lock()
					// Send current actions
					sendDoTurn(glClient, playerActions)
					// Clear actions array
					playerActions = playerActions[:0]
					playerActionsMutex.Unlock()
				}()
			} else {
				if doTurnAckMsg.WinnerPlayerID != -1 {
					log.WithFields(log.Fields{
						"winner player ID":      doTurnAckMsg.WinnerPlayerID,
						"winner nickname":       playersInfo[doTurnAckMsg.WinnerPlayerID].Nickname,
						"winner remote address": playersInfo[doTurnAckMsg.WinnerPlayerID].RemoteAddress,
					}).Info("Game is finished")
				} else {
					log.Info("Game is finished (no winner!)")
				}

				// Send GAME_ENDS to all clients
				globalState.Mutex.Lock()
				for _, player := range globalState.Players {
					player.gameEnds <- MessageGameEnds{
						MessageType:    "GAME_ENDS",
						WinnerPlayerID: doTurnAckMsg.WinnerPlayerID,
						GameState:      doTurnAckMsg.GameState,
					}
				}
				for _, visu := range globalState.Visus {
					visu.gameEnds <- MessageGameEnds{
						MessageType:    "GAME_ENDS",
						WinnerPlayerID: doTurnAckMsg.WinnerPlayerID,
						GameState:      doTurnAckMsg.GameState,
					}
				}

				globalState.Mutex.Unlock()

				// Leave the program
				Kick(glClient.client, "Game is finished")
				onexit <- 0
				return
			}
		}
	}
}

func sendDoInit(client *GameLogicClient, nbPlayers, nbTurnsMax int) error {
	msg := MessageDoInit{
		MessageType: "DO_INIT",
		NbPlayers:   nbPlayers,
		NbTurnsMax:  nbTurnsMax,
	}

	content, err := json.Marshal(msg)
	if err == nil {
		log.WithFields(log.Fields{
			"nickname":       client.client.nickname,
			"remote address": client.client.Conn.RemoteAddr(),
			"content":        string(content),
		}).Debug("Sending DO_INIT to game logic")
		err = sendMessage(client.client, content)
	}
	return err
}

func sendDoTurn(client *GameLogicClient,
	playerActions []MessageDoTurnPlayerAction) error {
	msg := MessageDoTurn{
		MessageType:   "DO_TURN",
		PlayerActions: playerActions,
	}

	content, err := json.Marshal(msg)
	if err == nil {
		log.WithFields(log.Fields{
			"nickname":       client.client.nickname,
			"remote address": client.client.Conn.RemoteAddr(),
			"content":        string(content),
		}).Debug("Sending DO_TURN to game logic")
		err = sendMessage(client.client, content)
	}
	return err
}
