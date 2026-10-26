package domain

import "fmt"

const LastLevel = 21

type Statistics struct {
	Treasure, Level, Kills, Food, Potions, Scrolls, Hits, Misses, HitsTaken, Steps int
}

// Game contains only serializable domain state. UI and storage are external.
type Game struct {
	ID            string
	Hero          Hero
	Level         *Level
	Stats         Statistics
	Random        Random
	Turn          int
	Finished, Won bool
	Difficulty    int
	FloorDamage   int
	Messages      []string
}

// Random keeps its state in the save so continuing preserves future rolls.
type Random struct{ State uint64 }

func (r *Random) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	if r.State == 0 {
		r.State = 0x9e3779b97f4a7c15
	}
	x := r.State
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	r.State = x
	return int(x % uint64(n))
}

func NewGame(seed int64) *Game {
	g := &Game{ID: fmt.Sprintf("%d", seed), Random: Random{State: uint64(seed)}, Hero: Hero{MaxHealth: 60, Health: 60, Strength: 7, Agility: 6}}
	g.Hero.Inventory = []Item{{Type: ItemTypeFood, Name: "Ration", Bonus: 20}, {Type: ItemTypeWeapon, Name: "Dagger", Bonus: 3}}
	g.enterLevel(1)
	g.Log("Find the stairs >. Survive all 21 floors.")
	return g
}

func (g *Game) Log(message string) {
	g.Messages = append(g.Messages, message)
	if len(g.Messages) > 5 {
		g.Messages = g.Messages[len(g.Messages)-5:]
	}
}

func (g *Game) enterLevel(number int) {
	g.Level = GenerateLevel(number, g.Difficulty, &g.Random)
	g.Hero.Position = g.Level.Start
	g.Hero.Keys = 0
	g.Stats.Level = number
	g.FloorDamage = 0
	g.UpdateVisibility()
}

// Move returns true only if an action consumed a turn.
func (g *Game) Move(direction Direction) bool {
	if g.Finished {
		return false
	}
	if g.Hero.Sleep > 0 {
		g.Hero.Sleep--
		g.Log("You wake up.")
		g.endTurn()
		return true
	}
	dest := g.Hero.Position.Add(Position{direction.DX, direction.DY})
	if !g.Level.Walkable(dest) {
		return false
	}
	if door := g.Level.LockAt(dest); door != nil && door.Color&g.Hero.Keys == 0 {
		g.Log("This door needs a matching key.")
		return false
	}
	if enemy := g.Level.EnemyAt(dest); enemy != nil {
		g.attack(enemy)
	} else {
		g.Hero.Position = dest
		if direction.DX != 0 || direction.DY != 0 {
			g.Stats.Steps++
		}
		g.pickup()
		if dest == g.Level.Exit {
			if g.Level.Number == LastLevel {
				g.Finished, g.Won = true, true
				g.Log("You escaped the dungeon!")
				return true
			}
			if g.FloorDamage < 8 && g.Hero.Health*4 > g.Hero.MaxHealth*3 {
				g.Difficulty = min(3, g.Difficulty+1)
			}
			if g.FloorDamage > 20 || g.Hero.Health*3 < g.Hero.MaxHealth {
				g.Difficulty = max(-2, g.Difficulty-1)
			}
			g.enterLevel(g.Level.Number + 1)
			g.Log(fmt.Sprintf("You descend to floor %d.", g.Level.Number))
			g.Turn++
			return true
		}
	}
	g.endTurn()
	return true
}

func (g *Game) pickup() {
	for i := 0; i < len(g.Level.Items); {
		item := g.Level.Items[i]
		if item.Position != g.Hero.Position {
			i++
			continue
		}
		if item.Type == ItemTypeKey {
			g.Hero.Keys |= item.Color
		} else {
			if len(g.Hero.Items(item.Type)) >= 9 {
				g.Log("Your backpack is full for this item type.")
				i++
				continue
			}
			g.Hero.Inventory = append(g.Hero.Inventory, item)
		}
		g.Log("Picked up " + item.Name + ".")
		g.Level.Items = append(g.Level.Items[:i], g.Level.Items[i+1:]...)
	}
}

func (g *Game) Use(kind ItemType, selection int) bool {
	if g.Finished {
		return false
	}
	if g.Hero.Sleep > 0 {
		return g.Move(Direction{})
	}
	if kind == ItemTypeWeapon && selection == 0 {
		if g.Hero.Weapon == nil {
			return false
		}
		if len(g.Hero.Items(kind)) >= 9 {
			g.Log("No room to store your weapon.")
			return false
		}
		name := g.Hero.Weapon.Name
		g.Hero.Inventory = append(g.Hero.Inventory, *g.Hero.Weapon)
		g.Hero.Weapon = nil
		g.Log(name + " stored in backpack. You are now unarmed.")
		g.endTurn()
		return true
	}
	index, count := -1, 0
	for i, item := range g.Hero.Inventory {
		if item.Type == kind {
			count++
			if count == selection {
				index = i
				break
			}
		}
	}
	if index < 0 {
		return false
	}
	item := g.Hero.Inventory[index]
	if kind == ItemTypeWeapon && g.Hero.Weapon != nil {
		placed := false
		for _, d := range Directions4 {
			p := g.Hero.Position.Add(Position{d.DX, d.DY})
			if g.Level.Walkable(p) && g.Level.EnemyAt(p) == nil && g.Level.ItemAt(p) == nil && g.Level.LockAt(p) == nil && p != g.Level.Exit {
				old := *g.Hero.Weapon
				old.Position = p
				g.Level.Items = append(g.Level.Items, old)
				g.Log(fmt.Sprintf("%s dropped nearby at (%d, %d).", old.Name, p.X, p.Y))
				placed = true
				break
			}
		}
		if !placed {
			g.Log("No space nearby to drop the old weapon.")
			return false
		}
	}
	g.Hero.Inventory = append(g.Hero.Inventory[:index], g.Hero.Inventory[index+1:]...)
	switch kind {
	case ItemTypeFood:
		g.Hero.Health = min(g.Hero.MaxHealth, g.Hero.Health+item.Bonus)
		g.Stats.Food++
	case ItemTypePotion:
		g.Hero.modify(item.Stat, item.Bonus)
		g.Hero.Effects = append(g.Hero.Effects, Effect{Stat: item.Stat, Bonus: item.Bonus, Turns: 21})
		g.Stats.Potions++
	case ItemTypeScroll:
		g.Hero.modify(item.Stat, item.Bonus)
		g.Stats.Scrolls++
	case ItemTypeWeapon:
		g.Hero.Weapon = &item
	}
	if kind == ItemTypeWeapon {
		low, high := g.Hero.DamageRange(g.Hero.Weapon)
		g.Log(fmt.Sprintf("Equipped %s. Damage: %d-%d.", item.Name, low, high))
	} else {
		g.Log("Used " + item.Name + ".")
	}
	g.endTurn()
	return true
}

func (g *Game) endTurn() {
	g.Turn++
	g.enemiesTurn()
	var effects []Effect
	for _, effect := range g.Hero.Effects {
		effect.Turns--
		if effect.Turns <= 0 {
			g.Hero.modify(effect.Stat, -effect.Bonus)
			g.Log("An elixir wears off.")
		} else {
			effects = append(effects, effect)
		}
	}
	g.Hero.Effects = effects
	if g.Hero.Health <= 0 {
		g.Finished = true
		g.Log("You died. Press N to start again.")
	}
	g.UpdateVisibility()
}
