package router

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type ARP struct {
	IP     net.IP
	MAC    net.HardwareAddr
	Mask   string
	Device string
}

func (a ARP) Interface() (*net.Interface, error) {
	return net.InterfaceByName(a.Device)
}

func (a ARP) String() string {
	return fmt.Sprintf(
		"IP=%s MAC=%s Mask=%s Device=%s",
		a.IP,
		a.MAC,
		a.Mask,
		a.Device,
	)
}

func NewARPFromProc(line []string) (*ARP, error) {
	if len(line) < 5 {
		return nil, fmt.Errorf("Line is too short")
	}

	arp := ARP{
		IP:     net.ParseIP(line[0]),
		Mask:   line[4],
		Device: line[5],
	}

	var err error
	arp.MAC, err = net.ParseMAC(line[3])
	if err != nil {
		return nil, err
	}

	return &arp, nil
}

type ARPTable []ARP

func (a ARPTable) Lookup(q net.IP) *ARP {
	for _, entry := range a {
		if entry.IP.Equal(q) {
			return &entry
		}
	}
	return nil
}

func NewARPTable(r io.Reader) (ARPTable, error) {
	table := ARPTable{}

	s := bufio.NewScanner(r)
	s.Scan()

	for s.Scan() {
		entry, err := NewARPFromProc(strings.Fields(s.Text()))
		if err != nil {
			return nil, err
		}
		table = append(table, *entry)
	}
	return table, nil
}

func LoadARPTable() (ARPTable, error) {
	// XXX: support non-linux or something I guess maybe?
	fd, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return NewARPTable(fd)
}
