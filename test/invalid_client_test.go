package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestInvalidClientLoginNotJson(t *testing.T) {
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

func TestInvalidClientLoginNoMessageType(t *testing.T) {
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

func TestInvalidClientLoginNoRole(t *testing.T) {
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

func TestInvalidClientLoginNoNickname(t *testing.T) {
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

func TestInvalidClientLoginBadRole(t *testing.T) {
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

func TestInvalidClientLoginBadNicknameShort(t *testing.T) {
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

func TestInvalidClientLoginBadNicknameLong(t *testing.T) {
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

func TestInvalidClientLoginBadNicknameBadCharacters(t *testing.T) {
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
