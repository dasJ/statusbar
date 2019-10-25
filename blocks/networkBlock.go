package blocks

import (
	"bufio"
	"fmt"
	"github.com/SlothOfAnarchy/statusbar"
	"os"
	"strings"
)

type NetworkBlock struct {
	block  *statusbar.I3Block
	failed bool
}

func (this *NetworkBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block
	return true
}

func (this NetworkBlock) Tick() {
	if this.failed {
		return
	}

	this.block.FullText = ""
	this.block.Color = ""

	f, err := os.Open("/proc/net/route")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "\t")
		// Only consider default routes (to 0.0.0.0)
		if parts[1] == "00000000" {
			iface := parts[0]
			if this.block.FullText == "" {
				this.block.FullText = iface
			} else if !strings.Contains(this.block.FullText, iface) {
				this.block.FullText += " " + iface
			}
		}
	}
	if this.block.FullText == "" {
		this.block.FullText = "No link"
		this.block.Color = "#ff0202"
	}
}

func (this NetworkBlock) Click(data statusbar.I3Click) {
}

func (this NetworkBlock) Block() *statusbar.I3Block {
	return this.block
}
