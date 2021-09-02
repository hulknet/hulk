package router

type RoutingTable struct {
	main Table
	apps map[string]Table
}

func NewRoutingTable() *RoutingTable {
	return &RoutingTable{}
}

func (t *RoutingTable) Find(idPrefix string) error {
	return nil
}
