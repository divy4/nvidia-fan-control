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

type Attribute struct {
	name  string
	value int
}

// Query Nvidia for any number of GPU attributes
func query_attributes(attributes []Attribute) {
	// Build command to query everything at once
	command := make([]string, 1+2*len(attributes))
	command[0] = "nvidia-settings"
	for i, attribute := range attributes {
		command[1+2*i] = "--query"
		command[2+2*i] = attribute.name
	}

	// Run command
	stdout, _ := run_command(command)

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
		match_int, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Panic(err)
		}

		// Save int
		attributes[i].value = match_int
		i++
	}
}

// Assign any number of GPU attributes to Nvidia
func assign_attributes(attributes []Attribute) {
	command := make([]string, 1+2*len(attributes))
	command[0] = BINARY_PATH
	for i, attribute := range attributes {
		command[1+2*i] = "--assign"
		command[2+2*i] = fmt.Sprintf("%s=%d", attribute.name, attribute.value)
	}
	run_command(command)
}
