package main

import (
	"fmt"
	"os"

	"github.com/jorgror/former-solver/game"
	"github.com/jorgror/former-solver/solver"
)

func main() {

	if len(os.Args) < 4 {
		fmt.Println("Usage: solve <input file> solver solver_param")
		return
	}

	file := os.Args[1]
	solverName := os.Args[2]

	grid := game.NewGrid()
	grid.LoadFromFile(file)

	switch solverName {
	case "random":
		randomParams := solver.ArgsToRandomParams(os.Args[3:])
		solver.Random(grid, randomParams)
	case "random_mt":
		randomParams := solver.ArgsToRandomParams(os.Args[3:])
		solver.RandomMultithread(grid, randomParams)
	case "beam":
		beamParams := solver.ArgsToBeamParams(os.Args[3:])
		solver.Beam(grid, beamParams)
	case "infinite":
		infiniteParams := solver.ArgsToInfiniteParams(os.Args[3:])
		solver.Infinite(grid, infiniteParams)
	default:
		fmt.Println("Unknown solver")
	}
}
