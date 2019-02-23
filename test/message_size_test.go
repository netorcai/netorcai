package test

import (
	"github.com/netorcai/netorcai"
	"github.com/netorcai/netorcai/client/go"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"regexp"
	"testing"
)

func RandomString(size int) string {
	alphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	buf := make([]rune, size)
	for i := range buf {
		buf[i] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(buf)
}

func TestFirstMessageTooBig(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendBytes([]byte(RandomString(1024-1)), false) // -1 for final '\n'
	assert.NoError(t, err, "Cannot send message")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, "InvalidClient",
		regexp.MustCompile("Received message size of first message is too big"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestSecondMessageTooBig(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	var client client.Client
	err := client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")
	defer client.Disconnect()

	err = client.SendLogin("player", "player", netorcai.Version)
	assert.NoError(t, err, "Cannot send LOGIN")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (LOGIN_ACK)")
	checkLoginAck(t, msg)

	err = client.SendBytes([]byte(RandomString(16777216-1)), false) // -1 for final '\n'

	msg, err = waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read client message (KICK)")
	checkKick(t, msg, "InvalidClient",
		regexp.MustCompile("Received message size is too big"))

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}
