package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestLoginNotJson(t *testing.T) {
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

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
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

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
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

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
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

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
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

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
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

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
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

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
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

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

func TestScenarioLoginPlayerAscii(t *testing.T) {
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

	player, err := connectClient(t, "player", "player", 1000)
	assert.NoError(t, err, "Cannot connect client")
	defer player.Disconnect()
}

func TestScenarioLoginPlayerArabic(t *testing.T) {
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

	player, err := connectClient(t, "player", "لاعب", 1000)
	assert.NoError(t, err, "Cannot connect client")
	defer player.Disconnect()
}

func TestScenarioLoginPlayerJapanese(t *testing.T) {
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

	player, err := connectClient(t, "player", "プレーヤー", 1000)
	assert.NoError(t, err, "Cannot connect client")
	defer player.Disconnect()
}

/**************************
 * More complex scenarios *
 **************************/

func subtestLoginMaxNbClientSequential(t *testing.T, loginRole string,
	nbConnections, expectedNbLogged int) {
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

	// Do many player connections sequentially
	var clients []Client
	nbLogged := 0

	assert.Condition(t, func() bool {
		return expectedNbLogged <= nbConnections
	})
	for i := 0; i < nbConnections; i++ {
		var client Client
		err := client.Connect("localhost", 4242)
		assert.NoError(t, err, "Cannot connect")
		clients = append(clients, client)

		err = client.SendLogin(loginRole, "клиент")
		assert.NoError(t, err, "Cannot send LOGIN")

		msg, err := waitReadMessage(&client, 1000)
		assert.NoError(t, err, "Cannot read message")

		if i < expectedNbLogged {
			checkLoginAck(t, msg)
			nbLogged += 1
		} else {
			checkKick(t, msg, regexp.MustCompile("Maximum number of players reached"))
		}
	}

	// Make sure only 4 clients could LOGIN successfully
	assert.Equal(t, expectedNbLogged, nbLogged,
		"Unexpected number of logged players")

	// Close all client sockets
	for i, client := range clients {
		err := client.Disconnect()
		if i < expectedNbLogged {
			assert.NoError(t, err, "Disconnection of connected client failed")
		}
	}

	// Connect the expected number of players
	for i := 0; i < expectedNbLogged; i++ {
		_, err := connectClient(t, loginRole, "клиент", 1000)
		assert.NoError(t, err, "Cannot connect client")
	}
}

func LoginMaxNbPlayerSequential(t *testing.T) {
	subtestLoginMaxNbClientSequential(t, "player", 100, 4)
}

func LoginMaxNbVisusSequential(t *testing.T) {
	subtestLoginMaxNbClientSequential(t, "visualization", 100, 1)
}

func LoginMaxNbGameLogicSequential(t *testing.T) {
	subtestLoginMaxNbClientSequential(t, "game logic", 100, 1)
}

func subtestLoginMaxNbClientParallel(t *testing.T, loginRole string,
	nbConnections, expectedNbLogged int) {
	_, _, _ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

	// Do many client connections in parallel
	var clients []*Client
	nbLogged := 0
	nbDisconnectSuccess := 0

	assert.Condition(t, func() bool {
		return expectedNbLogged <= nbConnections
	})

	clientsChan := make(chan *Client, nbConnections)
	clientLogged := make(chan int, nbConnections)
	for i := 0; i < nbConnections; i++ {
		go func() {
			var client Client
			err := client.Connect("localhost", 4242)
			assert.NoError(t, err, "Cannot connect")
			clientsChan <- &client

			err = client.SendLogin(loginRole, "клиент")
			assert.NoError(t, err, "Cannot send LOGIN")

			msg, err := waitReadMessage(&client, 1000)
			assert.NoError(t, err, "Cannot read message")
			switch msgType := msg["message_type"].(string); msgType {
			case "LOGIN_ACK":
				clientLogged <- 1
			case "KICK":
				assert.Regexp(t,
					regexp.MustCompile(`Maximum number of \S* reached`),
					msg["kick_reason"].(string), "Unexpected kick reason")
				clientLogged <- 0
			default:
				assert.Fail(t, "Unexpected message type %v", msgType)
			}
		}()
	}

	// Wait for all clients to finish their connection procedure
	for i := 0; i < nbConnections; i++ {
		clients = append(clients, <-clientsChan)
		nbLogged = nbLogged + <-clientLogged
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

	// Connect the expected number of players
	for i := 0; i < expectedNbLogged; i++ {
		go func() {
			player, err := connectClient(t, loginRole, "клиент", 1000)
			assert.NoError(t, err, "Cannot connect client")
			clientsChan <- &player
			clientLogged <- 1
		}()
	}

	for i := 0; i < expectedNbLogged; i++ {
		clients = append(clients, <-clientsChan)
		nbLogged = nbLogged + <-clientLogged
	}

	assert.Equal(t, expectedNbLogged, nbLogged,
		"Unexpected number of logged players")

	for _, client := range clients {
		err := client.Disconnect()
		assert.NoError(t, err, "Could not disconnect")
	}
}

func LoginMaxNbPlayerParallel(t *testing.T) {
	subtestLoginMaxNbClientParallel(t, "player", 100, 4)
}

func LoginMaxNbVisusParallel(t *testing.T) {
	subtestLoginMaxNbClientParallel(t, "visualization", 100, 1)
}

func LoginMaxNbGameLogicParallel(t *testing.T) {
	subtestLoginMaxNbClientParallel(t, "game logic", 100, 1)
}
