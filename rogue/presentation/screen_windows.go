//go:build windows

package presentation

import "github.com/gdamore/tcell/v2"

// Use Windows console events directly, including in a console opened by Explorer.
// This avoids depending on VT input protocol negotiation by the terminal host.
func newScreen() (tcell.Screen, error) {
	return tcell.NewConsoleScreen()
}
