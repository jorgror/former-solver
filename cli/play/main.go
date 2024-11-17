package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jorgror/former-solver/game"
)

func main() {
	grid := game.NewGrid()
	if len(os.Args) > 1 {
		err := grid.LoadFromFile(os.Args[1])
		if err != nil {
			fmt.Println("Error loading file:", err)
			return
		}
	}
	reader := bufio.NewReader(os.Stdin)

	clusters := grid.CountClusters()

	for clusters > 0 {
		grid.Print()
		fmt.Printf("Clusters: %d\n", clusters)
		fmt.Printf("Current numer of moves: %d\n", grid.GetNumActions())
		fmt.Print("Enter coordinates to clear (x y): ")
		input, _ := reader.ReadString('\n')
		var x, y int
		fmt.Sscanf(input, "%d %d", &x, &y)

		grid.ClearRegion(x, y)
		clusters = grid.CountClusters()
	}

	fmt.Println("Game Over!")
	fmt.Printf("Final Score: %d\n", grid.GetNumActions())
}
