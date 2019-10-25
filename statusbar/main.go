package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/SlothOfAnarchy/statusbar"
	"github.com/SlothOfAnarchy/statusbar/blocks"
	"os"
	"strings"
	"time"
)

func main() {
	// Header
	header, err := json.Marshal(statusbar.I3Header{1, 10, 12, true})
	if err != nil {
		panic(err)
	}
	fmt.Print(string(header))
	fmt.Println("[") // Begin stream

	// Initialize responder
	resp := statusbar.Responder{}

	// Initialize blocks
	resp.AppendBlock(&blocks.VolumeBlock{})
	resp.AppendBlock(&blocks.MemoryBlock{})
	resp.AppendBlock(&blocks.DiskBlock{})
	resp.AppendBlock(&blocks.BatteryBlock{})
	resp.AppendBlock(&blocks.NotmuchBlock{})
	resp.AppendBlock(&blocks.NetworkBlock{})
	resp.AppendBlock(&blocks.IpBlock{})
	resp.AppendBlock(&blocks.SsidBlock{})
	resp.AppendBlock(&blocks.LoadBlock{})
	resp.AppendBlock(&blocks.TempBlock{})
	resp.AppendBlock(&blocks.DateBlock{})

	// Click handler
	go func(resp statusbar.Responder) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			txt := scanner.Text()

			// Cleanup
			if txt == "[" {
				continue
			}
			txt = strings.TrimPrefix(txt, ",")

			go resp.HandleClick(txt)
		}
	}(resp)

	// Main loop
	for {
		resp.TickAll()
		resp.Output()
		time.Sleep(2000 * time.Millisecond)
	}
}
