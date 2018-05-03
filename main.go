package main

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
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

func initializeGlobalState(arguments map[string]interface{}) (GlobalState, error) {
	var gs GlobalState

	nbPlayersMax, err := readIntInString(arguments, "--nb-players-max",
		64, 1, 1024)
	if err != nil {
		return gs, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	nbVisusMax, err := readIntInString(arguments, "--nb-visus-max",
		64, 0, 1024)
	if err != nil {
		return gs, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	nbTurnsMax, err := readIntInString(arguments, "--nb-turns-max",
		64, 1, 65535)
	if err != nil {
		return gs, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	msBeforeFirstTurn, err := readFloatInString(arguments, "--delay-first-turn", 64, 50, 10000)
	if err != nil {
		return gs, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	msBetweenTurns, err := readFloatInString(arguments, "--delay-turns", 64, 50, 10000)
	if err != nil {
		return gs, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	gs = GlobalState{
		gameState:                   GAME_NOT_RUNNING,
		nbPlayersMax:                nbPlayersMax,
		nbVisusMax:                  nbVisusMax,
		nbTurnsMax:                  nbTurnsMax,
		millisecondsBeforeFirstTurn: msBeforeFirstTurn,
		millisecondsBetweenTurns:    msBetweenTurns,
	}

	return gs, nil
}

func setupGuards(onAbort chan int) {
	// Guard against SIGINT (ctrl+C) and SIGTERM (kill)
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		log.Warn("SIGTERM received. Aborting.")
		onAbort <- 3
	}()
}

func main() {
	os.Exit(mainReturnWithCode())
}

func mainReturnWithCode() int {
	usage := `NETwork ORChestrator for Artificial Intelligence games.

Usage:
  netorcai [--port=<port-number>]
           [--nb-turns-max=<nbt>]
           [--nb-players-max=<nbp>]
           [--nb-visus-max=<nbv>]
           [--delay-first-turn=<ms>]
           [--delay-turns=<ms>]
           [(--verbose | --quiet | --debug)] [--json-logs]
  netorcai -h | --help
  netorcai --version

Options:
  --port=<port-number>      The TCP port to listen incoming connections.
                            [default: 4242]
  --nb-turns-max=<nbt>      The maximum number of turns. [default: 100]
  --nb-players-max=<nbp>    The maximum number of players. [default: 4]
  --nb-visus-max=<nbv>      The maximum number of visualizations. [default: 1]
  --delay-first-turn=<ms>   The amount of time (in milliseconds) between the
                            GAME_STARTS message and the first TURN message.
                            [default: 1000]
  --delay-turns=<ms>		The amount of time (in milliseconds) between two
  							consecutive TURNs. [default: 1000]
  --quiet                   Only print critical information.
  --verbose                 Print information. Default verbosity mode.
  --debug                   Print debug information.
  --json-logs               Print log information in JSON.`

	netorcaiVersion := version
	if netorcaiVersion == "" {
		netorcaiVersion = "v0.1.0"
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

	port, err := readIntInString(arguments, "--port", 64, 1, 65535)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Invalid argument")
		return 1
	}

	globalState, err := initializeGlobalState(arguments)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Invalid argument")
		return 1
	}

	guardExit := make(chan int)
	serverExit := make(chan int)
	gameLogicExit := make(chan int)

	go setupGuards(guardExit)
	go server(int(port), &globalState, serverExit, gameLogicExit)

	run_prompt()

	select {
	case serverExitCode := <-serverExit:
		return serverExitCode
	case guardExitCode := <-guardExit:
		return guardExitCode
	case gameLogicExitCode := <-gameLogicExit:
		return gameLogicExitCode
	}
}
