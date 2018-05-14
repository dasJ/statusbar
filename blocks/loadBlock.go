package blocks

import (
	"fmt"
	"github.com/dasJ/statusbar"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type LoadBlock struct {
	block  *statusbar.I3Block
	failed bool
}

func (this *LoadBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block
	return true
}

func (this LoadBlock) Tick() {
	if this.failed {
		return
	}

	raw, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
		return
	}
	loadStr := strings.Split(string(raw), " ")[0]
	load, err := strconv.ParseFloat(loadStr, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
		return
	}
	if load > 10 {
		this.block.Color = "#ff0202"
	} else {
		this.block.Color = ""
	}

	this.block.FullText = fmt.Sprintf("%.2f", load)
}

func (this LoadBlock) Click(data statusbar.I3Click) {
}

func (this LoadBlock) Block() *statusbar.I3Block {
	return this.block
}
