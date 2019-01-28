package test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestKickallOnAbortKillSigterm(t *testing.T) {
	proc, clients, _, _, _, _ := runNetorcaiAndAllClients(t, []string{}, 1000, 0)
	defer killallNetorcaiSIGKILL()

	killallNetorcai()

	checkAllKicked(t, clients, regexp.MustCompile(`netorcai abort`), 1000)

	retCode, err := waitCompletionTimeout(proc.completion, 1000)
	assert.NoError(t, err, "netorcai did not complete")
	_, expRetCode := handleCoverage(t, 1)
	assert.Equal(t, expRetCode, retCode, "Unexpected netorcai return code")
}
