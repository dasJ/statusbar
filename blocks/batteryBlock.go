package blocks

import (
	"fmt"
	"github.com/SlothOfAnarchy/statusbar"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type BatteryBlock struct {
	block *statusbar.I3Block
}

func (this *BatteryBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block

	// Try to figure out if we have a battery
	if _, err := os.Stat("/sys/class/power_supply"); os.IsNotExist(err) {
		return false // No power supplies
	}
	supplies, err := ioutil.ReadDir("/sys/class/power_supply")
	if err != nil {
		return false // What?
	}
	for _, supply := range supplies {
		// Is this a battery?
		if strings.HasPrefix(supply.Name(), "BAT") {
			return true
		}
	}

	return false
}

func (this BatteryBlock) Tick() {
	supplies, err := ioutil.ReadDir("/sys/class/power_supply")
	if err != nil {
		return // What?
	}
	this.block.FullText = ""
	for _, supply := range supplies {
		// Is this a battery?
		if !strings.HasPrefix(supply.Name(), "BAT") {
			continue
		}
		if _, err = os.Stat("/sys/class/power_supply/" + supply.Name() + "/capacity"); err != nil {
			continue
		}
		// Read file
		rawCapacity, err := ioutil.ReadFile("/sys/class/power_supply/" + supply.Name() + "/capacity")
		if err != nil {
			continue
		}
		// Deal with the contents
		cap, err := strconv.Atoi(strings.TrimSpace(string(rawCapacity)))
		if err != nil {
			continue
		}
		if cap <= 15 {
			this.block.Color = "#ff0202"
		} else {
			this.block.Color = ""
		}
		if this.block.FullText == "" {
			this.block.FullText = fmt.Sprintf("%d%%", cap)
		} else {
			this.block.FullText += fmt.Sprintf(" %d%%", cap)
		}
	}

	for _, supply := range supplies {
		// Is this a power supply?
		if !strings.HasPrefix(supply.Name(), "AC") {
			continue
		}
		if _, err = os.Stat("/sys/class/power_supply/" + supply.Name() + "/online"); err != nil {
			continue
		}
		// Read file
		rawOnline, err := ioutil.ReadFile("/sys/class/power_supply/" + supply.Name() + "/online")
		if err != nil {
			continue
		}
		// Deal with the contents
		if string(rawOnline) == "1\n" {
			this.block.FullText += " +"
			this.block.Color = "#02ff02"
		}
	}
}

func (this BatteryBlock) Click(data statusbar.I3Click) {
}

func (this BatteryBlock) Block() *statusbar.I3Block {
	return this.block
}
