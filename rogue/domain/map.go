// map/generation/tiles
package domain

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
