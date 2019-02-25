package main

import (
	"fmt"
	"net"
	"os"

	"pault.ag/go/router"
)

func main() {
	r, err := router.New()
	if err != nil {
		panic(err)
	}

	ips := os.Args[1:]

	if len(ips) == 0 {
		for _, route := range r.RouteTable {
			fmt.Printf("%s\n", route.String())
		}
	}

	for _, ip := range ips {
		ipAddr := net.ParseIP(ip)

		route := r.RouteTable.Lookup(ipAddr)
		fmt.Printf("%s\n", route.String())
	}
}
