package main

import (
	"fmt"
	"github.com/c-bata/go-prompt"
	"regexp"
	"strings"
)

func stringInSlice(searchedValue string, slice []string) bool {
	for _, value := range slice {
		if value == searchedValue {
			return true
		}
	}
	return false
}

func executor(line string) {
	line = strings.TrimSpace(line)
	rStart, _ := regexp.Compile(`\Astart\z`)
	rPrint, _ := regexp.Compile(`\Aprint\s+(?P<variable>\S+)\z`)
	rSet, _ := regexp.Compile(`\Aset\s+(?P<variable>\S+)(?P<sep>\s|=)(?P<value>\S+)\z`)

	acceptedSetVariables := []string{
		"nb-turns-max",
		"nb-players-max",
		"nb-visus-max",
		"delay-first-turn",
		"delay-turns",
	}

	acceptedPrintVariables := append(acceptedSetVariables, "all")

	if rStart.MatchString(line) {
		fmt.Println("start")
	} else if rPrint.MatchString(line) {
		m := rPrint.FindStringSubmatch(line)
		names := rPrint.SubexpNames()
		matches := map[string]string{}
		for index, matchedString := range m {
			matches[names[index]] = matchedString
		}

		if stringInSlice(matches["variable"], acceptedPrintVariables) {
			fmt.Printf("print %v\n", matches["variable"])
		} else {
			fmt.Printf("Bad VARIABLE=%v. Accepted values: %v\n",
				matches["variable"],
				strings.Join(acceptedPrintVariables, " "))
		}
	} else if rSet.MatchString(line) {
		m := rSet.FindStringSubmatch(line)
		names := rSet.SubexpNames()
		matches := map[string]string{}
		for index, matchedString := range m {
			matches[names[index]] = matchedString
		}

		if stringInSlice(matches["variable"], acceptedSetVariables) {
			fmt.Printf("set %v\n", matches["variable"])
		} else {
			fmt.Printf("Bad VARIABLE=%v. Accepted values: %v\n",
				matches["variable"],
				strings.Join(acceptedSetVariables, " "))
		}
	} else {
		if strings.HasPrefix(line, "start") {
			fmt.Println("expected syntax: start")
		} else if strings.HasPrefix(line, "print") {
			fmt.Println("expected syntax: print VARIABLE")
		} else if strings.HasPrefix(line, "set") {
			fmt.Println("expected syntax: set VARIABLE=VALUE\n" +
				"   (alt syntax): set VARIABLE VALUE")
		}
	}
}

func completer(d prompt.Document) []prompt.Suggest {
	commandsSugestions := []prompt.Suggest{
		{Text: "start", Description: "Start the game"},
		{Text: "print", Description: "Print value of variable"},
		{Text: "set", Description: "Set value of variable"},
		{Text: "quit", Description: "Quit netorcai"},
	}

	setSuggestions := []prompt.Suggest{
		{Text: "nb-turns-max", Description: "Maximum number of turns"},
		{Text: "nb-players-max", Description: "Maximum number of players"},
		{Text: "nb-visus-max", Description: "Maximum number of visualizations"},
		{Text: "delay-first-turn", Description: "Time (ms) before 1st turn"},
		{Text: "delay-turns", Description: "Time (ms) between turns"},
	}

	printSuggestions := append(setSuggestions, prompt.Suggest{Text: "all",
		Description: "Print the value of all variables"})

	t := d.TextBeforeCursor()

	if strings.Count(t, " ") == 0 {
		return prompt.FilterHasPrefix(commandsSugestions, t, true)
	} else if strings.HasPrefix(t, "print") {
		return prompt.FilterHasPrefix(printSuggestions,
			strings.TrimPrefix(t, "print "), true)
	} else if strings.HasPrefix(t, "set") {
		return prompt.FilterHasPrefix(setSuggestions,
			strings.TrimPrefix(t, "set "), true)
	} else {
		return []prompt.Suggest{}
	}
}

func run_prompt() {
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(">>> "),
		prompt.OptionTitle(""),
	)
	p.Run()
}
