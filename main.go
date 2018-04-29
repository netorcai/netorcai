package main

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
)

var (
	version string
)

func setupLogging(arguments map[string]interface{}) {
	log.SetOutput(os.Stdout)

	if arguments["--json-logs"] == true {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		customFormatter := new(log.TextFormatter)
		customFormatter.TimestampFormat = "2006-01-02 15:04:05.000"
		customFormatter.FullTimestamp = true
		customFormatter.QuoteEmptyFields = true
		log.SetFormatter(customFormatter)
	}

	if arguments["--debug"] == true {
		log.SetLevel(log.DebugLevel)
	} else if arguments["--quiet"] == true {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	os.Exit(mainReturnWithCode())
}

func mainReturnWithCode() int {
	usage := `NETwork ORChestrator for Artificial Intelligence games.

Usage:
  netorcai [--port=<port-number>]
           [(--verbose | --quiet | --debug)] [--json-logs]
  netorcai -h | --help
  netorcai --version

Options:
  --port=<port-number>  The TCP port to listen incoming connections.
                        [default: 4242]
  --quiet               Only print critical information.
  --verbose             Print information. Default verbosity mode.
  --debug               Print debug information.
  --json-logs           Print log information in JSON.`

	netorcaiVersion := version
	if netorcaiVersion == "" {
		netorcaiVersion = "unreleased-yet"
	}

	ret := -1

	parser := &docopt.Parser{
		HelpHandler: func(err error, usage string) {
			fmt.Println(usage)
			if err != nil {
				ret = 1
			} else {
				ret = 0
			}
		},
		OptionsFirst: false,
	}

	arguments, _ := parser.ParseArgs(usage, os.Args[1:], netorcaiVersion)
	if ret != -1 {
		return ret
	}

	setupLogging(arguments)

	port, err := strconv.ParseInt(arguments["--port"].(string), 0, 16)
	if err != nil {
		log.WithFields(log.Fields{
			"err":  err,
			"port": arguments["--port"].(string),
		}).Error("Invalid port")
		return 1
	}

	serverExit := make(chan int)

	go server(int(port), serverExit)

	select {
	case serverExitCode := <-serverExit:
		return serverExitCode
	}

	return 0
}
