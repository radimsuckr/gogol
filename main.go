package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"time"
)

var clearFuncs map[string]func()

func init() {
	clearFuncs = make(map[string]func())
	clearFuncs["darwin"] = func() {
		cmd := exec.Command("clear") // Example for macOS, its tested
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
	clearFuncs["linux"] = func() {
		cmd := exec.Command("clear") // Linux example, its tested
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
	clearFuncs["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") // Windows example, its tested
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}

func callClear() error {
	f, ok := clearFuncs[runtime.GOOS] // Runtime.GOOS -> Linux, Windows, Darwin etc.
	if ok {
		// If we defined a clear function for that platform:
		f() // We execute it
	} else {
		// Unsupported platform
		return fmt.Errorf("terminal clearing is not supported for your platform %s", runtime.GOOS)
	}
	return nil
}

const (
	cols = 64
	rows = 32
)

func render(cells [rows][cols]bool) error {
	err := callClear()
	if err != nil {
		slog.Error(err.Error())
	}

	for i := range len(cells) {
		if i == 0 {
			for range cols + 2 {
				fmt.Print("#")
			}
			fmt.Println()
		}
		for j := range len(cells[i]) {
			if j == 0 {
				fmt.Print("#")
			}
			cell := cells[i][j]
			if cell {
				fmt.Print("O")
			} else {
				fmt.Print(" ")
			}
			if j == len(cells[i])-1 {
				fmt.Print("#")
			}
		}
		fmt.Println()
		if i == len(cells)-1 {
			for range cols + 2 {
				fmt.Print("#")
			}
			fmt.Println()
		}
	}

	return nil
}

func simulate(cells [rows][cols]bool) (newCells [rows][cols]bool) {
	for i := range len(cells) {
		for j := range len(cells[i]) {
			left := max(0, j-1)
			right := min(cols-1, j+1)
			top := max(0, i-1)
			bottom := min(rows-1, i+1)

			aliveNeighbors := 0
			for x := left; x <= right; x++ {
				if cells[top][x] {
					aliveNeighbors++
				}
			}
			for x := left; x <= right; x++ {
				if cells[bottom][x] {
					aliveNeighbors++
				}
			}
			if cells[i][left] {
				aliveNeighbors++
			}
			if cells[i][right] {
				aliveNeighbors++
			}

			if cells[i][j] {
				// Any live cell with fewer than two live neighbors dies, as if by underpopulation.
				if aliveNeighbors < 2 {
					newCells[i][j] = false
					continue
				}

				// Any live cell with two or three live neighbors lives on to the next generation.
				if aliveNeighbors == 2 || aliveNeighbors == 3 {
					newCells[i][j] = true
					continue
				}

				// Any live cell with more than three live neighbors dies, as if by overpopulation.
				if aliveNeighbors > 3 {
					newCells[i][j] = false
					continue
				}
			} else if aliveNeighbors == 3 {
				// Any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction.
				newCells[i][j] = true
				continue
			}
		}
	}
	return newCells
}

func main() {
	cells := [rows][cols]bool{}
	cells[9][9] = true
	cells[9][10] = true
	cells[9][11] = true

	cells[20][20] = true
	cells[20][21] = true
	cells[21][20] = true
	cells[21][21] = true

	cells[2][12] = true
	cells[3][13] = true
	cells[4][11] = true
	cells[4][12] = true
	cells[4][13] = true

	cells[2][17] = true
	cells[3][18] = true
	cells[4][16] = true
	cells[4][17] = true
	cells[4][18] = true

	for {
		cells = simulate(cells)
		err := render(cells)
		if err != nil {
			log.Fatal(err.Error())
		}
		time.Sleep(100 * time.Millisecond)
	}
}
