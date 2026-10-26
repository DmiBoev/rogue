package presenter

import (
	"fmt"
	"rogue/rogue/domain"
	ui "rogue/rogue/presentation"
)

func (p *Presenter) frame() ui.Frame {
	f := ui.Frame{}
	w, h := p.view.Size()
	if w < 60 || h < 20 {
		f.Text(0, 0, "Resize terminal to at least 60 x 20. Q: save and quit.", ui.Yellow)
		return f
	}
	if p.mode == "menu" {
		f.Heading("ROGUE / Go + tcell")
		f.Text(2, 3, "N / Т  New game", ui.White)
		f.Text(2, 5, "C / С  Continue saved game", ui.White)
		f.Text(2, 7, "R / К  Records", ui.White)
		f.Text(2, 9, "Q / Й  Quit", ui.White)
		f.Text(2, 11, "Enter: continue, or start if there is no save.", ui.Cyan)
		f.Text(2, 12, "21 floors. Nine rooms per floor. One way out.", ui.Gray)
		if p.game != nil && !p.game.Finished {
			f.Text(2, 14, "Esc  Return to current game", ui.Cyan)
		}
		f.Text(2, 16, "English and Russian keyboard layouts are supported.", ui.Gray)
	} else if p.mode == "records" {
		f.Heading("RECORDS / sorted by treasure")
		f.Text(1, 2, " #  Gold   Lv Kill Food Pot Scr Hit Miss Hurt Steps Result", ui.Cyan)
		for i, r := range p.records {
			if i >= h-7 {
				break
			}
			s := r.Stats
			result := "died"
			if r.Won {
				result = "won"
			}
			f.Text(1, 3+i, fmt.Sprintf("%2d %6d %2d %4d %4d %3d %3d %3d %4d %4d %5d %s", i+1, s.Treasure, s.Level, s.Kills, s.Food, s.Potions, s.Scrolls, s.Hits, s.Misses, s.HitsTaken, s.Steps, result), ui.White)
		}
		if len(p.records) == 0 {
			f.Text(2, 4, "No completed runs yet.", ui.Gray)
		}
		f.Text(1, h-3, "Press any key to return. Q: quit.", ui.Gray)
	} else {
		p.gameFrame(&f, w, h)
	}
	if p.selection != "" {
		f = ui.Frame{}
		f.Heading("BACKPACK / " + string(p.selection))
		f.Text(2, 2, "Choose an item (1-9). Esc: cancel.", ui.Cyan)
		if p.selection == domain.ItemTypeWeapon {
			weapon := p.game.Hero.Weapon
			name := "unarmed"
			if weapon != nil {
				name = weapon.Name
			}
			low, high := p.game.Hero.DamageRange(weapon)
			f.Text(2, 3, fmt.Sprintf("In hand: %s | Damage: %d-%d", name, low, high), ui.Yellow)
			if weapon != nil {
				f.Text(2, 4, "0  Unequip to backpack (keep weapon)", ui.White)
			}
			f.Text(2, 5, "Backpack weapons - choose a number to EQUIP:", ui.Cyan)
		}
		items := p.game.Hero.Items(p.selection)
		for i, item := range items {
			label := fmt.Sprintf("%d  %s (+%d)", i+1, item.Name, item.Bonus)
			if p.selection == domain.ItemTypeWeapon {
				low, high := p.game.Hero.DamageRange(&item)
				label = fmt.Sprintf("%d  Equip %s | Damage: %d-%d", i+1, item.Name, low, high)
			}
			f.Text(2, 6+i, label, ui.White)
		}
		if len(items) == 0 {
			f.Text(2, 6, "No items of this type.", ui.Gray)
		}
		if p.selection == domain.ItemTypeWeapon && p.game.Hero.Weapon != nil {
			f.Text(2, h-4, "Equipping another weapon drops the old one nearby.", ui.Yellow)
		}
		f.Text(2, h-3, "Esc: return to game. Equipping/using costs one turn.", ui.Gray)
		if len(p.game.Messages) > 0 {
			f.Text(2, h-2, p.game.Messages[len(p.game.Messages)-1], ui.White)
		}
	}
	if p.notice != "" {
		f.Text(1, h-1, p.notice, ui.Red)
	}
	return f
}

