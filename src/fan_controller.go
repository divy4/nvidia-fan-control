package main

import (
	"fmt"
	"strings"
	"time"
)

// Public types

type Fan struct {
	gpu_id        int
	min_speed     int
	max_speed     int
	control_curve FanCurve
}

type FanCurve map[int]int

type FanController struct {
	gpus  map[int]FanControllerGpu
	fans  map[int]FanControllerFan
	graph AsciiGraph
}

// Internal types

type FanControllerGpu struct {
	temperature int
}

type FanControllerFan struct {
	gpu_id        int
	min_speed     int
	max_speed     int
	current_speed int
	target_speed  int
	// Map of gpu temp -> target fan speed
	control_curve map[int]int
}

const GPU_TEMP_RUNE = '|'
const FAN_SPEED_RUNE = ':'
const GRAPH_RUNE_PRIORITY = "|:"
const GRID_SIZE = 10

func create_fan_controller(fans map[int]Fan, graph_min int, graph_max int) FanController {
	controller := FanController{
		gpus:  map[int]FanControllerGpu{},
		fans:  map[int]FanControllerFan{},
		graph: create_ascii_graph(graph_min, graph_max, GRID_SIZE, GRAPH_RUNE_PRIORITY),
	}
	for fan_id, fan := range fans {
		controller.gpus[fan.gpu_id] = FanControllerGpu{}
		controller.fans[fan_id] = FanControllerFan{
			gpu_id:        fan.gpu_id,
			min_speed:     fan.min_speed,
			max_speed:     fan.max_speed,
			target_speed:  fan.min_speed,
			control_curve: fan.control_curve,
		}
	}
	return controller
}

func (controller *FanController) run() {
	defer controller.disable_fan_control()
	controller.enable_fan_control()
	controller.print_stats_headers()
	for {
		controller.update_stats()
		controller.calculate_target_fan_speeds()
		controller.push_target_fan_speeds()
		controller.print_stats()
		time.Sleep(1 * time.Second)
	}
}

// I/O

func (controller *FanController) print_stats_headers() {
	headers := make([]string, len(controller.gpus)+len(controller.fans)*2+1)
	i := 0
	for id, _ := range controller.gpus {
		headers[i] = fmt.Sprintf("GPU %d TEMP", id)
		i++
	}
	for id, _ := range controller.fans {
		headers[i] = fmt.Sprintf("FAN %d CURRENT SPEED", id)
		i++
	}
	for id, _ := range controller.fans {
		headers[i] = fmt.Sprintf("FAN %d TARGET SPEED", id)
		i++
	}
	headers[i] = "GRAPH"
	fmt.Println(strings.Join(headers, ","))
}

func (controller *FanController) print_stats() {
	// Reset graph
	controller.graph.reset()

	values := make([]string, len(controller.gpus)+len(controller.fans)*2+1)
	i := 0

	// GPU temp
	for _, gpu := range controller.gpus {
		values[i] = fmt.Sprintf("%3d", gpu.temperature)
		controller.graph.set_rune(gpu.temperature, GPU_TEMP_RUNE)
		i++
	}
	// Fan current speed
	for _, fan := range controller.fans {
		values[i] = fmt.Sprintf("%3d", fan.current_speed)
		controller.graph.set_rune(fan.current_speed, FAN_SPEED_RUNE)
		i++
	}
	// Fan target speed
	for _, fan := range controller.fans {
		values[i] = fmt.Sprintf("%3d", fan.target_speed)
		i++
	}
	values[i] = controller.graph.String()
	// Join it all into 1 string
	fmt.Println(strings.Join(values, ","))
}

// Helpers

func (controller *FanController) enable_fan_control() {
	fmt.Println("Enabling fan control...")
	attributes := make([]Attribute, len(controller.gpus))
	i := 0
	for id, _ := range controller.gpus {
		attributes[i].name = fmt.Sprintf("[gpu:%d]/GPUFanControlState", id)
		attributes[i].value = 1
		i++
	}
	assign_attributes(attributes)
}

func (controller *FanController) disable_fan_control() {
	fmt.Println("Disabling fan control...")
	attributes := make([]Attribute, len(controller.gpus))
	i := 0
	for id, _ := range controller.gpus {
		attributes[id].name = fmt.Sprintf("[gpu:%d]/GPUFanControlState", id)
		attributes[id].value = 0
		i++
	}
	assign_attributes(attributes)
}

func (controller *FanController) update_stats() {
	// Build attributes to query
	attributes := make([]Attribute, len(controller.gpus)+len(controller.fans)*2)
	i := 0
	for id, _ := range controller.gpus {
		attributes[i].name = fmt.Sprintf("[gpu:%d]/GPUCoreTemp", id)
		i++
	}
	for id, _ := range controller.fans {
		attributes[i].name = fmt.Sprintf("[fan:%d]/GPUCurrentFanSpeed", id)
		attributes[i+1].name = fmt.Sprintf("[fan:%d]/GPUTargetFanSpeed", id)
		i += 2
	}

	query_attributes(attributes)

	// Map returned attributes to metrics
	i = 0
	for id, gpu := range controller.gpus {
		gpu.temperature = attributes[i].value
		controller.gpus[id] = gpu
		i++
	}
	for id, fan := range controller.fans {
		fan.current_speed = attributes[i].value
		fan.target_speed = attributes[i+1].value
		controller.fans[id] = fan
		i += 2
	}
}

func (controller *FanController) calculate_target_fan_speeds() {
	for id, fan := range controller.fans {
		gpu := controller.gpus[fan.gpu_id]

		speed := interpolate(gpu.temperature, fan.control_curve)
		// Limit speed reduction
		speed = max(fan.current_speed-1, speed)
		// Ensure speed is within 0-100
		speed = max(fan.min_speed, min(100, speed))
		// Set target speed
		fan.target_speed = speed
		controller.fans[id] = fan
	}
}

func (controller *FanController) push_target_fan_speeds() {
	attributes := make([]Attribute, len(controller.fans))
	i := 0
	for id, fan := range controller.fans {
		attributes[i].name = fmt.Sprintf("[fan:%d]/GPUTargetFanSpeed", id)
		attributes[i].value = fan.target_speed
		i++
	}
	assign_attributes(attributes)
}
