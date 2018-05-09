package test

import (
	"regexp"
	"testing"
)

func TestKickallOnAbortKillSigterm(t *testing.T) {
	_, clients, _, _, _ := runNetorcaiAndAllClients(t, 1000)

	killallNetorcai()

	checkAllKicked(t, clients, regexp.MustCompile(`netorcai abort`), 1000)
}