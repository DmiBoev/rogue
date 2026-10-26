package presentation

import (
	"github.com/gdamore/tcell/v2"
	"unicode"
)

type Input struct {
	Action string
	Choice int
}

func MapEvent(event tcell.Event) Input {
	switch e := event.(type) {
	case *tcell.EventResize:
		return Input{Action: "resize"}
	case *tcell.EventKey:
		switch e.Key() {
		case tcell.KeyCtrlC:
			return Input{Action: "quit"}
		case tcell.KeyEscape:
			return Input{Action: "cancel"}
		case tcell.KeyUp:
			return Input{Action: "up"}
		case tcell.KeyDown:
			return Input{Action: "down"}
		case tcell.KeyLeft:
			return Input{Action: "left"}
		case tcell.KeyRight:
			return Input{Action: "right"}
		case tcell.KeyEnter:
			return Input{Action: "confirm"}
		case tcell.KeyRune:
			r := unicode.ToLower(e.Rune())
			// The same physical keys work with a Russian keyboard layout.
			if latin, ok := map[rune]rune{'ц': 'w', 'ы': 's', 'ф': 'a', 'в': 'd', 'р': 'h', 'о': 'j', 'л': 'k', 'у': 'e', 'й': 'q', 'т': 'n', 'с': 'c', 'к': 'r', 'ь': 'm', 'ю': '.'}[r]; ok {
				r = latin
			}
			if r >= '0' && r <= '9' {
				return Input{Action: "select", Choice: int(r - '0')}
			}
			actions := map[rune]string{'w': "up", 's': "down", 'a': "left", 'd': "right", 'h': "weapon", 'j': "food", 'k': "potion", 'e': "scroll", 'q': "quit", 'n': "new", 'c': "continue", 'r': "records", 'm': "menu", '.': "wait", ' ': "wait"}
			return Input{Action: actions[r]}
		}
	}
	return Input{}
}
