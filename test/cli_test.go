package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setup() {
	fmt.Println("setup")
	killallNetorcai()
}

func teardown() {
	fmt.Println("teardown")
	killallNetorcai()
}

func TestCLINoArgs(t *testing.T) {
	args := []string{}
	nocIC := make(chan string)
	nocOC := make(chan string)
	nocCompletion := make(chan int)

	err := runNetorcaiCover(genCoverfile(t), args, nocIC, nocOC, nocCompletion)
	assert.NoError(t, err, "Cannot start netorcai")
	_, err = waitListening(nocOC, 1000)
	assert.NoError(t, err, "Netorcai is not listening")
}

func TestOther(t *testing.T) {
	fmt.Println("TestOther")
}

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}
