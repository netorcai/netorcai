package netorcai

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
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
	CLIENT_FINISHED = iota
	CLIENT_KICKED   = iota
)

type PlayerOrVisuClient struct {
	client     *Client
	playerID   int // TODO: generate them when the game is started
	isPlayer   bool
	gameStarts chan MessageGameStarts
	newTurn    chan MessageTurn
	gameEnds   chan MessageGameEnds
}

type GameLogicClient struct {
	client *Client
	// Messages to aggregate from player clients
	playerAction chan MessageDoTurnPlayerAction
	// Control messages
	start chan int
}

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
	MillisecondsBeforeFirstTurn float64
	MillisecondsBetweenTurns    float64
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

	globalState.Mutex.Lock()
	switch loginMessage.role {
	case "player":
		if globalState.GameState != GAME_NOT_RUNNING {
			globalState.Mutex.Unlock()
			Kick(client, "LOGIN denied: Game has been started")
		} else if len(globalState.Players) >= globalState.NbPlayersMax {
			globalState.Mutex.Unlock()
			Kick(client, "LOGIN denied: Maximum number of players reached")
		} else {
			err = sendLoginACK(client)
			if err != nil {
				globalState.Mutex.Unlock()
				Kick(client, "LOGIN denied: Could not send LOGIN_ACK")
			} else {
				pvClient := &PlayerOrVisuClient{
					client:   client,
					playerID: -1,
					isPlayer: true,
					newTurn:  make(chan MessageTurn),
				}

				globalState.Players = append(globalState.Players, pvClient)

				log.WithFields(log.Fields{
					"nickname":       client.nickname,
					"remote address": client.Conn.RemoteAddr(),
					"player count":   len(globalState.Players),
				}).Info("New player accepted")

				globalState.Mutex.Unlock()

				// Player behavior is handled in dedicated function.
				handlePlayerOrVisu(pvClient, globalState)
			}
		}
	case "visualization":
		if len(globalState.Visus) >= globalState.NbVisusMax {
			globalState.Mutex.Unlock()
			Kick(client, "LOGIN denied: Maximum number of visus reached")
		} else {
			err = sendLoginACK(client)
			if err != nil {
				globalState.Mutex.Unlock()
				Kick(client, "LOGIN denied: Could not send LOGIN_ACK")
			} else {
				pvClient := &PlayerOrVisuClient{
					client:   client,
					playerID: -1,
					isPlayer: false,
					newTurn:  make(chan MessageTurn),
				}

				globalState.Visus = append(globalState.Visus, pvClient)

				log.WithFields(log.Fields{
					"nickname":       client.nickname,
					"remote address": client.Conn.RemoteAddr(),
					"visu count":     len(globalState.Visus),
				}).Info("New visualization accepted")

				globalState.Mutex.Unlock()

				// Visu behavior is handled in dedicated function.
				handlePlayerOrVisu(pvClient, globalState)
			}
		}
	case "game logic":
		if globalState.GameState != GAME_NOT_RUNNING {
			globalState.Mutex.Unlock()
			Kick(client, "LOGIN denied: Game has been started")
		} else if len(globalState.GameLogic) >= 1 {
			globalState.Mutex.Unlock()
			Kick(client, "LOGIN denied: A game logic is already logged in")
		} else {
			err = sendLoginACK(client)
			if err != nil {
				globalState.Mutex.Unlock()
				Kick(client, "LOGIN denied: Could not send LOGIN_ACK")
			} else {
				glClient := &GameLogicClient{
					client: client,
				}

				globalState.GameLogic = append(globalState.GameLogic, glClient)

				log.WithFields(log.Fields{
					"nickname":       client.nickname,
					"remote address": client.Conn.RemoteAddr(),
				}).Info("Game logic accepted")

				globalState.Mutex.Unlock()

				handleGameLogic(glClient, globalState, gameLogicExit)
			}
		}
	default:
		globalState.Mutex.Unlock()
		Kick(client, fmt.Sprintf("LOGIN denied: Unknown role '%v'",
			loginMessage.role))
	}
}

