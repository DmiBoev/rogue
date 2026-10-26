package datalayer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"rogue/rogue/domain"
)

const saveVersion = 1

type JSONRepository struct{ path string }
type document struct {
	Version int             `json:"version"`
	Game    *domain.Game    `json:"game,omitempty"`
	Records []domain.Record `json:"records"`
}

func NewJSONRepository(directory string) *JSONRepository {
	return &JSONRepository{path: filepath.Join(directory, "rogue.json")}
}

func (r *JSONRepository) read() (document, error) {
	d := document{Version: saveVersion}
	f, err := os.Open(r.path)
	if errors.Is(err, os.ErrNotExist) {
		return d, nil
	}
	if err != nil {
		return d, err
	}
	defer f.Close()
	decoder := json.NewDecoder(io.LimitReader(f, 16<<20))
	if err = decoder.Decode(&d); err != nil {
		return d, fmt.Errorf("read save: %w", err)
	}
	var extra any
	if err = decoder.Decode(&extra); err != io.EOF {
		return d, fmt.Errorf("unexpected trailing save data")
	}
	if d.Version != saveVersion {
		return d, fmt.Errorf("unsupported save version: %d", d.Version)
	}
	return d, nil
}

// Current session and records commit together through an atomic replacement.
func (r *JSONRepository) write(d document) error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(r.path), ".rogue-*.json")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(d); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(name, r.path)
}
func (r *JSONRepository) Save(g *domain.Game) error {
	if err := validateGame(g); err != nil {
		return err
	}
	d, err := r.read()
	if err != nil {
		return err
	}
	d.Game = g
	return r.write(d)
}
func (r *JSONRepository) Load() (*domain.Game, error) {
	d, err := r.read()
	if err != nil {
		return nil, err
	}
	if d.Game == nil {
		return nil, os.ErrNotExist
	}
	if err = validateGame(d.Game); err != nil {
		return nil, err
	}
	return d.Game, nil
}
func validateGame(g *domain.Game) error {
	invalid := func() error { return errors.New("invalid game state in save") }
	if g == nil || g.ID == "" || g.Level == nil {
		return invalid()
	}
	l := g.Level
	if l.Number < 1 || l.Number > domain.LastLevel || l.Map.Width != domain.MapWidth || l.Map.Height != domain.MapHeight || len(l.Map.Tiles) != l.Map.Height || len(l.Map.Rooms) != 9 {
		return invalid()
	}
	for _, row := range l.Map.Tiles {
		if len(row) != l.Map.Width {
			return invalid()
		}
	}
	if !l.Walkable(g.Hero.Position) || !l.Walkable(l.Start) || !l.Walkable(l.Exit) || g.Hero.MaxHealth < 1 || g.Hero.Health > g.Hero.MaxHealth || !g.Finished && g.Hero.Health <= 0 {
		return invalid()
	}
	for _, room := range l.Map.Rooms {
		if room.Width < 3 || room.Height < 3 || room.X < 0 || room.Y < 0 || room.X+room.Width > l.Map.Width || room.Y+room.Height > l.Map.Height {
			return invalid()
		}
	}
	for _, e := range l.Enemies {
		if !l.Walkable(e.Position) || e.MaxHealth < 1 {
			return invalid()
		}
	}
	for _, item := range l.Items {
		if !l.Walkable(item.Position) {
			return invalid()
		}
	}
	for _, door := range l.Doors {
		if !l.Walkable(door.Position) || door.Color != 1 && door.Color != 2 && door.Color != 4 {
			return invalid()
		}
	}
	return nil
}
