package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

// Netorcai helpers
func runNetorcaiWaitListening(t *testing.T) (nocIC, nocOC chan string,
	nocCompletion chan int) {
	args := []string{}
	nocIC = make(chan string)
	nocOC = make(chan string)
	nocCompletion = make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")

	return nocIC, nocOC, nocCompletion
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
	client Client, err error) {
	err = client.Connect("localhost", 4242)
	assert.NoError(t, err, "Cannot connect")

	err = client.SendLogin(role, nickname)
	assert.NoError(t, err, "Cannot send LOGIN")

	msg, err := waitReadMessage(&client, 1000)
	assert.NoError(t, err, "Cannot read message")
	checkLoginAck(t, msg)
	return client, err
}

func checkKick(t *testing.T, msg map[string]interface{},
	reasonMatcher *regexp.Regexp) {
	assert.Equal(t, "KICK", msg["message_type"].(string),
		"Unexpected message type")
	assert.Regexp(t, reasonMatcher, msg["kick_reason"].(string),
		"Unexpected kick reason")
}

func checkLoginAck(t *testing.T, msg map[string]interface{}) {
	switch msgType := msg["message_type"].(string); msgType {
	case "LOGIN_ACK":
	case "KICK":
		assert.Failf(t, "Expected LOGIN_ACK, got KICK",
			msg["kick_reason"].(string))
	default:
		assert.Failf(t, "Expected LOGIN_ACK", msgType)
	}
}
