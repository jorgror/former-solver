package solver

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"

	"github.com/jorgror/former-solver/game"
)

type BeamResult struct {
	grid        *game.Grid
	totalScore  int
	bestScore   int
	bestActions []game.Action
}

type BeamParams struct {
	Itterations int
	KeepNumber  int
	Depth       int
	Cutoff      int
}

func ArgsToBeamParams(args []string) BeamParams {
	if len(args) < 4 {
		return BeamParams{Itterations: 20, KeepNumber: 10, Depth: 3, Cutoff: 50}
	}
	itterations, _ := strconv.Atoi(args[0])
	keepNumber, _ := strconv.Atoi(args[1])
	depth, _ := strconv.Atoi(args[2])
	cutoff, _ := strconv.Atoi(args[3])
	return BeamParams{Itterations: itterations, KeepNumber: keepNumber, Depth: depth, Cutoff: cutoff}
}

func Beam(grid *game.Grid, params BeamParams) {

	allTimeBest := params.Cutoff
	beamResults := []BeamResult{{grid: grid, bestScore: params.Cutoff}}

	for i := 0; ; i++ {
		for j := 0; j < params.Depth; j++ {
			tempStarts := make([]BeamResult, 0)
			for _, br := range beamResults {
				clusters := br.grid.GetClusters()
				for _, cluster := range clusters {
					gridCopy := br.grid.Clone()
					gridCopy.ClearRegion(cluster.X, cluster.Y)
					tempStarts = append(tempStarts, BeamResult{grid: gridCopy, bestScore: params.Cutoff})
				}
			}
			beamResults = tempStarts
		}

		numBeams := len(beamResults)

		// Create channels for work distribution and result collection
		workChan := make(chan int, numBeams)
		resultChan := make(chan BeamResult, numBeams)

		// Start worker goroutines
		numWorkers := runtime.NumCPU()
		for w := 0; w < numWorkers; w++ {
			go func() {
				for idx := range workChan {
					br := beamResults[idx]
					// Process beam result
					for j := 0; j < params.Itterations; j++ {
						ittCopy := br.grid.Clone()
						res := solveRandomRecursive(ittCopy, (i+1)*params.Depth, params.Cutoff)
						br.totalScore += res
						if res < br.bestScore {
							br.bestScore = res
							br.bestActions = ittCopy.GetActions()
						}
					}
					// Send updated beam result back
					resultChan <- br
				}
			}()
		}

		// Distribute work
		for idx := 0; idx < numBeams; idx++ {
			workChan <- idx
		}
		close(workChan)

		// Collect results
		updatedResults := make([]BeamResult, numBeams)
		for idx := 0; idx < numBeams; idx++ {
			updatedResults[idx] = <-resultChan
		}
		close(resultChan)

		beamResults = updatedResults

		// Use sort.Slice to sort beamResults
		sort.Slice(beamResults, func(i, j int) bool {
			return beamResults[i].totalScore < beamResults[j].totalScore
		})

		toKeep := numBeams
		if numBeams > params.KeepNumber {
			toKeep = params.KeepNumber
		}

		for _, res := range beamResults {
			if res.bestScore < allTimeBest {
				allTimeBest = res.bestScore
				fmt.Printf("New best score: %d\n", allTimeBest)
				game.PrintActions(res.bestActions)
			}
		}

		if toKeep > 1 {
			fmt.Printf("Best avg score kept: %.2f, worst avg score kept: %.2f, worst avg score: %.2f\n",
				float64(beamResults[0].totalScore)/float64(params.Itterations),
				float64(beamResults[toKeep-1].totalScore)/float64(params.Itterations),
				float64(beamResults[numBeams-1].totalScore)/float64(params.Itterations))
			beamResults = beamResults[:toKeep]
		} else {
			return
		}
	}

}