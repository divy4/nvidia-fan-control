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
var FAN_QUERY_LINE_REGEX = regexp.MustCompile(`([0-9]+)\[fan:([0-9]+)\]`)

type Attribute struct {
	name  string
	value int
}

// Queries Nvidia hardware for any number of GPU attributes.
func queryAttributes(attributes []Attribute) {
	// Build command to query everything at once
	command := make([]string, 1+2*len(attributes))
	command[0] = BINARY_PATH
	for i, attribute := range attributes {
		command[1+2*i] = "--query"
		command[2+2*i] = attribute.name
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
func assignAttributes(attributes []Attribute) {
	command := make([]string, 1+2*len(attributes))
	command[0] = BINARY_PATH
	for i, attribute := range attributes {
		command[1+2*i] = "--assign"
		command[2+2*i] = fmt.Sprintf("%s=%d", attribute.name, attribute.value)
	}
	runCommand(command)
}

// Gets a map of all fan IDs and what GPU IDs they belong to.
func getFans() map[int]int {
	command := []string{BINARY_PATH, "--query", "fans"}
	stdout, _ := runCommand(command)
	lines := strings.Split(stdout, "\n")

	fans := map[int]int{}

	for _, line := range lines {
		// Look for match in line
		matches := FAN_QUERY_LINE_REGEX.FindStringSubmatch(line)
		if len(matches) == 0 {
			continue
		}

		// Get GPU ID from first matching group
		gpuId, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Panic(err)
		}

		// Get fan ID from second matching group
		fanId, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Panic(err)
		}

		fans[fanId] = gpuId
	}

	return fans
}
