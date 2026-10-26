package presentation

type Color int

const (
	White Color = iota
	Gray
	Green
	Red
	Yellow
	Blue
	Cyan
)

type Cell struct {
	X, Y   int
	Symbol rune
	Color  Color
}
type Frame struct{ Cells []Cell }

func (f *Frame) Text(x, y int, s string, color Color) {
	for _, r := range s {
		f.Cells = append(f.Cells, Cell{x, y, r, color})
		x++
	}
}

type View interface {
	Size() (int, int)
	Render(Frame)
	ReadInput() Input
}
