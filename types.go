package statusbar

type I3Header struct {
	Version     int16 `json:"version"`
	StopSignal  int16 `json:"stop_signal"`
	ContSignal  int16 `json:"cont_signal"`
	ClickEvents bool  `json:"click_events"`
}

type I3Block struct {
	FullText  string `json:"full_text"`
	ShortText string `json:"short_text,omitempty"`
	Color     string `json:"color,omitempty"`
	Urgent    bool   `json:"urgent,omitempty"`
	Name      string `json:"name"`
}

type I3Click struct {
	Name      string `json:"name"`
	AbsoluteX int    `json:"x"`
	AbsoluteY int    `json:"y"`
	RelativeX int    `json:"relative_x"`
	RelativeY int    `json:"relative_y"`
	Button    int    `json:"button"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}

type Block interface {
	Init(ret *I3Block, resp *Responder) bool
	Tick()
	Click(data I3Click)
	Block() *I3Block
}
