package test

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func runNetorcai(command string, arguments []string,
	inputControl, outputControl chan string, completion chan int) error {
	cmd := exec.Command(command)
	cmd.Args = append([]string{command}, arguments...)

	stdinPipe, errIn := cmd.StdinPipe()
	stdoutPipe, errOut := cmd.StdoutPipe()
	if errIn != nil || errOut != nil {
		return fmt.Errorf("Could not setup process input/output pipes")
	}

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Cannot start process. %v", err)
	}

	go lineReader(bufio.NewReader(stdoutPipe), outputControl)
	go lineWriter(bufio.NewWriter(stdinPipe), inputControl)
	go waitCompletion(cmd, completion)
	return nil
}

func runNetorcaiCover(coverFile string, arguments []string,
	inputControl, outputControl chan string, completion chan int) error {

	if coverFile != "" {
		// Bypass arguments
		for _, arg := range arguments {
			if strings.HasPrefix(arg, "-") {
				arg = "__bypass" + arg
			}
		}

		arguments = append([]string{"-test.coverprofile=" + coverFile},
			arguments...)

		return runNetorcai("netorcai.cover", arguments, inputControl,
			outputControl, completion)
	} else {
		return runNetorcai("netorcai", arguments, inputControl,
			outputControl, completion)
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
	}
}

func waitCompletion(cmd *exec.Cmd, onCompletion chan int) {
	err := cmd.Wait()
	if err != nil {
		onCompletion <- 1
	}
	onCompletion <- 0
}
