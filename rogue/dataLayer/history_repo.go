package datalayer

import (
	"rogue/rogue/domain"
	"sort"
)

func (r *JSONRepository) Finish(record domain.Record) error {
	d, err := r.read()
	if err != nil {
		return err
	}
	found := false
	for _, old := range d.Records {
		if old.ID == record.ID {
			found = true
			break
		}
	}
	if !found {
		d.Records = append(d.Records, record)
	}
	if d.Game != nil && d.Game.ID == record.ID {
		d.Game = nil
	}
	return r.write(d)
}
func (r *JSONRepository) History() ([]domain.Record, error) {
	d, err := r.read()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(d.Records, func(i, j int) bool { return d.Records[i].Stats.Treasure > d.Records[j].Stats.Treasure })
	return d.Records, nil
}
