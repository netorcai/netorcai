package test

import (
	"fmt"
	"github.com/netorcai/netorcai"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestHelloGLOnly(t *testing.T) {
	proc, _, players, _, visus, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=2",
			"--delay-turns=500", "--debug"}, 1000, 0)
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
	go helloGameLogic(t, gl[0], 0, 0, 2, 2, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		regexp.MustCompile(`Game is finished`))

	// Start the game
	proc.inputControl <- "start"

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func TestHelloGLIdleClients(t *testing.T) {
	proc, _, _, _, _, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=2",
			"--delay-turns=500", "--debug"}, 1000, 0)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], 4, 0, 2, 2, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		regexp.MustCompile(`Game is finished`))

	// Start the game
	proc.inputControl <- "start"

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func TestHelloGLIdleClientsSpecial(t *testing.T) {
	proc, _, _, _, _, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=2",
			"--delay-turns=500", "--debug", "--nb-splayers-max=1"}, 1000, 1)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], 4, 1, 2, 2, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		regexp.MustCompile(`Game is finished`))

	// Start the game
	proc.inputControl <- "start"

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func TestHelloGLActiveVisu(t *testing.T) {
	proc, _, players, _, visus, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=3",
			"--delay-turns=500", "--debug", "--json-logs"}, 1000, 0)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], 0, 0, 3, 3, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		regexp.MustCompile(`Game is finished`))

	// Disconnect players
	for _, player := range players {
		player.Disconnect()
		waitOutputTimeout(regexp.MustCompile(`Remote endpoint closed`),
			proc.outputControl, 1000, false)
	}

	// Run visu clients
	for visuID, visu := range visus {
		go helloClient(t, visu, fmt.Sprintf("Visu%v", visuID),
			0, 0, 3, 3, 0, 500, 500, false, false, true, true,
			DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn,
			DefaultHelloClientCheckGameEnds,
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
	proc, _, players, _, visus, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=3",
			"--delay-turns=500", "--debug", "--json-logs"}, 1000, 0)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], 1, 0, 3, 3, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		regexp.MustCompile(`Game is finished`))

	// Run an active player
	go helloClient(t, players[0], "Player0", 1, 0, 3, 3, 0, 500, 500, true, false, true, true,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds,
		DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`))

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

func TestHelloGLActiveSpecialPlayer(t *testing.T) {
	proc, _, players, specialPlayers, visus, gl := runNetorcaiAndAllClients(
		t, []string{"--delay-first-turn=500", "--nb-turns-max=3",
			"--delay-turns=500", "--debug", "--json-logs",
			"--nb-splayers-max=1"}, 1000, 1)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], 0, 1, 3, 3, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		regexp.MustCompile(`Game is finished`))

	// Run an active special special
	go helloClient(t, specialPlayers[0], "SpecialPlayer0", 0, 1, 3, 3, 0, 500, 500, true, false, true, true,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds,
		DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`))

	// Disconnect classical players
	for _, player := range players {
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
	netorcaiAdditionalArgs []string,
	nbPlayers, nbSpecialPlayers, nbVisus int,
	nbTurnsNetorcai, nbTurnsGL, nbTurnsPlayer, nbTurnsVisu int,
	nbTurnsToSkipPlayer, nbTurnsToSkipVisu int,
	allowTurnSkipPlayer, allowTurnSkipVisu bool,
	checkGameStartsFunc ClientGameStartsCheckFunc,
	playerCheckTurnFunc, visuCheckTurnFunc ClientTurnCheckFunc,
	checkGameEndsFunc ClientGameEndsCheckFunc,
	checkDoTurnFunc GLCheckDoTurnFunc,
	doInitAckFunc GLDoInitAckFunc, doTurnAckFunc GLDoTurnAckFunc,
	playerTurnAckFunc, visuTurnAckFunc ClientTurnAckFunc,
	glKickReasonMatcher, playerKickReasonMatcher,
	visuKickReasonMatcher *regexp.Regexp) {
	proc, _, players, specialPlayers, visus, gl := runNetorcaiAndClients(
		t, append([]string{"--delay-first-turn=500",
			fmt.Sprintf("--nb-turns-max=%v", nbTurnsNetorcai),
			fmt.Sprintf("--nb-players-max=%v", nbPlayers),
			fmt.Sprintf("--nb-splayers-max=%v", nbSpecialPlayers),
			fmt.Sprintf("--nb-visus-max=%v", nbVisus),
			"--delay-turns=500", "--debug", "--json-logs", "--autostart"},
			netorcaiAdditionalArgs...),
		1000, nbPlayers, nbSpecialPlayers, nbVisus)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], nbPlayers, nbSpecialPlayers, nbTurnsNetorcai, nbTurnsGL,
		checkDoTurnFunc, doInitAckFunc, doTurnAckFunc,
		glKickReasonMatcher)

	// Run player clients
	for playerID, player := range players {
		go helloClient(t, player, fmt.Sprintf("Player%v", playerID),
			nbPlayers, nbSpecialPlayers, nbTurnsNetorcai, nbTurnsPlayer,
			nbTurnsToSkipPlayer, 500, 500, true, allowTurnSkipPlayer,
			nbTurnsPlayer == nbTurnsNetorcai, nbTurnsGL > 0,
			checkGameStartsFunc, playerCheckTurnFunc, checkGameEndsFunc,
			playerTurnAckFunc, playerKickReasonMatcher)
	}

	// Run special player clients
	for splayerID, splayer := range specialPlayers {
		go helloClient(t, splayer, fmt.Sprintf("SpecialPlayer%v", splayerID),
			nbPlayers, nbSpecialPlayers, nbTurnsNetorcai, nbTurnsPlayer,
			nbTurnsToSkipPlayer, 500, 500, true, allowTurnSkipPlayer,
			nbTurnsPlayer == nbTurnsNetorcai, nbTurnsGL > 0,
			checkGameStartsFunc, playerCheckTurnFunc, checkGameEndsFunc,
			playerTurnAckFunc, playerKickReasonMatcher)
	}

	// Run visu clients
	for visuID, visu := range visus {
		go helloClient(t, visu, fmt.Sprintf("Visu%v", visuID),
			nbPlayers, nbSpecialPlayers, nbTurnsNetorcai, nbTurnsVisu,
			nbTurnsToSkipVisu, 500, 500, false, allowTurnSkipVisu,
			nbTurnsVisu == nbTurnsNetorcai, nbTurnsGL > 0,
			checkGameStartsFunc, visuCheckTurnFunc, checkGameEndsFunc,
			visuTurnAckFunc, visuKickReasonMatcher)
	}

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func TestHelloGLActiveClients(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 4, 0, 1,
		3, 3, 3, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

func TestHelloGLActiveClientsFast(t *testing.T) {
	subtestHelloGlActiveClients(t, []string{"--fast"}, 4, 0, 1,
		1000, 1000, 1000, 1000,
		0, 0,
		false, true,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

func TestHelloGLActiveClientsSpecial(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 4, 1, 1,
		3, 3, 3, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}

// Invalid DO_INIT_ACK
func doInitAckNoMsgType(nbPlayers, nbSpecialPlayers, nbTurns int) string {
	return `{"initial_game_state":{"all_clients":{}}}`
}

func doInitAckNoInitialGameState(nbPlayers, nbSpecialPlayers, nbTurns int) string {
	return `{"message_type": "DO_INIT_ACK"}`
}

func doInitAckBadMsgType(nbPlayers, nbSpecialPlayers, nbTurns int) string {
	return `{"message_type": "DO_INIT_ACKz",
		"initial_game_state":{"all_clients":{}}}`
}

func doInitAckBadInitialGameStateNotObject(nbPlayers, nbSpecialPlayers, nbTurns int) string {
	return `{"message_type":"DO_INIT_ACK", "initial_game_state":0}`
}

func doInitAckBadInitialGameStateNoAllClients(nbPlayers, nbSpecialPlayers, nbTurns int) string {
	return `{"message_type":"DO_INIT_ACK", "initial_game_state":{}}`
}

func TestInvalidDoInitAckNoMsgType(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 0, 1, 1,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		doInitAckNoMsgType, DefaultHelloGlDoTurnAck,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Invalid DO_INIT_ACK message. `+
			`Field 'message_type' is missing`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

func TestInvalidDoInitAckNoInitialGameState(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 0, 1, 1,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		doInitAckNoInitialGameState, DefaultHelloGlDoTurnAck,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Invalid DO_INIT_ACK message. `+
			`Field 'initial_game_state' is missing`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

func TestInvalidDoInitAckBadMsgType(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 0, 1, 1,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		doInitAckBadMsgType, DefaultHelloGlDoTurnAck,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`DO_INIT_ACK was expected`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

func TestInvalidDoInitAckBadInitialGameStateNotObject(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 0, 1, 1,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		doInitAckBadInitialGameStateNotObject, DefaultHelloGlDoTurnAck,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Non-object value for field 'initial_game_state'`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

func TestInvalidDoInitAckBadInitialGameStateNoAllClients(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 0, 1, 1,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		doInitAckBadInitialGameStateNoAllClients, DefaultHelloGlDoTurnAck,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Field 'all_clients' is missing`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

// Invalid DO_TURN_ACK
func doTurnAckNoMsgType(turn int, actions []interface{}) string {
	return `{"winner_player_id":-1, "game_state":{"all_clients":{}}}`
}

func doTurnAckBadMsgType(turn int, actions []interface{}) string {
	return `{"message_type":"DO_TURN_ACKz", "winner_player_id":-1,` +
		`"game_state":{"all_clients":{}}}`
}

func doTurnAckNoWinner(turn int, actions []interface{}) string {
	return `{"message_type":"DO_TURN_ACK", "game_state":{"all_clients":{}}}`
}

func doTurnAckNoGameState(turn int, actions []interface{}) string {
	return `{"message_type":"DO_TURN_ACK", "winner_player_id":-1}`
}

func doTurnAckNoAllClients(turn int, actions []interface{}) string {
	return `{"message_type":"DO_TURN_ACK", "winner_player_id":-1, ` +
		`"game_state":{}}`
}

func doTurnAckBadWinner(turn int, actions []interface{}) string {
	return `{"message_type":"DO_TURN_ACK", "winner_player_id": 42,` +
		`"game_state":{"all_clients":{}}}`
}

func TestInvalidDoTurnAckNoMsgType(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 1, 0, 0,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, doTurnAckNoMsgType,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Field 'message_type' is missing`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

func TestInvalidDoTurnAckBadMsgType(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 1, 0, 0,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, doTurnAckBadMsgType,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`DO_TURN_ACK was expected`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

func TestInvalidDoTurnAckNoWinner(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 1, 0, 0,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, doTurnAckNoWinner,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Field 'winner_player_id' is missing`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

func TestInvalidDoTurnAckNoGameState(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 1, 0, 0,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, doTurnAckNoGameState,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Field 'game_state' is missing`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

func TestInvalidDoTurnAckNoAllClients(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 1, 0, 0,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, doTurnAckNoAllClients,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Field 'all_clients' is missing`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

func TestInvalidDoTurnAckBadWinner(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 1, 0, 0,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, doTurnAckBadWinner,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Invalid winner_player_id: Not in \[-1, 1\[`),
		regexp.MustCompile(`netorcai abort`),
		regexp.MustCompile(`netorcai abort`))
}

// Invalid TURN_ACK
func turnAckNoMsgType(turn, playerID int) string {
	return fmt.Sprintf(`{"turn_number": %v, "actions": []}`, turn)
}

func turnAckNoTurnNumber(turn, playerID int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK", "actions": []}`)
}

func turnAckNoActions(turn, playerID int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": %v}`, turn)
}

func turnAckBadMsgType(turn, playerID int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACKz",
		"turn_number": %v, "actions": []}`, turn)
}

func turnAckBadTurnNumberValue(turn, playerID int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": %v, "actions": []}`, turn+1)
}

func turnAckBadTurnNumberNotInt(turn, playerID int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": "nope", "actions": []}`)
}

func turnAckBadActions(turn, playerID int) string {
	return fmt.Sprintf(`{"message_type": "TURN_ACK",
		"turn_number": %v, "actions": {}}`, turn)
}

func TestInvalidTurnAckNoMsgType(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 3, 2, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckNoMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Field 'message_type' is missing`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckNoTurnNumber(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 3, 2, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckNoTurnNumber, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Field 'turn_number' is missing`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckNoActions(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 3, 2, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckNoActions, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Field 'actions' is missing`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckBadMsgType(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 3, 2, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckBadMsgType, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`TURN_ACK was expected`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckBadTurnNumberValue(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 3, 2, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckBadTurnNumberValue, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Invalid value \(turn_number=1\)`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckBadTurnNumberNotInt(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 3, 2, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckBadTurnNumberNotInt, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Non-integral value for field 'turn_number'`),
		regexp.MustCompile(`Game is finished`))
}

func TestInvalidTurnAckBadActions(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 1, 0, 1,
		3, 3, 2, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		DefaultHelloClientCheckGameEnds, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		turnAckBadActions, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Non-array value for field 'actions'`),
		regexp.MustCompile(`Game is finished`))
}

// Winner
func doTurnAckWinner(turn int, actions []interface{}) string {
	return `{"message_type":"DO_TURN_ACK",
		"winner_player_id":0,
		"game_state":{"all_clients":{}}}`
}

func checkGameEndsWinner(t *testing.T, msg map[string]interface{}) {
	checkGameEnds(t, msg)

	winner, err := netorcai.ReadInt(msg, "winner_player_id")
	assert.NoError(t, err, "Cannot read 'winner_player_id'")
	assert.Equal(t, 0, winner, "Unexpected 'winner_player_id' value")
}

func TestHelloWinner(t *testing.T) {
	subtestHelloGlActiveClients(t, nil, 4, 0, 1,
		3, 3, 3, 3,
		0, 0,
		false, false,
		DefaultHelloClientCheckGameStarts, DefaultHelloClientCheckTurn, DefaultHelloClientCheckTurn,
		checkGameEndsWinner, DefaultHelloGLCheckDoTurn,
		DefaultHelloGLDoInitAck, doTurnAckWinner,
		DefaultHelloClientTurnAck, DefaultHelloClientTurnAck,
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`),
		regexp.MustCompile(`Game is finished`))
}
