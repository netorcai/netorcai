package test

import (
	"fmt"
	"github.com/netorcai/netorcai"
	"github.com/netorcai/netorcai/client/go"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

func readFloat(data map[string]interface{}, field string) (float64, error) {
	value, exists := data[field]
	if !exists {
		return 0, fmt.Errorf("Field '%v' is missing", field)
	}

	switch value.(type) {
	default:
		return 0, fmt.Errorf("Non-float value for field '%v'", field)
	case float64:
		return value.(float64), nil
	}
}

func readBool(data map[string]interface{}, field string) (bool, error) {
	value, exists := data[field]
	if !exists {
		return false, fmt.Errorf("Field '%v' is missing", field)
	}

	switch value.(type) {
	default:
		return false, fmt.Errorf("Non-bool value for field '%v'", field)
	case bool:
		return value.(bool), nil
	}
}

// Netorcai helpers
func runNetorcaiWaitListening(t *testing.T,
	arguments []string) *NetorcaiProcess {
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, arguments)
	assert.NoError(t, err, "Cannot start netorcai")

	_, err = waitListening(proc.outputControl, 1000)
	if err != nil {
		killallNetorcai()
		assert.NoError(t, err, "Netorcai is not listening")
	}

	return proc
}

func waitCompletionTimeout(completion chan int, timeoutMS int) (
	exitCode int, err error) {
	select {
	case exitCode := <-completion:
		return exitCode, nil
	case <-time.After(time.Duration(timeoutMS) * time.Millisecond):
		return -1, fmt.Errorf("Timeout reached")
	}
}

func waitOutputTimeout(re *regexp.Regexp, output chan string,
	timeoutMS int, leaveOnNonMatch bool) (matchingLine string, err error) {
	timeoutReached := make(chan int)
	stopTimeout := make(chan int)
	defer close(timeoutReached)
	defer close(stopTimeout)
	go func() {
		select {
		case <-stopTimeout:
		case <-time.After(time.Duration(timeoutMS) * time.Millisecond):
			timeoutReached <- 0
		}
	}()

	for {
		select {
		case line := <-output:
			if re.MatchString(line) {
				stopTimeout <- 0
				return line, nil
			} else {
				if leaveOnNonMatch {
					stopTimeout <- 0
					return line, fmt.Errorf("Non-matching line read: %v", line)
				}
			}
		case <-timeoutReached:
			return "", fmt.Errorf("Timeout reached")
		}
	}
}

func waitListening(output chan string, timeoutMS int) (
	matchingLine string, err error) {
	re := regexp.MustCompile("Listening incoming connections")
	return waitOutputTimeout(re, output, timeoutMS, true)
}

func killallNetorcai() error {
	cmd := exec.Command("killall")
	cmd.Args = []string{"killall", "--quiet", "netorcai", "netorcai.cover"}
	return cmd.Run()
}

func killallNetorcaiSIGKILL() error {
	cmd := exec.Command("killall")
	cmd.Args = []string{"killall", "-KILL", "--quiet", "netorcai", "netorcai.cover"}
	return cmd.Run()
}

func handleCoverage(t *testing.T, expRetCode int) (coverFilename string,
	expectedReturnCode int) {
	_, exists := os.LookupEnv("DO_COVERAGE")
	if exists {
		coverFilename = t.Name() + ".covout"
		expectedReturnCode = 0
		return
	} else {
		coverFilename = ""
		expectedReturnCode = expRetCode
	}

	return coverFilename, expectedReturnCode
}

func isTravis() bool {
	_, exists := os.LookupEnv("TRAVIS")
	return exists
}

// Client helpers
func waitReadMessage(client *client.Client, timeoutMS int) (
	msg map[string]interface{}, err error) {
	msgChan := make(chan int)
	go func() {
		msg, err = client.ReadMessage()
		msgChan <- 0
	}()

	select {
	case <-msgChan:
		close(msgChan)
		return msg, err
	case <-time.After(time.Duration(timeoutMS) * time.Millisecond):
		return msg, fmt.Errorf("Timeout reached")
	}
}

