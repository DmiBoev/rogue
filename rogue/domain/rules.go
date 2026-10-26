package domain

import "fmt"

func distance(a, b Position) int { return abs(a.X-b.X) + abs(a.Y-b.Y) }
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (g *Game) attack(e *Enemy) {
	e.Engaged = true
	e.Hidden = false
	if e.Type == EnemyTypeVampire && !e.FirstDodged {
		e.FirstDodged = true
		g.Stats.Misses++
		g.Log("The vampire dodges your first strike.")
		return
	}
	if !g.hit(g.Hero.Agility, e.Agility) {
		g.Stats.Misses++
		g.Log("You miss " + e.Name + ".")
		return
	}
	base, _ := g.Hero.DamageRange(g.Hero.Weapon)
	damage := base + g.Random.Intn(3)
	e.Health -= damage
	g.Stats.Hits++
	g.Log(fmt.Sprintf("You hit %s for %d.", e.Name, damage))
	if e.Health <= 0 {
		gold := e.MaxHealth + e.Strength + e.Agility + e.Spite
		gold = gold/2 + g.Random.Intn(max(1, gold/2)) + g.Level.Number*2
		g.Stats.Treasure += gold
		g.Stats.Kills++
		g.Log(fmt.Sprintf("%s defeated. +%d treasure.", e.Name, gold))
	}
}

func (g *Game) hit(attack, defense int) bool {
	return g.Random.Intn(100) < max(10, min(95, 65+(attack-defense)*4))
}

func (g *Game) enemyAttack(e *Enemy) {
	e.Engaged = true
	e.Hidden = false
	hit := e.Counter || g.hit(e.Agility, g.Hero.Agility)
	e.Counter = false
	if e.Type == EnemyTypeOgre {
		e.Rest = 1
		e.Counter = true
	}
	if !hit {
		g.Log(e.Name + " misses.")
		return
	}
	damage := e.Strength + g.Random.Intn(3)
	g.Hero.Health -= damage
	g.Stats.HitsTaken++
	g.FloorDamage += damage
	g.Log(fmt.Sprintf("%s hits you for %d.", e.Name, damage))
	if e.Type == EnemyTypeVampire {
		loss := min(2, g.Hero.MaxHealth-1)
		g.Hero.MaxHealth -= loss
		g.Hero.Health = min(g.Hero.Health, g.Hero.MaxHealth)
	}
	if e.Type == EnemyTypeSnakeMage && g.Random.Intn(4) == 0 {
		g.Hero.Sleep = 1
		g.Log("The serpent puts you to sleep.")
	}
}

func (g *Game) canEnemyMove(p Position) bool {
	if !g.Level.Walkable(p) {
		return false
	}
	if d := g.Level.LockAt(p); d != nil && d.Color&g.Hero.Keys == 0 {
		return false
	}
	return p != g.Hero.Position && g.Level.EnemyAt(p) == nil
}

// chaseStep is BFS over adjacent cells. Occupied and locked cells block paths.
func (g *Game) chaseStep(start Position) (Position, bool) {
	queue := []Position{start}
	parent := map[Position]Position{start: start}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		if p == g.Hero.Position {
			for parent[p] != start {
				p = parent[p]
			}
			return p, true
		}
		for _, d := range Directions4 {
			n := p.Add(Position{d.DX, d.DY})
			if _, ok := parent[n]; ok {
				continue
			}
			if n != g.Hero.Position && !g.canEnemyMove(n) {
				continue
			}
			parent[n] = p
			queue = append(queue, n)
		}
	}
	return Position{}, false
}

