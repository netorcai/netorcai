package test

import (
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
	msBeforeFirstTurn, msBetweenTurns float64, isPlayer bool) {
	checkGameStarts(t, msg, nbPlayers, nbTurnsGL, msBeforeFirstTurn,
		msBetweenTurns, isPlayer)

	initialGS, err := netorcai.ReadObject(msg, "initial_game_state")
	assert.NoError(t, err, "Cannot read 'initial_game_state' in msg")
	subCheckFlattenedObject(t, initialGS)
}

func checkGameStartsNastyNested(t *testing.T,
	msg map[string]interface{}, nbPlayers, nbTurnsGL int,
	msBeforeFirstTurn, msBetweenTurns float64, isPlayer bool) {
	checkGameStarts(t, msg, nbPlayers, nbTurnsGL, msBeforeFirstTurn,
		msBetweenTurns, isPlayer)

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
}

func TestForwardInitialGameStateFlattened(t *testing.T) {
	subtestHelloGlActiveClients(t, 4, 1,
		3, 3, 3, 3,
		checkGameStartsFlattened, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds,
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
		DefaultHelloClientCheckGameEnds,
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
		DefaultHelloClientCheckGameEnds,
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
		DefaultHelloClientCheckGameEnds,
		DefaultHelloGLDoInitAck, doTurnAckNastyNested,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}
