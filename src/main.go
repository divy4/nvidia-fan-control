package main

import (
	"fmt"
	"os"
)

func main() {
	command, configFile := parseArgs()
	config := loadConfig(configFile)

	// TODO: figure out signal interrupts
	controller := createFanController(&config)

	switch command {
	case "run":
		controller.run()
	case "stop":
		controller.disableFanControl()
	default:
		printHelp()
		os.Exit(1)
	}
}

// Parses command line arguments.
func parseArgs() (string, string) {
	if len(os.Args) != 3 {
		printHelp()
		os.Exit(1)
	}

	// Command
	command := os.Args[1]
	if command != "run" && command != "stop" {
		printHelp()
		os.Exit(1)
	}

	// Config file
	configFile := os.Args[2]

	return command, configFile
}

// Prints the help message.
func printHelp() {
	fmt.Println("Usage: nvidia-fan-control run|stop <config-file>")
}
