package solver

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"sync"

	"github.com/jorgror/former-solver/game"
)

var (
	cache      = make(map[string][]*NextGrid)
	cacheLock  sync.Mutex
	cacheDepth = 5
)

func GetNextGrids(grid *game.Grid, clusters []game.Cluster) []*NextGrid {

	numActions := grid.GetNumActions()
	gridHash := ""
	if numActions <= cacheDepth {
		cacheLock.Lock()
		defer cacheLock.Unlock()

		gridHash = grid.HashGrid()
		if cachedResult, found := cache[gridHash]; found {
			return cachedResult
		}
	}

	nextGrids := make([]*NextGrid, 0, len(clusters))
	for i := range clusters {
		gridCopy := grid.Clone()
		gridCopy.ClearRegion(clusters[i].X, clusters[i].Y)
		clusters := gridCopy.GetClusters()
		nextGrids = append(nextGrids, &NextGrid{grid: gridCopy, clusters: clusters})
	}

	if numActions <= cacheDepth {
		cache[gridHash] = nextGrids
	}
	return nextGrids
}

type InfiniteParams struct {
	Cutoff     int
	CacheDepth int
}

func ArgsToInfiniteParams(args []string) InfiniteParams {
	if len(args) < 2 {
		return InfiniteParams{Cutoff: 50, CacheDepth: 5}
	}
	cutoff, _ := strconv.Atoi(args[0])
	cacheDepth, _ := strconv.Atoi(args[1])
	return InfiniteParams{Cutoff: cutoff, CacheDepth: cacheDepth}
}

type Task struct {
	grid     *game.Grid
	clusters []game.Cluster
	actions  []int
}

var (
	taskQueue = make(chan Task, 100)
)

var bestLock sync.Mutex

func Infinite(grid *game.Grid, params InfiniteParams) {

	allTimeBest := params.Cutoff

	cacheDepth = params.CacheDepth

	clusters := grid.GetClusters()
	maxWidth := len(clusters)

	level := 0

	actions := []int{}

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for task := range taskQueue {
				solveInfiniteRecursive(task.grid, task.clusters, task.actions, &allTimeBest)
			}
		}()
	}

	go func() {
		CreateInfiniteTasks(level, maxWidth, &allTimeBest, actions, grid, clusters)
		close(taskQueue)
	}()

	wg.Wait()
}

func CreateInfiniteTasks(level int, maxWidth int, allTimeBest *int, actions []int, grid *game.Grid, clusters []game.Cluster) {

	for i := 0; i < maxWidth-level; i++ {
		if level > *allTimeBest-1 {
			taskQueue <- Task{
				grid:     grid,
				clusters: clusters,
				actions:  actions}
		} else {
			actions = append([]int{i}, actions...)
			CreateInfiniteTasks(level+1, maxWidth, allTimeBest, actions, grid, clusters)
		}
	}
}

type NextGrid struct {
	grid     *game.Grid
	clusters []game.Cluster
}

func solveInfiniteRecursive(grid *game.Grid, clusters []game.Cluster, actions []int, best *int) {
	if grid.GetNumActions() >= *best {
		return
	}

	lenClusters := len(clusters)

	if lenClusters == 0 {
		bestLock.Lock()
		if grid.GetNumActions() < *best {
			*best = grid.GetNumActions()
			fmt.Printf("New best: %d\n", *best)
			grid.PrintMoves()
		}
		bestLock.Unlock()
		return
	}

	if len(actions) == 0 {
		return
	}

	if actions[0] >= lenClusters {
		return
	}

	uniqueColors := make(map[game.Color]struct{})
	for _, cluster := range clusters {
		uniqueColors[cluster.Color] = struct{}{}
	}
	numUniqueColors := len(uniqueColors)

	if numUniqueColors > len(actions) {
		return
	}

	nextGrids := GetNextGrids(grid, clusters)

	sort.Slice(nextGrids, func(i, j int) bool {
		return len(nextGrids[i].clusters) < len(nextGrids[j].clusters)
	})

	action := actions[0]
	solveInfiniteRecursive(nextGrids[action].grid, nextGrids[action].clusters, actions[1:], best)
}
