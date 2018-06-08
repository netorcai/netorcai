package test

import (
	"fmt"
	"github.com/mpoquet/netorcai"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

// Initial game state
func doInitAckFlattened(nbPlayers, nbTurns int) string {
	return `{
	  "message_type": "DO_INIT_ACK",
	  "initial_game_state": {
	    "all_clients": {
	      "string": "hello",
	      "integer": 42,
	      "float": 0.5,
	      "object": {},
	      "array": [],
	      "bool": true
	    }
	  }
	}`
}

func doInitAckNastyNested(nbPlayers, nbTurns int) string {
	return `{
	  "message_type": "DO_INIT_ACK",
	  "initial_game_state": {
	    "all_clients": {
	      "string": "hello",
	      "integer": 42,
	      "float": 0.5,
	      "object": {},
	      "array": [],
	      "bool": true,
	      "nested_object": {
	        "string": "hello",
	        "integer": 42,
	        "float": 0.5,
	        "object": {},
	        "array": [],
	        "bool": true,
	        "nested_array": [
	          {
	            "string": "hello",
	            "integer": 42,
	            "float": 0.5,
	            "object": {},
	            "array": [],
	            "bool": true
	          }
	        ]
	      },
	      "nested_array": [
	        {
	          "string": "hello",
	          "integer": 42,
	          "float": 0.5,
	          "object": {},
	          "array": [],
	          "bool": true,
	          "nested_object": {
	            "string": "hello",
	            "integer": 42,
	            "float": 0.5,
	            "object": {},
	            "array": [],
	            "bool": true
	          }
	        }
	      ]
	    }
	  }
	}`
}

func subCheckFlattenedObject(t *testing.T, object map[string]interface{}) {
	s, err := netorcai.ReadString(object, "string")
	assert.NoError(t, err, "Cannot read 'string' field in obj")
	assert.Equal(t, "hello", s, "Unexpected value for 'string' field in obj")

	i, err := netorcai.ReadInt(object, "integer")
	assert.NoError(t, err, "Cannot read 'integer' field in obj")
	assert.Equal(t, 42, i, "Unexpected value for 'integer' field in obj")

	f, err := readFloat(object, "float")
	assert.NoError(t, err, "Cannot read 'float' field in obj")
	assert.Equal(t, 0.5, f, "Unexpected value for 'float' field in obj")

	o, err := netorcai.ReadObject(object, "object")
	assert.NoError(t, err, "Cannot read 'object' field in obj")
	assert.Equal(t, 0, len(o), "Unexpected length for 'object' field in obj")

	a, err := netorcai.ReadArray(object, "array")
	assert.NoError(t, err, "Cannot read 'array' field in obj")
	assert.Equal(t, 0, len(a), "Unexpected length for 'array' field in obj")

	b, err := readBool(object, "bool")
	assert.NoError(t, err, "Cannot read 'bool' field in obj")
	assert.Equal(t, true, b, "Unexpected value for 'bool' field in obj")
}

func checkGameStartsFlattened(t *testing.T,
	msg map[string]interface{}, nbPlayers, nbTurnsGL int,
	msBeforeFirstTurn, msBetweenTurns float64, isPlayer bool) int {
	playerID := checkGameStarts(t, msg, nbPlayers, nbTurnsGL,
		msBeforeFirstTurn, msBetweenTurns, isPlayer)

	initialGS, err := netorcai.ReadObject(msg, "initial_game_state")
	assert.NoError(t, err, "Cannot read 'initial_game_state' in msg")
	subCheckFlattenedObject(t, initialGS)

	return playerID
}

func checkGameStartsNastyNested(t *testing.T,
	msg map[string]interface{}, nbPlayers, nbTurnsGL int,
	msBeforeFirstTurn, msBetweenTurns float64, isPlayer bool) int {
	playerID := checkGameStarts(t, msg, nbPlayers, nbTurnsGL,
		msBeforeFirstTurn, msBetweenTurns, isPlayer)

	initialGS, err := netorcai.ReadObject(msg, "initial_game_state")
	assert.NoError(t, err, "Cannot read 'initial_game_state' in msg")
	subCheckFlattenedObject(t, initialGS)

	nestedObject1, err := netorcai.ReadObject(initialGS, "nested_object")
	assert.NoError(t, err, "Cannot read 'nested_object' in msg")
	subCheckFlattenedObject(t, nestedObject1)

	nestedArray2, err := netorcai.ReadArray(nestedObject1, "nested_array")
	assert.NoError(t, err, "Cannot read 'nested_array' field in obj")
	assert.Equal(t, 1, len(nestedArray2),
		"Unexpected length for 'nested_array' field in obj")
	nestedArray2FirstElement := nestedArray2[0].(map[string]interface{})
	subCheckFlattenedObject(t, nestedArray2FirstElement)

	nestedArray1, err := netorcai.ReadArray(initialGS, "nested_array")
	assert.NoError(t, err, "Cannot read 'nested_array' field in obj")
	assert.Equal(t, 1, len(nestedArray1),
		"Unexpected length for 'nested_array' field in obj")
	nestedArray1FirstElement := nestedArray1[0].(map[string]interface{})
	subCheckFlattenedObject(t, nestedArray1FirstElement)

	nestedObject2, err := netorcai.ReadObject(nestedArray1FirstElement,
		"nested_object")
	assert.NoError(t, err, "Cannot read 'nested_object' in msg")
	subCheckFlattenedObject(t, nestedObject2)

	return playerID
}

func TestForwardInitialGameStateFlattened(t *testing.T) {
	subtestHelloGlActiveClients(t, 4, 1,
		3, 3, 3, 3,
		checkGameStartsFlattened, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		doInitAckFlattened, DefaultHelloGlDoTurnAck,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

func TestForwardInitialGameStateNastyNested(t *testing.T) {
	subtestHelloGlActiveClients(t, 4, 1,
		3, 3, 3, 3,
		checkGameStartsNastyNested, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		doInitAckNastyNested, DefaultHelloGlDoTurnAck,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

// Game state
func doTurnAckFlattened(turn int, actions []interface{}) string {
	return `{
	  "message_type": "DO_TURN_ACK",
	  "winner_player_id":-1,
	  "game_state": {
	    "all_clients": {
	      "string": "hello",
	      "integer": 42,
	      "float": 0.5,
	      "object": {},
	      "array": [],
	      "bool": true
	    }
	  }
	}`
}

func doTurnAckNastyNested(turn int, actions []interface{}) string {
	return `{
	  "message_type": "DO_TURN_ACK",
	  "winner_player_id":-1,
	  "game_state": {
	    "all_clients": {
	      "string": "hello",
	      "integer": 42,
	      "float": 0.5,
	      "object": {},
	      "array": [],
	      "bool": true,
	      "nested_object": {
	        "string": "hello",
	        "integer": 42,
	        "float": 0.5,
	        "object": {},
	        "array": [],
	        "bool": true,
	        "nested_array": [
	          {
	            "string": "hello",
	            "integer": 42,
	            "float": 0.5,
	            "object": {},
	            "array": [],
	            "bool": true
	          }
	        ]
	      },
	      "nested_array": [
	        {
	          "string": "hello",
	          "integer": 42,
	          "float": 0.5,
	          "object": {},
	          "array": [],
	          "bool": true,
	          "nested_object": {
	            "string": "hello",
	            "integer": 42,
	            "float": 0.5,
	            "object": {},
	            "array": [],
	            "bool": true
	          }
	        }
	      ]
	    }
	  }
	}`
}

func checkTurnFlattened(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int, isPlayer bool) {
	checkTurn(t, msg, expectedNbPlayers, expectedTurnNumber, isPlayer)

	gs, err := netorcai.ReadObject(msg, "game_state")
	assert.NoError(t, err, "Cannot read 'game_state' in msg")
	subCheckFlattenedObject(t, gs)
}

func checkTurnNastyNested(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int, isPlayer bool) {
	checkTurn(t, msg, expectedNbPlayers, expectedTurnNumber, isPlayer)

	gs, err := netorcai.ReadObject(msg, "game_state")
	assert.NoError(t, err, "Cannot read 'game_state' in msg")
	subCheckFlattenedObject(t, gs)

	nestedObject1, err := netorcai.ReadObject(gs, "nested_object")
	assert.NoError(t, err, "Cannot read 'nested_object' in msg")
	subCheckFlattenedObject(t, nestedObject1)

	nestedArray2, err := netorcai.ReadArray(nestedObject1, "nested_array")
	assert.NoError(t, err, "Cannot read 'nested_array' field in obj")
	assert.Equal(t, 1, len(nestedArray2),
		"Unexpected length for 'nested_array' field in obj")
	nestedArray2FirstElement := nestedArray2[0].(map[string]interface{})
	subCheckFlattenedObject(t, nestedArray2FirstElement)

	nestedArray1, err := netorcai.ReadArray(gs, "nested_array")
	assert.NoError(t, err, "Cannot read 'nested_array' field in obj")
	assert.Equal(t, 1, len(nestedArray1),
		"Unexpected length for 'nested_array' field in obj")
	nestedArray1FirstElement := nestedArray1[0].(map[string]interface{})
	subCheckFlattenedObject(t, nestedArray1FirstElement)

	nestedObject2, err := netorcai.ReadObject(nestedArray1FirstElement,
		"nested_object")
	assert.NoError(t, err, "Cannot read 'nested_object' in msg")
	subCheckFlattenedObject(t, nestedObject2)
}

func TestForwardGameStateFlattened(t *testing.T) {
	subtestHelloGlActiveClients(t, 4, 1,
		3, 3, 3, 3,
		DefaultHelloClientCheckGameStarts, checkTurnFlattened,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, doTurnAckFlattened,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

func TestForwardGameStateNastyNested(t *testing.T) {
	subtestHelloGlActiveClients(t, 4, 1,
		3, 3, 3, 3,
		DefaultHelloClientCheckGameStarts, checkTurnNastyNested,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, doTurnAckNastyNested,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

// Actions
func turnAckFlattened(turn, playerID int) string {
	actions := "[]"

	if playerID >= 0 {
		actions = fmt.Sprintf(`[{
		  "whoami": %v,
		  "sent_at_turn": %v,
		  "string": "hello",
		  "integer": 42,
		  "float": 0.5,
		  "object": {},
		  "array": [],
		  "bool": true
		}]`, playerID, turn)
	}
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": %v,
		"actions": %v}`, turn, actions)
}

func turnAckNastyNested(turn, playerID int) string {
	actions := "[]"

	if playerID >= 0 {
		actions = fmt.Sprintf(`[{
		  "whoami": %v,
		  "sent_at_turn": %v,
		  "string": "hello",
		  "integer": 42,
		  "float": 0.5,
		  "object": {},
		  "array": [],
		  "bool": true,
		  "nested_object": {
		    "string": "hello",
		    "integer": 42,
		    "float": 0.5,
		    "object": {},
		    "array": [],
		    "bool": true,
		    "nested_array": [
		      {
		        "string": "hello",
		        "integer": 42,
		        "float": 0.5,
		        "object": {},
		        "array": [],
		        "bool": true
		      }
		    ]
		  },
		  "nested_array": [
		    {
		      "string": "hello",
		      "integer": 42,
		      "float": 0.5,
		      "object": {},
		      "array": [],
		      "bool": true,
		      "nested_object": {
		        "string": "hello",
		        "integer": 42,
		        "float": 0.5,
		        "object": {},
		        "array": [],
		        "bool": true
		      }
		    }
		  ]
		}]`, playerID, turn)
	}
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": %v,
		"actions": %v}`, turn, actions)
}

func checkDoTurnFlattened(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int) []interface{} {
	pActions := checkDoTurn(t, msg, expectedNbPlayers, expectedTurnNumber)

	if expectedTurnNumber >= 0 {
		assert.Equal(t, expectedNbPlayers, len(pActions),
			"Unexpected number of player actions received")

		for _, pAction := range pActions {
			pAsObj := pAction.(map[string]interface{})

			actions, err := netorcai.ReadArray(pAsObj, "actions")
			assert.NoError(t, err, "Cannot read 'actions' field in obj")
			assert.Equal(t, 1, len(actions), "Unexpected 'actions' length")

			firstElement := actions[0].(map[string]interface{})
			subCheckFlattenedObject(t, firstElement)

			// Check player_id consistency
			playerID, err := netorcai.ReadInt(pAsObj, "player_id")
			assert.NoError(t, err, "Cannot read 'player_id' field in obj")

			whoami, err := netorcai.ReadInt(firstElement, "whoami")
			assert.NoError(t, err, "Cannot read 'whoami' field in obj")
			assert.Equal(t, playerID, whoami, "Unexpected 'whoami' value")

			// Check turn_number consistency
			turnNumber, err := netorcai.ReadInt(pAsObj, "turn_number")
			assert.NoError(t, err, "Cannot read 'turn_number' field in obj")

			sentAtTurn, err := netorcai.ReadInt(firstElement, "sent_at_turn")
			assert.NoError(t, err, "Cannot read 'sent_at_turn' field in obj")
			assert.Equal(t, turnNumber, sentAtTurn,
				"Unexpected 'sent_at_turn' value")
		}
	}

	return pActions
}

func checkDoTurnNastyNested(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedTurnNumber int) []interface{} {
	pActions := checkDoTurn(t, msg, expectedNbPlayers, expectedTurnNumber)

	if expectedTurnNumber >= 0 {
		assert.Equal(t, expectedNbPlayers, len(pActions),
			"Unexpected number of player actions received")

		for _, pAction := range pActions {
			pAsObj := pAction.(map[string]interface{})

			actions, err := netorcai.ReadArray(pAsObj, "actions")
			assert.NoError(t, err, "Cannot read 'actions' field in obj")
			assert.Equal(t, 1, len(actions), "Unexpected 'actions' length")

			firstElement := actions[0].(map[string]interface{})
			subCheckFlattenedObject(t, firstElement)

			// Check player_id consistency
			playerID, err := netorcai.ReadInt(pAsObj, "player_id")
			assert.NoError(t, err, "Cannot read 'player_id' field in obj")

			whoami, err := netorcai.ReadInt(firstElement, "whoami")
			assert.NoError(t, err, "Cannot read 'whoami' field in obj")
			assert.Equal(t, playerID, whoami, "Unexpected 'whoami' value")

			// Check turn_number consistency
			turnNumber, err := netorcai.ReadInt(pAsObj, "turn_number")
			assert.NoError(t, err, "Cannot read 'turn_number' field in obj")

			sentAtTurn, err := netorcai.ReadInt(firstElement, "sent_at_turn")
			assert.NoError(t, err, "Cannot read 'sent_at_turn' field in obj")
			assert.Equal(t, turnNumber, sentAtTurn,
				"Unexpected 'sent_at_turn' value")

			// Nesting check
			nestedObject1, err := netorcai.ReadObject(firstElement,
				"nested_object")
			assert.NoError(t, err, "Cannot read 'nested_object' in msg")
			subCheckFlattenedObject(t, nestedObject1)

			nestedArray2, err := netorcai.ReadArray(nestedObject1, "nested_array")
			assert.NoError(t, err, "Cannot read 'nested_array' field in obj")
			assert.Equal(t, 1, len(nestedArray2),
				"Unexpected length for 'nested_array' field in obj")
			nestedArray2FirstElement := nestedArray2[0].(map[string]interface{})
			subCheckFlattenedObject(t, nestedArray2FirstElement)

			nestedArray1, err := netorcai.ReadArray(firstElement,
				"nested_array")
			assert.NoError(t, err, "Cannot read 'nested_array' field in obj")
			assert.Equal(t, 1, len(nestedArray1),
				"Unexpected length for 'nested_array' field in obj")
			nestedArray1FirstElement := nestedArray1[0].(map[string]interface{})
			subCheckFlattenedObject(t, nestedArray1FirstElement)

			nestedObject2, err := netorcai.ReadObject(nestedArray1FirstElement,
				"nested_object")
			assert.NoError(t, err, "Cannot read 'nested_object' in msg")
			subCheckFlattenedObject(t, nestedObject2)
		}
	}

	return pActions
}

func TestForwardActionsFlattened(t *testing.T) {
	subtestHelloGlActiveClients(t, 4, 1,
		3, 3, 3, 3,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, checkDoTurnFlattened,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckFlattened, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

func TestForwardActionsNastyNested(t *testing.T) {
	subtestHelloGlActiveClients(t, 4, 1,
		3, 3, 3, 3,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, checkDoTurnNastyNested,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckNastyNested, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}
