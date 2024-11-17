package game

import (
	"fmt"
	"math/rand"
)

const (
	width  = 7
	height = 9
)

type Color int

const (
	Orange Color = iota
	Blue
	Green
	Pink
	Empty
)

func (c Color) Icon() string {
	var icon string
	switch c {
	case Orange:
		icon = "🟧"
	case Blue:
		icon = "🟦"
	case Green:
		icon = "🟩"
	case Pink:
		icon = "🟪"
	default:
		icon = "⬜"
	}
	return icon
}

type Cell struct {
	color Color
}
type Action struct {
	X, Y  int
	Color Color
}

type Grid struct {
	cells   [height][width]Cell
	actions []Action
}

func NewGrid() *Grid {
	grid := &Grid{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			grid.cells[y][x] = Cell{color: Color(rand.Intn(4))}
		}
	}
	return grid
}

func (g *Grid) Print() {

	for y := height - 1; y >= 0; y-- {
		// Print row header
		fmt.Printf("%2d ", y)
		for x := 0; x < width; x++ {
			fmt.Print(g.cells[y][x].color.Icon(), " ")
		}
		fmt.Println()
	}
	// Print column headers
	fmt.Print(" ")
	for x := 0; x < width; x++ {
		fmt.Printf("%3d", x)
	}
	fmt.Println()
}

func (g *Grid) ClearRegion(x, y int) {
	if x < 0 || x >= width || y < 0 || y >= height {
		return
	}

	targetColor := g.cells[y][x].color
	if targetColor == Empty {
		return
	}

	g.clearCell(x, y, targetColor)
	g.ApplyGravity()
	g.actions = append(g.actions, Action{X: x, Y: y, Color: targetColor})
}

func (g *Grid) clearCell(x, y int, targetColor Color) {
	if x < 0 || x >= width || y < 0 || y >= height {
		return
	}

	if g.cells[y][x].color != targetColor {
		return
	}

	g.cells[y][x].color = Empty

	// Recursively clear adjacent cells and count them
	g.clearCell(x+1, y, targetColor)
	g.clearCell(x-1, y, targetColor)
	g.clearCell(x, y+1, targetColor)
	g.clearCell(x, y-1, targetColor)

}

func (g *Grid) ApplyGravity() {
	for x := 0; x < width; x++ {
		writeIdx := 0
		for y := 0; y < height; y++ {
			if g.cells[y][x].color != Empty {
				g.cells[writeIdx][x] = g.cells[y][x]
				if writeIdx != y {
					g.cells[y][x].color = Empty
				}
				writeIdx++
			}
		}
		for y := writeIdx; y < height; y++ {
			g.cells[y][x].color = Empty
		}
	}
}

func (g *Grid) CountClusters() int {
	clusters := g.GetClusters()
	count := len(clusters)
	return count
}

type Cluster struct {
	Color Color
	X, Y  int
}

func (g *Grid) GetClusters() []Cluster {
	visited := [height][width]bool{}

	var clusters []Cluster
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if !visited[y][x] && g.cells[y][x].color != Empty {
				cluster := Cluster{Color: g.cells[y][x].color, X: x, Y: y}
				clusters = append(clusters, cluster)
				g.markAdjacent(x, y, g.cells[y][x].color, &visited)
			}
		}
	}
	return clusters
}

func (g *Grid) markAdjacent(x, y int, targetColor Color, visited *[height][width]bool) {
	if x < 0 || x >= width || y < 0 || y >= height {
		return
	}

	if visited[y][x] || g.cells[y][x].color != targetColor {
		return
	}

	visited[y][x] = true

	g.markAdjacent(x+1, y, targetColor, visited)
	g.markAdjacent(x-1, y, targetColor, visited)
	g.markAdjacent(x, y+1, targetColor, visited)
	g.markAdjacent(x, y-1, targetColor, visited)
}

func (g *Grid) Clone() *Grid {
	grid := &Grid{
		cells:   [height][width]Cell{},
		actions: make([]Action, len(g.actions)),
	}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			grid.cells[y][x] = g.cells[y][x]
		}
	}
	for idx, action := range g.actions {
		grid.actions[idx] = action
	}
	return grid
}

func (g *Grid) GetActions() []Action {
	return g.actions
}

func (g *Grid) GetNumActions() int {
	return len(g.actions)
}

func (g *Grid) PrintMoves() {
	actions := g.GetActions()
	PrintActions(actions)
}

func PrintActions(actions []Action) {
	for i, action := range actions {
		fmt.Printf("Step %2d x: %d, y: %d %s\n", i+1, action.X, action.Y, action.Color.Icon())
	}
}
