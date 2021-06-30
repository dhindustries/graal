package pathfinder

import "sort"

type Astar struct {
	Bias HeuristicFn
}

type astarNode struct {
	Node
	weight  float64
	prev    *astarNode
	visited bool
}

type queueItem struct {
	node     *astarNode
	priority float64
}

type queue struct {
	items []queueItem
}

type byPriority []queueItem

func (a byPriority) Len() int           { return len(a) }
func (a byPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool { return a[i].priority < a[j].priority }

func (q *queue) push(f float64, n *astarNode) {
	if q.items == nil {
		q.items = make([]queueItem, 0)
	}
	q.items = append(q.items, queueItem{n, f})
}

func (q *queue) pop() *astarNode {
	q.sort()
	if q.items != nil && len(q.items) > 0 {
		i := q.items[0]
		q.items = q.items[1:]
		return i.node
	}
	return nil
}

func (q *queue) sort() {
	if q.items != nil && len(q.items) > 0 {
		sort.Sort(byPriority(q.items))
	}
}

func (q *queue) empty() bool {
	return q.items == nil || len(q.items) == 0
}

func (finder *Astar) FindPath(src, dst Node) []Node {
	if f := finder.findFinal(src, dst); f != nil {
		return finder.buildQueue(f)
	}
	return []Node{}
}

func (finder *Astar) findFinal(src, dst Node) *astarNode {
	q := new(queue)
	nodes := make(map[Node]*astarNode)
	var bias HeuristicFn = finder.Bias
	if bias == nil {
		bias = func(a, b Node) float64 {
			return 0
		}
	}

	nodeOf := func(node Node) *astarNode {
		v, ok := nodes[node]
		if !ok {
			v = &astarNode{Node: node}
			nodes[node] = v
		}
		return v
	}

	q.push(0, nodeOf(src))
	for !q.empty() {
		node := q.pop()
		node.visited = true
		if node.Node == dst {
			return node
		}
		for _, conn := range node.Neighbors() {
			bval := bias(node.Node, conn.Node)
			weight := node.weight + conn.Weight
			neighbor := nodeOf(conn.Node)
			if !neighbor.visited {
				neighbor.weight = weight
				neighbor.prev = node
				neighbor.visited = true
				q.push(bval, neighbor)
			} else if neighbor.weight > weight {
				neighbor.weight = weight
				neighbor.prev = node
			}
		}
	}
	return nil
}

func (*Astar) buildQueue(node *astarNode) []Node {
	list := make([]Node, 0)
	for node != nil {
		list = append([]Node{node.Node}, list...)
		node = node.prev
	}
	return list[1:]
}
