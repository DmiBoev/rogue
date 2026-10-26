package domain

const MapWidth, MapHeight = 78, 30

type Door struct {
	Position Position
	Color    uint8
}
type Corridor struct{ Cells []Position }
type Level struct {
	Number      int
	Map         Mapper
	Start, Exit Position
	Corridors   []Corridor
	Doors       []Door
	Enemies     []Enemy
	Items       []Item
}

func (l *Level) InBounds(p Position) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < l.Map.Width && p.Y < l.Map.Height
}
func (l *Level) Walkable(p Position) bool { return l.InBounds(p) && l.Map.Tiles[p.Y][p.X].IsWalkable() }
func (l *Level) LockAt(p Position) *Door {
	for i := range l.Doors {
		if l.Doors[i].Position == p {
			return &l.Doors[i]
		}
	}
	return nil
}
func (l *Level) EnemyAt(p Position) *Enemy {
	for i := range l.Enemies {
		if l.Enemies[i].Health > 0 && l.Enemies[i].Position == p {
			return &l.Enemies[i]
		}
	}
	return nil
}
func (l *Level) ItemAt(p Position) *Item {
	for i := range l.Items {
		if l.Items[i].Position == p {
			return &l.Items[i]
		}
	}
	return nil
}
func (l *Level) RoomAt(p Position) *Room {
	for i := range l.Map.Rooms {
		if l.Map.Rooms[i].Contains(p) {
			return &l.Map.Rooms[i]
		}
	}
	return nil
}

// GenerateLevel uses a randomized DFS spanning tree on a 3x3 grid.
// Corridors stay inside the two adjacent sections they connect.
func GenerateLevel(number, difficulty int, rng RNG) *Level {
	for {
		l := generateLayout(number, rng)
		if l.Validate() {
			l.populate(difficulty, rng)
			return l
		}
	}
}

func generateLayout(number int, rng RNG) *Level {
	l := &Level{Number: number, Map: Mapper{Width: MapWidth, Height: MapHeight}}
	l.Map.Tiles = make([][]Tile, MapHeight)
	for y := range l.Map.Tiles {
		l.Map.Tiles[y] = make([]Tile, MapWidth)
	}
	for i := 0; i < 9; i++ {
		w, h := 12+rng.Intn(10), 5+rng.Intn(4)
		r := Room{ID: RoomID(i), X: (i%3)*26 + 1 + rng.Intn(25-w), Y: (i/3)*10 + 1 + rng.Intn(9-h), Width: w, Height: h}
		l.Map.Rooms = append(l.Map.Rooms, r)
		for y := r.Y; y < r.Y+h; y++ {
			for x := r.X; x < r.X+w; x++ {
				t := TileFloor
				if y == r.Y || y == r.Y+h-1 || x == r.X || x == r.X+w-1 {
					t = TileWall
				}
				l.Map.SetTile(Position{x, y}, t)
			}
		}
	}
	start := rng.Intn(9)
	l.Map.StartRoom = RoomID(start)
	r := l.Map.Rooms[start]
	l.Start = interior(r, rng)
	visited := map[int]bool{start: true}
	edges := 0
	var visit func(int)
	visit = func(a int) {
		neighbors := []int{}
		for _, d := range Directions4 {
			x, y := a%3+d.DX, a/3+d.DY
			if x >= 0 && x < 3 && y >= 0 && y < 3 {
				neighbors = append(neighbors, y*3+x)
			}
		}
		for i := len(neighbors) - 1; i > 0; i-- {
			j := rng.Intn(i + 1)
			neighbors[i], neighbors[j] = neighbors[j], neighbors[i]
		}
		for _, b := range neighbors {
			if visited[b] {
				continue
			}
			visited[b] = true
			path := l.connect(l.Map.Rooms[a], l.Map.Rooms[b])
			if edges < 3 {
				color := uint8(1 << edges)
				// The matching key is placed on the already reachable side of the edge.
				p := l.Map.Rooms[b].Center()
				for _, cell := range path {
					if l.Map.Rooms[b].Contains(cell) {
						p = cell
						break
					}
				}
				l.Doors = append(l.Doors, Door{Position: p, Color: color})
				keyPos := interior(l.Map.Rooms[a], rng)
				for keyPos == l.Start || l.ItemAt(keyPos) != nil {
					keyPos = interior(l.Map.Rooms[a], rng)
				}
				if RoomID(a) == l.Map.StartRoom {
					for _, cell := range path {
						if !l.Map.Rooms[a].Contains(cell) {
							keyPos = cell
							break
						}
					}
				}
				l.Items = append(l.Items, Item{Type: ItemTypeKey, Name: []string{"Red key", "Blue key", "Yellow key"}[edges], Color: color, Position: keyPos})
			}
			edges++
			l.Map.EndRoom = RoomID(b)
			visit(b)
		}
	}
	visit(start)
	l.Exit = interior(l.Map.Rooms[l.Map.EndRoom], rng)
	l.Map.SetTile(l.Exit, TileStairsDown)
	return l
}

func interior(r Room, rng RNG) Position {
	return Position{r.X + 1 + rng.Intn(r.Width-2), r.Y + 1 + rng.Intn(r.Height-2)}
}

