package test

import (
	"fmt"
	"github.com/mpoquet/netorcai"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

// Play once every two turns
func turnAckSkipOneTurnOverTwo(turn, playerID int) string {
	time.Sleep(time.Duration(1100) * time.Millisecond)
	msg := fmt.Sprintf(`{"message_type": "TURN_ACK",
        "turn_number": %v,
        "actions": []}`, turn)
	return msg
}

func checkTurnSkipOneTurnOverTwo(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int, isPlayer bool) int {

	turn, err := netorcai.ReadInt(msg, "turn_number")
	assert.NoError(t, err, "Cannot read 'turn_number'")
	assert.Equal(t, 0, turn%2, "Unexpected turn_number parity")

	return turn
}

func checkDoTurnSkipOneTurnOverTwo(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int) []interface{} {

	actions, err := netorcai.ReadArray(msg, "player_actions")
	assert.NoError(t, err, "Cannot read 'player_actions'")

	expectedPActionsLength := 0
	if (expectedTurnNumber != 0) && (expectedTurnNumber%2 == 0) {
		expectedPActionsLength = 1
	}

	assert.Equal(t, expectedPActionsLength, len(actions),
		"Unexpected array length for 'player_actions'")
	return actions
}

func TestSkipOneTurnOverTwo(t *testing.T) {
	subtestHelloGlActiveClients(t, 1, 0,
		7, 7, 7, 7,
		1, 0,
		DefaultHelloClientCheckGameStarts, checkTurnSkipOneTurnOverTwo,
		DefaultHelloClientCheckGameEnds, checkDoTurnSkipOneTurnOverTwo,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckSkipOneTurnOverTwo, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

// Skip first turn
func turnAckSkipFirstTurn(turn, playerID int) string {
	if turn == 0 {
		time.Sleep(time.Duration(600) * time.Millisecond)
	}
	msg := fmt.Sprintf(`{"message_type": "TURN_ACK",
        "turn_number": %v,
        "actions": []}`, turn)
	return msg
}

func checkTurnSkipFirstTurn(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int, isPlayer bool) int {

	turn, err := netorcai.ReadInt(msg, "turn_number")
	assert.NoError(t, err, "Cannot read 'turn_number'")

	return turn
}

func checkDoTurnSkipFirstTurn(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int) []interface{} {

	actions, err := netorcai.ReadArray(msg, "player_actions")
	assert.NoError(t, err, "Cannot read 'player_actions'")

	expectedPActionsLength := 1
	if expectedTurnNumber < 1 {
		expectedPActionsLength = 0
	}

	assert.Equal(t, expectedPActionsLength, len(actions),
		"Unexpected array length for 'player_actions'. turn=%v",
		expectedTurnNumber)
	return actions
}

func TestSkipFirstTurn(t *testing.T) {
	subtestHelloGlActiveClients(t, 1, 0,
		4, 4, 4, 4,
		0, 0,
		DefaultHelloClientCheckGameStarts, checkTurnSkipFirstTurn,
		DefaultHelloClientCheckGameEnds, checkDoTurnSkipFirstTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckSkipFirstTurn, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}
