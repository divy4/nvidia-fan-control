package main

import (
	"fmt"
	"os"
)

func main() {
	command := parse_args()

	control_curve := map[int]int{
		30: 35,
		60: 50,
		70: 100,
	}
	fans := map[int]Fan{
		0: Fan{
			gpu_id:        0,
			min_speed:     35,
			max_speed:     100,
			control_curve: control_curve,
		},
		1: Fan{
			gpu_id:        0,
			min_speed:     35,
			max_speed:     100,
			control_curve: control_curve,
		},
	}

	// TODO: figure out signal interrupts

	controller := create_fan_controller(fans, 30, 100)

	switch command {
	case "run":
		controller.run()
	case "stop":
		controller.disable_fan_control()
	default:
		print_help()
		os.Exit(1)
	}
}

func parse_args() string {
	if len(os.Args) != 3 {
		print_help()
		os.Exit(1)
	}

	command := os.Args[1]
	if command != "run" && command != "stop" {
		print_help()
		os.Exit(1)
	}

	return command
}

func print_help() {
	fmt.Println("Usage: nvidia-fan-control run|stop")
}