func (l *Level) connect(a, b Room) []Position {
	p, end := a.Center(), b.Center()
	var points []Position
	line := func(target Position) {
		for p != target {
			points = append(points, p)
			if p.X != target.X {
				if p.X < target.X {
					p.X++
				} else {
					p.X--
				}
			} else if p.Y < target.Y {
				p.Y++
			} else {
				p.Y--
			}
		}
	}
	if a.ID/3 == b.ID/3 {
		mid := (int(a.ID)%3 + int(b.ID)%3 + 1) * 13
		line(Position{mid, p.Y})
		line(Position{mid, end.Y})
		line(end)
	} else {
		mid := (int(a.ID)/3 + int(b.ID)/3 + 1) * 5
		line(Position{p.X, mid})
		line(Position{end.X, mid})
		line(end)
	}
	points = append(points, end)
	for _, cell := range points {
		t := l.Map.Tiles[cell.Y][cell.X].Type
		if t == TileFloor || t == TileDoor {
			continue
		}
		if t == TileWall {
			l.Map.SetTile(cell, TileDoor)
		} else {
			l.Map.SetTile(cell, TileCorridor)
		}
	}
	l.Corridors = append(l.Corridors, Corridor{Cells: points})
	return points
}

// Validate explores (position, key mask) states, so a reachable exit cannot
// accidentally hide an inaccessible room or a circular key dependency.
func (l *Level) Validate() bool {
	if len(l.Map.Rooms) != 9 || !l.Walkable(l.Start) || !l.Walkable(l.Exit) {
		return false
	}
	type state struct {
		P    Position
		Keys uint8
	}
	queue := []state{{P: l.Start}}
	seen := map[state]bool{}
	reached := map[Position]bool{}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		for _, i := range l.Items {
			if i.Position == s.P && i.Type == ItemTypeKey {
				s.Keys |= i.Color
			}
		}
		if seen[s] {
			continue
		}
		seen[s] = true
		reached[s.P] = true
		for _, d := range Directions4 {
			p := s.P.Add(Position{d.DX, d.DY})
			if !l.Walkable(p) {
				continue
			}
			if lock := l.LockAt(p); lock != nil && lock.Color&s.Keys == 0 {
				continue
			}
			queue = append(queue, state{p, s.Keys})
		}
	}
	if !reached[l.Exit] {
		return false
	}
	for _, r := range l.Map.Rooms {
		if !reached[r.Center()] {
			return false
		}
	}
	for _, i := range l.Items {
		if i.Type == ItemTypeKey && !reached[i.Position] {
			return false
		}
	}
	return true
}

func (l *Level) populate(difficulty int, rng RNG) {
	free := func(r Room) Position {
		for {
			p := interior(r, rng)
			if p != l.Start && p != l.Exit && l.EnemyAt(p) == nil && l.ItemAt(p) == nil && l.LockAt(p) == nil {
				return p
			}
		}
	}
	types := []EnemyType{EnemyTypeZomby, EnemyTypeVampire, EnemyTypeGhost, EnemyTypeOgre, EnemyTypeSnakeMage, EnemyTypeMimic}
	for _, r := range l.Map.Rooms {
		if r.ID == l.Map.StartRoom {
			continue
		}
		count := max(1, 1+l.Number/6+difficulty)
		for n := 0; n < count; n++ {
			kind := types[rng.Intn(len(types))]
			e := Enemy{Type: kind, Position: free(r)}
			switch kind {
			case EnemyTypeZomby:
				e.Name = "Zombie"
				e.Symbol = 'z'
				e.Health = 20
				e.Strength = 4
				e.Agility = 2
				e.Spite = 6
			case EnemyTypeVampire:
				e.Name = "Vampire"
				e.Symbol = 'v'
				e.Health = 24
				e.Strength = 4
				e.Agility = 8
				e.Spite = 9
			case EnemyTypeGhost:
				e.Name = "Ghost"
				e.Symbol = 'g'
				e.Health = 8
				e.Strength = 2
				e.Agility = 8
				e.Spite = 4
			case EnemyTypeOgre:
				e.Name = "Ogre"
				e.Symbol = 'O'
				e.Health = 35
				e.Strength = 9
				e.Agility = 2
				e.Spite = 6
			case EnemyTypeSnakeMage:
				e.Name = "Serpent mage"
				e.Symbol = 's'
				e.Health = 14
				e.Strength = 3
				e.Agility = 11
				e.Spite = 10
			case EnemyTypeMimic:
				e.Name = "Mimic"
				e.Symbol = 'm'
				e.Health = 26
				e.Strength = 2
				e.Agility = 8
				e.Spite = 3
			}
			e.Health += max(0, l.Number-1+difficulty) * 2
			e.MaxHealth = e.Health
			e.Strength += max(0, (l.Number-1)/4+difficulty)
			l.Enemies = append(l.Enemies, e)
		}
		items := max(1, 3-l.Number/8-difficulty)
		for n := 0; n < items; n++ {
			kind := rng.Intn(4)
			item := Item{Position: free(r)}
			switch kind {
			case 0:
				item.Type = ItemTypeFood
				item.Name = "Ration"
				item.Bonus = 20
			case 1, 2:
				item.Stat = []Stat{StatHealth, StatStrength, StatAgility}[rng.Intn(3)]
				item.Bonus = 2
				if item.Stat == StatHealth {
					item.Bonus = 10
				}
				if kind == 1 {
					item.Type = ItemTypePotion
					item.Name = string(item.Stat) + " elixir"
				} else {
					item.Type = ItemTypeScroll
					item.Name = string(item.Stat) + " scroll"
				}
			case 3:
				item.Type = ItemTypeWeapon
				item.Name = "Sword"
				item.Bonus = 4 + l.Number/3
			}
			// Struggling players receive more healing, not just more random loot.
			if difficulty < 0 && n == 0 {
				item.Type = ItemTypeFood
				item.Name = "Ration"
				item.Stat = ""
				item.Bonus = 25
			}
			l.Items = append(l.Items, item)
		}
	}
}
