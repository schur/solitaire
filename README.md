Peg Solitaire Solver
Copyright (C) 2014 blackflux.com <pegsolitaire@blackflux.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License version 3 as
published by the Free Software Foundation.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

Solver for the English peg solitaire.
This program finds a random solution for peg solitaire game by using brute force.

-- Runtime
A solution is typically found in less than two seconds, but the time does highly
fluctuate (I've seen everything from a few milliseconds to several seconds).

-- Implementation
The implementation is highly optimized and uses bit operators to efficiently find
a solution. The idea is as following: Since there exists 33 slots on the board, it
can not be represented by using an integer (32 bits), but we can use a long (64 bits).
The first 49 bits (7 x 7) of the long represent the board. However there are some bits
that are not valid and never used, i.e. 0,1,5,6,7,8,12,13 and so on. Checking of
possible moves and applying them can be done by using simple bit operations.

A recursive function is then used to check all possible moves for a given board,
applying each valid move and calling itself with the resulting board. The recursion is
done "in reverse", starting from the goal board. While this is not conceptually faster [a],
it allows for a minimum amount of bit operations in the recursion:

To reverse a move we can simply check
- (board & twoBalls) == 0 and
- (board & oneBall) != 0
where "twoBalls" indicates the two ball that would need to be added for this reversed move.
If we instead used the intuitive search direction, the same check would require additional
binary operations, since a simple inversion of the check would not work [b].

Paper [1] shows how the moves can be ordered to almost instantly find a solution.
Website [2] gives a nice overview of binary operations and some tricks that
can be applied.

[a] Playing the game in reverse is simply the inversion of the original game - just remove all
balls from the board and place ball where there were none before and you'll understand
what I mean.
[b] There is no "single" binary operation to check if two specific bits are set, but there
is one to check if they are both zero. There is further a binary operation to check if a specific
bit is set.

[1] http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.6.4826 (download at the top)
[2] http://graphics.stanford.edu/~seander/bithacks.html
