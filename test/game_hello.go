package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func hello_game_logic(t *testing.T, glClient *Client,
	nbPlayers, nbTurns int, onexit chan int) {
	// Wait DO_INIT
	msg, err := waitReadMessage(glClient, 1000)
	assert.NoError(t, err, "Could not read GLClient message (DO_INIT)")
	checkDoInit(t, msg, nbPlayers, nbTurns)

	// Send DO_INIT_ACK
	data := `{"message_type":"DO_INIT_ACK", "initial_game_state":{"all_clients":{}}}`
	err = glClient.SendString(data)
	assert.NoError(t, err, "GLClient could not send DO_INIT_ACK")

	// Wait for DO_TURN
	for turn := 0; turn < nbTurns; turn++ {
		msg, err := waitReadMessage(glClient, 1000)
		assert.NoError(t, err, "Could not read GLClient message (DO_TURN) "+
			"%v/%v", turn, nbTurns)
		checkDoTurn(t, msg, nbPlayers, turn)

		// Send DO_TURN_ACK
		data = `{"message_type":"DO_TURN_ACK",
			"winner_player_id":-1,
			"game_state":{"all_clients":{}}}`
		err = glClient.SendString(data)
		assert.NoError(t, err, "GLClient could not send DO_TURN_ACK")
	}

	msg, err = waitReadMessage(glClient, 1000)
	assert.NoError(t, err, "Could not read GLClient message (KICK)")
	checkKick(t, msg, regexp.MustCompile(".*"))

	// Close socket
	glClient.Disconnect()
	onexit <- 1
}
