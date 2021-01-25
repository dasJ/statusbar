package blocks

import (
	"bufio"
	"fmt"
	"github.com/SlothOfAnarchy/statusbar"
	"os"
	"os/exec"
	"strings"
)

type SsidBlock struct {
	block  *statusbar.I3Block
	failed bool
}

func (this *SsidBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block
	return true
}

func (this SsidBlock) Tick() {
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
			ssid := ""
			if this.CommandExists("nmcli") {
				ssid = this.GetNmSsid(iface)
			} else if this.CommandExists("wpa_cli") {
				ssid = this.GetWpaSupplSsid()
			}
			if this.block.FullText == "" {
				this.block.FullText = ssid
			} else {
				this.block.FullText += " " + ssid
			}
		}
	}
	if this.block.FullText == "" {
		this.block.FullText = "-"
	}
}

func (this SsidBlock) GetWpaSupplSsid() string {
	out, err := exec.Command("bash", "-c", "wpa_cli status | grep 'ssid=' | cut -d'=' -f2").Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
		return ""
	}
	outString := string(out)
	if len([]rune(outString)) > 0 {
		return outString
	}
	return ""
}
func (this SsidBlock) GetNmSsid(iface string) string {
	out, err := exec.Command("bash", "-c", "nmcli connection show | grep "+iface+" | grep wifi").Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
		return ""
	}
	outString := string(out)
	if strings.Contains(outString, iface) {
		return strings.Fields(outString)[0]
	}
	return ""
}

func (this SsidBlock) CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (this SsidBlock) Click(data statusbar.I3Click) {
}

func (this SsidBlock) Block() *statusbar.I3Block {
	return this.block
}
