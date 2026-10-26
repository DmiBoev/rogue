package domain

type Record struct {
	ID         string
	FinishedAt string
	Won        bool
	Stats      Statistics
}

type Repository interface {
	Save(*Game) error
	Load() (*Game, error)
	Finish(Record) error
	History() ([]Record, error)
}
