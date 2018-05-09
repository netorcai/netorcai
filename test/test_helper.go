package test

import (
	"fmt"
	"github.com/mpoquet/netorcai"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

// Netorcai helpers
func runNetorcaiWaitListening(t *testing.T) *NetorcaiProcess {
	coverFile, _ := handleCoverage(t, 0)
	args := []string{}

	proc, err := runNetorcaiCover(coverFile, args)
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
	go func() {
		time.Sleep(time.Duration(timeoutMS) * time.Millisecond)
		timeoutReached <- 0
	}()

	for {
		select {
		case line := <-output:
			if re.MatchString(line) {
				return line, nil
			} else {
				if leaveOnNonMatch {
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

// Client helpers
func waitReadMessage(client *Client, timeoutMS int) (
	msg map[string]interface{}, err error) {
	msgChan := make(chan int)
	go func() {
		msg, err = client.ReadMessage()
		msgChan <- 0
	}()

	select {
	case <-msgChan:
		return msg, err
	case <-time.After(time.Duration(timeoutMS) * time.Millisecond):
		return msg, fmt.Errorf("Timeout reached")
	}
}

func connectClient(t *testing.T, role, nickname string, timeoutMS int) (
	*Client, error) {
	client := &Client{}
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")

	err = client.SendLogin(role, nickname)
	assert.NoError(t, err, "Cannot send LOGIN")

	msg, err := waitReadMessage(client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkLoginAck(t, msg)
	return client, nil
}

func runNetorcaiAndAllClients(t *testing.T, timeoutMS int) (
	proc *NetorcaiProcess, clients, playerClients, visuClients,
	glClients []*Client) {
	proc = runNetorcaiWaitListening(t)

	// 4 players
	for i := 0; i < 4; i++ {
		player, err := connectClient(t, "player", "player", timeoutMS)
		if err != nil {
			killallNetorcai()
			assert.NoError(t, err, "Cannot connect client")
		}
		clients = append(clients, player)
		playerClients = append(playerClients, player)
	}

	// 1 visu
	for i := 0; i < 1; i++ {
		visu, err := connectClient(t, "visualization", "visu", timeoutMS)
		if err != nil {
			killallNetorcai()
			assert.NoError(t, err, "Cannot connect client")
		}
		clients = append(clients, visu)
		visuClients = append(visuClients, visu)
	}

	// 1 game logic
	for i := 0; i < 1; i++ {
		gl, err := connectClient(t, "game logic", "game_logic", timeoutMS)
		if err != nil {
			killallNetorcai()
			assert.NoError(t, err, "Cannot connect client")
		}
		clients = append(clients, gl)
		glClients = append(glClients, gl)
	}

	return proc, clients, playerClients, visuClients, glClients
}

func checkAllKicked(t *testing.T, clients []*Client,
	reasonMatcher *regexp.Regexp, timeoutMS int) {
	timeoutChan := make(chan int)
	go func() {
		time.Sleep(time.Duration(timeoutMS) * time.Millisecond)
		timeoutChan <- 0
	}()

	// All clients should receive a KICK
	kickChan := make(chan int, len(clients))
	for _, client := range clients {
		go func(c *Client) {
			msg, err := waitReadMessage(c, timeoutMS)
			assert.NoError(t, err, "Cannot read message")
			checkKick(t, msg, reasonMatcher)
			kickChan <- 0
		}(client)
	}

	for _ = range clients {
		select {
		case <-kickChan:
		case <-timeoutChan:
			assert.FailNow(t, "Timeout reached")
		}
	}
}

func checkKick(t *testing.T, msg map[string]interface{},
	reasonMatcher *regexp.Regexp) {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err, "Cannot read message_type")
	assert.Equal(t, "KICK", messageType, "Unexpected message type")

	kickReason, err := netorcai.ReadString(msg, "kick_reason")
	assert.NoError(t, err, "Cannot read kick_reason")
	assert.Regexp(t, reasonMatcher, kickReason, "Unexpected kick reason")
}

func checkLoginAck(t *testing.T, msg map[string]interface{}) {
	messageType, err := netorcai.ReadString(msg, "message_type")
	assert.NoError(t, err, "Cannot read message_type")

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
	assert.NoError(t, err, "Cannot read message_type")

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

		assert.FailNow(t, "Expected LOGIN_ACK, got KICK", kickReason)
	default:
		assert.FailNowf(t, "Expected LOGIN_ACK, got another message type",
			messageType)
	}
}
