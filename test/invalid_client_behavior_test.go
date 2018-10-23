package test

import (
	"github.com/netorcai/netorcai/client/go"
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

func TestInvalidPlayerMessageBeforeStart(t *testing.T) {
	proc, _, playerClients, _, glClients := runNetorcaiAndClients(t,
		[]string{}, 1000, 1, 0)
	defer killallNetorcaiSIGKILL()

	playerClients[0].SendString(`{"message_type": "TURN_ACK", ` +
		`"turn_number": -1, "actions":[]}`)
	checkAllKicked(t, playerClients, regexp.MustCompile(`Received a TURN_ACK `+
		`but the client state is not THINKING`), 1000)

	proc.inputControl <- `quit`
	waitCompletionTimeout(proc.completion, 1000)
	checkAllKicked(t, glClients, regexp.MustCompile(`netorcai abort`), 1000)
}

func TestInvalidVisuMessageBeforeStart(t *testing.T) {
	proc, _, _, visuClients, glClients := runNetorcaiAndClients(t,
		[]string{}, 1000, 0, 1)
	defer killallNetorcaiSIGKILL()

	visuClients[0].SendString(`{"message_type": "TURN_ACK", ` +
		`"turn_number": -1, "actions":[]}`)
	checkAllKicked(t, visuClients, regexp.MustCompile(`Received a TURN_ACK `+
		`but the client state is not THINKING`), 1000)

	proc.inputControl <- `quit`
	waitCompletionTimeout(proc.completion, 1000)
	checkAllKicked(t, glClients, regexp.MustCompile(`netorcai abort`), 1000)
}

func TestInvalidGlNoDoInitAck(t *testing.T) {
	proc, _, playerClients, visuClients, glClients := runNetorcaiAndAllClients(
		t, []string{}, 1000)
	defer killallNetorcaiSIGKILL()

	go func(glClient *client.Client) {
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
