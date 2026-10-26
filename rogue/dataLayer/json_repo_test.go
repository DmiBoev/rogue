package datalayer

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"rogue/rogue/domain"
	"testing"
)

func TestSaveRestoresExactStateAndRandomSequence(t *testing.T) {
	r := NewJSONRepository(t.TempDir())
	g := domain.NewGame(99)
	g.Move(domain.Direction{})
	g.Hero.Effects = []domain.Effect{{Stat: domain.StatStrength, Bonus: 2, Turns: 7}}
	if err := r.Save(g); err != nil {
		t.Fatal(err)
	}
	loaded, err := r.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(g, loaded) {
		t.Fatal("saved state changed after round trip")
	}
	for i := 0; i < 5; i++ {
		g.Move(domain.Direction{})
		loaded.Move(domain.Direction{})
	}
	if !reflect.DeepEqual(g, loaded) {
		t.Fatal("loaded RNG changes future turns")
	}
	// Replacing an existing save must work on Windows too.
	if err = r.Save(loaded); err != nil {
		t.Fatal(err)
	}
}

func TestFinishIsIdempotentAndSortsRecords(t *testing.T) {
	r := NewJSONRepository(t.TempDir())
	g := domain.NewGame(8)
	if err := r.Save(g); err != nil {
		t.Fatal(err)
	}
	for _, record := range []domain.Record{{ID: g.ID, Stats: domain.Statistics{Treasure: 10}}, {ID: "high", Stats: domain.Statistics{Treasure: 100}}, {ID: g.ID, Stats: domain.Statistics{Treasure: 10}}} {
		if err := r.Finish(record); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := r.Load(); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("finished save still resumable")
	}
	records, err := r.History()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 || records[0].Stats.Treasure != 100 {
		t.Fatal("duplicate or unsorted records")
	}
}

func TestCorruptSaveIsRejectedAndPreserved(t *testing.T) {
	dir := t.TempDir()
	r := NewJSONRepository(dir)
	path := filepath.Join(dir, "rogue.json")
	bad := []byte("{broken")
	if err := os.WriteFile(path, bad, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Load(); err == nil {
		t.Fatal("accepted corrupt save")
	}
	if err := r.Save(domain.NewGame(1)); err == nil {
		t.Fatal("silently overwrote corrupt save")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string(bad) {
		t.Fatal("damaged save was overwritten")
	}
}
