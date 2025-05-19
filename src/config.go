package main

import (
	"encoding/json"
	"log"
	"os"
)

const CURVE_SPEED_MAX = 100
const CURVE_SPEED_MIN = 0
const CURVE_TEMP_MAX = 90
const CURVE_TEMP_MIN = 30

type Config struct {
	Fans  map[int]ConfigFan `json:"fans"`
	Graph ConfigGraph       `json:"graph"`
}

type ConfigFan struct {
	GpuId        int         `json:"gpu_id"`
	ControlCurve map[int]int `json:"control_curve"`
}

type ConfigGraph struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// Loads a config file
func loadConfig(configFile string) Config {
	configText, err := os.ReadFile(configFile)
	if err != nil {
		log.Panic(err)
	}

	var config Config

	err = json.Unmarshal(configText, &config)
	if err != nil {
		log.Panic(err)
	}

	sanityCheckConfig(&config)

	return config
}

// Verifies a config is set correctly
func sanityCheckConfig(config *Config) {
	if config.Graph.Max < config.Graph.Min {
		log.Panicf(
			"Graph maximum '%d' is less than the graph minimum '%d'.",
			config.Graph.Min,
			config.Graph.Max,
		)
	} else if config.Graph.Max-config.Graph.Min < 10 {
		log.Panicf(
			"Graph minimum '%d' and maximum '%d' are less than 10C away from each other.",
			config.Graph.Min,
			config.Graph.Max,
		)
	}

	// Build a map of every fan and gpu according to nvidia-settings
	nvidiaFans := getFans()
	nvidiaGpus := map[int]bool{}
	for _, nvidiaGpuId := range nvidiaFans {
		nvidiaGpus[nvidiaGpuId] = true
	}

	// For each fan, according to nvidia
	for nvidiaFanId, _ := range nvidiaFans {
		// Verify every fan is configured
		_, ok := config.Fans[nvidiaFanId]
		if !ok {
			log.Panicf(
				"Missing config for fan %d, which exists according to nvidia-settings.",
				nvidiaFanId,
			)
		}
	}

	// For each fan, according to the config
	for configFanId, configFan := range config.Fans {
		// Verify every fan configured exists
		nvidiaGpuId, ok := nvidiaFans[configFanId]
		if !ok {
			log.Panicf(
				"Configuration for fan %d found, but fan %d does not exist according to nvidia-settings.",
				configFanId,
				configFanId,
			)
		}

		// Verify the fan is configured to monitor the correct GPU's temp
		if configFan.GpuId != nvidiaGpuId {
			log.Panicf(
				"Fan %d is configured to monitor GPU %d's temperature, but belongs to GPU %d according to nvidia-settings.",
				configFanId,
				configFan.GpuId,
				nvidiaGpuId,
			)
		}

		// Verify at least 1 point has been added to the control curve
		if len(configFan.ControlCurve) == 0 {
			log.Panicf(
				"Fan %d's control curve doesn't contain any points.",
				configFanId,
			)
		}

		// For each point in the control curve, sorted by temperature
		lastTemp := CURVE_TEMP_MIN
		lastSpeed := CURVE_SPEED_MIN
		for _, temp := range getSortedIndexes(&configFan.ControlCurve) {
			speed, ok := configFan.ControlCurve[temp]
			if !ok {
				log.Fatalf("Failed to read control curve data at %d", temp)
			}

			// Temperature checks
			if temp < CURVE_TEMP_MIN {
				log.Fatalf(
					"Fan %d's control curve contains (%dC, %d%%), which is below %d degrees.",
					configFanId,
					temp,
					speed,
					CURVE_TEMP_MIN,
				)
			} else if temp > CURVE_TEMP_MAX {
				log.Fatalf(
					"Fan %d's control curve contains (%dC, %d%%), which is above %d degrees.",
					configFanId,
					temp,
					speed,
					CURVE_TEMP_MAX,
				)
			} else if temp < lastTemp {
				// Technically this shouldn't happen because we sorted by temp.
				// ...but just in case...
				log.Fatalf(
					"Fan %d's control curve contains (%dC, %d%%), which is below the temperature of the previous point, (%dC, %d%%).",
					configFanId,
					temp,
					speed,
					lastTemp,
					lastSpeed,
				)
			}

			// Speed checks
			if speed < CURVE_SPEED_MIN {
				log.Fatalf(
					"Fan %d's control curve contains (%dC, %d%%), which is below %d%% fan speed.",
					configFanId,
					temp,
					speed,
					CURVE_SPEED_MIN,
				)
			} else if speed > CURVE_SPEED_MAX {
				log.Fatalf(
					"Fan %d's control curve contains (%dC, %d%%), which is above %d%% fan speed.",
					configFanId,
					temp,
					speed,
					CURVE_SPEED_MAX,
				)
			} else if speed < lastSpeed {
				log.Fatalf(
					"Fan %d's control curve contains (%dC, %d%%), which is below the speed of the previous point, (%dC, %d%%).",
					configFanId,
					temp,
					speed,
					lastTemp,
					lastSpeed,
				)
			}

			lastTemp, lastSpeed = temp, speed
		}
	}
}
