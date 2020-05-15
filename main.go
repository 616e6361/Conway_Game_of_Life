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
	"math/rand"
	"time"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

const (
	screenWidth  = 2000
	screenHeight = 1000
	all_cells    = screenWidth * screenHeight
)

var (
	cells            = make([]byte, screenWidth*screenHeight*4)
	flag, frames int = 0, 0
)

//
// FIRST RANDOM LIVE CELLS
//
func firstCells() {
	// Used for generating more random numbers.
	//
	rand.Seed(time.Now().UnixNano())

	for k := 0; k < all_cells; k++ {
		b := [2]byte{0, 0xff}
		c := rand.Intn(1-0+1) + 0

		// Assign each cell.
		//
		cells[k*4] = b[c]
		cells[k*4+1] = b[c]
		cells[k*4+2] = b[c]
		cells[k*4+3] = b[c]
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
				ty := nkey / screenWidth
				tx := nkey - (ty * screenWidth)
				if tx >= 0 && tx < screenWidth && ty >= 0 && ty < screenHeight {
					// If neighbor cell is alive, then incremend total.
					//
					if cells[nkey*4] == 0xff {
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
	if cells[key*4] == 0xff {
		total--
	}

	return total
}

//
// UPDATE AND DRAW CELLS
//
func Update(screen *ebiten.Image) {
	switch flag {
	case 0:
		firstCells()

		// Set flag to 1, so we don't repeat this process.
		//
		flag = 1

	case 1:
		var tempCells = make([]byte, screenWidth*screenHeight*4)
		
		for k := 0; k < all_cells; k++ {
			// Get number of neighbors.
			//
			neighbors := countNeighbors(k)

			// Handle dead and alive cells differently.
			//
			if cells[k*4] == 0xff {
				// See if this cell should die.
				//
				if neighbors == 2 || neighbors == 3 {
					// Extant.
					//
					tempCells[k*4] = cells[k*4]
					tempCells[k*4+1] = cells[k*4+1]
					tempCells[k*4+2] = cells[k*4+2]
					tempCells[k*4+3] = cells[k*4+3]
				} else {
					// Exctinct. Do nothing.
					//
				}
			} else {
				// See if this empty cell should come alive.
				//
				if neighbors == 3 {
					// Birth.
					//
					tempCells[k*4] = 0xff
					tempCells[k*4+1] = 0xff
					tempCells[k*4+2] = 0xff
					tempCells[k*4+3] = 0xff
				} else {
					// Stay unborn. Do nothing.
					//
				}
			}
		}

		// Now we assign all the results to the original cells array.
		//
		cells = tempCells
	}
}

// Default 60 FPS
//
func update(screen *ebiten.Image) error {
	// Update our pixel matrix.
	//
	Update(screen)

	// Clear screen / New generation.
	//
	screen.ReplacePixels(cells)

	ebitenutil.DebugPrint(screen, "Game of Life")
	return nil
}

func main() {
	//
	// RUN THE GAME
	//
	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Hello world!"); err != nil {
		panic(err)
	}
}
