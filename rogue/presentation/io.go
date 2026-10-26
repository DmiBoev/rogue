package presentation

import "github.com/gdamore/tcell/v2"

type Terminal struct{ screen tcell.Screen }

func NewTerminal() (*Terminal, error) {
	screen, err := newScreen()
	if err != nil {
		return nil, err
	}
	if err = screen.Init(); err != nil {
		return nil, err
	}
	screen.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack))
	screen.HideCursor()
	return &Terminal{screen: screen}, nil
}
func (t *Terminal) Close()           { t.screen.Fini() }
func (t *Terminal) Size() (int, int) { return t.screen.Size() }
func (t *Terminal) ReadInput() Input {
	e := t.screen.PollEvent()
	if e == nil {
		return Input{Action: "quit"}
	}
	if _, ok := e.(*tcell.EventResize); ok {
		t.screen.Sync()
	}
	return MapEvent(e)
}
func (t *Terminal) Render(frame Frame) {
	t.screen.Clear()
	w, h := t.Size()
	colors := []tcell.Color{tcell.ColorWhite, tcell.ColorGray, tcell.ColorGreen, tcell.ColorRed, tcell.ColorYellow, tcell.ColorBlue, tcell.ColorAqua}
	for _, c := range frame.Cells {
		if c.X >= 0 && c.Y >= 0 && c.X < w && c.Y < h {
			t.screen.SetContent(c.X, c.Y, c.Symbol, nil, tcell.StyleDefault.Foreground(colors[c.Color]).Background(tcell.ColorBlack))
		}
	}
	t.screen.Show()
}
