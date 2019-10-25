package blocks

import (
	"fmt"
	"github.com/dasJ/statusbar"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type MemoryBlock struct {
	block  *statusbar.I3Block
	failed bool
}

func (this *MemoryBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block
	return true
}

func (this MemoryBlock) Tick() {
	if this.failed {
		return
	}

	raw, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
		return
	}
	meminfoLines := strings.Split(string(raw), "\n")
	totalMemString := strings.Fields(meminfoLines[0])[1]
	availableMemString := strings.Fields(meminfoLines[2])[1]
	total, err := strconv.Atoi(totalMemString)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
		return
	}
	available, err := strconv.Atoi(availableMemString)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
		return
	}
	if float64(available) < float64(total)*float64(0.1) {
		this.block.Color = "#ff0202"
	} else {
		this.block.Color = ""
	}

	// meminfo output is in kB
	this.block.FullText = *statusbar.ByteSize(uint64(available * 1024))
}

func (this MemoryBlock) Click(data statusbar.I3Click) {
}

func (this MemoryBlock) Block() *statusbar.I3Block {
	return this.block
}
