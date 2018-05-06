package test

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"regexp"
	"testing"
)

func TestPromptStartNoClient(t *testing.T) {
	proc := runNetorcaiWaitListening(t)
	defer killallNetorcai()

	proc.inputControl <- "start"
	_, err := waitOutputTimeout(regexp.MustCompile(`Cannot start`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read line")
}

func TestPromptQuitNoClient(t *testing.T) {
	proc := runNetorcaiWaitListening(t)
	defer killallNetorcai()

	proc.inputControl <- "quit"
	_, err := waitOutputTimeout(regexp.MustCompile(`Shell exit`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read line")

	exitCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "Cannot wait netorcai completion")
	assert.Equal(t, 0, exitCode, "Invalid netorcai exit code")
}

func TestPromptQuitAllClient(t *testing.T) {
	proc := runNetorcaiWaitListening(t)
	defer killallNetorcai()

	proc.inputControl <- "quit"
	_, err := waitOutputTimeout(regexp.MustCompile(`Shell exit`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read line")

	exitCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "Cannot wait netorcai completion")
	assert.Equal(t, 0, exitCode, "Invalid netorcai exit code")
}

func TestControlProcessInputCatNoInut(t *testing.T) {
	cmd := exec.Command("cat")
	cmd.Args = []string{"cat"}

	stdinPipe, err := cmd.StdinPipe()
	assert.NoError(t, err, "Cannot get cat's command stdin pipe")

	err = cmd.Start()
	assert.NoError(t, err, "Cannot start cat command")

	err = stdinPipe.Close()
	assert.NoError(t, err, "Cannot close cat's stdin pipe")

	err = cmd.Wait()
	assert.NoError(t, err, "Could not wait cat's termination")
}

func TestControlProcessInputCatHelloWorld(t *testing.T) {
	cmd := exec.Command("cat")
	cmd.Args = []string{"cat"}

	stdinPipe, err := cmd.StdinPipe()
	assert.NoError(t, err, "Cannot get cat's command stdin pipe")
	writer := bufio.NewWriter(stdinPipe)

	stdoutPipe, err := cmd.StdoutPipe()
	assert.NoError(t, err, "Cannot get cat's command stdout pipe")
	reader := bufio.NewReader(stdoutPipe)

	err = cmd.Start()
	assert.NoError(t, err, "Cannot start cat command")

	text := "Hello world!\n"
	_, err = writer.WriteString(text)
	assert.NoError(t, err, "Cannot write text on cat's stdin pipe")

	err = writer.Flush()
	assert.NoError(t, err, "Cannot flush cat's stdin pipe")

	lineRead, err := reader.ReadString('\n')
	assert.NoError(t, err, "Cannot read on cat's stdout pipe")
	assert.Equal(t, text, lineRead, "Cat did not printed its input")

	err = stdinPipe.Close()
	assert.NoError(t, err, "Cannot close cat's stdin pipe")

	err = cmd.Wait()
	assert.NoError(t, err, "Could not wait cat's termination")
}
