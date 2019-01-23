package test

import (
	"fmt"
	"github.com/netorcai/netorcai"
	"github.com/stretchr/testify/assert"
	"regexp"
	"sort"
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
	subtestHelloGlActiveClients(t, nil, 1, 0,
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
	subtestHelloGlActiveClients(t, nil, 1, 0,
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

// Skip one turn, but not the same depending on the client
func turnAckSkipOneTurnMultiClient(turn, playerID int) string {
	if turn == playerID {
		time.Sleep(time.Duration(500) * time.Millisecond)
	}
	time.Sleep(time.Duration((playerID+1)*50) * time.Millisecond)
	msg := fmt.Sprintf(`{"message_type": "TURN_ACK",
        "turn_number": %v,
        "actions": [{"whoami": %v, "sent_at_turn": %v}]}`,
		turn, playerID, turn)
	return msg
}

func checkTurnSkipOneTurnMultiClient(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int, isPlayer bool) int {

	turn, err := netorcai.ReadInt(msg, "turn_number")
	assert.NoError(t, err, "Cannot read 'turn_number'")

	return turn
}

func subCheckPlayerActionsObject(t *testing.T, obj map[string]interface{},
	expectedPlayerID, expectedTurn int) {
	playerID, err := netorcai.ReadInt(obj, "player_id")
	assert.NoError(t, err, "Cannot read 'player_id'")
	assert.Equal(t, expectedPlayerID, playerID, "Unexpected 'player_id' value")

	turnNumber, err := netorcai.ReadInt(obj, "turn_number")
	assert.NoError(t, err, "Cannot read 'turn_number'")
	assert.Equal(t, expectedTurn, turnNumber,
		"Unexpected 'turn_number' value")

	actions, err := netorcai.ReadArray(obj, "actions")
	assert.NoError(t, err, "Cannot read 'actions'")
	assert.Len(t, actions, 1, "Unexpected 'actions' length")

	firstElement := actions[0].(map[string]interface{})

	whoami, err := netorcai.ReadInt(firstElement, "whoami")
	assert.NoError(t, err, "Cannot read 'whoami'")
	assert.Equal(t, expectedPlayerID, whoami, "Unexpected 'whoami' value")

	sendAtTurn, err := netorcai.ReadInt(firstElement, "sent_at_turn")
	assert.NoError(t, err, "Cannot read 'sent_at_turn'")
	assert.Equal(t, expectedTurn, sendAtTurn,
		"Unexpected 'sent_at_turn' value")
}

type ByPlayerID []interface{}

func (a ByPlayerID) Len() int      { return len(a) }
func (a ByPlayerID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPlayerID) Less(i, j int) bool {
	return a[i].(map[string]interface{})["player_id"].(float64) <
		a[j].(map[string]interface{})["player_id"].(float64)
}

func checkDoTurnSkipOneTurnMultiClient(t *testing.T,
	msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int) []interface{} {

	actions, err := netorcai.ReadArray(msg, "player_actions")
	assert.NoError(t, err, "Cannot read 'player_actions'")

	sort.Sort(ByPlayerID(actions))

	if expectedTurnNumber == -1 {
		assert.Equal(t, 0, len(actions),
			"Unexpected array length for 'player_actions' at turn=%v",
			expectedTurnNumber)
	} else if expectedTurnNumber == 0 {
		assert.Equal(t, 3, len(actions),
			"Unexpected array length for 'player_actions' at turn=%v",
			expectedTurnNumber)
		subCheckPlayerActionsObject(t, actions[0].(map[string]interface{}),
			1, 0)
		subCheckPlayerActionsObject(t, actions[1].(map[string]interface{}),
			2, 0)
		subCheckPlayerActionsObject(t, actions[2].(map[string]interface{}),
			3, 0)
	} else if expectedTurnNumber == 1 {
		assert.Equal(t, 3, len(actions),
			"Unexpected array length for 'player_actions' at turn=%v",
			expectedTurnNumber)
		subCheckPlayerActionsObject(t, actions[0].(map[string]interface{}),
			0, 1)
		subCheckPlayerActionsObject(t, actions[1].(map[string]interface{}),
			2, 1)
		subCheckPlayerActionsObject(t, actions[2].(map[string]interface{}),
			3, 1)
	} else if expectedTurnNumber == 2 {
		assert.Equal(t, 3, len(actions),
			"Unexpected array length for 'player_actions' at turn=%v",
			expectedTurnNumber)
		subCheckPlayerActionsObject(t, actions[0].(map[string]interface{}),
			0, 2)
		subCheckPlayerActionsObject(t, actions[1].(map[string]interface{}),
			1, 2)
		subCheckPlayerActionsObject(t, actions[2].(map[string]interface{}),
			3, 2)
	} else if expectedTurnNumber == 3 {
		assert.Equal(t, 3, len(actions),
			"Unexpected array length for 'player_actions' at turn=%v",
			expectedTurnNumber)
		subCheckPlayerActionsObject(t, actions[0].(map[string]interface{}),
			0, 3)
		subCheckPlayerActionsObject(t, actions[1].(map[string]interface{}),
			1, 3)
		subCheckPlayerActionsObject(t, actions[2].(map[string]interface{}),
			2, 3)
	} else {
		assert.Equal(t, 4, len(actions),
			"Unexpected array length for 'player_actions' at turn=%v",
			expectedTurnNumber)
		subCheckPlayerActionsObject(t, actions[0].(map[string]interface{}),
			0, expectedTurnNumber)
		subCheckPlayerActionsObject(t, actions[1].(map[string]interface{}),
			1, expectedTurnNumber)
		subCheckPlayerActionsObject(t, actions[2].(map[string]interface{}),
			2, expectedTurnNumber)
		subCheckPlayerActionsObject(t, actions[3].(map[string]interface{}),
			3, expectedTurnNumber)
	}

	return actions
}

func TestSkipOneTurnMultiClient(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 4, 0,
		7, 7, 7, 7,
		0, 0,
		DefaultHelloClientCheckGameStarts, checkTurnSkipOneTurnMultiClient,
		DefaultHelloClientCheckGameEnds, checkDoTurnSkipOneTurnMultiClient,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckSkipOneTurnMultiClient, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}
