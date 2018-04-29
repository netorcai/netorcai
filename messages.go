package main

import (
	"fmt"
	"regexp"
)

type MessageLogin struct {
	nickname string
	role     string
}

type MessageLoginAck struct {
}

type MessageGameStarts struct {
	playerID  int                    `json:"player_id"`
	gameState map[string]interface{} `json:"game_state"`
}

type MessageTurn struct {
	turnNumber int                    `json:"turn_number"`
	gameState  map[string]interface{} `json:"game_state"`
}

type MessageTurnAck struct {
	turnNumber int
	actions    map[string]interface{}
}

type MessageKick struct {
	MessageType string `json:"message_type"`
	KickReason  string `json:"kick_reason"`
}

func checkMessageType(data map[string]interface{}, expectedMessageType string) error {
	messageType, err := readString(data, "message_type")
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
	readMessage.nickname, err = readString(data, "nickname")
	if err != nil {
		return readMessage, err
	}

	// Check nickname
	r, _ := regexp.Compile(`\A\S{1,10}\z`)
	if !r.MatchString(readMessage.nickname) {
		return readMessage, fmt.Errorf("Invalid nickname")
	}

	// Read role
	readMessage.role, err = readString(data, "role")
	if err != nil {
		return readMessage, err
	}

	// Check role
	switch readMessage.role {
	case "player",
		"visualization",
		"game_logic":
		return readMessage, nil
	default:
		return readMessage, fmt.Errorf("Invalid role '%v'",
			readMessage.role)
	}
}
