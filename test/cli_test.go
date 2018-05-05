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
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgHelp(t *testing.T) {
	args := []string{"--help"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgHelpShort(t *testing.T) {
	args := []string{"-h"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgVersion(t *testing.T) {
	args := []string{"--version"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitOutputTimeout(regexp.MustCompile(`\Av\d+\.\d+\.\d+\S*\z`),
		nocOC, 1000, false)
	assert.NoError(t, err, "Cannot read version")

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgVerbose(t *testing.T) {
	args := []string{"--verbose"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgQuiet(t *testing.T) {
	args := []string{"--quiet", "--port=-1"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDebug(t *testing.T) {
	args := []string{"--debug"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgJsonLogs(t *testing.T) {
	args := []string{"--json-logs"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIInvalidVerbosityCombination(t *testing.T) {
	args := []string{"--debug", "--verbose"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIUnknownArg(t *testing.T) {
	args := []string{"--this-option-should-not-exist"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

/********************
 * --nb-players-max *
 ********************/
func TestCLIArgNbPlayersMaxNotInteger(t *testing.T) {
	args := []string{"--nb-players-max=meh"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbPlayersMaxTooSmall(t *testing.T) {
	args := []string{"--nb-players-max=0"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbPlayersMaxTooBig(t *testing.T) {
	args := []string{"--nb-players-max=1025"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbPlayersMaxSmall(t *testing.T) {
	args := []string{"--nb-players-max=1"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgNbPlayersMaxBig(t *testing.T) {
	args := []string{"--nb-players-max=1024"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/**********
 * --port *
 **********/
func TestCLIArgPortNotInteger(t *testing.T) {
	args := []string{"--port=meh"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgPortTooSmall(t *testing.T) {
	args := []string{"--port=0"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgPortTooBig(t *testing.T) {
	args := []string{"--port=65536"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgPortSmall(t *testing.T) {
	args := []string{"--port=1025"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgPortBig(t *testing.T) {
	args := []string{"--port=65535"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/******************
 * --nb-turns-max *
 ******************/
func TestCLIArgNbTurnsMaxNotInteger(t *testing.T) {
	args := []string{"--nb-turns-max=meh"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbTurnsMaxTooSmall(t *testing.T) {
	args := []string{"--nb-turns-max=0"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbTurnsMaxTooBig(t *testing.T) {
	args := []string{"--nb-turns-max=65536"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbTurnsMaxSmall(t *testing.T) {
	args := []string{"--nb-turns-max=1"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgNbTurnsMaxBig(t *testing.T) {
	args := []string{"--nb-turns-max=65535"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/******************
 * --nb-visus-max *
 ******************/
func TestCLIArgNbVisusMaxNotInteger(t *testing.T) {
	args := []string{"--nb-visus-max=meh"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbVisusMaxTooSmall(t *testing.T) {
	args := []string{"--nb-visus-max=-1"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbVisusMaxTooBig(t *testing.T) {
	args := []string{"--nb-visus-max=1025"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgNbVisusMaxSmall(t *testing.T) {
	args := []string{"--nb-visus-max=0"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgNbVisusMaxBig(t *testing.T) {
	args := []string{"--nb-visus-max=1024"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/**********************
 * --delay-first-turn *
 **********************/
func TestCLIArgDelayFirstTurnNotFloat(t *testing.T) {
	args := []string{"--delay-first-turn=meh"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayFirstTurnTooSmall(t *testing.T) {
	args := []string{"--delay-first-turn=49.999"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayFirstTurnTooBig(t *testing.T) {
	args := []string{"--delay-first-turn=10000.001"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayFirstTurnSmall(t *testing.T) {
	args := []string{"--delay-first-turn=50"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgDelayFirstTurnBig(t *testing.T) {
	args := []string{"--delay-first-turn=10000"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

/*****************
 * --delay-turns *
 *****************/
func TestCLIArgDelayTurnsNotFloat(t *testing.T) {
	args := []string{"--delay-turns=meh"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayTurnsTooSmall(t *testing.T) {
	args := []string{"--delay-turns=49.999"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayTurnsTooBig(t *testing.T) {
	args := []string{"--delay-turns=10000.001"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, expRetCode := handleCoverage(t, 1)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	retCode, err := waitCompletionTimeout(nocCompletion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}

func TestCLIArgDelayTurnsSmall(t *testing.T) {
	args := []string{"--delay-turns=50"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestCLIArgDelayTurnsBig(t *testing.T) {
	args := []string{"--delay-turns=10000"}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)
	coverFile, _ := handleCoverage(t, 0)

	err := runNetorcaiCover(coverFile, args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	defer killallNetorcai()

	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}
