package test

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"regexp"
	"strconv"
	"testing"
)

func promptReadValue(promptLine, variableName string) (string, error) {
	re := regexp.MustCompile(`\A` + variableName + `=(\d+)\z`)

	res := re.FindStringSubmatch(promptLine)
	if res == nil {
		return "", fmt.Errorf("No match")
	} else {
		strValue := res[1]
		return strValue, nil
	}
}

func TestPromptStartNoClient(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	proc.inputControl <- "start"
	_, err := waitOutputTimeout(regexp.MustCompile(`Cannot start`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read line")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestPromptDoubleStart(t *testing.T) {
	proc, _, _, _, _, _ := runNetorcaiAndAllClients(t, []string{}, 1000, 0)
	defer killallNetorcaiSIGKILL()

	proc.inputControl <- "start"
	proc.inputControl <- "start"
	_, err := waitOutputTimeout(
		regexp.MustCompile(`Game has already been started`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read line")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestPromptDoubleStartSpecial(t *testing.T) {
	proc, _, _, _, _, _ := runNetorcaiAndAllClients(t, []string{"--nb-splayers-max=1"}, 1000, 1)
	defer killallNetorcaiSIGKILL()

	proc.inputControl <- "start"
	proc.inputControl <- "start"
	_, err := waitOutputTimeout(
		regexp.MustCompile(`Game has already been started`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read line")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestPromptQuitNoClient(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	proc.inputControl <- "quit"
	_, err := waitOutputTimeout(regexp.MustCompile(`Shell exit`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read line")

	exitCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "Cannot wait netorcai completion")
	assert.Equal(t, 0, exitCode, "Invalid netorcai exit code")
}

func TestPromptQuitAllClient(t *testing.T) {
	proc, clients, _, _, _, _ := runNetorcaiAndAllClients(t, []string{}, 1000, 0)
	defer killallNetorcaiSIGKILL()

	proc.inputControl <- "quit"
	_, err := waitOutputTimeout(regexp.MustCompile(`Shell exit`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read line")

	checkAllKicked(t, clients, regexp.MustCompile(`netorcai abort`), 1000)
}

func TestPromptQuitAllClientSpecial(t *testing.T) {
	proc, clients, _, _, _, _ := runNetorcaiAndAllClients(t, []string{"--nb-splayers-max=1"}, 1000, 1)
	defer killallNetorcaiSIGKILL()

	proc.inputControl <- "quit"
	_, err := waitOutputTimeout(regexp.MustCompile(`Shell exit`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read line")

	checkAllKicked(t, clients, regexp.MustCompile(`netorcai abort`), 1000)
}

func subtestPromptVariablePrintSet(t *testing.T, variableName,
	invalidTypeValue, initialValue,
	tooSmallValue, okValue, tooBigValue string) {
	proc := runNetorcaiWaitListening(t, []string{})
	currentValue := initialValue
	defer killallNetorcaiSIGKILL()

	// Set invalid value (bad type)
	proc.inputControl <- "set " + variableName + "=" + invalidTypeValue
	line, err := waitOutputTimeout(regexp.MustCompile(`Bad VALUE`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err,
		"Cannot read prompt 'Bad VALUE' output (invalid type value)")

	// Initial value must still be there
	proc.inputControl <- "print " + variableName
	line, err = waitOutputTimeout(regexp.MustCompile(variableName+"="),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read prompt 'print' output (initial value)")
	value, err := promptReadValue(line, variableName)
	assert.NoError(t, err,
		"Cannot extract value from prompt print output (initial value)")
	assert.Equal(t, currentValue, value,
		"Unexpected value from prompt print output (initial value)")

	// Set a valid value, then check that the printed value is the expected one
	currentValue = okValue
	proc.inputControl <- "set " + variableName + "=" + okValue
	proc.inputControl <- "print " + variableName
	line, err = waitOutputTimeout(regexp.MustCompile(variableName+"="),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read prompt 'print' output (ok value)")
	value, err = promptReadValue(line, variableName)
	assert.NoError(t, err,
		"Cannot extract value from prompt print output (ok value)")
	assert.Equal(t, currentValue, value,
		"Unexpected value from prompt print output (ok value)")

	// Set invalid value (too small)
	proc.inputControl <- "set " + variableName + "=" + tooSmallValue
	line, err = waitOutputTimeout(regexp.MustCompile(`Bad VALUE`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err,
		"Cannot read prompt 'Bad VALUE' output (too small value)")

	// Set invalid value (too big)
	proc.inputControl <- "set " + variableName + "=" + tooBigValue
	line, err = waitOutputTimeout(regexp.MustCompile(`Bad VALUE`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err,
		"Cannot read prompt 'Bad VALUE' output (too big value)")

	// Previous value must still be there
	proc.inputControl <- "print " + variableName
	line, err = waitOutputTimeout(regexp.MustCompile(variableName+"="),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read prompt 'print' output (at end)")
	value, err = promptReadValue(line, variableName)
	assert.NoError(t, err,
		"Cannot extract value from prompt print output (at end)")
	assert.Equal(t, currentValue, value,
		"Unexpected value from prompt print output (at end)")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func subtestPromptIntVariablePrintSet(t *testing.T,
	variableName, invalidTypeValue string,
	initialValue, tooSmallValue, okValue, tooBigValue int) {
	subtestPromptVariablePrintSet(t, variableName, invalidTypeValue,
		strconv.Itoa(initialValue),
		strconv.Itoa(tooSmallValue),
		strconv.Itoa(okValue),
		strconv.Itoa(tooBigValue))
}

func TestPromptNbTurnsMax(t *testing.T) {
	subtestPromptIntVariablePrintSet(t, "nb-turns-max", "50.5", 100,
		0, 42, 65536)
}

func TestPromptNbPlayersMax(t *testing.T) {
	subtestPromptIntVariablePrintSet(t, "nb-players-max", "4.5", 4, 0, 2, 1025)
}

func TestPromptNbSpecialPlayersMax(t *testing.T) {
	subtestPromptIntVariablePrintSet(t, "nb-splayers-max", "4.5", 0, -1, 2, 1025)
}

func TestPromptNbVisusMax(t *testing.T) {
	subtestPromptIntVariablePrintSet(t, "nb-visus-max", "1.5", 1, -1, 10, 1025)
}

func subtestPromptFloatVariablePrintSet(t *testing.T,
	variableName, invalidTypeValue string,
	initialValue, tooSmallValue, okValue, tooBigValue float64) {
	subtestPromptVariablePrintSet(t, variableName, invalidTypeValue,
		strconv.FormatFloat(initialValue, 'f', -1, 64),
		strconv.FormatFloat(tooSmallValue, 'f', -1, 64),
		strconv.FormatFloat(okValue, 'f', -1, 64),
		strconv.FormatFloat(tooBigValue, 'f', -1, 64))
}

func TestPromptDelayFirstTurn(t *testing.T) {
	subtestPromptFloatVariablePrintSet(t, "delay-first-turn", "meh", 1000,
		49.999, 500, 10000.001)
}

func TestPromptDelayTurns(t *testing.T) {
	subtestPromptFloatVariablePrintSet(t, "delay-turns", "meh", 1000,
		49.999, 500, 10000.001)
}

func TestPromptPrintAll(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	proc.inputControl <- "print all"

	_, err := waitOutputTimeout(regexp.MustCompile(`nb-turns-max=100`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read print nb-turns-max")

	_, err = waitOutputTimeout(regexp.MustCompile(`nb-players-max=4`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read print nb-players-max")

	_, err = waitOutputTimeout(regexp.MustCompile(`nb-splayers-max=0`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read print nb-splayers-max")

	_, err = waitOutputTimeout(regexp.MustCompile(`nb-visus-max=1`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read print nb-visus-max")

	_, err = waitOutputTimeout(regexp.MustCompile(`delay-first-turn=1000`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read print delay-first-turn")

	_, err = waitOutputTimeout(regexp.MustCompile(`delay-turns=1000`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read print delay-turns")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestPromptPrintBadVariable(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	proc.inputControl <- "print unknown-var"
	_, err := waitOutputTimeout(regexp.MustCompile(`Bad VARIABLE=unknown-var`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read Bad VARIABLE")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestPromptSetBadVariable(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()

	proc.inputControl <- "set unknown-var=3"
	_, err := waitOutputTimeout(regexp.MustCompile(`Bad VARIABLE=unknown-var`),
		proc.outputControl, 1000, true)
	assert.NoError(t, err, "Cannot read Bad VARIABLE")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestPromptInvalidSyntaxPrint(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()
	re := regexp.MustCompile(`expected syntax: print VARIABLE`)

	proc.inputControl <- "print"
	_, err := waitOutputTimeout(re, proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read 'expected syntax [...]' after print")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestPromptInvalidSyntaxQuit(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()
	re := regexp.MustCompile(`expected syntax: quit`)

	proc.inputControl <- "quit meh"
	_, err := waitOutputTimeout(re, proc.outputControl, 1000, false)
	assert.NoError(t, err,
		"Cannot read 'expected syntax [...]' after quit meh")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestPromptInvalidSyntaxStart(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()
	re := regexp.MustCompile(`expected syntax: start`)

	proc.inputControl <- "start meh"
	_, err := waitOutputTimeout(re, proc.outputControl, 1000, false)
	assert.NoError(t, err,
		"Cannot read 'expected syntax [...]' after start meh")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
}

func TestPromptInvalidSyntaxSet(t *testing.T) {
	proc := runNetorcaiWaitListening(t, []string{})
	defer killallNetorcaiSIGKILL()
	re := regexp.MustCompile(`expected syntax: set VARIABLE=VALUE`)

	proc.inputControl <- "set"
	_, err := waitOutputTimeout(re, proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read 'expected syntax [...]' after set")

	proc.inputControl <- "set nb-turns-max"
	_, err = waitOutputTimeout(re, proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read 'expected syntax [...]' after set VAR")

	err = killNetorcaiGently(proc, 1000)
	assert.NoError(t, err, "Netorcai could not be killed gently")
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
