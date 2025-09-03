package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type FanController struct {
	gpus     map[int]FanControllerGpu
	fans     map[int]FanControllerFan
	graph    AsciiGraph
	xDisplay int
}

type FanControllerGpu struct {
	temperature int
}

type FanControllerFan struct {
	gpuId        int
	currentSpeed int
	targetSpeed  int
	// Map of gpu temp -> target fan speed
	controlCurve map[int]int
}

const GPU_TEMP_RUNE = '|'
const FAN_SPEED_RUNE = ':'
const GRAPH_RUNE_PRIORITY = "|:"
const GRID_SIZE = 10

// Creates a FanController object.
func createFanController(config *Config) FanController {
	controller := FanController{
		gpus:     map[int]FanControllerGpu{},
		fans:     map[int]FanControllerFan{},
		graph:    createAsciiGraph(config.Graph.Min, config.Graph.Max, GRID_SIZE, GRAPH_RUNE_PRIORITY),
		xDisplay: config.XDisplay,
	}
	for fanId, fan := range config.Fans {
		controller.gpus[fan.GpuId] = FanControllerGpu{}
		controller.fans[fanId] = FanControllerFan{
			gpuId:        fan.GpuId,
			controlCurve: fan.ControlCurve,
		}
	}
	return controller
}

// Runs a FanController to control fan speed.
// Note: This function runs infinitely!
func (controller *FanController) run() {
	defer controller.disableFanControl()
	controller.enableFanControl()
	controller.printStatsHeaders()
	for {
		controller.updateStats()
		controller.calculateTargetFanSpeeds()
		controller.pushTargetFanSpeeds()
		controller.printStats()
		time.Sleep(1 * time.Second)
	}
}

// I/O

// Prints metrics headers about a FanController.
func (controller *FanController) printStatsHeaders() {
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

// Prints metrics about a FanController.
func (controller *FanController) printStats() {
	// Reset graph
	controller.graph.clear()

	values := make([]string, len(controller.gpus)+len(controller.fans)*2+1)
	i := 0

	// GPU temp
	for _, gpu := range controller.gpus {
		values[i] = fmt.Sprintf("%3d", gpu.temperature)
		controller.graph.setRune(gpu.temperature, GPU_TEMP_RUNE)
		i++
	}
	// Fan current speed
	for _, fan := range controller.fans {
		values[i] = fmt.Sprintf("%3d", fan.currentSpeed)
		controller.graph.setRune(fan.currentSpeed, FAN_SPEED_RUNE)
		i++
	}
	// Fan target speed
	for _, fan := range controller.fans {
		values[i] = fmt.Sprintf("%3d", fan.targetSpeed)
		i++
	}
	values[i] = controller.graph.String()
	// Join it all into 1 string
	fmt.Println(strings.Join(values, ","))
}

// Helpers

// Enables control over all fans within a FanController.
func (controller *FanController) enableFanControl() {
	fmt.Println("Enabling fan control...")
	attributes := make([]Attribute, len(controller.gpus))
	i := 0
	for id, _ := range controller.gpus {
		attributes[i].name = fmt.Sprintf("[gpu:%d]/GPUFanControlState", id)
		attributes[i].value = 1
		i++
	}
	assignAttributes(attributes, controller.xDisplay)
}

// Disables control over all fans within a FanController.
func (controller *FanController) disableFanControl() {
	fmt.Println("Disabling fan control...")
	attributes := make([]Attribute, len(controller.gpus))
	i := 0
	for id, _ := range controller.gpus {
		attributes[id].name = fmt.Sprintf("[gpu:%d]/GPUFanControlState", id)
		attributes[id].value = 0
		i++
	}
	assignAttributes(attributes, controller.xDisplay)
}

// Pulls hardware stats from the real hardware.
func (controller *FanController) updateStats() {
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

	queryAttributes(attributes, controller.xDisplay)

	// Map returned attributes to metrics
	i = 0
	for id, gpu := range controller.gpus {
		gpu.temperature = attributes[i].value
		controller.gpus[id] = gpu
		i++
	}
	for id, fan := range controller.fans {
		fan.currentSpeed = attributes[i].value
		fan.targetSpeed = attributes[i+1].value
		controller.fans[id] = fan
		i += 2
	}
}

// Calculates what the target fan speed should be for all fans based on current
// metrics.
func (controller *FanController) calculateTargetFanSpeeds() {
	// For each gpu...
	for id, fan := range controller.fans {
		gpu := controller.gpus[fan.gpuId]
		// Compute the fan speed based on temp
		speed := interpolate(gpu.temperature, fan.controlCurve)
		// Limit speed slow down to 1% per tick
		speed = max(speed, fan.currentSpeed-1)
		// Limit speed to the minimum and maximum possible values on the curve
		// (because the current fan speed can be outside of this range)
		min_speed := interpolate(math.MinInt, fan.controlCurve)
		max_speed := interpolate(math.MaxInt, fan.controlCurve)
		speed = minMax(min_speed, speed, max_speed)
		// Set target speed
		fan.targetSpeed = speed
		controller.fans[id] = fan
	}
}

// Pushes target fan speeds to hardware.
func (controller *FanController) pushTargetFanSpeeds() {
	attributes := make([]Attribute, len(controller.fans))
	i := 0
	for id, fan := range controller.fans {
		attributes[i].name = fmt.Sprintf("[fan:%d]/GPUTargetFanSpeed", id)
		attributes[i].value = fan.targetSpeed
		i++
	}
	assignAttributes(attributes, controller.xDisplay)
}
