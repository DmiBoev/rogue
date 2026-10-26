package presentation

import (
	"github.com/gdamore/tcell/v2"
	"testing"
)

func TestTerminalRendersOnSimulationScreen(t *testing.T) {
	s := tcell.NewSimulationScreen("UTF-8")
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}
	defer s.Fini()
	s.SetSize(80, 24)
	terminal := &Terminal{screen: s}
	frame := Frame{}
	frame.Text(2, 3, "@zvO", Green)
	terminal.Render(frame)
	r, _, style, _ := s.GetContent(2, 3)
	fg, _, _ := style.Decompose()
	if r != '@' || fg != tcell.ColorGreen {
		t.Fatal("rendered glyph or color lost")
	}
	for _, c := range []struct {
		key    tcell.Key
		r      rune
		action string
	}{{tcell.KeyRune, 'W', "up"}, {tcell.KeyRune, 'h', "weapon"}, {tcell.KeyRune, 'j', "food"}, {tcell.KeyRune, 'k', "potion"}, {tcell.KeyRune, 'e', "scroll"}, {tcell.KeyRune, 'm', "menu"}, {tcell.KeyEscape, 0, "cancel"}, {tcell.KeyCtrlC, 0, "quit"}} {
		if got := MapEvent(tcell.NewEventKey(c.key, c.r, tcell.ModNone)); got.Action != c.action {
			t.Fatalf("%v: got %s", c, got.Action)
		}
	}
	if got := MapEvent(tcell.NewEventKey(tcell.KeyRune, '9', tcell.ModNone)); got.Action != "select" || got.Choice != 9 {
		t.Fatal("inventory selection lost")
	}
}

func TestRussianKeyboardLayout(t *testing.T) {
	for _, c := range []struct {
		r      rune
		action string
	}{
		{'т', "new"}, {'Т', "new"}, {'с', "continue"}, {'к', "records"}, {'й', "quit"},
		{'ц', "up"}, {'ф', "left"}, {'ы', "down"}, {'в', "right"},
		{'р', "weapon"}, {'о', "food"}, {'л', "potion"}, {'у', "scroll"}, {'ь', "menu"}, {'ю', "wait"},
	} {
		if got := MapEvent(tcell.NewEventKey(tcell.KeyRune, c.r, tcell.ModNone)); got.Action != c.action {
			t.Fatalf("%c: got %q, want %q", c.r, got.Action, c.action)
		}
	}
	if got := MapEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)); got.Action != "confirm" {
		t.Fatal("Enter not recognized")
	}
}
