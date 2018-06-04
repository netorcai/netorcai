package test

import (
	//"github.com/stretchr/testify/assert"
	"fmt"
	"regexp"
	"testing"
)

func TestHelloGLOnly(t *testing.T) {
	proc, _, players, visus, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=2",
			"--delay-turns=500", "--debug"}, 1000)
	defer killallNetorcaiSIGKILL()

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
	go helloGameLogic(t, gl[0], 0, 2,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck)

	// Start the game
	proc.inputControl <- "start"

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func TestHelloGLIdleClients(t *testing.T) {
	proc, _, _, _, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=2",
			"--delay-turns=500", "--debug"}, 1000)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], 4, 2,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck)

	// Start the game
	proc.inputControl <- "start"

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func TestHelloGLActiveVisu(t *testing.T) {
	proc, _, players, visus, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=3",
			"--delay-turns=500", "--debug", "--json-logs"}, 1000)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], 0, 3,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck)

	// Disconnect players
	for _, player := range players {
		player.Disconnect()
		waitOutputTimeout(regexp.MustCompile(`Remote endpoint closed`),
			proc.outputControl, 1000, false)
	}

	// Run visu clients
	for _, visu := range visus {
		go helloClient(t, visu, 0, 3, 3, 500, 500, false, true,
			DefaultHelloClientTurnAck, regexp.MustCompile(`Game is finished`))
	}

	// Start the game
	proc.inputControl <- "start"

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func TestHelloGLActivePlayer(t *testing.T) {
	proc, _, players, visus, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=3",
			"--delay-turns=500", "--debug", "--json-logs"}, 1000)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], 1, 3,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck)

	// Run an active player
	go helloClient(t, players[0], 1, 3, 3, 500, 500, true, true,
		DefaultHelloClientTurnAck, regexp.MustCompile(`Game is finished`))

	// Disconnect other players
	for _, player := range players[1:] {
		player.Disconnect()
		waitOutputTimeout(regexp.MustCompile(`Remote endpoint closed`),
			proc.outputControl, 1000, false)
	}

	// Disconnect visus
	for _, visu := range visus {
		visu.Disconnect()
		waitOutputTimeout(regexp.MustCompile(`Remote endpoint closed`),
			proc.outputControl, 1000, false)
	}

	// Start the game
	proc.inputControl <- "start"

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func subtestHelloGlActiveClients(t *testing.T,
	nbPlayers, nbVisus int,
	nbTurnsGL, nbTurnsPlayer, nbTurnsVisu int,
	playerTurnAckFunc, visuTurnAckFunc ClientTurnAckFunc,
	playerKickReasonMatcher, visuKickReasonMatcher *regexp.Regexp) {
	proc, _, players, visus, gl := runNetorcaiAndClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=3",
			"--delay-turns=500", "--debug", "--json-logs"}, 1000, nbPlayers,
		nbVisus)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], nbPlayers, nbTurnsGL,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck)

	// Run player clients
	for _, player := range players {
		go helloClient(t, player, nbPlayers, 3, nbTurnsPlayer, 500, 500, true,
			nbTurnsPlayer == nbTurnsGL, playerTurnAckFunc,
			playerKickReasonMatcher)
	}

	// Run visu clients
	for _, visu := range visus {
		go helloClient(t, visu, nbPlayers, 3, nbTurnsVisu, 500, 500, false,
			nbTurnsVisu == nbTurnsGL, visuTurnAckFunc,
			visuKickReasonMatcher)
	}

	// Start the game
	proc.inputControl <- "start"

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func TestHelloGLActiveClients(t *testing.T) {
	subtestHelloGlActiveClients(t, 4, 1, 3, 3, 3,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

// Invalid DO_INIT_ACK

// Invalid TURN_ACK
func turnAckNoMsgType(turn int) string {
	return fmt.Sprintf(`{"turn_number": %v, "actions": []}`, turn)
}

func turnAckNoTurnNumber(turn int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK", "actions": []}`)
}

func turnAckNoActions(turn int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": %v}`, turn)
}

func turnAckBadMsgType(turn int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACKz",
		"turn_number": %v, "actions": []}`, turn)
}

func turnAckBadTurnNumberValue(turn int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": %v, "actions": []}`, turn+1)
}

func turnAckBadTurnNumberNotInt(turn int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": "nope", "actions": []}`)
}

func turnAckBadActions(turn int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": %v, "actions": {}}`, turn)
}

func TestInvalidTurnAckNoMsgType(t *testing.T) {
	subtestHelloGlActiveClients(t, 1, 1, 3, 2, 3,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Field 'message_type' is missing`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckNoTurnNumber(t *testing.T) {
	subtestHelloGlActiveClients(t, 1, 1, 3, 2, 3,
		turnAckNoTurnNumber, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Field 'turn_number' is missing`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckNoActions(t *testing.T) {
	subtestHelloGlActiveClients(t, 1, 1, 3, 2, 3,
		turnAckNoActions, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Field 'actions' is missing`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckBadMsgType(t *testing.T) {
	subtestHelloGlActiveClients(t, 1, 1, 3, 2, 3,
		turnAckBadMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`TURN_ACK was expected`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckBadTurnNumberValue(t *testing.T) {
	subtestHelloGlActiveClients(t, 1, 1, 3, 2, 3,
		turnAckBadTurnNumberValue, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Invalid value \(turn_number=1\)`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckBadTurnNumberNotInt(t *testing.T) {
	subtestHelloGlActiveClients(t, 1, 1, 3, 2, 3,
		turnAckBadTurnNumberNotInt, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Non-integral value for field 'turn_number'`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckBadActions(t *testing.T) {
	subtestHelloGlActiveClients(t, 1, 1, 3, 2, 3,
		turnAckBadActions, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Non-array value for field 'actions'`),
		regexp.MustCompile(`Game is finished`))
}
