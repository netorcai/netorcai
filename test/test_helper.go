package test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

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
	timeoutMS int) (matchingLine string, err error) {
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
				log.Printf("Read non-matching line: %v\n", line)
			}
		case <-timeoutReached:
			return "", fmt.Errorf("Timeout reached")
		}
	}
}

func waitListening(output chan string, timeoutMS int) (
	matchingLine string, err error) {
	re := regexp.MustCompile("Listening incoming connections")
	return waitOutputTimeout(re, output, timeoutMS)
}

func killallNetorcai() error {
	cmd := exec.Command("killall")
	cmd.Args = []string{"killall", "--quiet", "netorcai", "netorcai.cover"}
	return cmd.Run()
}

func genCoverfile(t *testing.T) string {
	_, exists := os.LookupEnv("DO_COVERAGE")
	if exists {
		return t.Name() + ".covout"
	} else {
		return ""
	}
}