func (p *Presenter) gameFrame(f *ui.Frame, w, h int) {
	g := p.game
	hero := g.Hero
	f.Heading(fmt.Sprintf("ROGUE  Floor %d/21  HP %d/%d  STR %d  AGI %d  Gold %d", g.Level.Number, hero.Health, hero.MaxHealth, hero.Strength, hero.Agility, g.Stats.Treasure))
	weapon := "unarmed"
	if hero.Weapon != nil {
		weapon = fmt.Sprintf("%s +%d", hero.Weapon.Name, hero.Weapon.Bonus)
	}
	f.Text(1, 1, fmt.Sprintf("Weapon: %s  Turn: %d  Keys: ", weapon, g.Turn), ui.Gray)
	keyX := len(fmt.Sprintf("Weapon: %s  Turn: %d  Keys: ", weapon, g.Turn)) + 1
	for i, mask := range []uint8{1, 2, 4} {
		color := ui.Gray
		if hero.Keys&mask != 0 {
			color = keyColor(mask)
		}
		f.Text(keyX+i*2, 1, "k", color)
	}
	mapW, mapH := min(w-2, g.Level.Map.Width), min(h-11, g.Level.Map.Height)
	ox := max(0, min(g.Level.Map.Width-mapW, hero.Position.X-mapW/2))
	oy := max(0, min(g.Level.Map.Height-mapH, hero.Position.Y-mapH/2))
	paint := func(pos domain.Position, symbol rune, color ui.Color) {
		if pos.X >= ox && pos.X < ox+mapW && pos.Y >= oy && pos.Y < oy+mapH {
			f.Cells = append(f.Cells, ui.Cell{X: 1 + pos.X - ox, Y: 3 + pos.Y - oy, Symbol: symbol, Color: color})
		}
	}
	for y := oy; y < oy+mapH; y++ {
		for x := ox; x < ox+mapW; x++ {
			pos := domain.Position{X: x, Y: y}
			tile := g.Level.Map.Tiles[y][x]
			if !tile.Exploered {
				continue
			}
			symbol, color := rune(' '), ui.Gray
			switch tile.Type {
			case domain.TileWall:
				symbol = '#'
			case domain.TileDoor:
				symbol = '+'
			case domain.TileFloor:
				if tile.Visible {
					symbol = '.'
				}
			case domain.TileCorridor:
				symbol = '#'
			case domain.TileStairsDown:
				if tile.Visible {
					symbol = '>'
					color = ui.Cyan
				}
			}
			if lock := g.Level.LockAt(pos); lock != nil {
				color = keyColor(lock.Color)
				if hero.Keys&lock.Color != 0 {
					symbol = '/'
				}
			}
			paint(pos, symbol, color)
		}
	}
	for _, item := range g.Level.Items {
		if g.Level.Map.Tiles[item.Position.Y][item.Position.X].Visible {
			color := ui.Yellow
			if item.Type == domain.ItemTypeKey {
				color = keyColor(item.Color)
			}
			paint(item.Position, item.Symbol(), color)
		}
	}
	for _, e := range g.Level.Enemies {
		if e.Health <= 0 || e.Hidden || !g.Level.Map.Tiles[e.Position.Y][e.Position.X].Visible {
			continue
		}
		symbol, color := e.Symbol, ui.White
		switch e.Type {
		case domain.EnemyTypeZomby:
			color = ui.Green
		case domain.EnemyTypeVampire:
			color = ui.Red
		case domain.EnemyTypeOgre:
			color = ui.Yellow
		case domain.EnemyTypeMimic:
			if !e.Engaged {
				symbol = '%'
			}
		}
		paint(e.Position, symbol, color)
	}
	paint(hero.Position, '@', ui.Cyan)
	y := mapH + 4
	low, high := hero.DamageRange(hero.Weapon)
	f.Text(1, y-1, fmt.Sprintf("Damage: %d-%d | H: equipment and backpack weapons", low, high), ui.Gray)
	for i, message := range g.Messages {
		f.Text(1, y+i, message, ui.White)
	}
	if g.Finished {
		text := "YOU DIED"
		if g.Won {
			text = "YOU WON"
		}
		f.Text(1, y+5, text+" | N: new game  R: records  Q: quit", ui.Yellow)
	} else {
		f.Text(1, y+5, "WASD/arrows: move  H: weapon J: food K: elixir E: scroll", ui.Cyan)
		f.Text(1, y+6, "Space/.: wait  R: records  M: menu  Q: save and quit", ui.Gray)
	}
}

func keyColor(mask uint8) ui.Color {
	switch mask {
	case 1:
		return ui.Red
	case 2:
		return ui.Blue
	default:
		return ui.Yellow
	}
}
