package statusbar

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Responder struct {
	blocks    []*Block
	currentID int
}

func (this *Responder) AppendBlock(block Block) {
	if block.Init(&I3Block{}, this) {
		this.blocks = append(this.blocks, &block)
		block.Block().Name = strconv.Itoa(this.currentID)
		this.currentID++
	}
}

func (this *Responder) TickAll() {
	for _, block := range this.blocks {
		(*block).Tick()
	}
}

func (this *Responder) HandleClick(raw string) {
	var data I3Click
	err := json.Unmarshal([]byte(raw), &data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Dropping invalid click event from i3\n")
		return
	}
	blockID, err := strconv.Atoi(data.Name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Dropping invalid block name from i3\n")
		return
	}
	(*this.blocks[blockID]).Click(data)
}

func (this *Responder) Output() {
	fmt.Print("[")

	first := true
	for _, block := range this.blocks {
		if first {
			first = false
		} else {
			fmt.Print(",")
		}

		ret, err := json.Marshal((*block).Block())
		if err != nil {
			panic(err)
		}

		fmt.Print(string(ret))
	}
	fmt.Println("],")
}
