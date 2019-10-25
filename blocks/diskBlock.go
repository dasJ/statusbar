package blocks

import (
	"github.com/SlothOfAnarchy/statusbar"
	"syscall"
)

type DiskBlock struct {
	block  *statusbar.I3Block
	failed bool
}

func (this *DiskBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block
	return true
}

func (this DiskBlock) Tick() {
	if this.failed {
		return
	}

	var stat syscall.Statfs_t
	syscall.Statfs("/", &stat)

	if float64(stat.Bavail) < float64(stat.Bsize)*float64(0.1) {
		this.block.Color = "#ff0202"
	} else {
		this.block.Color = ""
	}

	this.block.FullText = ByteSize(stat.Bavail * uint64(stat.Bsize))
}

func (this DiskBlock) Click(data statusbar.I3Click) {
}

func (this DiskBlock) Block() *statusbar.I3Block {
	return this.block
}