func connectClient(t *testing.T, role, nickname string, timeoutMS int) (
	*client.Client, error) {
	client := &client.Client{}
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")

	err = client.SendLogin(role, nickname)
	assert.NoError(t, err, "Cannot send LOGIN")

	msg, err := waitReadMessage(client, 1000)
	assert.NoError(t, err, "Cannot read client message (LOGIN_ACK)")
	checkLoginAck(t, msg)
	return client, nil
}

func runNetorcaiAndClients(t *testing.T, arguments []string,
	timeoutMS int, nbPlayers, nbSpecialPlayers, nbVisus int) (
	proc *NetorcaiProcess, clients, playerClients, specialPlayerClients, visuClients,
	glClients []*client.Client) {
	proc = runNetorcaiWaitListening(t, arguments)

	// Players
	for i := 0; i < nbPlayers; i++ {
		player, err := connectClient(t, "player", "player", timeoutMS)
		if err != nil {
			killallNetorcai()
			assert.NoError(t, err, "Cannot connect client")
		}
		clients = append(clients, player)
		playerClients = append(playerClients, player)
	}

	// Special players
	for i := 0; i < nbSpecialPlayers; i++ {
		splayer, err := connectClient(t, "special player", "splayer", timeoutMS)
		if err != nil {
			killallNetorcai()
			assert.NoError(t, err, "Cannot connect client")
		}
		clients = append(clients, splayer)
		specialPlayerClients = append(specialPlayerClients, splayer)
	}

	// Visus
	for i := 0; i < nbVisus; i++ {
		visu, err := connectClient(t, "visualization", "visu", timeoutMS)
		if err != nil {
			killallNetorcai()
			assert.NoError(t, err, "Cannot connect client")
		}
		clients = append(clients, visu)
		visuClients = append(visuClients, visu)
	}

	// Game Logic
	for i := 0; i < 1; i++ {
		gl, err := connectClient(t, "game logic", "game_logic", timeoutMS)
		if err != nil {
			killallNetorcai()
			assert.NoError(t, err, "Cannot connect client")
		}
		clients = append(clients, gl)
		glClients = append(glClients, gl)
	}

	return proc, clients, playerClients, specialPlayerClients, visuClients, glClients
}

func runNetorcaiAndAllClients(t *testing.T, arguments []string,
	timeoutMS int, nbSpecialPlayers int) (
	proc *NetorcaiProcess, clients, playerClients, specialPlayerClients, visuClients,
	glClients []*client.Client) {
	return runNetorcaiAndClients(t, arguments, timeoutMS, 4, nbSpecialPlayers, 1)
}

func checkAllKicked(t *testing.T, clients []*client.Client,
	reasonMatcher *regexp.Regexp, timeoutMS int) {
	timeoutReached := make(chan int)
	stopTimeout := make(chan int)
	defer close(timeoutReached)
	defer close(stopTimeout)
	go func() {
		select {
		case <-stopTimeout:
		case <-time.After(time.Duration(timeoutMS) * time.Millisecond):
			timeoutReached <- 0
		}
	}()

	// All clients should receive a KICK
	kickChan := make(chan int, len(clients))
	for _, cli := range clients {
		go func(c *client.Client) {
			msg, err := waitReadMessage(c, timeoutMS)
			assert.NoError(t, err, "Cannot read client message (KICK)")
			checkKick(t, msg, reasonMatcher)
			kickChan <- 0
		}(cli)
	}

	for _ = range clients {
		select {
		case <-kickChan:
		case <-timeoutReached:
			assert.FailNow(t, "Timeout reached")
		}
	}

	stopTimeout <- 0
	close(kickChan)
}

func checkKick(t *testing.T, msg map[string]interface{},
	reasonMatcher *regexp.Regexp) {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err,
		"Cannot read 'message_type' field in received client message (KICK)")
	assert.Equal(t, "KICK", messageType, "Unexpected message type")

	kickReason, err := netorcai.ReadString(msg, "kick_reason")
	assert.NoError(t, err, "Cannot read 'kick_reason' in received client message (KICK)")
	assert.Regexp(t, reasonMatcher, kickReason, "Unexpected kick reason")
}

func checkLoginAck(t *testing.T, msg map[string]interface{}) {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err, "Cannot read 'message_type' field in "+
		"received client message (LOGIN_ACK)")

	switch messageType {
	case "LOGIN_ACK":
	case "KICK":
		kickReason, err := netorcai.ReadString(msg, "kick_reason")
		assert.NoError(t, err, "Cannot read kick_reason")

		assert.FailNow(t, "Expected LOGIN_ACK, got KICK", kickReason)
	default:
		assert.FailNowf(t, "Expected LOGIN_ACK, got another message type",
			messageType)
	}
}

