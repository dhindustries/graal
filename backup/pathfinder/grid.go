package pathfinder

import "math"

type Grid struct {
	algorithm  Finder
	avails     []bool
	weights    []float64
	dirs       [][2]int
	minx, maxx int
	miny, maxy int
	cells      map[int]map[int]*gridCell
}

type gridCell struct {
	grid *Grid
	x, y int
}

func NewGrid(minX, minY, maxX, maxY int) *Grid {
	w, h := maxX-minX, maxY-minY
	s := w * h
	weights := make([]float64, s)
	for i := range weights {
		weights[i] = 1
	}
	return &Grid{
		avails:  make([]bool, s),
		weights: weights,
		dirs: [][2]int{
			[2]int{0, -1},
			[2]int{-1, 0},
			[2]int{0, 1},
			[2]int{1, 0},
		},
		minx: minX, maxx: maxX,
		miny: minY, maxy: maxY,
	}
}

func (grid *Grid) FindPath(sx, sy, dx, dy int) [][2]int {
	algorithm := grid.algorithm
	if algorithm == nil {
		algorithm = &Astar{grid.Heuristics}
	}
	nodes := algorithm.FindPath(grid.node(sx, sy), grid.node(dx, dy))
	pos := make([][2]int, len(nodes))
	for i, node := range nodes {
		cell := node.(*gridCell)
		pos[i] = [2]int{cell.x, cell.y}
	}
	return pos
}

func (grid *Grid) Heuristics(a, b Node) float64 {
	if s, ok := a.(*gridCell); ok {
		if e, ok := b.(*gridCell); ok {
			dx, dy := float64(e.x-s.x), float64(e.y-s.y)
			return math.Sqrt(dx*dx + dy*dy)
		}
	}
	return 0
}

func (grid *Grid) SetAvailable(x, y int, v bool) {
	if i, ok := grid.index(x, y); ok {
		grid.avails[i] = v
	}
}

func (grid *Grid) SetWeight(x, y int, v float64) {
	if i, ok := grid.index(x, y); ok {
		grid.weights[i] = v
	}
}

func (grid *Grid) node(x, y int) *gridCell {
	if grid.cells == nil {
		grid.cells = make(map[int]map[int]*gridCell)
	}
	if _, ok := grid.cells[y]; !ok {
		grid.cells[y] = make(map[int]*gridCell)
	}
	cell, ok := grid.cells[y][x]
	if !ok {
		cell = &gridCell{
			grid: grid,
			x:    x,
			y:    y,
		}
		grid.cells[y][x] = cell
	}
	return cell
}

func (grid *Grid) available(x, y int) bool {
	if i, ok := grid.index(x, y); ok {
		return grid.avails[i]
	}
	return false
}

func (grid *Grid) weight(x, y int) float64 {
	if i, ok := grid.index(x, y); ok {
		return grid.weights[i]
	}
	return 0
}

func (grid *Grid) index(x, y int) (int, bool) {
	if x >= grid.minx && x <= grid.maxx && y >= grid.miny && y <= grid.maxy {
		w := grid.maxx - grid.minx
		return y*w + x, true
	}
	return 0, false
}

func (cell *gridCell) Equals(node Node) bool {
	n := node.(*gridCell)
	return n.x == cell.x && n.y == cell.y
}

func (cell *gridCell) Neighbors() []Connection {
	l := make([]Connection, 0, len(cell.grid.dirs))
	for _, dir := range cell.grid.dirs {
		x, y := cell.x+dir[0], cell.y+dir[1]
		if cell.grid.available(x, y) {
			l = append(l, Connection{
				Node:   cell.grid.node(x, y),
				Weight: cell.grid.weight(x, y),
			})
		}
	}
	return l
}
