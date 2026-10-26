package domain

type BSPNode struct {
	X, Y, Width, Height int
	Left, Right         *BSPNode
	Room                *Room
}
