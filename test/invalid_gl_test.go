package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestInvalidGlMessageBeforeStart(t *testing.T) {
	proc, _, playerClients, visuClients, glClients := runNetorcaiAndAllClients(
		t, []string{}, 1000)
	defer killallNetorcaiSIGKILL()

	glClients[0].SendString(`{}`)
	checkAllKicked(t, glClients, regexp.MustCompile(
		`Received a game logic message but the game has not started`), 1000)

	_, err := waitOutputTimeout(regexp.MustCompile(`Game logic failed`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err,
		"Cannot read `Game logic failed` in netorcai output")

	_, expRetCode := handleCoverage(t, 1)
	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")

	checkAllKicked(t, playerClients, regexp.MustCompile(`netorcai abort`),
		1000)
	checkAllKicked(t, visuClients, regexp.MustCompile(`netorcai abort`), 1000)
}

func TestInvalidGlNoDoInitAck(t *testing.T) {
	proc, _, playerClients, visuClients, glClients := runNetorcaiAndAllClients(
		t, []string{}, 1000)
	defer killallNetorcaiSIGKILL()

	go func(glClient *Client) {
		msg, err := waitReadMessage(glClient, 1000)
		assert.NoError(t, err, "Could not read GLClient message (DO_INIT)")
		checkDoInit(t, msg, 4, 100)

		// Do not send DO_INIT_ACK on purpose
		msg, err = waitReadMessage(glClient, 4000)
		checkKick(t, msg, regexp.MustCompile(`Did not receive DO_INIT_ACK after 3 seconds`))
	}(glClients[0])

	proc.inputControl <- `start`
	_, err := waitOutputTimeout(regexp.MustCompile(`Game logic failed`),
		proc.outputControl, 4000, false)
	assert.NoError(t, err,
		"Cannot read `Game logic failed` in netorcai output")

	_, expRetCode := handleCoverage(t, 1)
	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")

	checkAllKicked(t, playerClients, regexp.MustCompile(`netorcai abort`),
		1000)
	checkAllKicked(t, visuClients, regexp.MustCompile(`netorcai abort`), 1000)
}

func subtestBadDoInitAck(t *testing.T, doInitAckMessage string,
	kickReasonMatcher *regexp.Regexp) {
	proc, _, playerClients, visuClients, glClients := runNetorcaiAndAllClients(
		t, []string{}, 1000)
	defer killallNetorcaiSIGKILL()

	go func(glClient *Client) {
		msg, err := waitReadMessage(glClient, 1000)
		assert.NoError(t, err, "Could not read GLClient message (DO_INIT)")
		checkDoInit(t, msg, 4, 100)

		glClient.SendString(doInitAckMessage)

		// Should be kicked
		msg, err = waitReadMessage(glClient, 1000)
		checkKick(t, msg, kickReasonMatcher)
	}(glClients[0])

	proc.inputControl <- `start`
	_, err := waitOutputTimeout(regexp.MustCompile(`Game logic failed`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err,
		"Cannot read `Game logic failed` in netorcai output")

	_, expRetCode := handleCoverage(t, 1)
	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")

	checkAllKicked(t, playerClients, regexp.MustCompile(`netorcai abort`),
		1000)
	checkAllKicked(t, visuClients, regexp.MustCompile(`netorcai abort`), 1000)
}

func TestInvalidGlBadDoInitAckNoMessageType(t *testing.T) {
	subtestBadDoInitAck(t, `{}`,
		regexp.MustCompile(`Field 'message_type' is missing`))
}

func TestInvalidGlBadDoInitAckBadMessageType(t *testing.T) {
	subtestBadDoInitAck(t, `{"message_type":"LOGIN"}`,
		regexp.MustCompile(`Received 'LOGIN' message type, `+
			`while DO_INIT_ACK was expected`))
}

func TestInvalidGlBadDoInitAckNoGameState(t *testing.T) {
	subtestBadDoInitAck(t, `{"message_type":"DO_INIT_ACK"}`,
		regexp.MustCompile(`Field 'initial_game_state' is missing`))
}

func TestInvalidGlBadDoInitAckNonObjectGameState(t *testing.T) {
	subtestBadDoInitAck(t, `{"message_type":"DO_INIT_ACK", `+
		`"initial_game_state":0}`,
		regexp.MustCompile(`Non-object value for field 'initial_game_state'`))
}

func TestInvalidGlBadDoInitAckNoAllClients(t *testing.T) {
	subtestBadDoInitAck(t, `{"message_type":"DO_INIT_ACK", `+
		`"initial_game_state":{}}`,
		regexp.MustCompile(`Field 'all_clients' is missing`))
}
