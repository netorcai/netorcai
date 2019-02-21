package test

import (
	"github.com/netorcai/netorcai"
	"github.com/netorcai/netorcai/client/go"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
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
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Non-JSON"))

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
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Field 'message_type' is missing"))

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
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Field 'role' is missing"))

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
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Field 'nickname' is missing"))

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
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Non-string value for field 'role'"))

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
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Invalid role"))

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
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Invalid nickname"))

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
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Invalid nickname"))

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
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Invalid nickname"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginNoMetaprotocolVersion(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":"valid"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Field 'metaprotocol_version' is missing"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginBadMetaprotocolVersionNotString(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":"valid", "metaprotocol_version": false}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Non-string value for field 'metaprotocol_version'"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginBadMetaprotocolVersionNotSemver(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":"valid", "metaprotocol_version": "42"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Invalid metaprotocol version: Not MAJOR.MINOR.PATCH"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginBadMetaprotocolVersionDifferentMajor(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendString(`{"message_type":"LOGIN", "role":"player", "nickname":"valid", "metaprotocol_version": "0.1.0"}`)
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, "InvalidClient", regexp.MustCompile("Metaprotocol version mismatch. Major version must be identical"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

/************
 * LOGIN ok *
 ************/

func TestLoginPlayerAscii(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	player, err := connectClient(t, "player", "player", netorcai.Version, 1000)
	assert.NoError(t, err, "Cannot connect client")
	player.Disconnect()

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginPlayerArabic(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	player, err := connectClient(t, "player", "لاعب", netorcai.Version, 1000)
	assert.NoError(t, err, "Cannot connect client")
	player.Disconnect()

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestLoginPlayerJapanese(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	player, err := connectClient(t, "player", "プレーヤー", netorcai.Version, 1000)
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
	proc := runNetorcaiWaitListening(t, []string{"--nb-splayers-max=2"})
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

		err = client.SendLogin(loginRole, "клиент", netorcai.Version)
		assert.NoError(t, err, "Cannot send LOGIN")

		msg, err := waitReadMessage(client, 1000)
		assert.NoError(t, err, "Cannot read client message (LOGIN_ACK|KICK)")

		if i < expectedNbLogged {
			checkLoginAck(t, msg)
			clients = append(clients, client)
		} else {
			checkKick(t, msg, loginRole, kickReasonMatcher)
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
			_, err := connectClient(t, loginRole, "клиент", netorcai.Version, 1000)
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

func TestLoginMaxNbSpecialPlayerSequential(t *testing.T) {
	subtestLoginMaxNbClientSequential(t, "special player", 100, 2,
		regexp.MustCompile(`Maximum number of special players reached`))
}

func TestLoginMaxNbVisuSequential(t *testing.T) {
	subtestLoginMaxNbClientSequential(t, "visualization", 100, 1,
		regexp.MustCompile(`Maximum number of visus reached`))
}

func TestLoginMaxNbGameLogicSequential(t *testing.T) {
	subtestLoginMaxNbClientSequential(t, "game logic", 100, 1,
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

	err = client.SendLogin(loginRole, "client", netorcai.Version)
	assert.NoError(t, err, "Cannot send LOGIN")

	if shouldConnect {
		msg, err := waitReadMessage(client, 1000)
		assert.NoError(t, err, "Cannot read client message (LOGIN_ACK)")
		checkLoginAck(t, msg)
	} else {
		msg, err := waitReadMessage(client, 1000)
		assert.NoError(t, err, "Cannot read client message (KICK)")
		checkKick(t, msg, loginRole,
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
