package router

type Router struct {
	RouteTable RouteTable
	ARPTable   ARPTable
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
