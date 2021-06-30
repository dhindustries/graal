package pathfinder

type Node interface {
	Equals(node Node) bool
	Neighbors() []Connection
}

type Connection struct {
	Node   Node
	Weight float64
}
