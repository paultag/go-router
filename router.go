package router

import (
	"net"
)

type Router struct {
	RouteTable RouteTable
	ARPTable   ARPTable
}

//
func (r Router) Lookup(dest net.IP) (*Route, *ARP) {
	route := r.RouteTable.Lookup(dest)
	if route == nil {
		return nil, nil
	}
	arp := r.ARPTable.Lookup(route.Gateway)
	return route, arp
}

func NewRouter(routeTable RouteTable, arpTable ARPTable) Router {
	return Router{
		RouteTable: routeTable,
		ARPTable:   arpTable,
	}
}

func New() (*Router, error) {
	arp, err := LoadARPTable()
	if err != nil {
		return nil, err
	}
	route, err := LoadRouteTable()
	if err != nil {
		return nil, err
	}
	r := NewRouter(route, arp)
	return &r, nil
}
