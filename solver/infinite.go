package solver

import (
	"fmt"
	"runtime"
	"sort"
	"strconv"
	"sync"
	"time"

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

		gridHash = grid.HashGrid()
		cachedResult, found := cache[gridHash]
		cacheLock.Unlock()
		if found {
			return cachedResult
		}
	}

	nextGrids := make([]*NextGrid, 0, len(clusters))
	for i := range clusters {
		gridCopy := grid.Clone()
		gridCopy.ClearRegion(clusters[i].X, clusters[i].Y)
		clusters := gridCopy.GetClusters()
		score := ScoringFunc(gridCopy, clusters)
		nextGrids = append(nextGrids, &NextGrid{grid: gridCopy, clusters: clusters, score: score})
	}
	sort.Slice(nextGrids, func(i, j int) bool {
		return nextGrids[i].score < nextGrids[j].score
	})

	if numActions <= cacheDepth {
		cacheLock.Lock()
		cache[gridHash] = nextGrids
		cacheLock.Unlock()
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

	numWorkers := runtime.NumCPU()

	lastTime = time.Now()
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for task := range taskQueue {
				solveInfiniteRecursive(task.grid, task.clusters, task.actions, &allTimeBest)
				LogTasks(task.actions)
			}
		}()
	}

	go func() {
		InfiniteTaskController(grid, &allTimeBest)
		close(taskQueue)
	}()

	wg.Wait()
}

var (
	store   [20][40]int
	counter int = 0
)

var lastTime time.Time

func LogTasks(actions []int) {
	// for i, action := range actions {
	// 	store[i][action]++
	// }
	counter++

	if counter%100000 == 0 {
		interval := time.Since(lastTime)
		rate := float64(100000) / interval.Seconds()
		fmt.Printf("Counter: %d Rate: %.0f / sec\n", counter, rate)
		lastTime = time.Now()
		// for i := 0; i < 20; i++ {
		// 	for j := 0; j < 30; j++ {
		// 		if store[i][j] > 0 {
		// 			fmt.Printf("%d ", store[i][j])
		// 		}
		// 	}
		// 	fmt.Println()
		// }
	}
}

type RunPlan struct {
	FullLevels  int
	SearchWidth int
}

func InfiniteTaskController(grid *game.Grid, allTimeBest *int) {

	runPlan := []RunPlan{
		{FullLevels: 1, SearchWidth: 1},
		{FullLevels: 2, SearchWidth: 1},
		{FullLevels: 3, SearchWidth: 1},
		{FullLevels: 3, SearchWidth: 3},
		{FullLevels: 4, SearchWidth: 1},
		{FullLevels: 4, SearchWidth: 3},
		{FullLevels: 5, SearchWidth: 1},
		{FullLevels: 4, SearchWidth: 3},
	}

	clusters := grid.GetClusters()
	maxWidth := len(clusters)

	for _, plan := range runPlan {
		actions := make([]int, 0)
		fmt.Println("Running plan:", plan)
		CreateInfiniteTasks(0, maxWidth, allTimeBest, actions, grid, clusters, plan.FullLevels, plan.SearchWidth)
	}
}

func CreateInfiniteTasks(level int, maxWidth int, allTimeBest *int, actions []int, grid *game.Grid, clusters []game.Cluster, fullLevels int, searchWidth int) {

	width := maxWidth
	if *allTimeBest-level-1 > fullLevels {
		width = searchWidth
	}

	for i := 0; i < width; i++ {
		actionsCopy := make([]int, len(actions))
		copy(actionsCopy, actions)
		actionsCopy = append([]int{i}, actionsCopy...)

		if level > *allTimeBest-2 {
			taskQueue <- Task{
				grid:     grid,
				clusters: clusters,
				actions:  actionsCopy}
		} else {
			CreateInfiniteTasks(level+1, maxWidth, allTimeBest, actionsCopy, grid, clusters, fullLevels, searchWidth)
		}
	}
}

type NextGrid struct {
	grid     *game.Grid
	clusters []game.Cluster
	score    int
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

	action := actions[0]
	solveInfiniteRecursive(nextGrids[action].grid, nextGrids[action].clusters, actions[1:], best)
}

func ScoringFunc(grid *game.Grid, clusters []game.Cluster) int {
	score := len(clusters)
	colors := make([]int, 4)
	for _, cluster := range clusters {
		colors[cluster.Color]++
	}

	sort.Slice(colors, func(i, j int) bool {
		return colors[i] < colors[j]
	})

	score += colors[0] * 2
	score += colors[1] * 1
	score -= colors[2] * 1
	score -= colors[3] * 2

	return score
}
