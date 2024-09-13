package main

import (
	"math"
	"math/rand/v2"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Any live cell with fewer than two live neighbours dies, as if by underpopulation.
// Any live cell with two or three live neighbours lives on to the next generation.
// Any live cell with more than three live neighbours dies, as if by overpopulation.
// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.

const (
	debug                   = false // show debug information
	screenDimension int32   = 1000  // the size of the window
	gridDimension   int     = 100   // the size of the grid
	gridLifeTime    int     = 10000 // the number of generations to run
	probability     float32 = .1    // how probable it should be for a cell to be alive at seeding, higher is more likely
	tickRate        float32 = 100   // tick rate in milliseconds
	useRandom       bool    = true  // use random seeding
	startPaused     bool    = true  // start paused
)

type Cell struct {
	x          int
	y          int
	state      bool // true = alive, false = dead
	neighbours int
}

type Grid struct {
	cells           [][]Cell
	dimension       int
	lifetime        int
	paused          bool
	tickRate        float32
	initGrid        func(*Grid) *Grid
	countNeighbours func(*Grid) *Grid
	updateGrid      func(*Grid) *Grid
	drawGrid        func(*Grid)
}

func initGrid(grid *Grid) *Grid {
	grid.dimension = gridDimension
	grid.lifetime = gridLifeTime
	grid.tickRate = tickRate
	grid.paused = startPaused
	grid.cells = make([][]Cell, grid.dimension)
	// Initialize the Grid with random values
	for i := 0; i < grid.dimension; i++ {
		grid.cells[i] = make([]Cell, grid.dimension)
		for j := 0; j < grid.dimension; j++ {
			grid.cells[i][j].x = i
			grid.cells[i][j].y = j
			if useRandom {
				grid.cells[i][j].state = rand.Float32() < probability
			} else {
				grid.cells[i][j].state = false
			}
		}
	}
	return grid
}

func countNeighbours(grid *Grid) *Grid {
	// Count the number of live neigbours
	for i := 0; i < grid.dimension; i++ {
		for j := 0; j < grid.dimension; j++ {
			liveNeighbours := 0
			for x := i - 1; x <= i+1; x++ {
				for y := j - 1; y <= j+1; y++ {
					if x >= 0 && x < grid.dimension && y >= 0 && y < grid.dimension {
						if grid.cells[x][y].state {
							liveNeighbours++
						}
					}
				}
			}
			// Subtract the cell itself from the count
			if grid.cells[i][j].state {
				liveNeighbours--
			}
			grid.cells[i][j].neighbours = liveNeighbours
		}
	}
	return grid
}

func updateGrid(grid *Grid) *Grid {
	// Update the Grid based on the rules of the Game of Life
	for i := 0; i < grid.dimension; i++ {
		for j := 0; j < grid.dimension; j++ {
			currCell := grid.cells[i][j]
			if currCell.state {
				if currCell.neighbours < 2 || currCell.neighbours > 3 {
					grid.cells[i][j].state = false
				} else {
					if currCell.neighbours == 3 || currCell.neighbours == 2 {
						grid.cells[i][j].state = true
					}
				}
			} else {
				if currCell.neighbours == 3 {
					grid.cells[i][j].state = true
				}
			}
		}
	}
	return grid
}

func drawGrid(grid *Grid) {
	cyclesRemaining := grid.lifetime
	cellSize := screenDimension / int32(gridDimension)
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	for i := 0; i < grid.dimension; i++ {
		for j := 0; j < grid.dimension; j++ {
			currentCell := grid.cells[i][j]
			if currentCell.state {
				rl.DrawRectangle(int32(currentCell.x)*cellSize, int32(currentCell.y)*cellSize, cellSize, cellSize, rl.Black)
				if debug {
					rl.DrawText(strconv.Itoa(currentCell.neighbours), int32(currentCell.x)*cellSize+2, int32(currentCell.y)*cellSize+2, 10, rl.White)
				}
			} else {
				rl.DrawRectangle(int32(currentCell.x)*cellSize, int32(currentCell.y)*cellSize, cellSize, cellSize, rl.White)
			}
		}
	}

	rl.DrawText("Cycles Remaining: "+strconv.Itoa(cyclesRemaining), 10, int32(screenDimension)+10, 32, rl.Black)
	rl.EndDrawing()

}

func main() {
	rl.InitWindow(screenDimension, screenDimension+50, "Game of Life")
	defer rl.CloseWindow()
	rl.SetTargetFPS(120)

	grid := initGrid(new(Grid))

	// work out neighbours for initial seed
	grid = countNeighbours(grid)
	cellSize := float64(screenDimension / int32(gridDimension))

	// main loop
	for !rl.WindowShouldClose() {
		if grid.lifetime <= 0 {
			grid = initGrid(grid)
			grid = countNeighbours(grid)
		}
		if rl.IsKeyPressed(rl.KeyR) {
			grid = initGrid(grid)
			grid = countNeighbours(grid)
		}

		if rl.IsKeyPressed(rl.KeyQ) {
			return
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			grid.paused = !grid.paused
		}

		if rl.IsKeyPressed(rl.KeyMinus) {
			grid.tickRate += 100
		}
		if rl.IsKeyPressed(rl.KeyEqual) {
			grid.tickRate -= 100
		}

		mousePos := rl.GetMousePosition()
		// calculate cell position
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			x := int(math.Floor(float64(mousePos.X) / cellSize))
			y := int(math.Floor(float64(mousePos.Y) / cellSize))
			grid.cells[x][y].state = !grid.cells[x][y].state
			grid = countNeighbours(grid)
		}

		if grid.paused && rl.IsKeyPressed(rl.KeyRight) {
			updateGrid(grid)
			countNeighbours(grid)
			grid.lifetime -= 1
		}

		drawGrid(grid)

		time.Sleep(time.Millisecond * time.Duration(grid.tickRate))

		// update grid
		if !grid.paused {
			updateGrid(grid)
			countNeighbours(grid)
			grid.lifetime -= 1
		}

	}
}
