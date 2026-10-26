//go:build !windows

package presentation

import "github.com/gdamore/tcell/v2"

func newScreen() (tcell.Screen, error) {
	return tcell.NewScreen()
}