func checkDoInit(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedNbTurnsMax int) {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err, "Cannot read 'message_type' field in "+
		"received client message (DO_INIT)")

	switch messageType {
	case "DO_INIT":
		nbPlayers, err := netorcai.ReadInt(msg, "nb_players")
		assert.NoError(t, err, "Cannot read nb_players")
		assert.Equal(t, expectedNbPlayers, nbPlayers,
			"Unexpected value for nb_players in received DO_INIT message")

		nbTurnsMax, err := netorcai.ReadInt(msg, "nb_turns_max")
		assert.NoError(t, err, "Cannot read nb_turns_max")
		assert.Equal(t, expectedNbTurnsMax, nbTurnsMax,
			"Unexpected value for nb_turns_max in received DO_INIT message")
	case "KICK":
		kickReason, err := netorcai.ReadString(msg, "kick_reason")
		assert.NoError(t, err, "Cannot read kick_reason")

		assert.FailNow(t, "Expected DO_INIT, got KICK", kickReason)
	default:
		assert.FailNowf(t, "Expected DO_INIT, got another message type",
			messageType)
	}
}

func checkDoTurn(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedNbSpecialPlayers, expectedTurnNumber int) []interface{} {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err, "Cannot read 'message_type' field in "+
		"received client message (DO_TURN)")

	switch messageType {
	case "DO_TURN":
		playerActions, err := netorcai.ReadArray(msg, "player_actions")
		assert.NoError(t, err, "Cannot read player_actions in DO_TURN message")
		assert.Condition(t, func() bool {
			return len(playerActions) <= expectedNbPlayers+expectedNbSpecialPlayers
		}, "Invalid player_actions array in DO_TURN message: Size=%v while "+
			"nb_players=%v and nb_special_players=%v",
			len(playerActions), expectedNbPlayers, expectedNbSpecialPlayers)

		for playerIndex, pActions := range playerActions {
			obj := pActions.(map[string]interface{})

			playerID, err := netorcai.ReadInt(obj, "player_id")
			assert.NoError(t, err, "Invalid player_actions in DO_TURN "+
				"message: Cannot read player_id in array element %v",
				playerIndex)
			assert.Condition(t, func() bool {
				return playerID >= 0 && playerID < expectedNbPlayers+expectedNbSpecialPlayers
			}, "Invalid player_id=%v in player_actions[%v] in DO_TURN "+
				"message: Should be in [0,%v[",
				playerID, playerIndex, expectedNbPlayers+expectedNbSpecialPlayers)

			turnNumber, err := netorcai.ReadInt(obj, "turn_number")
			assert.NoError(t, err, "Invalid player_actions in DO_TURN "+
				"message: Cannot read turn_number in array element %v",
				playerIndex)
			assert.Equal(t, expectedTurnNumber, turnNumber,
				"Unexpected turn_number in DO_TURN player action %v",
				playerIndex)

			_, err = netorcai.ReadArray(obj, "actions")
			assert.NoError(t, err, "Invalid player_actions in DO_TURN "+
				"message: Cannot read the actions array in player action %v",
				playerIndex)

			return playerActions
		}
	case "KICK":
		kickReason, err := netorcai.ReadString(msg, "kick_reason")
		assert.NoError(t, err, "Cannot read kick_reason")

		assert.FailNow(t, "Expected DO_TURN, got KICK", kickReason)
	default:
		assert.FailNowf(t, "Expected DO_TURN, got another message type",
			messageType)
	}

	return []interface{}{}
}

