package test

import (
	"fmt"
	"github.com/netorcai/netorcai"
	"github.com/netorcai/netorcai/client/go"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

type ClientGameStartsCheckFunc func(*testing.T, map[string]interface{}, int,
	int, int, float64, float64, bool) int
type ClientTurnCheckFunc func(*testing.T, map[string]interface{}, int, int, int,
	bool) int
type ClientGameEndsCheckFunc func(*testing.T, map[string]interface{})
type GLCheckDoTurnFunc func(*testing.T, map[string]interface{},
	int, int, int) []interface{}
type ClientTurnAckFunc func(int, int) string
type GLDoInitAckFunc func(int, int, int) string
type GLDoTurnAckFunc func(int, []interface{}) string

func DefaultHelloClientCheckGameStarts(t *testing.T,
	msg map[string]interface{}, nbPlayers, nbSpecialPlayers, nbTurnsGL int,
	msBeforeFirstTurn, msBetweenTurns float64, isPlayer bool) int {
	playerID := checkGameStarts(t, msg, nbPlayers, nbSpecialPlayers, nbTurnsGL,
		msBeforeFirstTurn, msBetweenTurns, isPlayer)
	return playerID
}

func DefaultHelloClientCheckTurn(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedNbSpecialPlayers, expectedTurnNumber int, isPlayer bool) int {
	return checkTurn(t, msg, expectedNbPlayers, expectedNbSpecialPlayers, expectedTurnNumber, isPlayer)
}

func DefaultHelloClientCheckGameEnds(t *testing.T,
	msg map[string]interface{}) {
	checkGameEnds(t, msg)
}

func DefaultHelloGLCheckDoTurn(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedNbSpecialPlayers, expectedTurnNumber int) []interface{} {
	actions := checkDoTurn(t, msg, expectedNbPlayers, expectedNbSpecialPlayers, expectedTurnNumber)
	return actions
}

func DefaultHelloClientTurnAck(turn, playerID int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": %v,
		"actions": []}`, turn)
}

func DefaultHelloGLDoInitAck(nbPlayers, nbSpecialPlayers, nbTurns int) string {
	return `{"message_type":"DO_INIT_ACK", "initial_game_state":{"all_clients":{}}}`
}

func DefaultHelloGlDoTurnAck(turn int, actions []interface{}) string {
	return `{"message_type":"DO_TURN_ACK",
		"winner_player_id":-1,
		"game_state":{"all_clients":{}}}`
}

func helloGameLogic(t *testing.T, glClient *client.Client,
	nbPlayers, nbSpecialPlayers, nbTurnsNetorcai, nbTurns int,
	checkDoTurnFunc GLCheckDoTurnFunc,
	doInitAckFunc GLDoInitAckFunc, doTurnAckFunc GLDoTurnAckFunc,
	kickReasonMatcher *regexp.Regexp) {
	// Wait DO_INIT
	msg, err := waitReadMessage(glClient, 1000)
	assert.NoError(t, err, "Could not read GLClient message (DO_INIT)")
	checkDoInit(t, msg, nbPlayers, nbSpecialPlayers, nbTurnsNetorcai)

	// Send DO_INIT_ACK
	data := doInitAckFunc(nbPlayers, nbSpecialPlayers, nbTurnsNetorcai)
	err = glClient.SendString(data)
	assert.NoError(t, err, "GLClient could not send DO_INIT_ACK")

	// Wait for DO_TURN
	for turn := 0; turn < nbTurns; turn++ {
		msg, err := waitReadMessage(glClient, 1000)
		assert.NoError(t, err, "Could not read GLClient message (DO_TURN) "+
			"%v/%v", turn, nbTurns)
		actions := checkDoTurnFunc(t, msg, nbPlayers, nbSpecialPlayers, turn-1)

		// Send DO_TURN_ACK
		data = doTurnAckFunc(turn, actions)
		err = glClient.SendString(data)
		assert.NoError(t, err, "GLClient could not send DO_TURN_ACK")
	}

	msg, err = waitReadMessage(glClient, 1000)
	assert.NoError(t, err, "Could not read GLClient message (KICK)")
	checkKick(t, msg, "GameLogic", kickReasonMatcher)

	// Close socket
	glClient.Disconnect()
}

func helloClient(t *testing.T, client *client.Client, clientName string,
	nbPlayers, nbSpecialPlayers, nbTurnsGL, nbTurnsClient, turnsToSkip int,
	msBeforeFirstTurn, msBetweenTurns float64,
	isPlayer, allowTurnSkip, shouldTurnAckBeValid, shouldDoInitAckBeValid bool,
	checkGameStartsFunc ClientGameStartsCheckFunc,
	checkTurnFunc ClientTurnCheckFunc,
	checkGameEndsFunc ClientGameEndsCheckFunc,
	turnAckFunc ClientTurnAckFunc, kickReasonMatcher *regexp.Regexp) {

	// Close socket when leaving this function
	defer client.Disconnect()

	if shouldDoInitAckBeValid {
		// Wait GAME_STARTS
		msg, err := waitReadMessage(client, 1000)
		assert.NoError(t, err, "%v could not read message (GAME_STARTS)", clientName)
		playerID := checkGameStartsFunc(t, msg, nbPlayers, nbSpecialPlayers, nbTurnsGL,
			msBeforeFirstTurn, msBetweenTurns, isPlayer)

		if !allowTurnSkip {
			for turn := 0; turn < nbTurnsClient-1; turn += 1 + turnsToSkip {
				// Wait TURN
				msg, err := waitReadMessage(client, 1000)
				assert.NoError(t, err, "%v could not read message (TURN) %v/%v",
					clientName, turn, nbTurnsClient)
				turnReceived := checkTurnFunc(t, msg, nbPlayers, nbSpecialPlayers, turn, isPlayer)

				// Send TURN_ACK
				data := turnAckFunc(turnReceived, playerID)
				err = client.SendString(data)
				assert.NoError(t, err, "%s cannot send TURN_ACK", clientName)
			}

			if shouldTurnAckBeValid {
				// Wait GAME_ENDS
				msg, err = waitReadMessage(client, 1000)
				assert.NoError(t, err, "%v could not read message (GAME_ENDS)", clientName)
				checkGameEndsFunc(t, msg)
			}
		} else {
		TurnLoop:
			for turn := 0; turn < nbTurnsClient; turn += 1 {
				msg, err := waitReadMessage(client, 1000)
				assert.NoError(t, err, "%v could not read message (TURN or GAME_ENDS) %v/%v"+
					clientName, turn, nbTurnsClient)
				turnReceived := checkTurnPotentialTurnsSkipped(t, msg, nbPlayers, nbSpecialPlayers, turn, isPlayer)

				messageType, _ := netorcai.ReadString(msg, "message_type")

				switch messageType {
				case "TURN":
					// Send TURN_ACK
					data := turnAckFunc(turnReceived, playerID)
					err = client.SendString(data)
					assert.NoError(t, err, "%s cannot send TURN_ACK", clientName)
				case "GAME_ENDS":
					break TurnLoop
				case "KICK":
					return
				}
			}
		}
	}

	// Wait Kick
	msg, err := waitReadMessage(client, 2000)
	assert.NoError(t, err, "Could not read %v message (KICK)", clientName)
	checkKick(t, msg, clientName, kickReasonMatcher)
}
