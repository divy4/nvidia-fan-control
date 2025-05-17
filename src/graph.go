package main

type AsciiGraph struct {
	gridSize     int
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

func createAsciiGraph(
	minLocation int,
	maxLocation int,
	gridSize int,
	runePriority string,
) AsciiGraph {
	graph := AsciiGraph{
		gridSize:     gridSize,
		iteration:    -2,
		min:          minLocation,
		max:          maxLocation,
		runes:        make([]rune, maxLocation-minLocation+1),
		runePriority: runePriority,
	}
	graph.reset()
	return graph
}

func (graph *AsciiGraph) reset() {
	// Keep track of how many times the line has been reset
	graph.iteration++
	isGridLine := graph.iteration%graph.gridSize == 0

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

func (graph AsciiGraph) String() string {
	return string(graph.runes)
}