func (g *Game) enemiesTurn() {
	for i := range g.Level.Enemies {
		e := &g.Level.Enemies[i]
		if e.Health <= 0 || g.Hero.Health <= 0 {
			continue
		}
		if e.Rest > 0 {
			e.Rest--
			continue
		}
		if distance(e.Position, g.Hero.Position) == 1 {
			g.enemyAttack(e)
			continue
		}
		if distance(e.Position, g.Hero.Position) <= e.Spite {
			if p, ok := g.chaseStep(e.Position); ok {
				e.Position = p
				continue
			}
		}
		room := g.Level.RoomAt(e.Position)
		if room == nil { // A monster stranded in a corridor wanders until it finds a room.
			d := Directions4[g.Random.Intn(4)]
			p := e.Position.Add(Position{d.DX, d.DY})
			if g.canEnemyMove(p) {
				e.Position = p
			}
			continue
		}
		if e.Type == EnemyTypeGhost && !e.Engaged {
			e.Hidden = g.Turn%3 != 0
			p := interior(*room, &g.Random)
			if g.canEnemyMove(p) {
				e.Position = p
			}
			continue
		}
		d := Directions4[g.Random.Intn(4)]
		delta := Position{d.DX, d.DY}
		if e.Type == EnemyTypeSnakeMage {
			diagonals := []Position{{1, 1}, {-1, 1}, {-1, -1}, {1, -1}}
			delta = diagonals[e.Pattern%4]
			e.Pattern++
		}
		steps := 1
		if e.Type == EnemyTypeOgre {
			steps = 2
		}
		for step := 0; step < steps; step++ {
			p := e.Position.Add(delta)
			if !room.Contains(p) || !g.canEnemyMove(p) {
				break
			}
			if delta.X != 0 && delta.Y != 0 && (!g.Level.Walkable(e.Position.Add(Position{delta.X, 0})) || !g.Level.Walkable(e.Position.Add(Position{0, delta.Y}))) {
				break
			}
			e.Position = p
		}
	}
}

func (g *Game) UpdateVisibility() {
	l := g.Level
	for y := range l.Map.Tiles {
		for x := range l.Map.Tiles[y] {
			l.Map.Tiles[y][x].Visible = false
		}
	}
	room := l.RoomAt(g.Hero.Position)
	for y := max(0, g.Hero.Position.Y-10); y < min(l.Map.Height, g.Hero.Position.Y+11); y++ {
		for x := max(0, g.Hero.Position.X-12); x < min(l.Map.Width, g.Hero.Position.X+13); x++ {
			p := Position{x, y}
			if room != nil && l.Map.Tiles[g.Hero.Position.Y][g.Hero.Position.X].Type != TileDoor {
				if !room.Contains(p) && distance(p, g.Hero.Position) > 3 {
					continue
				}
			} else if distance(p, g.Hero.Position) > 7 {
				continue
			}
			if g.lineOfSight(g.Hero.Position, p) {
				t := &l.Map.Tiles[y][x]
				t.Visible = true
				t.Exploered = true
			}
		}
	}
	// A room is revealed as a unit while the player is inside its interior.
	if room != nil && l.Map.Tiles[g.Hero.Position.Y][g.Hero.Position.X].Type == TileFloor || room != nil && g.Hero.Position == l.Exit {
		for y := room.Y; y < room.Y+room.Height; y++ {
			for x := room.X; x < room.X+room.Width; x++ {
				l.Map.Tiles[y][x].Visible = true
				l.Map.Tiles[y][x].Exploered = true
			}
		}
	}
}

// Bresenham rays stop at walls, void, and doors for which no key is held.
func (g *Game) lineOfSight(a, b Position) bool {
	dx, dy := abs(b.X-a.X), -abs(b.Y-a.Y)
	sx, sy := 1, 1
	if a.X > b.X {
		sx = -1
	}
	if a.Y > b.Y {
		sy = -1
	}
	err := dx + dy
	for {
		if a == b {
			return true
		}
		t := g.Level.Map.Tiles[a.Y][a.X].Type
		if t == TileWall || t == TileUnknown || t == "" {
			return false
		}
		if d := g.Level.LockAt(a); d != nil && d.Color&g.Hero.Keys == 0 {
			return false
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			a.X += sx
		}
		if e2 <= dx {
			err += dx
			a.Y += sy
		}
	}
}
