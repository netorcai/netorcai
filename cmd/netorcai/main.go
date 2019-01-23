package main

import (
	"fmt"
	docopt "github.com/docopt/docopt-go"
	"github.com/netorcai/netorcai"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
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

func initializeGlobalState(arguments map[string]interface{}) (
	*netorcai.GlobalState, error) {
	nbPlayersMax, err := netorcai.ReadIntInString(arguments,
		"--nb-players-max", 64, 1, 1024)
	if err != nil {
		return nil, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	nbVisusMax, err := netorcai.ReadIntInString(arguments,
		"--nb-visus-max", 64, 0, 1024)
	if err != nil {
		return nil, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	nbTurnsMax, err := netorcai.ReadIntInString(arguments,
		"--nb-turns-max", 64, 1, 65535)
	if err != nil {
		return nil, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	msBeforeFirstTurn, err := netorcai.ReadFloatInString(arguments, "--delay-first-turn", 64, 50, 10000)
	if err != nil {
		return nil, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	msBetweenTurns, err := netorcai.ReadFloatInString(arguments,
		"--delay-turns", 64, 50, 10000)
	if err != nil {
		return nil, fmt.Errorf("Invalid arguments: %v", err.Error())
	}

	autostart := arguments["--autostart"].(bool)

	gs := &netorcai.GlobalState{
		GameState:                   netorcai.GAME_NOT_RUNNING,
		NbPlayersMax:                nbPlayersMax,
		NbVisusMax:                  nbVisusMax,
		NbTurnsMax:                  nbTurnsMax,
		Autostart:                   autostart,
		MillisecondsBeforeFirstTurn: msBeforeFirstTurn,
		MillisecondsBetweenTurns:    msBetweenTurns,
	}

	return gs, nil
}

func setupGuards(gs *netorcai.GlobalState, onAbort chan int) {
	// Guard against SIGINT (ctrl+C) and SIGTERM (kill)
	sigterm := make(chan os.Signal, 2)
	signal.Notify(sigterm, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigterm
		onAbort <- 1
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
           [--autostart]
           [--simple-prompt]
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
  --delay-turns=<ms>        The amount of time (in milliseconds) between two
                            consecutive TURNs. [default: 1000]
  --autostart               Start game when all clients are connnected.
                            Set --nb-{players,visus}-max accordingly.
  --simple-prompt           Always use a simple prompt.
  --quiet                   Only print critical information.
  --verbose                 Print information. Default verbosity mode.
  --debug                   Print debug information.
  --json-logs               Print log information in JSON.`

	netorcaiVersion := version
	if netorcaiVersion == "" {
		netorcaiVersion = "v1.2.0"
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

	port, err := netorcai.ReadIntInString(arguments, "--port", 64, 1, 65535)
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
	shellExit := make(chan int)

	setupGuards(globalState, guardExit)
	go netorcai.RunServer(int(port), globalState, serverExit, gameLogicExit)

	interactivePrompt := true
	if arguments["--simple-prompt"] == true {
		interactivePrompt = false
	} else {
		interactivePrompt = terminal.IsTerminal(int(os.Stdout.Fd()))
	}

	go netorcai.RunPrompt(globalState, shellExit, interactivePrompt)

	select {
	case serverExitCode := <-serverExit:
		return serverExitCode
	case guardExitCode := <-guardExit:
		log.Warn("SIGTERM received. Aborting.")
		netorcai.Cleanup()
		return guardExitCode
	case gameLogicExitCode := <-gameLogicExit:
		if gameLogicExitCode != 0 {
			log.Warn("Game logic failed. Aborting.")
		}
		netorcai.Cleanup()
		return gameLogicExitCode
	case shellExitCode := <-shellExit:
		log.Warn("Shell exited. Aborting.")
		netorcai.Cleanup()
		return shellExitCode
	}
}
