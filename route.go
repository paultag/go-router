package router

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type Route struct {
	Interface   string
	Destination net.IP
	Gateway     net.IP
	// Flags
	// RefCnt
	// Use
	// Metric
	Mask net.IPMask
}

func (r Route) String() string {
	return fmt.Sprintf(
		"%s via %s dev %s",
		r.IPNet().String(),
		r.Gateway,
		r.Interface,
	)
}

func (r Route) IPNet() *net.IPNet {
	return &net.IPNet{IP: r.Destination, Mask: r.Mask}
}

// XXX: combine this with the IP helper below
func parseHexIPMask(h string) (net.IPMask, error) {
	r := make(net.IPMask, 4)
	num, err := strconv.ParseUint(h, 16, 64)
	if err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint32(r, uint32(num))
	return r, nil
}

func parseHexIP(h string) (net.IP, error) {
	r := make(net.IP, 4)
	num, err := strconv.ParseUint(h, 16, 64)
	if err != nil {
		return nil, err
	}
	binary.LittleEndian.PutUint32(r, uint32(num))
	return r, nil
}

func NewRouteFromProc(line []string) (*Route, error) {
	if len(line) < 8 {
		return nil, fmt.Errorf("Line is too short")
	}

	destination, err := parseHexIP(line[1])
	if err != nil {
		return nil, err
	}
	gateway, err := parseHexIP(line[2])
	if err != nil {
		return nil, err
	}

	mask, err := parseHexIPMask(line[7])
	if err != nil {
		return nil, err
	}

	route := Route{
		Interface:   line[0],
		Destination: destination,
		Gateway:     gateway,
		Mask:        mask,
	}

	return &route, nil
}

type RouteTable []Route

func (rt RouteTable) Lookup(q net.IP) *Route {
	for _, entry := range rt {
		net := entry.IPNet()
		if net.Contains(q) {
			return &entry
		}
	}
	return nil
}

func NewRouteTable(r io.Reader) (RouteTable, error) {
	table := RouteTable{}

	s := bufio.NewScanner(r)
	s.Scan()

	for s.Scan() {
		entry, err := NewRouteFromProc(strings.Fields(s.Text()))
		if err != nil {
			return nil, err
		}
		table = append(table, *entry)
	}
	return table, nil
}

func LoadRouteTable() (RouteTable, error) {
	// XXX: support non-linux or something I guess maybe?
	fd, err := os.Open("/proc/net/route")
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	return NewRouteTable(fd)
}
