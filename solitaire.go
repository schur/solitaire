package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

// the below constants are binary representations of the bitmaps that model the board
// a "1" represents a marble in the slot, a "0" rpresents an empty slot
// However, in VALID_BOARD_CELLS a "1" represents a valid slot

// the route via strconv is done to break the binary numbers into multiple lines to visualise the board

// Valid Cells that can contain a ball (i.e. thev available slots)
var VALID_BOARD_CELLS, _ = strconv.ParseUint("0"+
	"0011100"+
	"0011100"+
	"1111111"+
	"1111111"+
	"1111111"+
	"0011100"+
	"0011100", 2, 64)

// initial board (one marble free in center)
var INITIAL_BOARD, _ = strconv.ParseUint("0"+
	"0011100"+
	"0011100"+
	"1111111"+
	"1110111"+
	"1111111"+
	"0011100"+
	"0011100", 2, 64)

// goal board (one marble in center)
var GOAL_BOARD, _ = strconv.ParseUint("0"+
	"0000000"+
	"0000000"+
	"0001000"+
	"0000000"+
	"0000000"+
	"0000000"+
	"0000000", 2, 64)

// the structure represtenting a move is composed as follows:
// - first entry (after) holds the peg that is added by the move
// - second entry (before) holds the two pegs that are removed by the move
// - third entry holds all three involved pegs
type Move struct {
	after, before, all uint64
}

// Global Variables:

// list of seen boards - this is used to prevent rechecking of paths
var seenBoards = map[uint64]bool{}

// list of solution boards in ascending order - filled in once the solution is found - array capcity 32 is based on known max. number of moves
var Solution = make([]uint64, 0, 32)

// holds all 76 moves that are possible
var Moves = make([]Move, 0, 76)

func main() {

	// generate all possible moves

	// holds all starting positions in west-east direction
	var startsX = [19]int{2, 9, 14, 15, 16, 17, 18, 21, 22, 23, 24, 25, 28, 29, 30, 31, 32, 37, 44}
	for _, x := range startsX {
		Moves = createMoves(x, x+1, x+2, Moves)
	}
	// holds all starting positions in north-south direction
	var startsY = [19]int{2, 3, 4, 9, 10, 11, 14, 15, 16, 17, 18, 19, 20, 23, 24, 25, 30, 31, 32}
	for _, y := range startsY {
		Moves = createMoves(y, y+7, y+14, Moves)
	}

	// randomize the order of the moves (this highly influences the resulting runtime)
	rand.Shuffle(len(Moves), func(i, j int) { Moves[i], Moves[j] = Moves[j], Moves[i] })

	// add starting board (as this board is not added by the recursive function)
	Solution = append(Solution, INITIAL_BOARD)

	// start recursively search for the initial board from the goal (reverse direction!)
	search(GOAL_BOARD)

	// print the solution
	PrintSolution()

}

// do the calculation recursively by starting from
// the "GOAL_BOARD" and doing moves in reverse
func search(board uint64) bool {
	// for all possible moves
	for _, move := range Moves {
		// check if the move is valid
		// Note: we place "two ball" check first since it is more
		// likely to fail. This saves about 20% in run time (!)
		if (move.before&board) == 0 && (move.after&board) != 0 {
			// calculate the board after this move was applied
			newBoard := board ^ move.all
			// only continue processing if we have not seen this board before
			if !seenBoards[newBoard] {
				seenBoards[newBoard] = true
				// check if the initial board is reached
				if newBoard == INITIAL_BOARD || search(newBoard) {
					Solution = append(Solution, board)
					return true
				}
			}
		}
	}
	return false
}

// create the two possible moves for the three added pegs
// (this function assumes that the pegs are in one continuous line)
func createMoves(bit1 int, bit2 int, bit3 int, moves []Move) []Move {
	var newmove Move
	newmove.after = 1 << bit1
	newmove.before = (1 << bit2) | (1 << bit3)
	newmove.all = (1 << bit1) | (1 << bit2) | (1 << bit3)
	moves = append(moves, newmove)

	newmove.after = 1 << bit3
	newmove.before = (1 << bit2) | (1 << bit1)
	newmove.all = (1 << bit1) | (1 << bit2) | (1 << bit3)
	moves = append(moves, newmove)

	return moves
}

// print the found solution
func PrintSolution() {

	for i := 0; i < len(Solution); i++ {
		// loop over all 7 rows
		var k int
		for m := 0; m < 7; m++ {
			// print 16 steps in 1 row
			for k = 0; k < 16; k++ {
				//fmt.Printf("i: %d, m: %d, k: %d", i, m, k)
				previous := i + k - 1
				if previous < 0 {
					previous = 0
				}
				printLine(Solution[i+k], Solution[previous], m)
				if (i + k) == len(Solution)-1 {
					k++
					break
				}
				fmt.Printf("   ")
			}
			fmt.Println()
		}
		i = i + k - 1
		fmt.Println("-------------")
	}
}

// print one line of the board
// first argument: board to print
// second argument: previous board - the function will highlight any changes made by a move
// pass the board from the first argument again to not highlight any changes
// third argument: line number to print
func printLine(board uint64, prev_board uint64, line int) {
	const colorReset = "\033[0m"
	const colorRed = "\033[31m"
	const colorBlue = "\033[34m"
	const colorGrey = "\033[37m"
	const colorWhite = "\033[97m"

	// loop over all cells (the board is 7 x 7)
	var cell uint64 = 1 << (7 * line) // move to first cell in the line
	for i := 0; i < 7; i++ {
		validCell := (cell & VALID_BOARD_CELLS) != 0
		if validCell {
			if (cell & board) != 0 {
				if (cell & prev_board) == 0 {
					fmt.Printf(colorRed)
				} else {
					fmt.Printf(colorWhite)
				}
				fmt.Printf("X" + colorReset)
			} else {
				if (cell & prev_board) != 0 {
					fmt.Printf(colorBlue)
				} else {
					fmt.Printf(colorGrey)
				}
				fmt.Printf("0" + colorReset)
			}
		} else {
			fmt.Printf(" ")
		}
		cell = cell << 1 // move to next cell
		// print new line after 7 slots
	}
}
