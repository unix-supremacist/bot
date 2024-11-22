package main

import (
	"strings"
	"math/rand"
)

type rotation struct {
	X, Y int
}

func generateBoard(sizex, sizey int, input string) string {
	words := strings.Split(input, ",")
	return generateBoardWords(sizex, sizey, words)
}

func generateBoardWords(sizex, sizey int, words []string) string {
	board := ""
	grid := make([][]string, sizey)
	for i := 0; i < sizey; i++ {
		grid[i] = make([]string, sizex)
	}
	rotations := []rotation{}
	rotations = append(rotations, rotation{1,0})
	rotations = append(rotations, rotation{1,1})
	rotations = append(rotations, rotation{0,1})
	rotations = append(rotations, rotation{-1,1})
	rotations = append(rotations, rotation{-1,0})
	rotations = append(rotations, rotation{-1,-1})
	rotations = append(rotations, rotation{0,-1})
	rotations = append(rotations, rotation{1,-1})

	for height := range grid {
		for width := range grid[height] {
			grid[height][width] = "!"
		}
	}

	for wordi := range words {
		word := strings.Replace(strings.Replace(words[wordi], "'", "", -1), " ", "", -1)
		length := len(word)
		succeeded := false
		loops := 0
		for !succeeded {
			loops++
			if loops > sizex*sizey*len(words) {
				board = board+"word cannot fit with board size\n"
				break
			}
			invalid := false
			rot := rotations[rand.Intn(len(rotations))]
			
			x_pos := rand.Intn(sizex)
			y_pos := rand.Intn(sizey)
			x_end := x_pos+length*rot.X
			y_end := y_pos+length*rot.Y

			if x_end >= sizex {continue}
			if y_end >= sizey {continue}
			if x_end < 0 {continue}
			if y_end < 0 {continue}

			for ci, c := range word {
				i := ci+1
				char := string(c)
				x_pos_new := x_pos+i*rot.X
				y_pos_new := y_pos+i*rot.Y

				if !((grid[y_pos_new][x_pos_new] == "!") || (grid[y_pos_new][x_pos_new] == char)) {
					invalid=true
					continue
				}
			}

			if invalid {continue}

			for ci, c := range word {
				i := ci+1
				char := string(c)
				x_pos_new := x_pos+i*rot.X
				y_pos_new := y_pos+i*rot.Y
				grid[y_pos_new][x_pos_new] = char
			}

			succeeded = true
			loops = 0
		}
	}

	for height := range grid {
		for width := range grid[height] {
			if grid[height][width] == "!" {
				//grid[height][width] = string(rand.Intn(25)+97)
			}
			board = board+grid[height][width]
			if width+1 != len(grid[height]) {
				board = board+" "
			}
		}
		board = board+"\n"
	}

	return board
}