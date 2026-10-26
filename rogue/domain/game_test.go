package domain

import (
	"reflect"
	"testing"
)

func TestGeneratedLevelsHaveReachableRoomsAndKeys(t *testing.T) {
	for seed := int64(1); seed <= 20; seed++ {
		rng := Random{State: uint64(seed)}
		for floor := 1; floor <= LastLevel; floor++ {
			l := GenerateLevel(floor, int(seed%6)-2, &rng)
			if !l.Validate() {
				t.Fatalf("seed %d floor %d is inaccessible", seed, floor)
			}
			if len(l.Corridors) != 8 || len(l.Doors) != 3 {
				t.Fatal("expected a spanning tree and three locks")
			}
			occupied := map[Position]bool{l.Start: true, l.Exit: true}
			for _, e := range l.Enemies {
				if l.Map.Rooms[l.Map.StartRoom].Contains(e.Position) {
					t.Fatal("enemy in starting room")
				}
				if occupied[e.Position] {
					t.Fatal("overlapping actors")
				}
				occupied[e.Position] = true
			}
			for _, i := range l.Items {
				if l.Map.Rooms[l.Map.StartRoom].Contains(i.Position) {
					t.Fatal("item in starting room")
				}
				if occupied[i.Position] {
					t.Fatal("overlapping item")
				}
				occupied[i.Position] = true
			}
		}
	}
}

func TestSeedReproducesWholeSession(t *testing.T) {
	if !reflect.DeepEqual(NewGame(123), NewGame(123)) {
		t.Fatal("seed does not reproduce game")
	}
}

func TestValidationRejectsKeyBehindItsOwnDoor(t *testing.T) {
	rng := Random{State: 12}
	l := GenerateLevel(1, 0, &rng)
	for i := range l.Items {
		if l.Items[i].Type == ItemTypeKey && l.Items[i].Color == 1 {
			l.Items[i].Position = l.Doors[0].Position
		}
	}
	if l.Validate() {
		t.Fatal("accepted a key behind its own locked door")
	}
}

func TestTurnsInventoryAndEffects(t *testing.T) {
	g := NewGame(7)
	g.Level.Enemies = nil
	// A blocked action does not advance time.
	g.Hero.Position = Position{0, 0}
	if g.Move(DirWest) || g.Turn != 0 {
		t.Fatal("blocked move consumed turn")
	}
	g.Hero.Position = g.Level.Start
	g.Hero.Inventory = []Item{{Type: ItemTypePotion, Stat: StatHealth, Bonus: 10, Name: "Health elixir"}}
	if g.Use(ItemTypePotion, 9) || g.Turn != 0 {
		t.Fatal("invalid choice consumed turn")
	}
	if !g.Use(ItemTypePotion, 1) || g.Hero.MaxHealth != 70 || g.Stats.Potions != 1 {
		t.Fatal("elixir not applied")
	}
	g.Hero.Health = 2
	for i := 0; i < 20; i++ {
		g.Move(Direction{})
	}
	if g.Hero.Health != 1 || g.Hero.MaxHealth != 60 || len(g.Hero.Effects) != 0 {
		t.Fatal("elixir expiration must leave living player at least one HP")
	}
}

func TestVampireFirstStrikeAndDeathCannotBeReversedByExpiration(t *testing.T) {
	g := NewGame(4)
	g.Level.Enemies = nil
	e := Enemy{Type: EnemyTypeVampire, Name: "Vampire", Health: 20, MaxHealth: 20}
	g.attack(&e)
	if e.Health != 20 || !e.FirstDodged || g.Stats.Misses != 1 {
		t.Fatal("first strike must miss vampire")
	}
	g.Hero.Health = 0
	g.Hero.Effects = []Effect{{Stat: StatHealth, Bonus: 10, Turns: 1}}
	g.endTurn()
	if !g.Finished || g.Hero.Health > 0 {
		t.Fatal("an expired elixir revived a dead player")
	}
}

