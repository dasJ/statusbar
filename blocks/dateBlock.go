package blocks

import (
	"fmt"
	"github.com/dasJ/statusbar"
	"time"
)

type DateBlock struct {
	block *statusbar.I3Block
}

func (this *DateBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block
	return true
}

func (this DateBlock) Tick() {
	now := time.Now()
	_, kw := now.ISOWeek()
	this.block.FullText = fmt.Sprint("(KW", kw, ") ", now.Format("Mon, 02. Jan 2006 15:04"))
	this.block.ShortText = now.Format("15:04")
}

func (this DateBlock) Click(data statusbar.I3Click) {
}

func (this DateBlock) Block() *statusbar.I3Block {
	return this.block
}