func handlePlayerOrVisu(pvClient *PlayerOrVisuClient,
	globalState *GlobalState) {
	turnBuffer := make([]MessageTurn, 1)
	lastTurnNumberSent := -1

	for {
		select {
		case turn := <-pvClient.newTurn:
			// A new turn has been received.
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
				// The turn message therefore buffered.
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

			// Check client state
			if pvClient.client.state != CLIENT_THINKING {
				KickLoggedPlayerOrVisu(pvClient, globalState,
					"Received a TURN_ACK but the client state is not THINKING")
				return
			}

			// Check turnNumber value
			if turnAckMsg.turnNumber != lastTurnNumberSent {
				KickLoggedPlayerOrVisu(pvClient, globalState,
					fmt.Sprintf("Invalid TURN_ACK received: "+
						"Expected turn_number=%v, got %v", lastTurnNumberSent,
						turnAckMsg.turnNumber))
			}

			if pvClient.isPlayer {
				// Forward the player actions to the game logic
				globalState.Mutex.Lock()
				if len(globalState.GameLogic) == 1 {
					globalState.GameLogic[0].playerAction <- MessageDoTurnPlayerAction{
						PlayerID:   pvClient.playerID,
						TurnNumber: turnAckMsg.turnNumber,
						Actions:    turnAckMsg.actions,
					}
				}
				globalState.Mutex.Unlock()
			}

			// If a TURN is buffered, send it right now.
			if len(turnBuffer) > 0 {
				err := sendTurn(pvClient.client, turnBuffer[0])
				if err != nil {
					KickLoggedPlayerOrVisu(pvClient, globalState,
						fmt.Sprintf("Cannot send TURN. %v", err.Error()))
					return
				}

				// Empty turn buffer
				turnBuffer = turnBuffer[0:]
				pvClient.client.state = CLIENT_THINKING
			} else {
				pvClient.client.state = CLIENT_READY
			}
		}
	}
}

func handleGameLogic(glClient *GameLogicClient, globalState *GlobalState,
	onexit chan int) {
	// Wait for the game to start
	select {
	case <-glClient.start:
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
		return
	}

	// Generate randomized player identifiers
	globalState.Mutex.Lock()
	playerIDs := rand.Perm(len(globalState.Players))
	for playerIndex, player := range globalState.Players {
		player.playerID = playerIDs[playerIndex]
	}

	// Send DO_FIRST_TURN
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
	msg := <-glClient.client.incomingMessages
	if msg.err != nil {
		Kick(glClient.client,
			fmt.Sprintf("Cannot read DO_INIT_ACK. %v", msg.err.Error()))
		onexit <- 1
		return
	}

	doTurnAckMsg, err := readDoInitAckMessage(msg.content)
	if err != nil {
		Kick(glClient.client,
			fmt.Sprintf("Invalid DO_INIT_ACK message. %v", msg.err.Error()))
		onexit <- 1
		return
	}

	// Send GAME_STARTS to all clients
	globalState.Mutex.Lock()
	initialNbPlayers := len(globalState.Players)
	for _, player := range globalState.Players {
		player.gameStarts <- MessageGameStarts{
			PlayerID:         player.playerID,
			NbPlayers:        initialNbPlayers,
			NbTurnsMax:       globalState.NbTurnsMax,
			DelayFirstTurn:   globalState.MillisecondsBeforeFirstTurn,
			InitialGameState: doTurnAckMsg.InitialGameState,
		}
	}
	globalState.Mutex.Unlock()

	// Wait before really starting the game
	time.Sleep(time.Duration(globalState.MillisecondsBeforeFirstTurn) *
		time.Millisecond)

	// Order the game logic to compute a TURN (without any action)
	turnNumber := 0
	playerActions := make([]MessageDoTurnPlayerAction, 0)
	sendDoTurn(glClient, playerActions)

	for {
		select {
		case action := <-glClient.playerAction:
			// A client sent its actions.
			// Replace the current message from this player if it exists,
			// and place it at the end of the array.
			// This may happen if the client was late in a previous turn but
			// catched up in current turn by sending two TURN_ACK.
			actionFound := false
			for actionIndex, act := range playerActions {
				if act.PlayerID == action.PlayerID {
					playerActions[len(playerActions)-1], playerActions[actionIndex] = playerActions[actionIndex], playerActions[len(playerActions)-1]
					playerActions[len(playerActions)-1] = action
					break
				}
			}

			if !actionFound {
				// Append the action into the actions array
				playerActions = append(playerActions, action)
			}

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

			// Forward the TURN to the clients
			globalState.Mutex.Lock()
			for _, player := range globalState.Players {
				player.newTurn <- MessageTurn{
					TurnNumber: turnNumber,
					GameState:  doTurnAckMsg.GameState,
				}
			}
			for _, visu := range globalState.Visus {
				visu.newTurn <- MessageTurn{
					TurnNumber: turnNumber,
					GameState:  doTurnAckMsg.GameState,
				}
			}
			globalState.Mutex.Unlock()
			turnNumber = turnNumber + 1

			if turnNumber < globalState.NbTurnsMax {
				// Trigger a new DO_TURN in some time
				go func() {
					time.Sleep(time.Duration(globalState.MillisecondsBetweenTurns) *
						time.Millisecond)

					// Send current actions
					sendDoTurn(glClient, playerActions)
					// Clear actions array
					playerActions = playerActions[:0]
				}()
			} else {
				// Send GAME_ENDS to all clients
				globalState.Mutex.Lock()
				for _, player := range globalState.Players {
					player.gameEnds <- MessageGameEnds{
						WinnerPlayerID: doTurnAckMsg.WinnerPlayerID,
						GameState:      doTurnAckMsg.GameState,
					}
				}
				for _, visu := range globalState.Visus {
					visu.gameEnds <- MessageGameEnds{
						WinnerPlayerID: doTurnAckMsg.WinnerPlayerID,
						GameState:      doTurnAckMsg.GameState,
					}
				}

				globalState.Mutex.Unlock()
			}
		}
	}
}