func checkPlayersInfo(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedNbSpecialPlayers int, isPlayer bool) {
	playersInfo, err := netorcai.ReadArray(msg, "players_info")
	assert.NoError(t, err, "Cannot read players_info in GAME_STARTS")
	if isPlayer {
		assert.Equal(t, 0, len(playersInfo),
			"Unexpected players_info: Should be empty for players")
	} else {
		assert.Equal(t, expectedNbPlayers+expectedNbSpecialPlayers, len(playersInfo),
			"Unexpected player_info array size: "+
				"Should match number of players for visualization")
		playerIDs := make([]int, 0)
		for playerIndex, player := range playersInfo {
			obj := player.(map[string]interface{})

			pid, err := netorcai.ReadInt(obj, "player_id")
			assert.NoError(t, err, "Cannot read player_id in "+
				"players_info[%v] of GAME_STARTS message (as a visu)",
				playerIndex)
			playerIDs = append(playerIDs, pid)

			_, err = netorcai.ReadString(obj, "nickname")
			assert.NoError(t, err, "Cannot read nickname in "+
				"players_info[%v] of GAME_STARTS message (as a visu)",
				playerIndex)

			_, err = netorcai.ReadString(obj, "remote_address")
			assert.NoError(t, err, "Cannot read remote_address in "+
				"players_info[%v] of GAME_STARTS message (as a visu)",
				playerIndex)

			_, err = readBool(obj, "is_connected")
			assert.NoError(t, err, "Cannot read nickname in "+
				"players_info[%v] of GAME_STARTS message (as a visu)",
				playerIndex)
		}

		for i := 0; i < expectedNbPlayers+expectedNbSpecialPlayers; i++ {
			assert.Contains(t, playerIDs, i,
				"Invalid players_info in GAME_STARTS message (as a visu): "+
					"No info for player_id=%v while there should be "+
					"nb_players=%v", i, expectedNbPlayers)
		}
	}
}

func checkGameStarts(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedNbSpecialPlayers, expectedNbTurnsMax int,
	expectedMsBeforeFirstTurn, expectedMsBetweenTurns float64,
	isPlayer bool) (playerID int) {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err, "Cannot read 'message_type' field in "+
		"received client message (GAME_STARTS)")

	switch messageType {
	case "GAME_STARTS":
		nbPlayers, err := netorcai.ReadInt(msg, "nb_players")
		assert.NoError(t, err, "Cannot read nb_players in GAME_STARTS")
		assert.Equal(t, expectedNbPlayers, nbPlayers,
			"Unexpected value for nb_players in received GAME_STARTS message")

		nbSpecialPlayers, err := netorcai.ReadInt(msg, "nb_special_players")
		assert.NoError(t, err, "Cannot read nb_special_players in GAME_STARTS")
		assert.Equal(t, expectedNbSpecialPlayers, nbSpecialPlayers,
			"Unexpected value for nb_special_players in received GAME_STARTS message")

		nbTurnsMax, err := netorcai.ReadInt(msg, "nb_turns_max")
		assert.NoError(t, err, "Cannot read nb_turns_max")
		assert.Equal(t, expectedNbTurnsMax, nbTurnsMax,
			"Unexpected value for nb_turns_max in GAME_STARTS message")

		playerID, err := netorcai.ReadInt(msg, "player_id")
		assert.NoError(t, err, "Cannot read player_id in GAME_STARTS")
		if isPlayer {
			assert.Condition(t, func() bool {
				return playerID >= 0 && playerID < expectedNbPlayers+expectedNbSpecialPlayers
			}, "Invalid player_id=%v in GAME_STARTS message: "+
				"Should be in [0,%v[ for a player",
				playerID, expectedNbPlayers+expectedNbSpecialPlayers)
		} else {
			assert.Equal(t, -1, playerID, "Invalid player_id=%v in "+
				"GAME_STARTS message: Should be -1 for visualization",
				playerID)
		}

		msBeforeFirstTurn, err := readFloat(msg,
			"milliseconds_before_first_turn")
		assert.NoError(t, err,
			"Cannot read milliseconds_before_first_turn in GAME_STARTS")
		assert.InEpsilon(t, expectedMsBeforeFirstTurn, msBeforeFirstTurn,
			1e-3, "Unexpected value for milliseconds_before_first_turn "+
				"in GAME_STARTS message")

		msBetweenTurns, err := readFloat(msg,
			"milliseconds_before_first_turn")
		assert.NoError(t, err,
			"Cannot read milliseconds_before_first_turn in GAME_STARTS")
		assert.InEpsilon(t, expectedMsBetweenTurns, msBetweenTurns,
			1e-3, "Unexpected value for milliseconds_before_first_turn "+
				"in GAME_STARTS message")

		checkPlayersInfo(t, msg, expectedNbPlayers, expectedNbSpecialPlayers, isPlayer)
		return playerID
	case "KICK":
		kickReason, err := netorcai.ReadString(msg, "kick_reason")
		assert.NoError(t, err, "Cannot read kick_reason")

		assert.FailNow(t, "Expected GAME_STARTS, got KICK", kickReason)
	default:
		assert.FailNowf(t, "Expected GAME_STARTS, got another message type",
			messageType)
	}
	return -2
}

