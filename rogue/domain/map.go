// map/generation/tiles
package domain

import "math/rand"

type TileType string

const (
	TileUnknown    TileType = "unknown_tile"  //не исследованно
	TileWall       TileType = "wall_tile"     //стена
	TileFloor      TileType = "floor_tile"    //пол
	TileCorridor   TileType = "corridor_tile" //коридор
	TileDoor       TileType = "door_tile"     //дверь
	TileStairsDown TileType = "stairs_down"   //лестница вниз пока что
)

type Tile struct {
	Type      TileType
	Exploered bool //видел ли игрок вообщше
	Visible   bool //видит прямо сейчас
}

//проходимость
func (t Tile) IsWalkable() bool {
	switch t.Type {
	case TileFloor, TileCorridor, TileDoor, TileStairsDown:
		return true
	default:
		return false
	}
}

//стены и закрытые двери блокируют обзор
func (t Tile) IsOralcul() bool {
	return t.Type == TileWall
}

//room

type RoomID int

type Room struct {
	ID     RoomID
	X, Y   int //левый верхний угол
	Width  int
	Height int
}

func (r Room) Center() Position {
	return Position{r.X + r.Width/2, r.Y + r.Height/2}
}

func (r Room) Contains(p Position) bool {
	return p.X >= r.X && p.X < r.X+r.Width &&
		p.Y >= r.Y && p.Y < r.Y+r.Height
}

func (r Room) RandomPoint(rng RNG) Position {
	return Position{
		r.X + rng.Intn(r.Width),
		r.Y + rng.Intn(r.Height),
	}
}

type Mapper struct {
	Width, Height int
	Tiles         [][]Tile
	Rooms         []Room
	StartRoom     RoomID
	EndRoom       RoomID
}

func (m *Mapper) SetTile(p Position, t TileType) {
	if p.X >= 0 && p.X < m.Width && p.Y >= 0 && p.Y < m.Height {
		m.Tiles[p.Y][p.X].Type = t
	}
}

func (m *Mapper) CarveRoom(r Room) {
	for y := r.Y; y < r.Y+r.Height; y++ {
		for x := r.X; x < r.X+r.Width; x++ {
			m.SetTile(Position{x, y}, TileFloor)
		}
	}
}

func (m *Mapper) CarveCorridor(x1, y1, x2, y2 int) {
	rn := rand.Float32()
	if rn >= 0.5 {
		for x := min(x1, x2); x <= max(x1, x2); x++ {
			m.SetTile(Position{x, y1}, TileCorridor)
		}
		for y := min(y1, y2); y <= max(y1, y2); y++ {
			m.SetTile(Position{x2, y}, TileCorridor)
		}
	} else {
		for y := min(y1, y2); y <= max(y1, y2); y++ {
			m.SetTile(Position{x2, y}, TileCorridor)
		}
		for x := min(x1, x2); x <= max(x1, x2); x++ {
			m.SetTile(Position{x, y1}, TileCorridor)
		}
	}
}

func (n *BSPNode) Split(rng RNG, depth int) {
	if depth <= 0 {
		return
	}

	const MinSize = 10
	if n.Width < MinSize*2 || n.Height < MinSize*2 {
		return
	}

	horizontal := rng.Intn(2) == 0
	if float32(n.Width) > float32(n.Height)*1.25 {
		horizontal = false // режем вертикально (Left/Right)
	} else if float32(n.Height) > float32(n.Width)*1.25 {
		horizontal = true // режем горизонтально (Top/Bottom)
	}

	if horizontal {
		splitY := n.Y + MinSize + rng.Intn(n.Height-MinSize*2)
		n.Left = &BSPNode{X: n.X, Y: n.Y, Width: n.Width, Height: splitY - n.Y}
		n.Right = &BSPNode{X: n.X, Y: splitY, Width: n.Width, Height: n.Y + n.Height - splitY}
	} else {
		splitX := n.X + MinSize + rng.Intn(n.Width-MinSize*2)
		n.Left = &BSPNode{X: n.X, Y: n.Y, Width: splitX - n.X, Height: n.Height}
		n.Right = &BSPNode{X: splitX, Y: n.Y, Width: n.X + n.Width - splitX, Height: n.Height}
	}

	n.Left.Split(rng, depth-1)
	n.Right.Split(rng, depth-1)
}
