package presenter

import (
	"reflect"
	datalayer "rogue/rogue/dataLayer"
	"rogue/rogue/domain"
	ui "rogue/rogue/presentation"
	"strings"
	"testing"
)

type scriptedView struct {
	inputs []ui.Input
	frames []ui.Frame
}

func (v *scriptedView) Size() (int, int)  { return 80, 41 }
func (v *scriptedView) Render(f ui.Frame) { v.frames = append(v.frames, f) }
func (v *scriptedView) ReadInput() ui.Input {
	if len(v.inputs) == 0 {
		return ui.Input{Action: "quit"}
	}
	e := v.inputs[0]
	v.inputs = v.inputs[1:]
	return e
}

func TestGameLoopMenusDoNotAdvanceTurnsAndResumeIsExact(t *testing.T) {
	repo := datalayer.NewJSONRepository(t.TempDir())
	v := &scriptedView{inputs: []ui.Input{
		{Action: "new"}, {Action: "weapon"}, {Action: "select", Choice: 1},
		{Action: "food"}, {Action: "cancel"}, {Action: "records"}, {Action: "cancel"},
		{Action: "menu"}, {Action: "continue"}, {Action: "wait"}, {Action: "quit"},
	}}
	p := New(v, repo, 123)
	if err := p.Run(); err != nil {
		t.Fatal(err)
	}
	if p.game.Turn != 2 || p.game.Hero.Weapon == nil {
		t.Fatal("menu altered time or equipping failed")
	}
	if len(v.frames) < 10 {
		t.Fatal("screens were not rendered")
	}
	resume := New(&scriptedView{inputs: []ui.Input{{Action: "continue"}, {Action: "quit"}}}, repo, 0)
	if err := resume.Run(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(p.game, resume.game) {
		t.Fatal("continue did not restore session")
	}
}

func TestEscapeFromEveryBackpackReturnsToGame(t *testing.T) {
	for _, action := range []string{"weapon", "food", "potion", "scroll"} {
		t.Run(action, func(t *testing.T) {
			view := &scriptedView{inputs: []ui.Input{{Action: "new"}, {Action: action}, {Action: "cancel"}, {Action: "cancel"}}}
			p := New(view, datalayer.NewJSONRepository(t.TempDir()), 123)
			if err := p.Run(); err != nil {
				t.Fatal(err)
			}
			if p.mode != "game" || p.selection != "" {
				t.Fatal("Esc or repeated Esc left the game")
			}
			if p.game.Turn != 0 || p.game.Hero.Weapon != nil {
				t.Fatal("cancel changed player state")
			}
		})
	}
}

func TestWeaponMenuShowsEquippedWeaponAndEscapeKeepsIt(t *testing.T) {
	view := &scriptedView{inputs: []ui.Input{{Action: "new"}, {Action: "weapon"}, {Action: "select", Choice: 1}, {Action: "weapon"}, {Action: "cancel"}, {Action: "cancel"}}}
	p := New(view, datalayer.NewJSONRepository(t.TempDir()), 123)
	if err := p.Run(); err != nil {
		t.Fatal(err)
	}
	if p.game.Hero.Weapon == nil || p.game.Hero.Weapon.Name != "Dagger" || p.game.Turn != 1 || p.mode != "game" {
		t.Fatal("cancel unequipped weapon or left game")
	}
	found := false
	for _, frame := range view.frames {
		var text strings.Builder
		for _, cell := range frame.Cells {
			text.WriteRune(cell.Symbol)
		}
		if strings.Contains(text.String(), "In hand: Dagger | Damage: 12-14") && strings.Contains(text.String(), "0  Unequip to backpack (keep weapon)") {
			found = true
		}
	}
	if !found {
		t.Fatal("weapon screen does not explain current equipment")
	}
}

func TestMenuHasDedicatedKeyAndEscapeResumes(t *testing.T) {
	view := &scriptedView{inputs: []ui.Input{{Action: "new"}, {Action: "menu"}, {Action: "cancel"}}}
	p := New(view, datalayer.NewJSONRepository(t.TempDir()), 123)
	if err := p.Run(); err != nil {
		t.Fatal(err)
	}
	if p.mode != "game" || p.game.Turn != 0 {
		t.Fatal("menu/resume changed game")
	}
}

func TestEnterStartsOrResumesGame(t *testing.T) {
	repo := datalayer.NewJSONRepository(t.TempDir())
	first := New(&scriptedView{inputs: []ui.Input{{Action: "confirm"}, {Action: "wait"}}}, repo, 123)
	if err := first.Run(); err != nil {
		t.Fatal(err)
	}
	if first.game == nil || first.game.Turn != 1 {
		t.Fatal("Enter did not start game")
	}
	second := New(&scriptedView{inputs: []ui.Input{{Action: "confirm"}}}, repo, 0)
	if err := second.Run(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first.game, second.game) {
		t.Fatal("Enter did not restore existing save")
	}
}

func TestDeathThroughLoopCreatesSingleRecordAndAllowsRestart(t *testing.T) {
	repo := datalayer.NewJSONRepository(t.TempDir())
	g := domain.NewGame(55)
	g.Hero.Health = 1
	r := g.Level.RoomAt(g.Hero.Position)
	p := g.Hero.Position.Add(domain.Position{X: 1})
	if !r.Contains(p) || !g.Level.Walkable(p) {
		p = g.Hero.Position.Add(domain.Position{X: -1})
	}
	g.Level.Enemies = []domain.Enemy{{Name: "Ogre", Type: domain.EnemyTypeOgre, Position: p, Health: 10, MaxHealth: 10, Strength: 100, Counter: true}}
	if err := repo.Save(g); err != nil {
		t.Fatal(err)
	}
	app := New(&scriptedView{inputs: []ui.Input{{Action: "continue"}, {Action: "wait"}, {Action: "records"}, {Action: "cancel"}, {Action: "new"}, {Action: "quit"}}}, repo, 0)
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
	records, err := repo.History()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].Won || records[0].ID != g.ID {
		t.Fatal("death was not recorded exactly once")
	}
	if app.game.ID == g.ID || app.game.Finished {
		t.Fatal("restart did not reset game")
	}
}