func checkTurn(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedNbSpecialPlayers, expectedTurnNumber int, isPlayer bool) int {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err, "Cannot read 'message_type' field in "+
		"received client message (TURN)")

	switch messageType {
	case "TURN":
		turnNumber, err := netorcai.ReadInt(msg, "turn_number")
		assert.NoError(t, err, "Cannot read turn_number in TURN")
		assert.Equal(t, expectedTurnNumber, turnNumber,
			"Unexpected value for turn_number in received TURN message")

		_, err = netorcai.ReadObject(msg, "game_state")
		assert.NoError(t, err, "Cannot read game_state in TURN")

		checkPlayersInfo(t, msg, expectedNbPlayers, expectedNbSpecialPlayers, isPlayer)
		return turnNumber
	case "KICK":
		kickReason, err := netorcai.ReadString(msg, "kick_reason")
		assert.NoError(t, err, "Cannot read kick_reason")

		assert.FailNow(t, "Expected TURN, got KICK", kickReason)
	default:
		assert.FailNowf(t, "Expected TURN, got another message type",
			messageType)
	}
	return expectedTurnNumber
}

func checkTurnPotentialTurnsSkipped(t *testing.T, msg map[string]interface{},
	expectedNbPlayers, expectedNbSpecialPlayers, expectedMinimalTurnNumber int, isPlayer bool) int {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err, "Cannot read 'message_type' field in "+
		"received client message (TURN or GAME_ENDS)")

	switch messageType {
	case "TURN":
		turnNumber, err := netorcai.ReadInt(msg, "turn_number")
		assert.NoError(t, err, "Cannot read turn_number in TURN")
		assert.Condition(t, func() bool {
			return turnNumber >= expectedMinimalTurnNumber
		})

		_, err = netorcai.ReadObject(msg, "game_state")
		assert.NoError(t, err, "Cannot read game_state in TURN")

		checkPlayersInfo(t, msg, expectedNbPlayers, expectedNbSpecialPlayers, isPlayer)
		return turnNumber
	case "GAME_ENDS":
		_, err := netorcai.ReadInt(msg, "winner_player_id")
		assert.NoError(t, err, "Cannot read winner_player_id in GAME_ENDS")

		_, err = netorcai.ReadObject(msg, "game_state")
		assert.NoError(t, err, "Cannot read game_state in GAME_ENDS")
	case "KICK":
		kickReason, err := netorcai.ReadString(msg, "kick_reason")
		assert.NoError(t, err, "Cannot read kick_reason")

		assert.FailNow(t, "Expected (TURN or GAME_ENDS), got KICK", kickReason)
	default:
		assert.FailNowf(t, "Expected (TURN or GAME_ENDS), got another message type",
			messageType)
	}
	return expectedMinimalTurnNumber
}

func checkGameEnds(t *testing.T, msg map[string]interface{}) {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err, "Cannot read 'message_type' field in "+
		"received client message (GAME_ENDS)")

	switch messageType {
	case "GAME_ENDS":
		_, err := netorcai.ReadInt(msg, "winner_player_id")
		assert.NoError(t, err, "Cannot read winner_player_id in GAME_ENDS")

		_, err = netorcai.ReadObject(msg, "game_state")
		assert.NoError(t, err, "Cannot read game_state in GAME_ENDS")
	case "KICK":
		kickReason, err := netorcai.ReadString(msg, "kick_reason")
		assert.NoError(t, err, "Cannot read kick_reason")

		assert.FailNow(t, "Expected GAME_ENDS, got KICK", kickReason)
	default:
		assert.FailNowf(t, "Expected GAME_ENDS, got another message type",
			messageType)
	}
}

func killNetorcaiGently(proc *NetorcaiProcess, timeoutMS int) error {
	killallNetorcai()

	_, err := waitCompletionTimeout(proc.completion, timeoutMS)
	return err
}
