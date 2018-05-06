package test

import (
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"testing"
)

func TestMain(m *testing.M) {
	killallNetorcai()
	retCode := m.Run()
	killallNetorcai()
	os.Exit(retCode)
}

func TestCLINoArgs(t *testing.T) {
	args := []string{}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgHelp(t *testing.T) {
	args := []string{"--help"}
	coverFile, expRetCode := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgHelpShort(t *testing.T) {
	args := []string{"-h"}
	coverFile, expRetCode := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgVersion(t *testing.T) {
	args := []string{"--version"}
	coverFile, expRetCode := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitOutputTimeout(regexp.MustCompile(`\Av\d+\.\d+\.\d+\S*\z`),
		proc.outputControl, 1000, false)
	assert.NoError(t, err, "Cannot read version")

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgVerbose(t *testing.T) {
	args := []string{"--verbose"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgQuiet(t *testing.T) {
	args := []string{"--quiet", "--port=-1"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDebug(t *testing.T) {
	args := []string{"--debug"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgJsonLogs(t *testing.T) {
	args := []string{"--json-logs"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIInvalidVerbosityCombination(t *testing.T) {
	args := []string{"--debug", "--verbose"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIUnknownArg(t *testing.T) {
	args := []string{"--this-option-should-not-exist"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

/********************
 * --nb-players-max *
 ********************/
func TestCLIArgNbPlayersMaxNotInteger(t *testing.T) {
	args := []string{"--nb-players-max=meh"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbPlayersMaxTooSmall(t *testing.T) {
	args := []string{"--nb-players-max=0"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbPlayersMaxTooBig(t *testing.T) {
	args := []string{"--nb-players-max=1025"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbPlayersMaxSmall(t *testing.T) {
	args := []string{"--nb-players-max=1"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgNbPlayersMaxBig(t *testing.T) {
	args := []string{"--nb-players-max=1024"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/**********
 * --port *
 **********/
func TestCLIArgPortNotInteger(t *testing.T) {
	args := []string{"--port=meh"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgPortTooSmall(t *testing.T) {
	args := []string{"--port=0"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgPortTooBig(t *testing.T) {
	args := []string{"--port=65536"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgPortSmall(t *testing.T) {
	args := []string{"--port=1025"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgPortBig(t *testing.T) {
	args := []string{"--port=65535"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/******************
 * --nb-turns-max *
 ******************/
func TestCLIArgNbTurnsMaxNotInteger(t *testing.T) {
	args := []string{"--nb-turns-max=meh"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbTurnsMaxTooSmall(t *testing.T) {
	args := []string{"--nb-turns-max=0"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbTurnsMaxTooBig(t *testing.T) {
	args := []string{"--nb-turns-max=65536"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbTurnsMaxSmall(t *testing.T) {
	args := []string{"--nb-turns-max=1"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgNbTurnsMaxBig(t *testing.T) {
	args := []string{"--nb-turns-max=65535"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/******************
 * --nb-visus-max *
 ******************/
func TestCLIArgNbVisusMaxNotInteger(t *testing.T) {
	args := []string{"--nb-visus-max=meh"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbVisusMaxTooSmall(t *testing.T) {
	args := []string{"--nb-visus-max=-1"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbVisusMaxTooBig(t *testing.T) {
	args := []string{"--nb-visus-max=1025"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbVisusMaxSmall(t *testing.T) {
	args := []string{"--nb-visus-max=0"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgNbVisusMaxBig(t *testing.T) {
	args := []string{"--nb-visus-max=1024"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/**********************
 * --delay-first-turn *
 **********************/
func TestCLIArgDelayFirstTurnNotFloat(t *testing.T) {
	args := []string{"--delay-first-turn=meh"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayFirstTurnTooSmall(t *testing.T) {
	args := []string{"--delay-first-turn=49.999"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayFirstTurnTooBig(t *testing.T) {
	args := []string{"--delay-first-turn=10000.001"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayFirstTurnSmall(t *testing.T) {
	args := []string{"--delay-first-turn=50"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgDelayFirstTurnBig(t *testing.T) {
	args := []string{"--delay-first-turn=10000"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/*****************
 * --delay-turns *
 *****************/
func TestCLIArgDelayTurnsNotFloat(t *testing.T) {
	args := []string{"--delay-turns=meh"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayTurnsTooSmall(t *testing.T) {
	args := []string{"--delay-turns=49.999"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayTurnsTooBig(t *testing.T) {
	args := []string{"--delay-turns=10000.001"}
	coverFile, expRetCode := handleCoverage(t, 1)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayTurnsSmall(t *testing.T) {
	args := []string{"--delay-turns=50"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgDelayTurnsBig(t *testing.T) {
	args := []string{"--delay-turns=10000"}
	coverFile, _ := handleCoverage(t, 0)

	proc, err := runNetorcaiCover(coverFile, args)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(proc.outputControl, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}
