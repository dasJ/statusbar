package blocks

import (
	"github.com/SlothOfAnarchy/statusbar"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type TempBlock struct {
	block      *statusbar.I3Block
	sensorPath string
	highTemp   int
}

func (this *TempBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block

	var preferredSensors = [...]string{ /* Thinkpad */ "coretemp"}

	// Try to figure out the correct sensor
	if _, err := os.Stat("/sys/class/hwmon"); os.IsNotExist(err) {
		return false // No hwmon
	}
	sensors, err := ioutil.ReadDir("/sys/class/hwmon")
	if err != nil {
		return false // What?
	}
SensorSearch:
	for _, sensor := range sensors {
		// Check if a we have a temperature sensor
		if _, err := os.Stat("/sys/class/hwmon/" + sensor.Name() + "/temp1_input"); os.IsNotExist(err) {
			continue // No temperature sensor here
		}

		// Use this sensor
		this.sensorPath = "/sys/class/hwmon/" + sensor.Name() + "/temp1_input"

		// Try to find high temperature
		if _, err := os.Stat("/sys/class/hwmon/" + sensor.Name() + "/temp1_max"); os.IsNotExist(err) {
			rawTemp, err := ioutil.ReadFile("/sys/class/hwmon/" + sensor.Name() + "/temp1_max")
			if err != nil {
				this.highTemp = 0
			} else {
				this.highTemp, _ = strconv.Atoi(string(rawTemp))
			}
		} else {
			this.highTemp = 0
		}

		// Handle name (preferred sensors)
		rawName, err := ioutil.ReadFile(this.sensorPath)
		if err != nil {
			continue
		}
		name := string(rawName)
		for _, cur := range preferredSensors {
			if name == cur {
				break SensorSearch
			}
		}
	}

	return true
}

func (this TempBlock) Tick() {
	rawTemp, err := ioutil.ReadFile(this.sensorPath)
	if err != nil {
		this.block.FullText = "ERROR"
		this.block.Color = "#ff0202"
		return
	}
	temp, err := strconv.Atoi(strings.TrimSpace(string(rawTemp)))
	if err != nil {
		this.block.FullText = "ERROR"
		this.block.Color = "#ff0202"
		return
	}
	if this.highTemp > 0 && temp >= this.highTemp {
		this.block.Color = "#ff0202"
	} else {
		this.block.Color = ""
	}

	this.block.FullText = strconv.Itoa(temp/1000) + "°C"
}

func (this TempBlock) Click(data statusbar.I3Click) {
}

func (this TempBlock) Block() *statusbar.I3Block {
	return this.block
}
