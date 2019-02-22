package test

import (
	"fmt"
	"github.com/netorcai/netorcai/client/go"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

type ClientDisconnectionTurnFunc func(int) int

func ClientDisconnectionWhenTurnIsGreaterThanPlayerID(playerID int) int {
	return playerID
}

func disconnectingClient(t *testing.T, client *client.Client, clientName string,
	nbPlayers, nbSpecialPlayers int,
	nbTurnsGL int,
	isPlayer bool,
	msBeforeFirstTurn, msBetweenTurns float64,
	clientDisconnectionTurnFunc ClientDisconnectionTurnFunc) {

	msg, err := waitReadMessage(client, 2000)
	assert.NoError(t, err, "%v could not read message (GAME_STARTS)", clientName)
	playerID := checkGameStarts(t, msg, nbPlayers, nbSpecialPlayers, nbTurnsGL,
		msBeforeFirstTurn, msBetweenTurns, isPlayer)
	disconnectionTurn := clientDisconnectionTurnFunc(playerID)

	for turn := 0; turn < disconnectionTurn; turn += 1 {
		// Wait TURN
		msg, err := waitReadMessage(client, 2000)
		assert.NoError(t, err, "%v could not read message (TURN) %v/%v",
			clientName, turn, disconnectionTurn)
		turnReceived := checkTurn(t, msg, nbPlayers, nbSpecialPlayers, turn, isPlayer)

		// Send TURN_ACK
		if turn != disconnectionTurn-1 {
			data := DefaultHelloClientTurnAck(turnReceived, playerID)
			err = client.SendString(data)
			assert.NoError(t, err, "%s cannot send TURN_ACK", clientName)
		}
	}

	client.Disconnect()
}

func subtestDisconnectingClients(t *testing.T,
	netorcaiAdditionalArgs []string,
	nbPlayers, nbSpecialPlayers, nbVisus int,
	nbTurns int,
	playerDiscoTurnFunc, splayerDiscoTurnFunc, visuDiscoTurnFunc ClientDisconnectionTurnFunc) {

	proc, _, players, specialPlayers, visus, gl := runNetorcaiAndClients(
		t, append([]string{"--delay-first-turn=500",
			fmt.Sprintf("--nb-turns-max=%v", nbTurns),
			fmt.Sprintf("--nb-players-max=%v", nbPlayers),
			fmt.Sprintf("--nb-splayers-max=%v", nbSpecialPlayers),
			fmt.Sprintf("--nb-visus-max=%v", nbVisus),
			"--delay-turns=500", "--debug", "--autostart"},
			netorcaiAdditionalArgs...),
		1000, nbPlayers, nbSpecialPlayers, nbVisus)
	defer killallNetorcaiSIGKILL()

	// Run a game client
	go helloGameLogic(t, gl[0], nbPlayers, nbSpecialPlayers, nbTurns, nbTurns,
		DefaultHelloGLCheckDoTurn, DefaultHelloGLDoInitAck, DefaultHelloGlDoTurnAck,
		regexp.MustCompile(`Game is finished`))

	// Run player clients
	for playerID, player := range players {
		go disconnectingClient(t, player, fmt.Sprintf("Player%v", playerID),
			nbPlayers, nbSpecialPlayers, nbTurns, true, 500, 500, playerDiscoTurnFunc)
	}

	// Run special player clients
	for splayerID, splayer := range specialPlayers {
		go disconnectingClient(t, splayer, fmt.Sprintf("SpecialPlayer%v", splayerID),
			nbPlayers, nbSpecialPlayers, nbTurns, true, 500, 500, playerDiscoTurnFunc)
	}

	// Run visu clients
	for visuID, visu := range visus {
		go disconnectingClient(t, visu, fmt.Sprintf("Visu%v", visuID),
			nbPlayers, nbSpecialPlayers, nbTurns, false, 500, 500, playerDiscoTurnFunc)
	}

	// Wait for game end
	waitOutputTimeout(regexp.MustCompile(`Game is finished`),
		proc.outputControl, 5000, false)
	waitCompletionTimeout(proc.completion, 1000)
}

func TestPlayerDisconnectionDuringGame(t *testing.T) {
	subtestDisconnectingClients(t, nil, 4, 1, 1,
		7,
		ClientDisconnectionWhenTurnIsGreaterThanPlayerID,
		ClientDisconnectionWhenTurnIsGreaterThanPlayerID,
		ClientDisconnectionWhenTurnIsGreaterThanPlayerID)
}

func TestPlayerDisconnectionDuringGameFast(t *testing.T) {
	subtestDisconnectingClients(t, []string{"--fast"}, 4, 1, 1,
		7,
		ClientDisconnectionWhenTurnIsGreaterThanPlayerID,
		ClientDisconnectionWhenTurnIsGreaterThanPlayerID,
		ClientDisconnectionWhenTurnIsGreaterThanPlayerID)
}
