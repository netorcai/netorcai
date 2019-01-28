package test

import (
	"github.com/netorcai/netorcai/client/go"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestLoginNotJson(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`definitely not JSON`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, regexp.MustCompile("Non-JSON"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginNoMessageType(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"nickname":"bot", "role":"player"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, regexp.MustCompile("Field 'message_type' is missing"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginNoRole(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "nickname":"bot"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, regexp.MustCompile("Field 'role' is missing"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginNoNickname(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, regexp.MustCompile("Field 'nickname' is missing"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginRoleNotString(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":1, "nickname":"bot"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, regexp.MustCompile("Non-string value for field 'role'"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginBadRole(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"¿Qué?", "nickname":"bot"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, regexp.MustCompile("Invalid role"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginBadNicknameShort(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":""}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, regexp.MustCompile("Invalid nickname"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginBadNicknameLong(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":"1234567890a"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, regexp.MustCompile("Invalid nickname"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginBadNicknameBadCharacters(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":"hi world"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, regexp.MustCompile("Invalid nickname"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

/************
 * LOGIN ok *
 ************/

func TestLoginPlayerAscii(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	player, err := connectClient(t, "player", "player", 1000)
	assert.NoError(t, err, "Cannot connect client")
	player.Disconnect()

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginPlayerArabic(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	player, err := connectClient(t, "player", "لاعب", 1000)
	assert.NoError(t, err, "Cannot connect client")
	player.Disconnect()

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginPlayerJapanese(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	player, err := connectClient(t, "player", "プレーヤー", 1000)
	assert.NoError(t, err, "Cannot connect client")
	player.Disconnect()

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

/**************************
 * More complex scenarios *
 **************************/

func subtestLoginMaxNbClientSequential(t *testing.T, loginRole string,
	nbConnections, expectedNbLogged int, kickReasonMatcher *regexp.Regexp) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	// Do many player connections sequentially
	var clients []*client.Client

	assert.Condition(t, func() bool {
		return expectedNbLogged <= nbConnections
	})
	for i := 0; i < nbConnections; i++ {
		client := &client.Client{}
		err := client.Connect("localhost", 4242)
		assert.NoError(t, err, "Cannot connect")

		err = client.SendLogin(loginRole, "клиент")
		assert.NoError(t, err, "Cannot send LOGIN")

		msg, err := waitReadMessage(client, 1000)
		assert.NoError(t, err, "Cannot read client message (LOGIN_ACK|KICK)")

		if i < expectedNbLogged {
			checkLoginAck(t, msg)
			clients = append(clients, client)
		} else {
			checkKick(t, msg, kickReasonMatcher)
			err = client.Disconnect()
			assert.NoError(t, err, "Kicked client could not disconnect")
		}
	}

	// Make sure the expected number of clients could LOGIN successfully
	assert.Equal(t, expectedNbLogged, len(clients),
		"Unexpected number of logged players")

	// Close all client sockets
	for _, client := range clients {
		err := client.Disconnect()
		assert.NoError(t, err, "Logged client could not disconnect")

		// Check netorcai awareness of the disconnection
		_, err = waitOutputTimeout(
			regexp.MustCompile(`Remote endpoint closed`), proc.outputControl,
			500, false)
		assert.NoError(t, err,
			"Could not read disconnection discovery in netorcai output")
	}

	if loginRole != "game logic" {
		// Connect the expected number of clients
		for i := 0; i < expectedNbLogged; i++ {
			_, err := connectClient(t, loginRole, "клиент", 1000)
			assert.NoError(t, err, "Cannot connect client")
		}
	}

	err := killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "First instance could not be killed gently")
}

func TestLoginMaxNbPlayerSequential(t *testing.T) {
	subtestLoginMaxNbClientSequential(t, "player", 100, 4,
		regexp.MustCompile(`Maximum number of players reached`))
}

func TestLoginMaxNbVisuSequential(t *testing.T) {
	subtestLoginMaxNbClientSequential(t, "visualization", 100, 1,
		regexp.MustCompile(`Maximum number of visus reached`))
}

func TestLoginMaxNbGameLogicSequential(t *testing.T) {
	subtestLoginMaxNbClientSequential(t, "game logic", 100, 1,
		regexp.MustCompile(`A game logic is already logged in`))
}

func subtestLoginMaxNbClientParallel(t *testing.T, loginRole string,
	nbConnections, expectedNbLogged int, kickReasonMatcher *regexp.Regexp) {
	if isTravis() {
		// Do not run this test on Travis
		// Cclients would be too slow to react between KICK and socket close.
		t.SkipNow()
	}

	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	// Do many client connections in parallel
	var clients []*client.Client
	nbLogged := 0
	nbDisconnectSuccess := 0

	assert.Condition(t, func() bool {
		return expectedNbLogged <= nbConnections
	})

	clientsChan := make(chan *client.Client, nbConnections)
	clientLogged := make(chan int, nbConnections)
	defer close(clientsChan)
	defer close(clientLogged)

	for i := 0; i < nbConnections; i++ {
		go func() {
			client := &client.Client{}
			err := client.Connect("localhost", 4242)
			assert.NoError(t, err, "Client cannot connect")
			clientsChan <- client

			err = client.SendLogin(loginRole, "клиент")
			assert.NoError(t, err, "Client cannot send LOGIN")

			msg, err := waitReadMessage(client, 1000)
			assert.NoError(t, err, "Client cannot read LOGIN_ACK|KICK")
			switch msgType := msg["message_type"].(string); msgType {
			case "LOGIN_ACK":
				clientLogged <- 1
			case "KICK":
				clientLogged <- 0
				assert.Regexp(t, kickReasonMatcher,
					msg["kick_reason"].(string),
					"Client kicked for an unexpected reason")
				client.Disconnect()
			default:
				assert.Fail(t, "Unexpected message type %v", msgType)
			}
		}()
	}

	timeoutReached := make(chan int)
	stopTimeout := make(chan int)
	defer close(timeoutReached)
	defer close(stopTimeout)
	go func(timeout, stop chan int) {
		select {
		case <-stopTimeout:
		case <-time.After(1000 * time.Millisecond):
			timeoutReached <- 0
		}
	}(timeoutReached, stopTimeout)

	// Wait for all clients to finish their connection procedure
	for i := 0; i < 2*nbConnections; i++ {
		select {
		case client := <-clientsChan:
			clients = append(clients, client)
		case logResult := <-clientLogged:
			nbLogged += logResult
		case <-timeoutReached:
			assert.FailNow(t, "Timeout reached while waiting all clients to finish their connection procedure (first phase)")
		}
	}

	stopTimeout <- 0

	// Make sure the right number of clients could LOGIN successfully
	assert.Equal(t, expectedNbLogged, nbLogged,
		"Unexpected number of logged clients")

	// Disconnect all clients
	for _, client := range clients {
		err := client.Disconnect()
		if err == nil {
			nbDisconnectSuccess += 1
		}
	}
	assert.Condition(t, func() bool {
		return nbDisconnectSuccess >= expectedNbLogged
	})
	clients = clients[:0]
	nbLogged = 0

	// Wait netorcai awareness of the disconnection
	for i := 0; i < expectedNbLogged; i++ {
		_, err := waitOutputTimeout(
			regexp.MustCompile(`Remote endpoint closed`), proc.outputControl,
			500, false)
		assert.NoError(t, err, "Could not read disconnection discovery in netorcai output")
	}

	// Connect the expected number of players
	for i := 0; i < expectedNbLogged; i++ {
		go func() {
			player, err := connectClient(t, loginRole, "клиент", 1000)
			assert.NoError(t, err, "Cannot connect client")
			clientsChan <- player
			clientLogged <- 1
		}()
	}

	go func(timeout, stop chan int) {
		select {
		case <-stopTimeout:
		case <-time.After(1000 * time.Millisecond):
			timeoutReached <- 0
		}
	}(timeoutReached, stopTimeout)

	// Wait for all clients to finish their connection procedure
	for i := 0; i < 2*expectedNbLogged; i++ {
		select {
		case client := <-clientsChan:
			clients = append(clients, client)
		case logResult := <-clientLogged:
			nbLogged += logResult
		case <-timeoutReached:
			assert.FailNow(t, "Timeout reached while waiting all clients to finish their connection procedure (second phase)")
		}
	}

	stopTimeout <- 0

	assert.Equal(t, expectedNbLogged, nbLogged,
		"Unexpected number of logged clients")

	for _, client := range clients {
		err := client.Disconnect()
		assert.NoError(t, err, "Logged client could not disconnect")
	}

	// Kill netorcai and wait for its port to be freed
	killallNetorcai()
	_, err := waitOutputTimeout(regexp.MustCompile(`Closing listening socket`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Could not read `Closing listening socket` in netorcai output")

	_, expRetCode := handleCoverage(t, 1)
	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestLoginMaxNbPlayerParallel(t *testing.T) {
	subtestLoginMaxNbClientParallel(t, "player", 100, 4,
		regexp.MustCompile(`Maximum number of players reached`))
}

func TestLoginMaxNbVisuParallel(t *testing.T) {
	subtestLoginMaxNbClientParallel(t, "visualization", 100, 1,
		regexp.MustCompile(`Maximum number of visus reached`))
}

func TestLoginMaxNbGameLogicParallel(t *testing.T) {
	subtestLoginMaxNbClientParallel(t, "game logic", 100, 1,
		regexp.MustCompile(`A game logic is already logged in`))
}

func subtestLoginGameAlreadyStarted(t *testing.T, loginRole string,
	shouldConnect bool) {
	proc, _, _, _, _, _ := runNetorcaiAndClients(t,
		[]string{}, 1000, 0, 0, 0)
	defer killallNetorcaiSIGKILL()

	waitOutputTimeout(regexp.MustCompile(`Game logic accepted`),
		proc.outputControl, 1000, false)

	proc.inputControl <- `start`
	waitOutputTimeout(regexp.MustCompile(`Game started`), proc.outputControl,
		1000, true)

	client := &client.Client{}
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")

	err = client.SendLogin(loginRole, "client")
	assert.NoError(t, err, "Cannot send LOGIN")

	if shouldConnect {
		msg, err := waitReadMessage(client, 1000)
		assert.NoError(t, err, "Cannot read client message (LOGIN_ACK)")
		checkLoginAck(t, msg)
	} else {
		msg, err := waitReadMessage(client, 1000)
		assert.NoError(t, err, "Cannot read client message (KICK)")
		checkKick(t, msg,
			regexp.MustCompile(`LOGIN denied: Game has been started`))
	}

	proc.inputControl <- `quit`
	waitCompletionTimeout(proc.completion, 1000)
}

func TestLoginPlayerGameAlreadyStarted(t *testing.T) {
	subtestLoginGameAlreadyStarted(t, "player", false)
}

func TestLoginVisuGameAlreadyStarted(t *testing.T) {
	subtestLoginGameAlreadyStarted(t, "visualization", true)
}

func TestLoginGLGameAlreadyStarted(t *testing.T) {
	subtestLoginGameAlreadyStarted(t, "game logic", false)
}
