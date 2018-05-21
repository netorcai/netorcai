package test

import (
	//"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestHelloGLOnly(t *testing.T) {
	proc, _, players, visus, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=2",
			"--delay-turns=500", "--debug"}, 1000)
	defer killallNetorcaiSIGKILL()
	proc.printOutput = true

	// Disconnect all players
	for _, player := range players {
		player.Disconnect()
		waitOutputTimeout(regexp.MustCompile(`Remote endpoint closed`),
			proc.outputControl, 1000, false)
	}

	// Disconnect all visus
	for _, visu := range visus {
		visu.Disconnect()
		waitOutputTimeout(regexp.MustCompile(`Remote endpoint closed`),
			proc.outputControl, 1000, false)
	}

	// Run a game client
	glExit := make(chan int)
	go hello_game_logic(t, gl[0], 0, 2, glExit)

	// Start the game
	proc.inputControl <- "start"

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}
