package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestKillKick(t *testing.T) {
	_ = runNetorcaiWaitListening(t)
	defer killallNetorcai()

	// Do client connections sequentially
	var clients []*Client

	// 4 players
	for i := 0; i < 4; i++ {
		player, err := connectClient(t, "player", "player", 1000)
		assert.NoError(t, err, "Cannot connect client")
		clients = append(clients, player)
	}

	// 1 visu
	for i := 0; i < 1; i++ {
		visu, err := connectClient(t, "visualization", "visu", 1000)
		assert.NoError(t, err, "Cannot connect client")
		clients = append(clients, visu)
	}

	// 1 game logic
	for i := 0; i < 1; i++ {
		gl, err := connectClient(t, "game logic", "game_logic", 1000)
		assert.NoError(t, err, "Cannot connect client")
		clients = append(clients, gl)
	}

	// Kill netorcai
	go func() {
		killallNetorcai()
	}()

	// All clients should receive a KICK
	kickChan := make(chan int, 6)
	for _, client := range clients {
		go func(c *Client) {
			msg, err := waitReadMessage(c, 1000)
			assert.NoError(t, err, "Cannot read message")
			checkKick(t, msg, regexp.MustCompile(`netorcai abort`))
			kickChan <- 0
		}(client)
	}

	for _, _ = range clients {
		<-kickChan
	}
}
