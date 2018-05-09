package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestInvalidGlMessageBeforeStart(t *testing.T) {
	proc, _, playerClients, visuClients, glClients := runNetorcaiAndAllClients(
		t, 1000)
	defer killallNetorcai()

	glClients[0].SendString(`{}`)
	checkAllKicked(t, glClients, regexp.MustCompile(
		`Received a game logic message but the game has not started`), 1000)

	_, err := waitOutputTimeout(regexp.MustCompile(`Game logic failed`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err,
		"Cannot read `Game logic failed` in netorcai output")

	checkAllKicked(t, playerClients, regexp.MustCompile(`netorcai abort`),
		1000)
	checkAllKicked(t, visuClients, regexp.MustCompile(`netorcai abort`), 1000)
}

func TestInvalidGlNoDoInitAck(t *testing.T) {
	proc, _, playerClients, visuClients, glClients := runNetorcaiAndAllClients(
		t, 1000)
	defer killallNetorcai()

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

	checkAllKicked(t, playerClients, regexp.MustCompile(`netorcai abort`),
		1000)
	checkAllKicked(t, visuClients, regexp.MustCompile(`netorcai abort`), 1000)
}