func Kick(client *Client, reason string) {
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
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot marshal JSON message")
	} else {
		_ = sendMessage(client, content)
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
}

func KickLoggedPlayerOrVisu(pvClient *PlayerOrVisuClient,
	gs *GlobalState, reason string) {
	// Remove the client from the global state
	gs.Mutex.Lock()

	if pvClient.isPlayer {
		// Locate the player in the array
		playerIndex := -1
		for index, player := range gs.Players {
			if player.client == pvClient.client {
				playerIndex = index
				break
			}
		}

		if playerIndex == -1 {
			log.Error("Could not remove player: Did not find it")
		} else {
			// Remove the player by placing it at the end of the slice,
			// then reducing the slice length
			gs.Players[len(gs.Players)-1], gs.Players[playerIndex] = gs.Players[playerIndex], gs.Players[len(gs.Players)-1]
			gs.Players = gs.Players[:len(gs.Players)-1]
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

		if visuIndex == -1 {
			log.Error("Could not remove visu: Did not find it")
		} else {
			// Remove the visu by placing it at the end of the slice,
			// then reducing the slice length
			gs.Visus[len(gs.Visus)-1], gs.Visus[visuIndex] = gs.Visus[visuIndex], gs.Visus[len(gs.Visus)-1]
			gs.Visus = gs.Visus[:len(gs.Visus)-1]
		}
	}

	gs.Mutex.Unlock()

	// Kick the client
	Kick(pvClient.client, reason)
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

func sendDoInit(client *GameLogicClient, nbPlayers, nbTurnsMax int) error {
	msg := MessageDoInit{
		MessageType: "DO_INIT",
		NbPlayers:   nbPlayers,
		NbTurnsMax:  nbTurnsMax,
	}

	content, err := json.Marshal(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot marshal JSON message")
		return err
	} else {
		err = sendMessage(client.client, content)
		return err
	}
}

func sendDoTurn(client *GameLogicClient,
	playerActions []MessageDoTurnPlayerAction) error {
	msg := MessageDoTurn{
		MessageType:   "DO_TURN",
		PlayerActions: playerActions,
	}

	content, err := json.Marshal(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Cannot marshal JSON message")
		return err
	} else {
		err = sendMessage(client.client, content)
		return err
	}
}

func Cleanup() {
	globalGS.Mutex.Lock()
	log.Warn("Closing listening socket.")
	globalGS.Listener.Close()

	nbClients := len(globalGS.Players) + len(globalGS.Visus) +
		len(globalGS.GameLogic)
	if nbClients > 0 {
		log.Warn("Sending KICK messages to clients")
		kickChan := make(chan int)
		for _, client := range globalGS.Players {
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

	globalGS.Mutex.Unlock()
}
