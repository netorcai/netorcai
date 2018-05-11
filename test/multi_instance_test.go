package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestTwoInstancesSamePort(t *testing.T) {
	args := []string{"--port=5151"}
	coverFile, expectedExitCode2 := handleCoverage(t, 1)

	proc1, err := runNetorcaiCover("", args) // never covered
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcaiSIGKILL()

	_, err = waitListening(proc1.outputControl, 1000)
	assert.NoError(t, err, "First instance is not listening")

	proc2, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")

	_, err = waitOutputTimeout(
		regexp.MustCompile(`Cannot listen incoming connections`),
		proc2.outputControl, 1000, false)
	assert.NoError(t, err, "Second instance is listening")

	exitCode, err := waitCompletionTimeout(proc2.completion, 1000)
	assert.NoError(t, err, "Second instance has not completed")
	assert.Equal(t, expectedExitCode2, exitCode,
		"Second instance bad exit code")

	err = killNetorcaiGently(proc1, 1000)
	assert.NoError(t, err, "First instance could not be killed gently")
}

func TestTwoInstancesDifferentPort(t *testing.T) {
	args := []string{"--port=5151"}
	coverFile, _ := handleCoverage(t, 1)

	proc1, err := runNetorcaiCover("", args) // never covered
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcaiSIGKILL()

	_, err = waitListening(proc1.outputControl, 1000)
	assert.NoError(t, err, "First instance is not listening")

	args = []string{"--port=5252"}
	proc2, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")

	_, err = waitListening(proc2.outputControl, 1000)
	assert.NoError(t, err, "Second instance is not listening")

	err = killNetorcaiGently(proc1, 1000)
	assert.NoError(t, err, "First instance could not be killed gently")

	err = killNetorcaiGently(proc2, 1000)
	assert.NoError(t, err, "Second instance could not be killed gently")
}
