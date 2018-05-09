package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestLoginNotJson(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	var client Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`definitely not JSON`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkKick(t, msg, regexp.MustCompile("Non-JSON"))
}

func TestLoginNoMessageType(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	var client Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"nickname":"bot", "role":"player"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkKick(t, msg, regexp.MustCompile("Field 'message_type' is missing"))
}

func TestLoginNoRole(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	var client Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "nickname":"bot"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkKick(t, msg, regexp.MustCompile("Field 'role' is missing"))
}

func TestLoginNoNickname(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	var client Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkKick(t, msg, regexp.MustCompile("Field 'nickname' is missing"))
}

func TestLoginBadRole(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	var client Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"¿Qué?", "nickname":"bot"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkKick(t, msg, regexp.MustCompile("Invalid role"))
}

func TestLoginBadNicknameShort(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	var client Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":""}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkKick(t, msg, regexp.MustCompile("Invalid nickname"))
}

func TestLoginBadNicknameLong(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	var client Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":"1234567890a"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkKick(t, msg, regexp.MustCompile("Invalid nickname"))
}

func TestLoginBadNicknameBadCharacters(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	var client Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":"hi world"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkKick(t, msg, regexp.MustCompile("Invalid nickname"))
}

/************
 * LOGIN ok *
 ************/

func TestLoginPlayerAscii(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	player, err := connectClient(t, "player", "player", 1000)
	assert.NoError(t, err, "Cannot connect client")
	defer player.Disconnect()
}

func TestLoginPlayerArabic(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	player, err := connectClient(t, "player", "لاعب", 1000)
	assert.NoError(t, err, "Cannot connect client")
	defer player.Disconnect()
}

func TestLoginPlayerJapanese(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	player, err := connectClient(t, "player", "プレーヤー", 1000)
	assert.NoError(t, err, "Cannot connect client")
	defer player.Disconnect()
}

/**************************
 * More complex scenarios *
 **************************/

func subtestLoginMaxNbClientSequential(t *testing.T, loginRole string,
	nbConnections, expectedNbLogged int, kickReasonMatcher *regexp.Regexp) {
	proc := runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	// Do many player connections sequentially
	var clients []*Client
	nbLogged := 0

	assert.Condition(t, func() bool {
		return expectedNbLogged <= nbConnections
	})
	for i := 0; i < nbConnections; i++ {
		client := &Client{}
		err := client.Connect("localhost", 4242)
		assert.NoError(t, err, "Cannot connect")
		clients = append(clients, client)

		err = client.SendLogin(loginRole, "клиент")
		assert.NoError(t, err, "Cannot send LOGIN")

		msg, err := waitReadMessage(client, 1000)
		assert.NoError(t, err, "Cannot read message")

		if i < expectedNbLogged {
			checkLoginAck(t, msg)
			nbLogged += 1
		} else {
			checkKick(t, msg, kickReasonMatcher)
		}
	}

	// Make sure the expected number of clients could LOGIN successfully
	assert.Equal(t, expectedNbLogged, nbLogged,
		"Unexpected number of logged players")

	// Close all client sockets
	for i, client := range clients {
		err := client.Disconnect()
		if i < expectedNbLogged {
			assert.NoError(t, err, "Disconnection of connected client failed")
		}
	}

	// Wait netorcai awareness of the disconnection
	for i := 0; i < expectedNbLogged; i++ {
		_, err := waitOutputTimeout(
			regexp.MustCompile(`Remote endpoint closed`), proc.outputControl,
			500, false)
		assert.NoError(t, err, "Could not read disconnection discovery in netorcai output")
	}

	// Connect the expected number of players
	for i := 0; i < expectedNbLogged; i++ {
		_, err := connectClient(t, loginRole, "клиент", 1000)
		assert.NoError(t, err, "Cannot connect client")
	}
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
	proc := runNetorcaiWaitListening(t)
	defer killallNetorcaiSIGKILL()

	// Do many client connections in parallel
	var clients []*Client
	nbLogged := 0
	nbDisconnectSuccess := 0

	assert.Condition(t, func() bool {
		return expectedNbLogged <= nbConnections
	})

	clientsChan := make(chan *Client, nbConnections)
	clientLogged := make(chan int, nbConnections)
	defer close(clientsChan)
	defer close(clientLogged)

	for i := 0; i < nbConnections; i++ {
		go func() {
			client := &Client{}
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
			default:
				assert.Fail(t, "Unexpected message type %v", msgType)
			}
		}()
	}

	timeoutChan := make(chan int, 2)
	defer close(timeoutChan)
	go func(c chan int) {
		time.Sleep(1000 * time.Millisecond)
		c <- 0
	}(timeoutChan)

	// Wait for all clients to finish their connection procedure
	for i := 0; i < 2*nbConnections; i++ {
		select {
		case client := <-clientsChan:
			clients = append(clients, client)
		case logResult := <-clientLogged:
			nbLogged += logResult
		case <-timeoutChan:
			assert.FailNow(t, "Timeout reached while waiting all clients to finish their connection procedure (first phase)")
		}
	}

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

	timeoutChan2 := make(chan int, 2)
	defer close(timeoutChan2)
	go func(c chan int) {
		time.Sleep(1000 * time.Millisecond)
		c <- 0
	}(timeoutChan)

	// Wait for all clients to finish their connection procedure
	for i := 0; i < 2*expectedNbLogged; i++ {
		select {
		case client := <-clientsChan:
			clients = append(clients, client)
		case logResult := <-clientLogged:
			nbLogged += logResult
		case <-timeoutChan2:
			assert.FailNow(t, "Timeout reached while waiting all clients to finish their connection procedure (second phase)")
		}
	}

	assert.Equal(t, expectedNbLogged, nbLogged,
		"Unexpected number of logged clients")

	for _, client := range clients {
		err := client.Disconnect()
		assert.NoError(t, err, "Logged client could not disconnect")
	}

	killallNetorcai()
	_, err := waitOutputTimeout(regexp.MustCompile(`Closing listening socket`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Could not read `Closing listening socket` in netorcai output")
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
