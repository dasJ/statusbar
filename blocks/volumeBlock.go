package blocks

import (
	"fmt"
	"github.com/dasJ/statusbar"
	"os/exec"

	//#cgo LDFLAGS: -lpulse
	//#include <volblock.h>
	"C"
)

type VolumeBlock struct {
	block *statusbar.I3Block
	responder *statusbar.Responder
	failed bool
	restart chan bool
}

var volumeBlock *VolumeBlock

func (this *VolumeBlock) autorestartPulse() {
	for (true) {
		// Wait for the restart channel
		<-this.restart
		// Initialize and check result
		out := C.GoString(C.initPulse())
		if out != "" {
			this.block.FullText = out
			this.block.Color = "#ff0202"
			this.failed = true
		}
		// Run the main loop
		go C.runPulse()
	}
}

func (this *VolumeBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block
	this.responder = resp
	this.failed = false
	this.restart = make(chan bool, 1)
	volumeBlock = this

	// Launch the autorestarter
	go this.autorestartPulse()
	// Trigger the initial autostart
	this.restart <- true

	return true
}

func (this VolumeBlock) Tick() {
}

func (this VolumeBlock) Click(data statusbar.I3Click) {
	if data.Button == 1 {
		exec.Command("pavucontrol").Start()
	} else if data.Button == 3 {
		C.toggleMute()
	} else if data.Button == 4 {
		C.setVolume(1, 5)
	} else if data.Button == 5 {
		C.setVolume(0, 5)
	}
}

func (this VolumeBlock) Block() *statusbar.I3Block {
	return this.block
}

//export goPulseRestart
func goPulseRestart() {
	volumeBlock.restart <- true
}

//export goPulseError
func goPulseError(c *C.char) {
	volumeBlock.failed = true
	goPulseMsg(c)
}

//export goPulseMsg
func goPulseMsg(c *C.char) {
	volumeBlock.block.FullText = C.GoString(c)
	volumeBlock.block.Color = "#ff0202"
	volumeBlock.responder.Output()
}

//export goPulseVol
func goPulseVol(mute C.char, vol C.int) {
	if mute == 1 {
		volumeBlock.block.FullText = "muted"
		volumeBlock.responder.Output()
		return
	}
	volumeBlock.block.Color = ""
	volumeBlock.block.FullText = fmt.Sprintf("%d%%", int(vol))
	volumeBlock.responder.Output()
}
