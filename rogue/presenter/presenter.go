package presenter

import (
	"errors"
	"fmt"
	"os"
	"rogue/rogue/domain"
	"rogue/rogue/presentation"
	"time"
)

type Presenter struct {
	view      presentation.View
	repo      domain.Repository
	game      *domain.Game
	mode      string
	selection domain.ItemType
	records   []domain.Record
	notice    string
	seed      int64
}

func New(view presentation.View, repo domain.Repository, seed int64) *Presenter {
	return &Presenter{view: view, repo: repo, mode: "menu", seed: seed}
}

func (p *Presenter) Run() error {
	for {
		p.view.Render(p.frame())
		input := p.view.ReadInput()
		if input.Action == "quit" {
			if p.game != nil {
				if err := p.persist(); err != nil {
					return fmt.Errorf("save before exit: %w", err)
				}
			}
			return nil
		}
		w, h := p.view.Size()
		if w < 60 || h < 20 {
			continue
		}
		if input.Action == "resize" {
			continue
		}
		if p.mode == "records" {
			if input.Action != "" {
				if p.game == nil {
					p.mode = "menu"
				} else {
					p.mode = "game"
				}
			}
			continue
		}
		if p.selection != "" {
			if input.Action == "cancel" {
				p.selection = ""
				p.mode = "game"
			}
			if input.Action == "select" {
				if p.game.Use(p.selection, input.Choice) {
					p.selection = ""
					p.saveNotice()
				}
			}
			continue
		}
		if input.Action == "records" {
			records, err := p.repo.History()
			if err != nil {
				p.notice = err.Error()
			} else {
				p.records = records
				p.mode = "records"
			}
			continue
		}
		if p.mode == "menu" || p.game != nil && p.game.Finished {
			switch input.Action {
			case "new":
				if p.game != nil && p.game.Finished {
					if err := p.persist(); err != nil {
						p.notice = err.Error()
						continue
					}
				}
				p.startGame()
			case "continue", "confirm", "cancel":
				if p.game != nil && !p.game.Finished {
					p.mode = "game"
					continue
				}
				if input.Action == "cancel" {
					continue
				}
				game, err := p.repo.Load()
				if input.Action == "confirm" && errors.Is(err, os.ErrNotExist) {
					p.startGame()
					continue
				}
				if err != nil {
					p.notice = "Cannot continue: " + err.Error()
				} else {
					p.game = game
					p.mode = "game"
					p.notice = ""
				}
			}
			continue
		}
		if p.game == nil {
			continue
		}
		if input.Action == "menu" {
			p.saveNotice()
			p.mode = "menu"
			continue
		}
		var acted bool
		switch input.Action {
		case "up":
			acted = p.game.Move(domain.DirNorth)
		case "down":
			acted = p.game.Move(domain.DirSouth)
		case "left":
			acted = p.game.Move(domain.DirWest)
		case "right":
			acted = p.game.Move(domain.DirEast)
		case "wait":
			acted = p.game.Move(domain.Direction{})
		case "weapon":
			p.selection = domain.ItemTypeWeapon
		case "food":
			p.selection = domain.ItemTypeFood
		case "potion":
			p.selection = domain.ItemTypePotion
		case "scroll":
			p.selection = domain.ItemTypeScroll
		}
		if acted {
			p.saveNotice()
		}
	}
}

func (p *Presenter) startGame() {
	seed := p.seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	p.seed = 0
	p.game = domain.NewGame(seed)
	p.mode = "game"
	p.notice = ""
	p.saveNotice()
}

func (p *Presenter) persist() error {
	if p.game.Finished {
		return p.repo.Finish(domain.Record{ID: p.game.ID, FinishedAt: time.Now().UTC().Format(time.RFC3339), Won: p.game.Won, Stats: p.game.Stats})
	}
	return p.repo.Save(p.game)
}
func (p *Presenter) saveNotice() {
	if err := p.persist(); err != nil {
		p.notice = "Save failed: " + err.Error()
	} else {
		p.notice = ""
	}
}
