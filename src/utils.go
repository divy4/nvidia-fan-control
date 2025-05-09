package main

import (
	"bytes"
	"log"
	"math"
	"os/exec"
	"strings"
)

// Naively interpolate (and extrapolate) the value of a curve at x.
func interpolate(x int, curve map[int]int) int {
	found_0, found_1 := false, false
	x0, y0 := math.MinInt, -1
	x1, y1 := math.MaxInt, -1

	if len(curve) == 0 {
		log.Panic("Cannot interpolate on a curve with no anchors.")
	}

	// Find anchors to left and right of x
	for curr_x, curr_y := range curve {
		// Return anchor value if we match the anchor exactly
		if curr_x == x {
			return curr_y

		} else if curr_x < x {
			if curr_x > x0 {
				found_0 = true
				x0, y0 = curr_x, curr_y
			}

			// curr_x > x
		} else {
			if curr_x < x1 {
				found_1 = true
				x1, y1 = curr_x, curr_y
			}
		}
	}

	// Return anchor if we only found 1
	if found_0 && !found_1 {
		return y0
	} else if !found_0 && found_1 {
		return y1
	}

	// https://en.wikipedia.org/wiki/Linear_interpolation
	return (y0*(x1-x) + y1*(x-x0)) / (x1 - x0)
}

// Min and max functions joined together. Good for keeping numbers within a range.
func minmax(min_x int, x int, max_x int) int {
	return max(min_x, min(x, max_x))
}

// Execute a command and return stdout, stderr
func run_command(command []string) (string, string) {
	// Setup command
	process := exec.Command(command[0], command[1:]...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	process.Stdout = &stdout
	process.Stderr = &stderr
	// Run it
	err := process.Run()
	if err != nil {
		log.Panic("Failed to execute %s: %s", command, err)
	}
	if strings.Contains(stderr.String(), "ERROR") {
		log.Panicf(
			"Failed to execute %s: stderr contains 'ERROR':\nstdout:\n%sstderr:\n%s",
			command,
			stdout.String(),
			stderr.String(),
		)
	}
	// Get output (TODO: Remove output)
	// fmt.Printf("command: %s\n", command)
	return stdout.String(), stderr.String()
}
