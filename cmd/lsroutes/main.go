package main

import (
	"fmt"

	"pault.ag/go/router"
)

func main() {
	r, err := router.New()
	if err != nil {
		panic(err)
	}
	for _, route := range r.RouteTable {
		arpString := "unknown mac"
		arp := r.ARPTable.Lookup(route.Gateway)
		if arp != nil {
			arpString = arp.MAC.String()
		}

		fmt.Printf(
			"%s %s via %s (%s)\n",
			route.Interface,
			route.IPNet().String(),
			route.Gateway,
			arpString,
		)
	}
}
