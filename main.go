// Title: Conway's Game of Life
// Description: Reproducing the Game of Life: https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life
//
// RULES:
// 1. Any live cell with < 2 live neighbours dies, as if by underpopulation.
// 2. Any live cell with 2 or 3 live neighbours lives on to the next generation.
// 3. Any live cell with > 3 live neighbours dies, as if by overpopulation.
// 4. Any dead cell with exactly 3 live neighbours becomes a live cell, as if by reproduction.
//
// Library used for making 2D graphics: Ebiten
// Link: https://github.com/hajimehoshi/ebiten
//

package main

import (
	"fmt"
	"math/rand"
	"time"
	//"image/color"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Info struct {
	x, y float64
	life bool
}

const (
	screenWidth  = 1000
	screenHeight = 500
	all_cells    = screenWidth * screenHeight
	first_cells  = 1000 * 300 // You can decide the number of live cells to begin with.
)

var (
	square      *ebiten.Image
	cells_map   [all_cells]Info
	cells               = make([]byte, screenWidth*screenHeight*4)
	yellowColor         = []byte{0xff, 0xff, 0, 0xff}
	x, y        float64 = 0, 0
	flag, f     int     = 0, 1
	frames      int     = 0
)

//
// FIRST RANDOM LIVE CELLS
//
func firstCells() {
	// Used for generating more random numbers.
	//
	rand.Seed(time.Now().UnixNano())

	for n := 0; n < first_cells; n++ {
		// Generate rand cell number.
		//
		c := rand.Intn(screenWidth * screenHeight)
		// Check if that number is repeating.
		//
		if cells_map[c].life == true {
			n--
		} else {
			cells_map[c].life = true
			// Test print.
			//
			//fmt.Println(cells_map[c].x, " - ", cells_map[c].y, " - ", cells_map[c].life)
		}
	}
}

//
// RETURN NR OF NEIGHBORS AROUND A CELL
//
func countNeighbors(key int) int {
	total := 0
	w := screenWidth * (-1)

	// Go through this counter 3 times.
	//
	for c := 0; c < 3; c++ {
		// Go through all 3 neighbors on each of the 3 rows.
		//
		for n := -1; n <= 1; n++ {
			// Check if neighbor is outside our window limits.
			//
			nkey := key + w + n

			// Check if nkey < 0 or > all_cells.
			//
			if nkey >= 0 && nkey < all_cells {
				// Check if neighbor is outside the window  limits.
				//
				tx := cells_map[nkey].x
				ty := cells_map[nkey].y
				if tx >= 0 && tx < screenWidth && ty >= 0 && ty < screenHeight {
					// If neighbor cell is alive, then incremend total.
					//
					if cells_map[nkey].life == true {
						total++
					}
				}
			}
		}
		w += screenWidth
	}

	// We also counted the current live cell.
	// If it's alive, we'll decrement that from the equation.
	//
	if cells_map[key].life == true {
		total--
	}

	return total
}

//
// RECALCULATE CELLS FOR THE NEXT GENERATION
//
func Update() {
	switch flag {
	case 0:
		//
		// DEFAULT CELLS AND FIRST LIVE CELLS
		//
		// Define each cell inside the map.
		// (Well, it's not a map, but an ARRAY really. But I call it a map 'cause I'm mapping the cells. Not very logical.)
		// First set x, y pos, with life.bool = false as a default.
		// Then create the first random live cells.
		//

		k := 0
		// Going through y pos.
		//
		for y < screenHeight {
			// Going through x pos.
			//
			for x < screenWidth {
				// Set cell as default dead.
				//
				cells_map[k] = Info{x, y, false}
				x++
				k++
			}
			x = 0
			y++
		}

		// Creating the first, random live cells.
		//
		firstCells()

		// Set flag to 1, so we don't repeat this process.
		//
		flag = 1

	case 1:
		//
		// CHECK FOR LIVE CELLS NEIGHBOURS
		//
		// We make a temp cells_map, so we can update the whole cells_map array at the end.
		// If we directly change the cells_map array, then the final result will be affected in the process -
		// - the detection of neighbors will be faulty.
		//
		var tempCellsMap [all_cells]Info = cells_map

		for k, v := range cells_map {
			//
			// x - width - 1	||	x - width	||	x - width + 1	/	y - 1
			//
			// x - 1			||		-		||	x + 1			/	y
			//
			// x + width - 1	||	x + width	||	x + width + 1	/	y + 1
			//

			// Get number of neighbors.
			//
			neighbors := countNeighbors(k)

			// Test print.
			//
			if f == 0 {
				//fmt.Println(v.x, " - ", v.y, " - ", v.life, " - ", neighbors)
			}

			// Handle dead and alive cells differently.
			//
			if v.life == true {
				// Test print.
				//
				if f == 0 {
					fmt.Println("OLD: ", v.x, " - ", v.y, " - ", v.life, " - ", neighbors)
				}

				// See if this cell should die.
				//
				if neighbors == 2 || neighbors == 3 {
					// Extant.
					// Do nothing.
					//
				} else {
					// Exctinct.
					//
					// We change the temp cells_map, not the original one.
					//
					tempCellsMap[k].life = false
				}
			} else {
				// See if this empty cell should come alive.
				//
				if neighbors == 3 {
					// Birth.
					//
					// We change the temp cells_map, not the original one.
					//
					tempCellsMap[k].life = true
				} else {
					// Stay unborn.
					// Do nothing.
					//
				}
			}
		}

		// Now we assign all the results to the original cells_map array.
		//
		cells_map = tempCellsMap

		// Used only for Test print.
		//
		f = 1
	}
}

//
// CHANGE PIXELS COLOR IN THE BYTE ARRAY / DRAW
//
func Draw(screen *ebiten.Image) {
	// Go through the cells map.
	//
	for k, v := range cells_map {
		if v.life == true {
			// Find cell position in the pixel matrix.
			//
			idx := 4 * k
			// Set color.
			//
			cells[idx] = yellowColor[0]
			cells[idx+1] = yellowColor[1]
			cells[idx+2] = yellowColor[2]
			cells[idx+3] = yellowColor[3]
		} else {
			idx := 4 * k
			// Set transparent color / remove pixel.
			//
			cells[idx] = 0
			cells[idx+1] = 0
			cells[idx+2] = 0
			cells[idx+3] = 0
		}
	}
}

// Default 60 FPS
//
func update(screen *ebiten.Image) error {
	// Update our cells map.
	//
	Update()

	// Update our pixel matrix.
	//
	Draw(screen)

	// Clear screen / Replace with the new generation cells.
	//
	screen.ReplacePixels(cells)

	// Timer.
	//
	time.Sleep(5 * time.Millisecond)

	//
	// INTRO
	//
	// We have it here, below the Game functions,
	// because it needs to be printed on top of the cells.
	//
	ebitenutil.DebugPrint(screen, "Game of Life")

	// If we want to count frames / generations.
	//
	frames++

	return nil
}

func main() {
	//
	// RUN THE GAME
	//
	if err := ebiten.Run(update, screenWidth, screenHeight, 3, "Hello world!"); err != nil {
		panic(err)
	}
}
