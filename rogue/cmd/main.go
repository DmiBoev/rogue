package main

import (
	"flag"
	"fmt"
	"os"
	datalayer "rogue/rogue/dataLayer"
	"rogue/rogue/presentation"
	"rogue/rogue/presenter"
)

func run() error {
	directory := flag.String("save-dir", "saves", "directory for JSON saves and records")
	seed := flag.Int64("seed", 0, "random seed (0 uses current time)")
	flag.Parse()
	view, err := presentation.NewTerminal()
	if err != nil {
		return fmt.Errorf("open terminal: %w", err)
	}
	defer view.Close()
	return presenter.New(view, datalayer.NewJSONRepository(*directory), *seed).Run()
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
