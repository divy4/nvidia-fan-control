package main

type AsciiGraph struct {
	grid_size     int
	iteration     int
	min           int
	max           int
	runes         []rune
	rune_priority string
}

const BORDER_RUNE = '#'
const GRID_RUNE = '.'
const DEFAULT_RUNE = ' '

var DIGIT_RUNES = [...]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

func create_ascii_graph(
	min_location int,
	max_location int,
	grid_size int,
	rune_priority string,
) AsciiGraph {
	graph := AsciiGraph{
		grid_size:     grid_size,
		iteration:     -2,
		min:           min_location,
		max:           max_location,
		runes:         make([]rune, max_location-min_location+1),
		rune_priority: rune_priority,
	}
	graph.reset()
	return graph
}

func (graph *AsciiGraph) reset() {
	// Keep track of how many times the line has been reset
	graph.iteration++
	is_grid_line := graph.iteration%graph.grid_size == 0

	size := graph.max - graph.min + 1

	// Reset grid to empty space, grid lines, and edges
	for i := range size {
		location := i + graph.min

		// Draw a grid line every graph.grid_size lines
		if is_grid_line {
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

func (graph *AsciiGraph) set_rune(location int, char rune) {
	index := minmax(graph.min, location, graph.max) - graph.min
	for _, c := range graph.rune_priority {
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
