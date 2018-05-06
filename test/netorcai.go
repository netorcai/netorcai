package test

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"strings"
)

type NetorcaiProcess struct {
	cmd           *exec.Cmd
	stdinPipe     io.WriteCloser
	stdoutPipe    io.ReadCloser
	inputControl  chan string // user can send messages on this channel
	outputControl chan string // user can receive messages on this channel
	completion    chan int    // user can receive an exit code on this channel
}

func runNetorcai(command string, arguments []string) (*NetorcaiProcess, error) {
	proc := &NetorcaiProcess{
		inputControl:  make(chan string),
		outputControl: make(chan string),
		completion:    make(chan int),
	}
	proc.cmd = exec.Command(command)
	proc.cmd.Args = append([]string{command}, arguments...)

	var errIn, errOut error
	proc.stdinPipe, errIn = proc.cmd.StdinPipe()
	proc.stdoutPipe, errOut = proc.cmd.StdoutPipe()

	if errIn != nil || errOut != nil {
		return proc, fmt.Errorf("Could not setup process input/output pipes")
	}

	log.WithFields(log.Fields{
		"args": proc.cmd.Args,
	}).Debug("Starting process")
	err := proc.cmd.Start()
	if err != nil {
		return proc, fmt.Errorf("Cannot start process. %v", err)
	}

	go lineReader(bufio.NewReader(proc.stdoutPipe), proc.outputControl)
	go lineWriter(bufio.NewWriter(proc.stdinPipe), proc.inputControl)
	go waitCompletion(proc.cmd, proc.completion)
	return proc, nil
}

func runNetorcaiCover(coverFile string, arguments []string) (
	*NetorcaiProcess, error) {
	if coverFile != "" {
		// Bypass arguments
		for index, arg := range arguments {
			if strings.HasPrefix(arg, "-") {
				arguments[index] = "__bypass" + arg
			}
		}

		arguments = append([]string{"-test.coverprofile=" + coverFile},
			arguments...)

		return runNetorcai("netorcai.cover", arguments)
	} else {
		return runNetorcai("netorcai", arguments)
	}
}

func lineReader(reader *bufio.Reader, lineRead chan string) {
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		} else {
			lineRead <- strings.TrimRight(line, "\n")
		}
	}
}

func lineWriter(writer *bufio.Writer, lineToWrite chan string) {
	for {
		line := <-lineToWrite
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return
		}
		err = writer.Flush()
		if err != nil {
			return
		}
	}
}

func waitCompletion(cmd *exec.Cmd, onCompletion chan int) {
	err := cmd.Wait()
	if err != nil {
		onCompletion <- 1
	}
	onCompletion <- 0
}
