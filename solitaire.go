package main

import "fmt"
import "math/rand"
import "time"

const VALID_BOARD_CELLS uint64 = 124141734710812
const INITIAL_BOARD uint64 = 124141717933596
const GOAL_BOARD uint64 = 16777216

// holds all 76 moves that are possible
// the inner array is structures as following:
// - first entry holds the peg that is added by the move
// - second entry holds the two pegs that are removed by the move
// - third entry holds all three involved pegs
var Moves [76][3]uint64

func main () {

	// generate all possible moves
//	moves := make ([]uint64,0,228)
	// holds all starting positions in west-east direction
	var pos = 0;
	var startsX = [19]int{2,9,14,15,16,17,18,21,22,23,24,25,28,29,30,31,32,37,44}
	for _,x := range startsX {
		createMoves(x, x + 1, x + 2, &pos)
	}
	// holds all starting positions in north-south direction
	var startsY = [19]int{2,3,4,9,10,11,14,15,16,17,18,19,20,23,24,25,30,31,32}
	for _,y := range startsY {
		createMoves(y, y + 7, y + 14, &pos)
	}
	fmt.Println(pos);

	fmt.Println(Moves)

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(Moves), func(i, j int) { Moves[i], Moves[j] = Moves[j], Moves[i] })

	fmt.Println(Moves)

//	printBoard(INITIAL_BOARD)
//	printBoard(GOAL_BOARD)
	
}


// create the two possible moves for the three added pegs
// (this function assumes that the pegs are in one continuous line)
func createMoves(bit1 int, bit2 int, bit3 int, pos *int) {
	Moves[*pos][0] = 1 << bit1
	Moves[*pos][1] = (1 << bit2) | (1 << bit3)
	Moves[*pos][2] = (1 << bit1) | (1 << bit2) | (1 << bit3)
	*pos++;
	Moves[*pos][0] = 1 << bit3
	Moves[*pos][1] = (1 << bit2) | (1 << bit1)
	Moves[*pos][2] = (1 << bit1) | (1 << bit2) | (1 << bit3)
	*pos++;
	}

// print the board
func printBoard (board uint64) {
	// loop over all cells (the board is 7 x 7)
	for i:=0; i<49; i++ {
		var character byte
		validCell := ((1<<i) & VALID_BOARD_CELLS) != 0 
		if validCell {
			if ((1<<i) & board) != 0 {
				character = 'X'
			} else {
				character = '0'
			}
		} else {
			character = ' '
		}
		// print new line after 7 slots
		fmt.Printf("%c",character)
		if i % 7 == 6 {
			fmt.Println()
		}
	}
	fmt.Println("-------------")
}