package blocks

import (
	"bufio"
	"fmt"
	"github.com/SlothOfAnarchy/statusbar"
	"net"
	"os"
	"strings"
)

type IpBlock struct {
	block  *statusbar.I3Block
	failed bool
}

func (this *IpBlock) Init(block *statusbar.I3Block, resp *statusbar.Responder) bool {
	this.block = block
	return true
}

func (this IpBlock) Tick() {
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
			ipAddress := this.GetInterfaceIp(parts[0])
			if this.block.FullText == "" {
				this.block.FullText = ipAddress
			} else {
				this.block.FullText += " " + ipAddress
			}
		}
	}
	if this.block.FullText == "" {
		this.block.FullText = "No IP"
		this.block.Color = "#ff0202"
	}
}

func (this IpBlock) GetInterfaceIp(ifaceName string) string {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
	}
	addrs, err := iface.Addrs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		this.failed = true
	}
	for _, address := range addrs {
		// check the address type and that it's not a loopback
		if ipnet, err := address.(*net.IPNet); err && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (this IpBlock) Click(data statusbar.I3Click) {
}

func (this IpBlock) Block() *statusbar.I3Block {
	return this.block
}
