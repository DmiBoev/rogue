package domain

type Position struct {
	X, Y int
}

func (p Position) Add(d Position) Position {
	return (Position{p.X + d.X, p.Y + d.Y})
}

func (p Position) Eq(d Position) bool {
	return p.X == d.X && p.Y == d.Y
}

type Direction struct {
	DX, DY int
}

var (
	DirNorth = Direction{0, -1}
	DirSouth = Direction{0, 1}
	DirWest  = Direction{-1, 0}
	DirEast  = Direction{1, 0}
)
var Directions4 = []Direction{DirNorth, DirSouth, DirWest, DirEast}
