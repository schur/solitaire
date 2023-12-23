/**
 * Peg Solitaire Solver
 * Copyright (C) 2014 blackflux.com <pegsolitaire@blackflux.com>
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License version 3 as
 *  published by the Free Software Foundation.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
 
/**
 * Solver for the English peg solitaire.
 * This program finds a random solution for peg solitaire game by using brute force.
 *
 * -- Runtime
 * A solution is typically found in less than two seconds, but the time does highly
 * fluctuate (I've seen everything from a few milliseconds to several seconds).
 *
 * -- Implementation
 *
 * The implementation is highly optimized and uses bit operators to efficiently find
 * a solution. The idea is as following: Since there exists 33 slots on the board, it
 * can not be represented by using an integer (32 bits), but we can use a long (64 bits).
 * The first 49 bits (7 x 7) of the long represent the board. However there are some bits
 * that are not valid and never used, i.e. 0,1,5,6,7,8,12,13 and so on. Checking of
 * possible moves and applying them can be done by using simple bit operations.
 *
 * A recursive function is then used to check all possible moves for a given board,
 * applying each valid move and calling itself with the resulting board. The recursion is
 * done "in reverse", starting from the goal board. While this is not conceptually faster [a],
 * it allows for a minimum amount of bit operations in the recursion:
 *
 * To reverse a move we can simply check
 * - (board & twoBalls) == 0 and
 * - (board & oneBall) != 0
 * where "twoBalls" indicates the two ball that would need to be added for this reversed move.
 * If we instead used the intuitive search direction, the same check would require additional
 * binary operations, since a simple inversion of the check would not work [b].
 *
 * Paper [1] shows how the moves can be ordered to almost instantly find a solution.
 * Website [2] gives a nice overview of binary operations and some tricks that
 * can be applied.
 *
 * [a] Playing the game in reverse is simply the inversion of the original game - just remove all
 * balls from the board and place ball where there were none before and you'll understand
 * what I mean.
 * [b] There is no "single" binary operation to check if two specific bits are set, but there
 * is one to check if they are both zero. There is further a binary operation to check if a specific
 * bit is set.
 *
 * [1] http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.6.4826 (download at the top)
 * [2] http://graphics.stanford.edu/~seander/bithacks.html
 */

package main

import (
	"fmt"
	"math/rand"
)

// board that contains a ball in every available slot
const VALID_BOARD_CELLS uint64 = 124141734710812

// initial board (one marble free in center)
const INITIAL_BOARD uint64 = 124141717933596

// goal board (one marble in center)
const GOAL_BOARD uint64 = 16777216

const colorReset = "\033[0m"
const colorRed = "\033[31m"

// the structure represtenting a move is composed as follows:
// - first entry (after) holds the peg that is added by the move
// - second entry (before) holds the two pegs that are removed by the move
// - third entry holds all three involved pegs
type Move struct {
	after, before, all uint64
}

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

	// print the found solution
	for i := 0; i < len(Solution); i++ {
		// loop over all 7 rows
		var k int
		for m := 0; m < 7; m++ {
			// print 10 steps in 1 row
			for k = 0; k < 16; k++ {
				//fmt.Printf("i: %d, m: %d, k: %d", i, m, k)
				printLine(Solution[i+k], m)
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

// print the board
func printBoard(board uint64) {
	// loop over all 7 rows
	for i := 0; i < 7; i++ {
		printLine(board, i)
		fmt.Println()
	}
	fmt.Println("-------------")
}

// print one line of the board
func printLine(board uint64, line int) {
	// loop over all cells (the board is 7 x 7)
	var cell uint64 = 1 << (7 * line) // move to first cell in the line
	for i := 0; i < 7; i++ {
		validCell := (cell & VALID_BOARD_CELLS) != 0
		if validCell {
			if (cell & board) != 0 {
				fmt.Printf(colorRed + "X" + colorReset)
			} else {
				fmt.Printf("0")
			}
		} else {
			fmt.Printf(" ")
		}
		cell = cell << 1 // move to next cell
		// print new line after 7 slots
	}
}

// print a move
func printMove(move Move) {
	// loop over all 7 rows
	for i := 0; i < 7; i++ {
		printLine(move.after, i)
		fmt.Printf("   ")
		printLine(move.before, i)
		fmt.Printf("   ")
		printLine(move.all, i)
		fmt.Println()

	}
	fmt.Println("-------------")
}
