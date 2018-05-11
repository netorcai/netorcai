package netorcai

import (
	"fmt"
	"regexp"
)

type MessageLogin struct {
	nickname string
	role     string
}

type MessageLoginAck struct {
	MessageType string `json:"message_type"`
}

type MessageGameStarts struct {
	PlayerID         int                    `json:"player_id"`
	NbPlayers        int                    `json:"nb_players"`
	NbTurnsMax       int                    `json:"nb_turns_max"`
	DelayFirstTurn   float64                `json:"milliseconds_before_first_turn"`
	InitialGameState map[string]interface{} `json:"initial_game_state"`
}

type MessageGameEnds struct {
	WinnerPlayerID int                    `json:"winner_player_id"`
	GameState      map[string]interface{} `json:"game_state"`
}

type MessageTurn struct {
	MessageType string                 `json:"message_type"`
	TurnNumber  int                    `json:"turn_number"`
	GameState   map[string]interface{} `json:"game_state"`
}

type MessageTurnAck struct {
	turnNumber int
	actions    []interface{}
}

type MessageDoInit struct {
	MessageType string `json:"message_type"`
	NbPlayers   int    `json:"nb_players"`
	NbTurnsMax  int    `json:"nb_turns_max"`
}

type MessageDoInitAck struct {
	InitialGameState map[string]interface{}
}

type MessageDoTurnPlayerAction struct {
	PlayerID   int           `json:"player_id"`
	TurnNumber int           `json:"turn_number"`
	Actions    []interface{} `json:"actions"`
}

type MessageDoTurn struct {
	MessageType   string                      `json:"message_type"`
	PlayerActions []MessageDoTurnPlayerAction `json:"player_actions"`
}

type MessageDoTurnAck struct {
	WinnerPlayerID int
	GameState      map[string]interface{}
}

type MessageKick struct {
	MessageType string `json:"message_type"`
	KickReason  string `json:"kick_reason"`
}

func checkMessageType(data map[string]interface{}, expectedMessageType string) error {
	messageType, err := ReadString(data, "message_type")
	if err != nil {
		return err
	}

	if messageType != expectedMessageType {
		return fmt.Errorf("Received '%v' message type, "+
			"while %v was expected", messageType, expectedMessageType)
	}

	return nil
}

func readLoginMessage(data map[string]interface{}) (MessageLogin, error) {
	var readMessage MessageLogin

	// Check message type
	err := checkMessageType(data, "LOGIN")
	if err != nil {
		return readMessage, err
	}

	// Read nickname
	readMessage.nickname, err = ReadString(data, "nickname")
	if err != nil {
		return readMessage, err
	}

	// Check nickname
	r, _ := regexp.Compile(`\A\S{1,10}\z`)
	if !r.MatchString(readMessage.nickname) {
		return readMessage, fmt.Errorf("Invalid nickname")
	}

	// Read role
	readMessage.role, err = ReadString(data, "role")
	if err != nil {
		return readMessage, err
	}

	// Check role
	switch readMessage.role {
	case "player",
		"visualization",
		"game logic":
		return readMessage, nil
	default:
		return readMessage, fmt.Errorf("Invalid role '%v'",
			readMessage.role)
	}
}

func readTurnAckMessage(data map[string]interface{}, expectedTurnNumber int) (
	MessageTurnAck, error) {
	var readMessage MessageTurnAck

	// Check message type
	err := checkMessageType(data, "TURN_ACK")
	if err != nil {
		return readMessage, err
	}

	// Read turn number
	readMessage.turnNumber, err = ReadInt(data, "turn_number")
	if err != nil {
		return readMessage, err
	}

	// Check turn number
	if readMessage.turnNumber != expectedTurnNumber {
		return readMessage, fmt.Errorf("Invalid value (turn_number=%v): "+
			"expecting %v", readMessage.turnNumber, expectedTurnNumber)
	}

	// Read actions
	readMessage.actions, err = ReadArray(data, "actions")
	if err != nil {
		return readMessage, err
	}

	return readMessage, nil
}

func readDoInitAckMessage(data map[string]interface{}) (
	MessageDoInitAck, error) {
	var readMessage MessageDoInitAck

	// Check message type
	err := checkMessageType(data, "DO_INIT_ACK")
	if err != nil {
		return readMessage, err
	}

	// Read game state
	gameState, err := ReadObject(data, "initial_game_state")
	if err != nil {
		return readMessage, err
	}

	// Read game state -> all clients
	readMessage.InitialGameState, err = ReadObject(gameState, "all_clients")
	if err != nil {
		return readMessage, err
	}

	return readMessage, nil
}

func readDoTurnAckMessage(data map[string]interface{}, nbPlayers int) (
	MessageDoTurnAck, error) {
	var readMessage MessageDoTurnAck

	// Check message type
	err := checkMessageType(data, "DO_TURN_ACK")
	if err != nil {
		return readMessage, err
	}

	// Read winner player id
	readMessage.WinnerPlayerID, err = ReadInt(data, "winner_player_id")
	if err != nil {
		return readMessage, err
	}

	// Check player id
	if readMessage.WinnerPlayerID < 0 ||
		readMessage.WinnerPlayerID >= nbPlayers {
		return readMessage, fmt.Errorf("Invalid winner_player_id: "+
			"Not in [0, %v[", nbPlayers)
	}

	// Read game state
	gameState, err := ReadObject(data, "game_state")
	if err != nil {
		return readMessage, err
	}

	// Read game state -> all clients
	readMessage.GameState, err = ReadObject(gameState, "all_clients")
	if err != nil {
		return readMessage, err
	}

	return readMessage, nil
}
