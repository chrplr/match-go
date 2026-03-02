package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		runGUI()
		return
	}

	scriptPath := os.Args[1]
	config, err := parseScript(scriptPath)
	if err != nil {
		fmt.Printf("Error parsing script: %v\n", err)
		os.Exit(1)
	}

	dataSets, err := readDataSets(config)
	if err != nil {
		fmt.Printf("Error reading data sets: %v\n", err)
		os.Exit(1)
	}

	preprocessData(config, dataSets)

	matcher := NewMatcher(config, dataSets)
	matcher.Run()
}
