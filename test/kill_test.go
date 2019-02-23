package test

import (
	"fmt"
	"github.com/netorcai/netorcai/client/go"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestKickallOnAbortKillSigterm(t *testing.T) {
	proc, clients, _, _, _, _ := runNetorcaiAndAllClients(t, []string{}, 1000, 0)
	defer killallNetorcaiSIGKILL()

	killallNetorcai()

	checkAllKicked(t, clients, regexp.MustCompile(`netorcai abort`), 1000)

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	_, expRetCode := handleCoverage(t, 1)
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestKickallOnAbortKillSigtermSpecial(t *testing.T) {
	proc, clients, _, _, _, _ := runNetorcaiAndAllClients(t, []string{"--nb-splayers-max=1"}, 1000, 1)
	defer killallNetorcaiSIGKILL()

	killallNetorcai()

	checkAllKicked(t, clients, regexp.MustCompile(`netorcai abort`), 1000)

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	_, expRetCode := handleCoverage(t, 1)
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

type GLKillEvent int

const (
	OnDoTurnReception GLKillEvent = 0
)

func subtestKillDuringGame(t *testing.T, netorcaiArgs []string,
	glKillEvent GLKillEvent,
	nbTurns int,
	msBeforeFirstTurn, msBetweenTurns float64) {
	proc, _, players, _, visus, gls := runNetorcaiAndAllClients(t,
		append([]string{
			fmt.Sprintf("--delay-first-turn=%v", msBeforeFirstTurn),
			fmt.Sprintf("--delay-turns=%v", msBetweenTurns),
			fmt.Sprintf("--nb-turns-max=%v", nbTurns)},
			netorcaiArgs...),
		1000, 0)
	defer killallNetorcaiSIGKILL()

	// Disconnect visus
	for _, visu := range visus {
		visu.Disconnect()
		waitOutputTimeout(regexp.MustCompile(`Remote endpoint closed`),
			proc.outputControl, 1000, false)
	}

	clientFinished := make(chan int, 5)
	// Game logic
	for _, glClient := range gls {
		go func(netorcai *NetorcaiProcess, gl *client.Client, onexit chan int) {
			// Read DO_INIT
			msg, err := waitReadMessage(gl, 1000)
			assert.NoError(t, err, "GL could not read message (DO_INIT)")
			checkDoInit(t, msg, 4, 0, nbTurns)

			// Answer DO_INIT_ACK
			doInitAck := DefaultHelloGLDoInitAck(4, 0, nbTurns)
			err = gl.SendString(doInitAck)
			assert.NoError(t, err, "GL could not send DO_INIT_ACK")

			// Read DO_TURN
			msg, err = waitReadMessage(gl, 1000)
			assert.NoError(t, err, "GL could not read message (DO_TURN)")
			checkDoTurn(t, msg, 4, 0, nbTurns)

			if glKillEvent == OnDoTurnReception {
				// Instead of answering DO_TURN_ACK, kill netorcai gently.
				err = killNetorcaiGently(proc, 1000)
				assert.NoError(t, err, "Netorcai could not be killed gently")

				// Read KICK
				msg, err = waitReadMessage(gl, 1000)
				assert.NoError(t, err, "GL could not read message (KICK)")
				checkKick(t, msg, "GL", regexp.MustCompile(`netorcai abort`))
			}
			onexit <- 1
		}(proc, glClient, clientFinished)
	}

	// Players: Expect KICK after GAME_STARTS
	for playerID, playerClient := range players {
		go func(clientName string, player *client.Client, onexit chan int) {
			// Read GAME_STARTS
			msg, err := waitReadMessage(player, 1000)
			assert.NoError(t, err, "%v could not read message (GAME_STARTS)", clientName)
			checkGameStarts(t, msg, 4, 0, nbTurns, msBeforeFirstTurn, msBetweenTurns, true)

			// Read KICK
			msg, err = waitReadMessage(player, 1000)
			assert.NoError(t, err, "%v could not read message (KICK)", clientName)
			checkKick(t, msg, clientName, regexp.MustCompile(`netorcai abort`))
			onexit <- 1
		}(fmt.Sprintf("Player%v", playerID), playerClient, clientFinished)
	}

	proc.inputControl <- `start`
	waitOutputTimeout(regexp.MustCompile(`Game started`), proc.outputControl, 1000, true)

	// Wait until completion of all clients or timeout
	timeoutReached := make(chan int)
	stopTimeout := make(chan int)
	defer close(timeoutReached)
	defer close(stopTimeout)
	go func() {
		select {
		case <-stopTimeout:
		case <-time.After(time.Duration(1000) * time.Millisecond):
			timeoutReached <- 0
		}
	}()
	nbClientFinished := 0
	for nbClientFinished < 5 {
		select {
		case <-clientFinished:
			nbClientFinished++
		case <-timeoutReached:
			assert.FailNow(t, "Timeout reached while waiting clients' finition")
		}
	}
}

func TestKillDuringTurn(t *testing.T) {
	subtestKillDuringGame(t, nil, OnDoTurnReception, 3, 50, 50)
}

func TestKillDuringTurnFast(t *testing.T) {
	subtestKillDuringGame(t, []string{"--fast"}, OnDoTurnReception, 3, 50, 50)
}