func TestBackpackCapacityAndWeaponSwap(t *testing.T) {
	g := NewGame(5)
	g.Level.Enemies = nil
	g.Hero.Inventory = nil
	for i := 0; i < 9; i++ {
		g.Hero.Inventory = append(g.Hero.Inventory, Item{Type: ItemTypeFood})
	}
	g.Level.Items = []Item{{Type: ItemTypeFood, Position: g.Hero.Position}}
	g.pickup()
	if len(g.Hero.Items(ItemTypeFood)) != 9 || len(g.Level.Items) != 1 {
		t.Fatal("capacity not enforced")
	}
	g.Hero.Inventory = []Item{{Type: ItemTypeWeapon, Name: "Sword", Bonus: 4}}
	g.Hero.Weapon = &Item{Type: ItemTypeWeapon, Name: "Dagger", Bonus: 2}
	if !g.Use(ItemTypeWeapon, 1) || g.Hero.Weapon.Name != "Sword" {
		t.Fatal("could not equip")
	}
	found := false
	for _, item := range g.Level.Items {
		if item.Name == "Dagger" && distance(item.Position, g.Hero.Position) == 1 {
			found = true
		}
	}
	if !found {
		t.Fatal("old weapon not dropped onto adjacent cell")
	}
	if !g.Use(ItemTypeWeapon, 0) || g.Hero.Weapon != nil || len(g.Hero.Items(ItemTypeWeapon)) != 1 {
		t.Fatal("unequip lost weapon")
	}
}

func TestEquippingFromBackpackHoldsWeaponInsteadOfDroppingIt(t *testing.T) {
	g := NewGame(5)
	g.Level.Enemies = nil
	before := len(g.Level.Items)
	if !g.Use(ItemTypeWeapon, 1) {
		t.Fatal("could not equip initial dagger")
	}
	if g.Hero.Weapon == nil || g.Hero.Weapon.Name != "Dagger" || len(g.Hero.Items(ItemTypeWeapon)) != 0 {
		t.Fatal("weapon not transferred to hand")
	}
	if len(g.Level.Items) != before {
		t.Fatal("equipping dropped selected weapon")
	}
	low, high := g.Hero.DamageRange(g.Hero.Weapon)
	if low != 12 || high != 14 {
		t.Fatalf("unexpected weapon damage %d-%d", low, high)
	}
	if !g.Use(ItemTypeWeapon, 0) || g.Hero.Weapon != nil || len(g.Hero.Items(ItemTypeWeapon)) != 1 {
		t.Fatal("unequipping lost weapon")
	}
	if len(g.Level.Items) != before {
		t.Fatal("unequipping dropped weapon")
	}
	low, high = g.Hero.DamageRange(nil)
	if low != 7 || high != 9 {
		t.Fatal("incorrect unarmed damage")
	}
}

// Play an entire campaign through the public turn API with enemies removed.
// This checks key pickup, locks, automatic transitions, and the final victory.
func TestCanCompleteAllTwentyOneFloors(t *testing.T) {
	g := NewGame(42)
	for !g.Finished {
		g.Level.Enemies = nil
		for _, color := range []uint8{1, 2, 4} {
			for _, item := range append([]Item(nil), g.Level.Items...) {
				if item.Type == ItemTypeKey && item.Color == color {
					walkTo(t, g, item.Position)
				}
			}
		}
		walkTo(t, g, g.Level.Exit)
	}
	if !g.Won || g.Stats.Level != LastLevel {
		t.Fatal("campaign did not finish at floor 21")
	}
}

func walkTo(t *testing.T, g *Game, target Position) {
	t.Helper()
	start := g.Hero.Position
	queue := []Position{start}
	parent := map[Position]Position{start: start}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		if p == target {
			break
		}
		for _, d := range Directions4 {
			n := p.Add(Position{d.DX, d.DY})
			if _, seen := parent[n]; seen || !g.Level.Walkable(n) {
				continue
			}
			if lock := g.Level.LockAt(n); lock != nil && g.Hero.Keys&lock.Color == 0 {
				continue
			}
			// Avoid descending while walking to a key.
			if n == g.Level.Exit && n != target {
				continue
			}
			parent[n] = p
			queue = append(queue, n)
		}
	}
	if _, found := parent[target]; !found {
		t.Fatalf("cannot reach %v on floor %d with keys %d", target, g.Level.Number, g.Hero.Keys)
	}
	var path []Position
	for p := target; p != start; p = parent[p] {
		path = append(path, p)
	}
	for i := len(path) - 1; i >= 0; i-- {
		p := path[i]
		d := Direction{DX: p.X - g.Hero.Position.X, DY: p.Y - g.Hero.Position.Y}
		if !g.Move(d) {
			t.Fatal("valid path move rejected")
		}
	}
}
