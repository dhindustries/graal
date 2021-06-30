package pathfinder

type HeuristicFn func(a, b Node) float64

type Finder interface {
	FindPath(source, destination Node) []Node
}
