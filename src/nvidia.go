package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const BINARY_PATH = "/usr/bin/nvidia-settings"

var ATTRIBUTE_QUERY_LINE_REGEX = regexp.MustCompile(`Attribute '.*' \(.*\): ([0-9]+)\.`)
var FAN_QUERY_LINE_REGEX = regexp.MustCompile(`\[fan:([0-9]+)\]`)

type Attribute struct {
	name  string
	value int
}

// Queries Nvidia hardware for any number of GPU attributes.
func queryAttributes(attributes []Attribute, xDisplay int) {
	// Build command to query everything at once
	command := make([]string, 3+2*len(attributes))
	command[0] = BINARY_PATH
	command[1] = "--display"
	command[2] = strconv.Itoa(xDisplay)
	for i, attribute := range attributes {
		command[3+2*i] = "--query"
		command[4+2*i] = attribute.name
	}

	// Run command
	stdout, _ := runCommand(command)

	// Parse output into ints
	lines := strings.Split(stdout, "\n")
	i := 0
	for _, line := range lines {
		// Look for match in line
		matches := ATTRIBUTE_QUERY_LINE_REGEX.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue
		}

		// Convert match group to int
		matchInt, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Panic(err)
		}

		// Save int
		attributes[i].value = matchInt
		i++
	}
}

// Assigns any number of GPU attributes to Nvidia hardware.
func assignAttributes(attributes []Attribute, xDisplay int) {
	command := make([]string, 3+2*len(attributes))
	command[0] = BINARY_PATH
	command[1] = "--display"
	command[2] = strconv.Itoa(xDisplay)
	for i, attribute := range attributes {
		command[3+2*i] = "--assign"
		command[4+2*i] = fmt.Sprintf("%s=%d", attribute.name, attribute.value)
	}
	runCommand(command)
}

// Gets a map of all fan IDs and what GPU IDs they belong to.
func getFans(xDisplay int) map[int]int {
	command := []string{
		BINARY_PATH,
		"--display",
		strconv.Itoa(xDisplay),
		"--query",
		"fans",
	}
	stdout, _ := runCommand(command)
	lines := strings.Split(stdout, "\n")

	fans := map[int]int{}

	for _, line := range lines {
		// Look for match in line
		matches := FAN_QUERY_LINE_REGEX.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue
		}

		// Get fan ID from matching group
		fanId, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Panic(err)
		}

		// TODO: Support multi-gpu
		fans[fanId] = 0
	}

	return fans
}
