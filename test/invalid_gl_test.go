package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestInvalidGameLogicMessageBeforeStart(t *testing.T) {
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
