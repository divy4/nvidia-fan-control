package main

type AsciiGraph struct {
	gridSize     int
	mergeSize    int
	iteration    int
	min          int
	max          int
	runes        []rune
	runePriority string
}

const BORDER_RUNE = '#'
const GRID_RUNE = '.'
const DEFAULT_RUNE = ' '

var DIGIT_RUNES = [...]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

// Creates an AsciiGraph.
func createAsciiGraph(
	minLocation int,
	maxLocation int,
	gridSize int,
	mergeSize int,
	runePriority string,
) AsciiGraph {
	graph := AsciiGraph{
		gridSize:     gridSize,
		mergeSize:    mergeSize,
		iteration:    -1,
		min:          minLocation,
		max:          maxLocation,
		runes:        make([]rune, maxLocation-minLocation+1),
		runePriority: runePriority,
	}
	graph.next_line()
	return graph
}

// Clears an AsciiGraph.
func (graph *AsciiGraph) next_line() {
	// Keep track of how many times the line has been reset
	graph.iteration++
	// Only clear the graph once very mergeSize times
	if graph.iteration % graph.mergeSize != 0 {
		return
	}

	// And make it a grid line once every mergeSize * gridSize times
	isGridLine := graph.iteration % (graph.mergeSize * graph.gridSize) == 0

	size := graph.max - graph.min + 1

	// Reset grid to empty space, grid lines, and edges
	for i := range size {
		location := i + graph.min

		// Draw a grid line every graph.gridSize lines
		if isGridLine {
			if i == 0 || i == size-1 {
				graph.runes[i] = BORDER_RUNE
			} else if location%10 == 0 {
				graph.runes[i] = DIGIT_RUNES[location/10%10]
			} else {
				graph.runes[i] = GRID_RUNE
			}

			// Draw a normal line
		} else {
			if i == 0 || i == size-1 {
				graph.runes[i] = BORDER_RUNE
			} else {
				graph.runes[i] = DEFAULT_RUNE
			}
		}
	}
}

// Sets a single rune on an AsciiGraph.
func (graph *AsciiGraph) setRune(location int, char rune) {
	index := minMax(graph.min, location, graph.max) - graph.min
	for _, c := range graph.runePriority {
		// Don't change anything if the current rune has higher priority
		switch c {
		// Skip if the current rune has higher priority
		case graph.runes[index]:
			return
		// Update char if it's a higher priority
		case char:
			{
				graph.runes[index] = char
				return
			}
		}
	}
}

// Check if the currnet line is a line where the graph should print
func (graph AsciiGraph) ready() bool {
	return graph.iteration % graph.mergeSize == graph.mergeSize - 1
}

// Convert an AsciiGraph into a string.
func (graph AsciiGraph) String() string {
	return string(graph.runes)
}
