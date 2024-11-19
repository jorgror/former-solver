package solver

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"

	"github.com/jorgror/former-solver/game"
)

type RandomParams struct {
	MaxTries int
}

func ArgsToRandomParams(args []string) RandomParams {
	if len(args) < 1 {
		return RandomParams{MaxTries: 10000}
	}
	maxTries, _ := strconv.Atoi(args[0])
	return RandomParams{MaxTries: maxTries}
}

func Random(grid *game.Grid, params RandomParams) {

	best := 50
	for i := 0; i < params.MaxTries; i++ {

		if i%(params.MaxTries/100) == 0 {
			fmt.Printf("Progress: %d%%\n", i*100/params.MaxTries)
		}

		gridCopy := grid.Clone() // Make a copy of grid
		res := solveRandomRecursive(gridCopy, 0, best)
		if res < best {
			best = res
			fmt.Printf("Best: %d\n", best)
			gridCopy.PrintMoves()
		}
	}

}

func RandomMultithread(grid *game.Grid, params RandomParams) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	best := 50

	solved := 0
	for t := 0; t < runtime.NumCPU(); t++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mu.Lock()
				if solved >= params.MaxTries {
					mu.Unlock()
					return
				}
				solved++
				if solved%(params.MaxTries/100) == 0 {
					fmt.Printf("Progress: %d%%\n", solved*100/params.MaxTries)
				}
				mu.Unlock()

				gridCopy := grid.Clone() // Make a copy of grid
				res := solveRandomRecursive(gridCopy, 0, best)
				if res < best {
					mu.Lock()
					best = res
					fmt.Printf("Best: %d\n", best)
					gridCopy.PrintMoves()
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
}

func solveRandomRecursive(grid *game.Grid, steps int, stop int) int {
	if steps >= stop {
		return stop + 1
	}
	clusters := grid.GetClusters()
	numClusters := len(clusters)
	if numClusters == 0 {
		return steps
	}
	randomIdx := rand.Intn(numClusters)
	cluster := clusters[randomIdx]
	grid.ClearRegion(cluster.X, cluster.Y)
	res := solveRandomRecursive(grid, steps+1, stop)

	return res

}
