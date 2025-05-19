package main

import (
	"bytes"
	"log"
	"math"
	"os/exec"
	"sort"
	"strings"
)

// Math

// Naively interpolate (and extrapolate) the value of a curve at x.
func interpolate(x int, curve map[int]int) int {
	found0, found1 := false, false
	x0, y0 := math.MinInt, -1
	x1, y1 := math.MaxInt, -1

	if len(curve) == 0 {
		log.Panic("Cannot interpolate on a curve with no anchors.")
	}

	// Find anchors to left and right of x
	for currX, currY := range curve {
		// Return anchor value if we match the anchor exactly
		if currX == x {
			return currY

		} else if currX < x {
			if currX > x0 {
				found0 = true
				x0, y0 = currX, currY
			}

			// currX > x
		} else {
			if currX < x1 {
				found1 = true
				x1, y1 = currX, currY
			}
		}
	}

	// Return anchor if we only found 1
	if found0 && !found1 {
		return y0
	} else if !found0 && found1 {
		return y1
	}

	// https://en.wikipedia.org/wiki/Linear_interpolation
	return (y0*(x1-x) + y1*(x-x0)) / (x1 - x0)
}

// Min and max functions joined together. Good for keeping numbers within a range.
func minMax(minX int, x int, maxX int) int {
	return max(minX, min(x, maxX))
}

// Lists/Maps

// Returns a sorted list of all keys in a map.
func getSortedIndexes(data *map[int]int) []int {
	indexes := make([]int, 0)
	for key := range *data {
		indexes = append(indexes, key)
	}
	sort.Ints(indexes)
	return indexes
}

// OS

// Executes a command and returns stdout, stderr.
func runCommand(command []string) (string, string) {
	// Setup command
	process := exec.Command(command[0], command[1:]...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	process.Stdout = &stdout
	process.Stderr = &stderr
	// Run it
	err := process.Run()
	if err != nil {
		log.Panicf("Failed to execute %s: %s", command, err)
	}
	if strings.Contains(stderr.String(), "ERROR") {
		log.Panicf(
			"Failed to execute %s: stderr contains 'ERROR':\nstdout:\n%sstderr:\n%s",
			command,
			stdout.String(),
			stderr.String(),
		)
	}
	return stdout.String(), stderr.String()
}
