package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	cols = 64
	rows = 32
)

func render(cells [rows][cols]bool) error {
	var buf strings.Builder
	// Pre-allocate buffer capacity to avoid reallocations
	// Rough estimate: (cols + 2) * (rows + 2) characters + ANSI codes
	buf.Grow((cols + 3) * (rows + 3))

	// Move cursor to home position (1,1) instead of clearing
	buf.WriteString("\033[H")

	// Top border
	for range cols + 2 {
		buf.WriteByte('#')
	}
	buf.WriteByte('\n')

	// Grid rows
	for i := range len(cells) {
		buf.WriteByte('#')
		for j := range len(cells[i]) {
			if cells[i][j] {
				buf.WriteByte('O')
			} else {
				buf.WriteByte(' ')
			}
		}
		buf.WriteByte('#')
		buf.WriteByte('\n')
	}

	// Bottom border
	for range cols + 2 {
		buf.WriteByte('#')
	}
	buf.WriteByte('\n')

	// Write entire buffer at once to minimize flickering
	_, err := os.Stdout.Write([]byte(buf.String()))
	return err
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
	// Enable ANSI escape codes on Windows
	if err := enableANSI(); err != nil {
		log.Printf("Warning: Failed to enable ANSI support: %v\n", err)
	}

	// Set up signal handler to restore cursor on CTRL+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		// Restore cursor and exit cleanly
		fmt.Print("\033[?25h")
		os.Exit(0)
	}()

	// Clear screen once and hide cursor for cleaner rendering
	fmt.Print("\033[2J\033[H\033[?25l")

	// Ensure cursor is shown again on normal exit
	defer fmt.Print("\033[?25h")

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
		time.Sleep(1000 / 60 * time.Millisecond)
	}
}
